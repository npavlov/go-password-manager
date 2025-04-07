package passwords

import (
	"context"
	"time"

	pb "github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// Client GRPCClient handles communication with the gRPC server
type Client struct {
	conn         *grpc.ClientConn
	client       pb.PasswordServiceClient
	tokenManager *auth.TokenManager
	log          *zerolog.Logger
}

// NewPasswordClient  creates a new gRPC connection
func NewPasswordClient(conn *grpc.ClientConn, tokenManager *auth.TokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		client:       pb.NewPasswordServiceClient(conn),
		tokenManager: tokenManager,
		log:          log,
	}
}

// GetPassword sends a register request to the server
func (as *Client) GetPassword(ctx context.Context, id string) (*pb.PasswordData, time.Time, error) {
	resp, err := as.client.GetPassword(ctx, &pb.GetPasswordRequest{
		PasswordId: id,
	})

	if err != nil {
		as.log.Error().Err(err).Msg("error getting password")

		return nil, time.Time{}, errors.Wrap(err, "error getting password")
	}

	return resp.Password, resp.LastUpdate.AsTime(), nil
}

func (as *Client) UpdatePassword(ctx context.Context, id, login, password string) error {
	_, err := as.client.UpdatePassword(ctx, &pb.UpdatePasswordRequest{
		PasswordId: id,
		Data: &pb.PasswordData{
			Login:    login,
			Password: password,
		},
	})

	if err != nil {
		as.log.Error().Err(err).Msg("error updating password")

		return errors.Wrap(err, "error updating password")
	}

	return nil
}

func (as *Client) StorePassword(ctx context.Context, login, password string) (string, error) {
	resp, err := as.client.StorePassword(ctx, &pb.StorePasswordRequest{
		Password: &pb.PasswordData{
			Login:    login,
			Password: password,
		},
	})

	if err != nil {
		as.log.Error().Err(err).Msg("error storing password")

		return "", errors.Wrap(err, "error storing password")
	}

	return resp.PasswordId, nil
}

func (as *Client) DeletePassword(ctx context.Context, id string) (bool, error) {
	resp, err := as.client.DeletePassword(ctx, &pb.DeletePasswordRequest{
		PasswordId: id,
	})

	if err != nil {
		as.log.Error().Err(err).Msg("error deleting password")

		return false, errors.Wrap(err, "error deleting password")
	}

	return resp.Ok, nil
}
