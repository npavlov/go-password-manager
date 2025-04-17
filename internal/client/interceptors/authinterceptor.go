//nolint:wrapcheck,ireturn,lll
package interceptors

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/npavlov/go-password-manager/internal/client/config"
)

type Client interface {
	RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest, opts ...grpc.CallOption) (*pb.RefreshTokenResponse, error)
}

type AuthInterceptor struct {
	mu           sync.Mutex
	AuthClient   Client
	config       config.Config
	tokenManager auth.ITokenManager
}

// NewAuthInterceptor initializes the interceptor with tokens.
func NewAuthInterceptor(cfg config.Config, tokenManager auth.ITokenManager) *AuthInterceptor {
	return &AuthInterceptor{
		config:       cfg,
		tokenManager: tokenManager,
		mu:           sync.Mutex{},
		AuthClient:   nil,
	}
}

func (ai *AuthInterceptor) SetAuthClient(client Client) {
	ai.AuthClient = client
}

// UnaryInterceptor checks for authentication errors and refreshes token if necessary.
func (ai *AuthInterceptor) UnaryInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	if method == pb.AuthService_Register_FullMethodName ||
		method == pb.AuthService_Login_FullMethodName ||
		method == pb.AuthService_RefreshToken_FullMethodName {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	// Add access token to request metadata
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", ai.tokenManager.GetAccessToken())

	// Perform the gRPC request
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}

	// Check if the error is due to an expired token
	st, ok := status.FromError(err)
	if ok && st.Code() == codes.Unauthenticated && strings.Contains(st.Message(), "invalid token") {
		// Try refreshing the token
		newToken, newRefresh, err := ai.refreshAccessToken(ctx)
		if err != nil {
			ai.tokenManager.HandleAuthFailure()

			return errors.Wrap(err, "failed to refresh access token")
		}

		// Save new tokens and retry request
		err = ai.tokenManager.UpdateTokens(newToken, newRefresh)
		if err != nil {
			return errors.Wrap(err, "failed to update access token")
		}

		// Retry the original request with the new token
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", newToken)

		return invoker(ctx, method, req, reply, cc, opts...)
	}

	return err
}

// StreamInterceptor attaches the token to streaming RPCs.
//

func (ai *AuthInterceptor) StreamInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	if ai.tokenManager == nil {
		return streamer(ctx, desc, cc, method, opts...)
	}

	token := ai.tokenManager.GetAccessToken()
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "no access token")
	}

	newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", token)

	return streamer(newCtx, desc, cc, method, opts...)
}

// refreshAccessToken calls the auth service to get a new token.
func (ai *AuthInterceptor) refreshAccessToken(ctx context.Context) (string, string, error) {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	// Call the RefreshToken gRPC endpoint
	resp, err := ai.AuthClient.RefreshToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: ai.tokenManager.GetRefreshToken(),
	})
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get new access token")
	}

	return resp.GetToken(), resp.GetRefreshToken(), nil
}
