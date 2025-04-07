package binary

import (
	"context"
	"io"

	pb "github.com/npavlov/go-password-manager/gen/proto/file"
	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// Client handles file (binary) operations over gRPC
type Client struct {
	conn         *grpc.ClientConn
	client       pb.FileServiceClient
	tokenManager *auth.TokenManager
	log          *zerolog.Logger
}

// NewBinaryClient creates a new FileService client
func NewBinaryClient(conn *grpc.ClientConn, tokenManager *auth.TokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		client:       pb.NewFileServiceClient(conn),
		tokenManager: tokenManager,
		log:          log,
	}
}

// UploadFile streams file data to the server
func (c *Client) UploadFile(ctx context.Context, filename string, reader io.Reader) (string, error) {
	stream, err := c.client.UploadFile(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to start upload stream")
	}

	// Send metadata first
	err = stream.Send(&pb.UploadFileRequest{
		Filename: filename,
		Data:     make([]byte, 0),
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to send file metadata")
	}

	// Send file chunks
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return "", errors.Wrap(err, "failed to read from input")
		}

		if n > 0 {
			err := stream.Send(&pb.UploadFileRequest{
				Filename: filename,
				Data:     buf[:n],
			})
			if err != nil && err != io.EOF {
				return "", errors.Wrap(err, "failed to send file chunk")
			}
		}

		if err == io.EOF {
			break
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", errors.Wrap(err, "failed to receive upload response")
	}

	c.log.Info().Str("file_id", resp.FileId).Msg("file uploaded successfully")
	return resp.FileId, nil
}

// DownloadFile streams file data from the server
func (c *Client) DownloadFile(ctx context.Context, fileID string, writer io.Writer) error {
	stream, err := c.client.DownloadFile(ctx, &pb.DownloadFileRequest{
		FileId: fileID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to start download stream")
	}

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to receive file chunk")
		}

		_, err = writer.Write(chunk.Data)
		if err != nil {
			return errors.Wrap(err, "failed to write to output")
		}
	}

	c.log.Info().Str("file_id", fileID).Msg("file downloaded successfully")
	return nil
}

// DeleteFile removes a file by ID
func (c *Client) DeleteFile(ctx context.Context, fileID string) (bool, error) {
	resp, err := c.client.DeleteFile(ctx, &pb.DeleteFileRequest{
		FileId: fileID,
	})
	if err != nil {
		return false, errors.Wrap(err, "failed to delete file")
	}

	return resp.Ok, nil
}

// GetFile retrieves metadata for a file by ID
func (c *Client) GetFile(ctx context.Context, fileID string) (*pb.FileMeta, error) {
	resp, err := c.client.GetFile(ctx, &pb.GetFileRequest{
		FileId: fileID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get file metadata")
	}

	return resp.File, nil
}
