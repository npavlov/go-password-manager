package interceptors

import (
	"context"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/npavlov/go-password-manager/internal/server/redis"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TokenInterceptor extracts a token from metadata and injects it into the context.
func TokenInterceptor(logger *zerolog.Logger, jwtSecret string, memStorage redis.MemStorage) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// List of methods that should bypass token verification
		skipMethods := map[string]bool{
			pb.AuthService_Register_FullMethodName: true,
			pb.AuthService_Login_FullMethodName:    true,
		}

		// Skip authentication for specified methods
		if skipMethods[info.FullMethod] {
			logger.Info().Str("method", info.FullMethod).Msg("skipping authentication for public method")

			return handler(ctx, req)
		}

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Error().Msg("missing metadata")
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			logger.Error().Msg("missing authorization token")

			return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
		}

		tokenString := tokens[0]

		userID, err := utils.ValidateJWT(tokenString, jwtSecret)
		if err != nil {
			logger.Error().Msg("invalid token")

			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		// Check if the token exists in Redis and match with User ID
		result, err := memStorage.Get(ctx, tokenString)
		if result != userID || err != nil {
			logger.Error().Msg("invalid token")

			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, "user_id", userID)

		logger.Info().Str("method", info.FullMethod).Msg("user_id extracted and added to context")

		return handler(ctx, req)
	}
}
