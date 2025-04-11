package notes

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/note"
	"github.com/npavlov/go-password-manager/internal/client/auth"
)

// Client GRPCClient handles communication with the gRPC server.
type Client struct {
	conn         *grpc.ClientConn
	client       pb.NoteServiceClient
	tokenManager *auth.TokenManager
	log          *zerolog.Logger
}

// NewNoteClient  creates a new gRPC connection.
func NewNoteClient(conn *grpc.ClientConn, tokenManager *auth.TokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		client:       pb.NewNoteServiceClient(conn),
		tokenManager: tokenManager,
		log:          log,
	}
}

// GetNote sends a register request to the server.
func (as *Client) GetNote(ctx context.Context, id string) (*pb.NoteData, time.Time, error) {
	resp, err := as.client.GetNote(ctx, &pb.GetNoteRequest{
		NoteId: id,
	})
	if err != nil {
		as.log.Error().Err(err).Msg("error getting password")

		return nil, time.Time{}, errors.Wrap(err, "error getting password")
	}

	return resp.GetNote(), resp.GetLastUpdate().AsTime(), nil
}

func (as *Client) StoreNote(ctx context.Context, content string) (string, error) {
	resp, err := as.client.StoreNote(ctx, &pb.StoreNoteRequest{
		Note: &pb.NoteData{
			Content: content,
		},
	})
	if err != nil {
		as.log.Error().Err(err).Msg("error storing note")

		return "", errors.Wrap(err, "error storing note")
	}

	return resp.GetNoteId(), nil
}

func (as *Client) DeleteNote(ctx context.Context, id string) (bool, error) {
	resp, err := as.client.DeleteNote(ctx, &pb.DeleteNoteRequest{
		NoteId: id,
	})
	if err != nil {
		as.log.Error().Err(err).Msg("error deleting note")

		return false, errors.Wrap(err, "error deleting note")
	}

	return resp.GetOk(), nil
}
