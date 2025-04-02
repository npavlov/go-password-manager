package auth

import (
	"context"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// Client GRPCClient handles communication with the gRPC server
type Client struct {
	conn         *grpc.ClientConn
	client       pb.AuthServiceClient
	tokenManager *auth.TokenManager
	log          *zerolog.Logger
}

// NewAuthClient  creates a new gRPC connection
func NewAuthClient(conn *grpc.ClientConn, tokenManager *auth.TokenManager, log *zerolog.Logger) *Client {
	return &Client{
		conn:         conn,
		client:       pb.NewAuthServiceClient(conn),
		tokenManager: tokenManager,
		log:          log,
	}
}

// Register sends a register request to the server
func (as *Client) Register(username, password, email string) (string, error) {
	resp, err := as.client.Register(context.Background(), &pb.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	})
	if err != nil {
		return "", err
	}

	err = as.tokenManager.UpdateTokens(resp.Token, resp.RefreshToken)
	if err != nil {
		as.log.Error().Err(err).Msg("failed to update tokens")

		return "", err
	}

	return resp.GetUserKey(), nil
}

// Login sends a login request to the server
func (as *Client) Login(username, password string) error {
	resp, err := as.client.Login(context.Background(), &pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}
	err = as.tokenManager.UpdateTokens(resp.Token, resp.RefreshToken)
	if err != nil {
		as.log.Error().Err(err).Msg("failed to update tokens")

		return err
	}

	return nil
}
