package auth

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/npavlov/go-password-manager/internal/client/auth"
)

// Client GRPCClient handles communication with the gRPC server.
type Client struct {
	conn         *grpc.ClientConn
	Client       pb.AuthServiceClient
	TokenManager auth.ITokenManager
	Log          *zerolog.Logger
}

// NewAuthClient  creates a new gRPC connection.
func NewAuthClient(conn *grpc.ClientConn, tokenManager auth.ITokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		Client:       pb.NewAuthServiceClient(conn),
		TokenManager: tokenManager,
		Log:          log,
	}
}

// Register sends a register request to the server.
func (as *Client) Register(username, password, email string) (string, error) {
	resp, err := as.Client.Register(context.Background(), &pb.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to register user")
	}

	err = as.TokenManager.UpdateTokens(resp.GetToken(), resp.GetRefreshToken())
	if err != nil {
		as.Log.Error().Err(err).Msg("failed to update tokens")

		return "", errors.New("failed to update tokens")
	}

	return resp.GetUserKey(), nil
}

// Login sends a login request to the server.
func (as *Client) Login(username, password string) error {
	resp, err := as.Client.Login(context.Background(), &pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		as.TokenManager.HandleAuthFailure()

		return errors.Wrap(err, "failed to login")
	}
	err = as.TokenManager.UpdateTokens(resp.GetToken(), resp.GetRefreshToken())
	if err != nil {
		as.Log.Error().Err(err).Msg("failed to update tokens")

		return errors.Wrap(err, "failed to update tokens")
	}

	return nil
}
