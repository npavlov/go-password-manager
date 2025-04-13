package passwords

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/client/auth"
)

// Client GRPCClient handles communication with the gRPC server.
type Client struct {
	conn         *grpc.ClientConn
	Client       pb.PasswordServiceClient
	TokenManager auth.ITokenManager
	Log          *zerolog.Logger
}

// NewPasswordClient  creates a new gRPC connection.
func NewPasswordClient(conn *grpc.ClientConn, tokenManager auth.ITokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		Client:       pb.NewPasswordServiceClient(conn),
		TokenManager: tokenManager,
		Log:          log,
	}
}

// GetPassword sends a register request to the server.
func (as *Client) GetPassword(ctx context.Context, id string) (*pb.PasswordData, time.Time, error) {
	resp, err := as.Client.GetPassword(ctx, &pb.GetPasswordRequest{
		PasswordId: id,
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error getting password")

		return nil, time.Time{}, errors.Wrap(err, "error getting password")
	}

	return resp.GetPassword(), resp.GetLastUpdate().AsTime(), nil
}

func (as *Client) UpdatePassword(ctx context.Context, id, login, password string) error {
	_, err := as.Client.UpdatePassword(ctx, &pb.UpdatePasswordRequest{
		PasswordId: id,
		Data: &pb.PasswordData{
			Login:    login,
			Password: password,
		},
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error updating password")

		return errors.Wrap(err, "error updating password")
	}

	return nil
}

func (as *Client) StorePassword(ctx context.Context, login, password string) (string, error) {
	resp, err := as.Client.StorePassword(ctx, &pb.StorePasswordRequest{
		Password: &pb.PasswordData{
			Login:    login,
			Password: password,
		},
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error storing password")

		return "", errors.Wrap(err, "error storing password")
	}

	return resp.GetPasswordId(), nil
}

func (as *Client) DeletePassword(ctx context.Context, id string) (bool, error) {
	resp, err := as.Client.DeletePassword(ctx, &pb.DeletePasswordRequest{
		PasswordId: id,
	})
	if err != nil {
		as.Log.Error().Err(err).Msg("error deleting password")

		return false, errors.Wrap(err, "error deleting password")
	}

	return resp.GetOk(), nil
}
