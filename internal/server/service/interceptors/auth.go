package interceptors

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/npavlov/go-password-manager/internal/server/redis"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
)

// TokenInterceptor extracts a token from metadata and injects it into the context.
func TokenInterceptor(log *zerolog.Logger, jwtSecret string, memSt redis.MemStorage) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// List of methods that should bypass token verification
		skipMethods := map[string]bool{
			pb.AuthService_Register_FullMethodName:     true,
			pb.AuthService_Login_FullMethodName:        true,
			pb.AuthService_RefreshToken_FullMethodName: true,
		}

		// Skip authentication for specified methods
		if skipMethods[info.FullMethod] {
			log.Info().Str("method", info.FullMethod).Msg("skipping authentication for public method")

			return handler(ctx, req)
		}

		userID, err := AuthenticateToken(ctx, jwtSecret, memSt)
		if err != nil {
			log.Info().Str("method", info.FullMethod).Msg("authentication failed")

			return nil, errors.Wrap(err, "authenticating token")
		}

		//nolint:revive,staticcheck
		ctx = context.WithValue(ctx, "user_id", userID)

		log.Info().Str("method", info.FullMethod).Msg("user_id extracted and added to context")

		return handler(ctx, req)
	}
}

type wrappedStream struct {
	grpc.ServerStream
	//nolint:containedctx
	ctx context.Context
}

// Context Override Context() to return the modified context.
func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func StreamTokenInterceptor(logger *zerolog.Logger,
	jwtSecret string,
	memStorage redis.MemStorage,
) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Skip token authentication if the request is for reflection
		if info.FullMethod == "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo" {
			return handler(srv, stream)
		}

		// Extract context from stream
		ctx := stream.Context()

		// Authenticate token
		userID, err := AuthenticateToken(ctx, jwtSecret, memStorage)
		if err != nil {
			logger.Error().Err(err).Msg("Unauthorized stream request")

			return status.Error(codes.Unauthenticated, "invalid token")
		}

		// Add user ID to context
		//nolint:revive,staticcheck
		ctx = context.WithValue(ctx, "user_id", userID)

		// Wrap the original stream with the new context
		wrapped := &wrappedStream{ServerStream: stream, ctx: ctx}

		// Pass the modified stream to the handler
		return handler(srv, wrapped)
	}
}

func AuthenticateToken(ctx context.Context, jwtSecret string, memStorage redis.MemStorage) (string, error) {
	// Extract token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "missing authorization token")
	}

	tokenString := tokens[0]

	userID, err := utils.ValidateJWT(tokenString, jwtSecret)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "invalid token")
	}

	// Check if the token exists in Redis and match with User ID
	result, err := memStorage.Get(ctx, tokenString)
	if result != userID || err != nil {
		return "", status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return userID, nil
}
