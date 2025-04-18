//nolint:wrapcheck,exhaustruct,revive
package file

import (
	"context"
	"io"

	"github.com/bufbuild/protovalidate-go"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/npavlov/go-password-manager/gen/proto/file"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	gu "github.com/npavlov/go-password-manager/internal/utils"
)

const (
	chunkSize = 1024
)

type Storage interface {
	StoreBinary(ctx context.Context, createBinary db.StoreBinaryEntryParams) (*db.BinaryEntry, error)
	DeleteBinary(ctx context.Context, arg db.DeleteBinaryEntryParams) error
	GetBinaries(ctx context.Context, userID string) ([]db.BinaryEntry, error)
	GetBinary(ctx context.Context, binaryID string, userID pgtype.UUID) (*db.BinaryEntry, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (*db.User, error)
}

type S3Storage interface {
	PutObject(ctx context.Context,
		bucketName string,
		objectName string,
		reader io.Reader,
		objectSize int64,
		opts minio.PutObjectOptions,
	) (info minio.UploadInfo, err error)
	GetObject(ctx context.Context, bucketName string, objName string, opts minio.GetObjectOptions) (io.ReadCloser, error)
	RemoveObject(ctx context.Context, bucketName string, objName string, opts minio.RemoveObjectOptions) error
}

type Service struct {
	pb.UnimplementedFileServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   Storage
	cfg       *config.Config
	minio     S3Storage
}

func NewFileService(log *zerolog.Logger, storage Storage, cfg *config.Config, minioClient S3Storage) *Service {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create validator")
	}

	return &Service{
		logger:    log,
		validator: validator,
		storage:   storage,
		cfg:       cfg,
		minio:     minioClient,
	}
}

func (fs *Service) RegisterService(grpcServer *grpc.Server) {
	pb.RegisterFileServiceServer(grpcServer, fs)
}

//nolint:cyclop,funlen,lll
func (fs *Service) UploadFileV1(stream grpc.ClientStreamingServer[pb.UploadFileV1Request, pb.UploadFileV1Response]) error {
	ctx := stream.Context()

	// First, receive metadata (filename)
	req, err := stream.Recv()
	if err != nil {
		fs.logger.Error().Err(err).Msg("failed to receive file metadata")

		return errors.Wrap(err, "failed to receive file metadata")
	}

	if err := fs.validator.Validate(req); err != nil {
		return errors.Wrap(err, "failed to validate file metadata")
	}

	userUUID, decryptedUserKey, err := utils.GetDecryptionKey(ctx, fs.storage, fs.cfg.SecuredMasterKey.Get())
	if err != nil {
		fs.logger.Error().Err(err).Msg("error getting user id")

		return errors.Wrap(err, "error getting user id")
	}

	// Prepare MinIO upload
	objectName := userUUID.String() + "-" + req.GetFilename()
	pipeReader, pipeWriter := io.Pipe()

	// Upload file asynchronously
	go func() {
		uploadCtx := context.Background() // Create a new independent context
		defer pipeReader.Close()          // Ensure reader is closed after upload

		_, err := fs.minio.PutObject(
			uploadCtx, fs.cfg.Bucket, objectName, pipeReader, -1,
			//nolint:exhaustruct
			minio.PutObjectOptions{ContentType: "application/octet-stream"},
		)
		if err != nil {
			fs.logger.Error().Err(err).Msg("Failed to upload to MinIO")
		}
		fs.logger.Info().Msg("Successfully uploaded file to MinIO")
	}()

	// Encrypt and write data to the pipe
	encryptor, err := utils.NewEncryptor(pipeWriter, decryptedUserKey)
	if err != nil {
		fs.logger.Error().Err(err).Msg("error creating encryptor")

		return errors.Wrap(err, "error creating encryptor")
	}
	defer pipeWriter.Close()

	// Stream and encrypt file data
	var totalSize int64 = 0
	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fs.logger.Error().Err(err).Msg("failed to receive file data")

			return errors.Wrap(err, "failed to receive file data")
		}

		if len(chunk.GetData()) == 0 {
			continue
		}

		blockSize, err := encryptor.Write(chunk.GetData())
		totalSize += int64(blockSize)
		if err != nil {
			fs.logger.Error().Err(err).Msg("failed to write chunk")

			return errors.Wrap(err, "failed to write chunk")
		}
	}

	binary, err := fs.storage.StoreBinary(ctx, db.StoreBinaryEntryParams{
		UserID:   userUUID,
		FileName: req.GetFilename(),
		FileSize: totalSize,
		FileUrl:  objectName,
	})
	if err != nil {
		fs.logger.Error().Err(err).Msg("failed to store binary")

		_ = fs.minio.RemoveObject(ctx, fs.cfg.Bucket, objectName, minio.RemoveObjectOptions{
			ForceDelete: true,
		})

		return errors.Wrap(err, "failed to store binary")
	}

	return stream.SendAndClose(&pb.UploadFileV1Response{FileId: binary.ID.String()})
}

func (fs *Service) GetFilesV1(ctx context.Context, req *pb.GetFilesV1Request) (*pb.GetFilesV1Response, error) {
	if err := fs.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	return &pb.GetFilesV1Response{}, nil
}

func (fs *Service) DeleteFileV1(ctx context.Context, req *pb.DeleteFileV1Request) (*pb.DeleteFileV1Response, error) {
	if err := fs.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserID(ctx)
	if err != nil {
		fs.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	err = fs.storage.DeleteBinary(ctx, db.DeleteBinaryEntryParams{
		ID:     gu.GetIDFromString(req.GetFileId()),
		UserID: userUUID,
	})
	if err != nil {
		fs.logger.Error().Err(err).Msg("error deleting file")

		return nil, errors.Wrap(err, "error deleting file")
	}

	return &pb.DeleteFileV1Response{
		Ok: true,
	}, nil
}

//nolint:cyclop,funlen,lll
func (fs *Service) DownloadFileV1(req *pb.DownloadFileV1Request, str grpc.ServerStreamingServer[pb.DownloadFileV1Response]) error {
	ctx := str.Context()

	userUUID, decryptedUserKey, err := utils.GetDecryptionKey(ctx, fs.storage, fs.cfg.SecuredMasterKey.Get())
	if err != nil {
		fs.logger.Error().Err(err).Msg("error getting user id")

		return errors.Wrap(err, "error getting user id")
	}

	// Retrieve file metadata from DB
	fileEntry, err := fs.storage.GetBinary(ctx, req.GetFileId(), userUUID)
	if err != nil {
		fs.logger.Error().Err(err).Msg("failed to fetch file metadata")

		//nolint:wrapcheck
		return status.Error(codes.NotFound, "file not found")
	}

	// Ensure the file belongs to the user
	if fileEntry.UserID != userUUID {
		fs.logger.Error().Msg("you do not have access to this file")

		//nolint:wrapcheck
		return status.Error(codes.PermissionDenied, "you do not have access to this file")
	}

	// Fetch encrypted file from MinIO
	objectName := fileEntry.FileUrl
	//nolint:exhaustruct
	reader, err := fs.minio.GetObject(ctx, fs.cfg.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		fs.logger.Error().Err(err).Msg("failed to fetch file from MinIO")

		//nolint:wrapcheck
		return status.Error(codes.Internal, "failed to retrieve file")
	}
	defer reader.Close()

	// Create a decryptor to decrypt the data on the fly
	decryptor, err := utils.NewDecryptor(reader, decryptedUserKey)
	if err != nil {
		fs.logger.Error().Err(err).Msg("error creating decryptor")

		//nolint:wrapcheck
		return status.Error(codes.Internal, "error creating decryptor")
	}

	// Stream decrypted data in blocks
	buffer := make([]byte, chunkSize) // Read in 1024-byte chunks

	for {
		cursor, err := decryptor.Read(buffer) // Read a chunk
		if cursor > 0 {
			// Send only the exact number of bytes read
			if err := str.Send(&pb.DownloadFileV1Response{
				Data:       buffer[:cursor], // Trim the buffer to actual size
				LastUpdate: timestamppb.New(fileEntry.UpdatedAt.Time),
			}); err != nil {
				return errors.Wrap(err, "error sending download response")
			}
		}
		if err != nil && !errors.Is(err, io.EOF) {
			fs.logger.Error().Err(err).Msg("error reading and decrypting file")

			return status.Error(codes.Internal, "error reading and decrypting file")
		}

		if cursor < len(buffer) || errors.Is(err, io.EOF) {
			break
		}
	}

	fs.logger.Info().Str("file", fileEntry.FileName).Msg("file successfully streamed")

	return nil
}

func (fs *Service) GetFileV1(ctx context.Context, req *pb.GetFileV1Request) (*pb.GetFileV1Response, error) {
	if err := fs.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserID(ctx)
	if err != nil {
		fs.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	file, err := fs.storage.GetBinary(ctx, req.GetFileId(), userUUID)
	if err != nil {
		fs.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	return &pb.GetFileV1Response{
		File: &pb.FileMeta{
			Id:       file.ID.String(),
			FileName: file.FileName,
			FileSize: file.FileSize,
			FileUrl:  file.FileUrl,
		},
		LastUpdate: timestamppb.New(file.UpdatedAt.Time),
	}, nil
}
