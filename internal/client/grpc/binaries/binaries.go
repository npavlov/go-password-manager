package binary

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/file"
	"github.com/npavlov/go-password-manager/internal/client/auth"
)

const (
	chunkSize = 1024
)

// Client handles file (binary) operations over gRPC.
type Client struct {
	conn         *grpc.ClientConn
	Client       pb.FileServiceClient
	TokenManager auth.ITokenManager
	Log          *zerolog.Logger
}

// NewBinaryClient creates a new FileService client.
func NewBinaryClient(conn *grpc.ClientConn, tokenManager auth.ITokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		Client:       pb.NewFileServiceClient(conn),
		TokenManager: tokenManager,
		Log:          log,
	}
}

// UploadFile streams file data to the server.
//
//nolint:cyclop
func (c *Client) UploadFile(ctx context.Context, filename string, reader io.Reader) (string, error) {
	stream, err := c.Client.UploadFileV1(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to start upload stream")
	}

	// Send metadata first
	err = stream.Send(&pb.UploadFileV1Request{
		Filename: filename,
		Data:     make([]byte, 0),
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to send file metadata")
	}

	// Send file chunks
	buf := make([]byte, chunkSize)
	for {
		cursor, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return "", errors.Wrap(err, "failed to read from input")
		}

		if cursor > 0 {
			err := stream.Send(&pb.UploadFileV1Request{
				Filename: filename,
				Data:     buf[:cursor],
			})
			if err != nil && !errors.Is(err, io.EOF) {
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

	c.Log.Info().Str("file_id", resp.GetFileId()).Msg("file uploaded successfully")

	return resp.GetFileId(), nil
}

// DownloadFile streams file data from the server.
func (c *Client) DownloadFile(ctx context.Context, fileID string, writer io.Writer) error {
	stream, err := c.Client.DownloadFileV1(ctx, &pb.DownloadFileV1Request{
		FileId: fileID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to start download stream")
	}

	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to receive file chunk")
		}

		_, err = writer.Write(chunk.GetData())
		if err != nil {
			return errors.Wrap(err, "failed to write to output")
		}
	}

	c.Log.Info().Str("file_id", fileID).Msg("file downloaded successfully")

	return nil
}

// DeleteFile removes a file by ID.
func (c *Client) DeleteFile(ctx context.Context, fileID string) (bool, error) {
	resp, err := c.Client.DeleteFileV1(ctx, &pb.DeleteFileV1Request{
		FileId: fileID,
	})
	if err != nil {
		return false, errors.Wrap(err, "failed to delete file")
	}

	return resp.GetOk(), nil
}

// GetFile retrieves metadata for a file by ID.
func (c *Client) GetFile(ctx context.Context, fileID string) (*pb.FileMeta, error) {
	resp, err := c.Client.GetFileV1(ctx, &pb.GetFileV1Request{
		FileId: fileID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get file metadata")
	}

	return resp.GetFile(), nil
}
