package file

import (
	"context"
	"io"

	"github.com/bufbuild/protovalidate-go"
	"github.com/minio/minio-go/v7"
	pb "github.com/npavlov/go-password-manager/gen/proto/file"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	"github.com/npavlov/go-password-manager/internal/server/storage"
	gu "github.com/npavlov/go-password-manager/internal/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	pb.UnimplementedFileServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   *storage.DBStorage
	cfg       *config.Config
	minio     *minio.Client
}

func NewFileService(log *zerolog.Logger, storage *storage.DBStorage, cfg *config.Config, minioClient *minio.Client) *Service {
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

func (fs *Service) UploadFile(stream grpc.ClientStreamingServer[pb.UploadFileRequest, pb.UploadFileResponse]) error {
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

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		fs.logger.Error().Err(err).Msg("error getting user id")

		return errors.Wrap(err, "error getting user id")
	}

	decryptedUserKey, err := utils.GetUserKey(ctx, fs.storage, userUUID, fs.cfg.SecuredMasterKey.Get())
	if err != nil {
		fs.logger.Error().Err(err).Msg("error getting user id")

		return errors.Wrap(err, "error getting user id")
	}

	// Prepare MinIO upload
	objectName := userUUID.String() + "-" + req.Filename
	pipeReader, pipeWriter := io.Pipe()

	// Upload file asynchronously
	go func() {
		uploadCtx := context.Background() // Create a new independent context
		defer pipeReader.Close()          // Ensure reader is closed after upload

		_, err := fs.minio.PutObject(
			uploadCtx, fs.cfg.Bucket, objectName, pipeReader, -1,
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
		if err == io.EOF {
			break
		}
		if err != nil {
			fs.logger.Error().Err(err).Msg("failed to receive file data")

			return errors.Wrap(err, "failed to receive file data")
		}

		if len(chunk.Data) == 0 {
			continue
		}

		blockSize, err := encryptor.Write(chunk.Data)
		totalSize += int64(blockSize)
		if err != nil {
			fs.logger.Error().Err(err).Msg("failed to write chunk")

			return errors.Wrap(err, "failed to write chunk")
		}
	}

	binary, err := fs.storage.StoreBinary(ctx, db.StoreBinaryEntryParams{
		UserID:   userUUID,
		FileName: req.Filename,
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

	return stream.SendAndClose(&pb.UploadFileResponse{FileId: binary.ID.String()})
}

func (fs *Service) GetFiles(ctx context.Context, req *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	if err := fs.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	return &pb.GetFilesResponse{}, nil
}

func (fs *Service) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	if err := fs.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		fs.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	err = fs.storage.DeleteBinary(ctx, db.DeleteBinaryEntryParams{
		ID:     gu.GetIdFromString(req.FileId),
		UserID: userUUID,
	})
	if err != nil {
		fs.logger.Error().Err(err).Msg("error deleting file")

		return nil, errors.Wrap(err, "error deleting file")
	}

	return &pb.DeleteFileResponse{
		Ok: true,
	}, nil
}

func (fs *Service) DownloadFile(req *pb.DownloadFileRequest, stream grpc.ServerStreamingServer[pb.DownloadFileResponse]) error {
	ctx := stream.Context()

	// Get user ID from context
	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		fs.logger.Error().Err(err).Msg("error getting user ID")
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	// Retrieve file metadata from DB
	fileEntry, err := fs.storage.GetBinary(ctx, req.FileId, userUUID)
	if err != nil {
		fs.logger.Error().Err(err).Msg("failed to fetch file metadata")
		return status.Error(codes.NotFound, "file not found")
	}

	// Ensure the file belongs to the user
	if fileEntry.UserID != userUUID {
		fs.logger.Error().Msg("you do not have access to this file")

		return status.Error(codes.PermissionDenied, "you do not have access to this file")
	}

	// Retrieve user encryption key
	decryptedUserKey, err := utils.GetUserKey(ctx, fs.storage, userUUID, fs.cfg.SecuredMasterKey.Get())
	if err != nil {
		fs.logger.Error().Err(err).Msg("failed to retrieve user encryption key")
		return status.Error(codes.Internal, "failed to retrieve encryption key")
	}

	// Fetch encrypted file from MinIO
	objectName := fileEntry.FileUrl
	reader, err := fs.minio.GetObject(ctx, fs.cfg.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		fs.logger.Error().Err(err).Msg("failed to fetch file from MinIO")
		return status.Error(codes.Internal, "failed to retrieve file")
	}
	defer reader.Close()

	// Create a decryptor to decrypt the data on the fly
	decryptor, err := utils.NewDecryptor(reader, decryptedUserKey)
	if err != nil {
		fs.logger.Error().Err(err).Msg("error creating decryptor")

		return status.Error(codes.Internal, "error creating decryptor")
	}

	// Stream decrypted data in blocks
	buffer := make([]byte, 1024) // Read in 1024-byte chunks

	for {
		n, err := decryptor.Read(buffer) // Read a chunk
		if n > 0 {
			// Send only the exact number of bytes read
			if err := stream.Send(&pb.DownloadFileResponse{
				Data:       buffer[:n], // Trim the buffer to actual size
				LastUpdate: timestamppb.New(fileEntry.UpdatedAt.Time),
			}); err != nil {
				return errors.Wrap(err, "error sending download response")
			}
		}

		if n < len(buffer) || err == io.EOF {
			break
		}

		if err != nil {
			fs.logger.Error().Err(err).Msg("error reading and decrypting file")

			return status.Error(codes.Internal, "error reading and decrypting file")
		}
	}

	fs.logger.Info().Str("file", fileEntry.FileName).Msg("file successfully streamed")

	return nil
}

func (fs *Service) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileResponse, error) {
	if err := fs.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		fs.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	file, err := fs.storage.GetBinary(ctx, req.FileId, userUUID)
	if err != nil {
		fs.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	return &pb.GetFileResponse{
		File: &pb.FileMeta{
			Id:       file.ID.String(),
			FileName: file.FileName,
			FileSize: file.FileSize,
			FileUrl:  file.FileUrl,
		},
		LastUpdate: timestamppb.New(file.UpdatedAt.Time),
	}, nil
}
