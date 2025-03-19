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

		// Process the chunk in 4096-byte blocks
		for offset := 0; offset < len(chunk.Data); offset += 4096 {
			end := offset + 4096
			if end > len(chunk.Data) {
				end = len(chunk.Data)
			}

			// Encrypt the 4096-byte block
			encryptedBlock, err := utils.Encrypt(string(chunk.Data[offset:end]), decryptedUserKey)
			if err != nil {
				fs.logger.Error().Err(err).Msg("failed to encrypt file block")
				return errors.Wrap(err, "failed to encrypt file block")
			}

			// Write the encrypted block to MinIO
			_, err = pipeWriter.Write([]byte(encryptedBlock))
			if err != nil {
				fs.logger.Error().Err(err).Msg("failed to write to MinIO")
				return errors.Wrap(err, "failed to write to MinIO")
			}

			totalSize += int64(len(encryptedBlock))
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
	fileEntry, err := fs.storage.GetBinary(ctx, req.FileId)
	if err != nil {
		fs.logger.Error().Err(err).Msg("failed to fetch file metadata")
		return status.Error(codes.NotFound, "file not found")
	}

	// Ensure the file belongs to the user
	if fileEntry.UserID != userUUID {
		fs.logger.Warn().Msg("unauthorized file access attempt")
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

	// Read and stream the file in chunks
	buf := make([]byte, 4096) // 4KB buffer for chunked streaming
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			fs.logger.Error().Err(err).Msg("failed to read file chunk")
			return status.Error(codes.Internal, "failed to read file")
		}

		if n == 0 {
			break
		}

		// Decrypt the chunk
		decryptedChunk, err := utils.Decrypt(string(buf[:n]), decryptedUserKey)
		if err != nil {
			fs.logger.Error().Err(err).Msg("failed to decrypt file chunk")
			return status.Error(codes.Internal, "failed to decrypt file")
		}

		// Send chunk to client
		err = stream.Send(&pb.DownloadFileResponse{Data: []byte(decryptedChunk)})
		if err != nil {
			fs.logger.Error().Err(err).Msg("failed to send file chunk")
			return status.Error(codes.Canceled, "stream interrupted")
		}
	}

	fs.logger.Info().Str("file", fileEntry.FileName).Msg("file successfully streamed")

	return nil
}
