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
	Client       pb.NoteServiceClient
	TokenManager auth.ITokenManager
	Log          *zerolog.Logger
}

// NewNoteClient  creates a new gRPC connection.
func NewNoteClient(conn *grpc.ClientConn, tokenManager auth.ITokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		Client:       pb.NewNoteServiceClient(conn),
		TokenManager: tokenManager,
		Log:          log,
	}
}

// GetNote sends a register request to the server.
func (as *Client) GetNote(ctx context.Context, id string) (*pb.NoteData, time.Time, error) {
	resp, err := as.Client.GetNoteV1(ctx, &pb.GetNoteV1Request{
		NoteId: id,
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error getting password")

		return nil, time.Time{}, errors.Wrap(err, "error getting password")
	}

	return resp.GetNote(), resp.GetLastUpdate().AsTime(), nil
}

func (as *Client) StoreNote(ctx context.Context, content string) (string, error) {
	resp, err := as.Client.StoreNoteV1(ctx, &pb.StoreNoteV1Request{
		Note: &pb.NoteData{
			Content: content,
		},
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error storing note")

		return "", errors.Wrap(err, "error storing note")
	}

	return resp.GetNoteId(), nil
}

func (as *Client) DeleteNote(ctx context.Context, id string) (bool, error) {
	resp, err := as.Client.DeleteNoteV1(ctx, &pb.DeleteNoteV1Request{
		NoteId: id,
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error deleting note")

		return false, errors.Wrap(err, "error deleting note")
	}

	return resp.GetOk(), nil
}
