package auth

import (
	"context"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/redis"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	"github.com/npavlov/go-password-manager/internal/server/storage"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

const (
	TokenExpiration        = time.Minute * 60
	RefreshTokenExpiration = time.Hour * 24 * 7 // 7 days
)

type Service struct {
	pb.UnimplementedAuthServiceServer
	validator  protovalidate.Validator
	logger     *zerolog.Logger
	storage    *storage.DBStorage
	cfg        *config.Config
	memStorage redis.MemStorage
}

func NewAuthService(log *zerolog.Logger, storage *storage.DBStorage, cfg *config.Config, memStorage redis.MemStorage) *Service {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create validator")
	}

	return &Service{
		logger:     log,
		validator:  validator,
		storage:    storage,
		cfg:        cfg,
		memStorage: memStorage,
	}
}

func (as *Service) RegisterService(grpcServer *grpc.Server) {
	pb.RegisterAuthServiceServer(grpcServer, as)
}

// Register a new user
func (as *Service) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := as.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "error hashing password")
	}

	// Generate a unique encryption key for this user
	userKey, err := utils.GenerateRandomKey()
	if err != nil {
		return nil, err
	}
	// Encrypt the user key using the master key
	encryptedKey, err := utils.Encrypt(userKey, as.cfg.SecuredMasterKey.Get())
	if err != nil {
		return nil, err
	}

	user, err := as.storage.RegisterUser(ctx, db.CreateUserParams{
		Username:      req.Username,
		Password:      string(hashedPassword),
		EncryptionKey: encryptedKey,
		Email:         req.Email,
	})

	if err != nil {
		as.logger.Error().Err(err).Msg("failed to register user")

		return nil, errors.Wrap(err, "error creating user")
	}

	as.logger.Info().Interface("user", user).Msg("user created")

	token, refreshToken, err := as.tokenGeneration(ctx, user.ID)
	if err != nil {
		as.logger.Error().Err(err).Msg("error generating token")

		return nil, errors.Wrap(err, "error generating token")
	}

	return &pb.RegisterResponse{Token: token, RefreshToken: refreshToken, UserKey: userKey}, nil
}

// Login user and return JWT token
func (as *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := as.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	user, err := as.storage.GetUser(ctx, req.Username)
	if err != nil {
		as.logger.Error().Err(err).Msg("failed to get user")

		return nil, errors.Wrap(err, "error getting user")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		as.logger.Error().Err(err).Msg("invalid password")

		return nil, errors.Wrap(err, "invalid password")
	}

	token, refreshToken, err := as.tokenGeneration(ctx, user.ID)
	if err != nil {
		as.logger.Error().Err(err).Msg("error generating token")

		return nil, errors.Wrap(err, "error generating token")
	}

	return &pb.LoginResponse{Token: token, RefreshToken: refreshToken}, nil
}

func (as *Service) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	if err := as.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	// Get refresh token from DB
	tokenRow, err := as.storage.GetToken(ctx, req.RefreshToken)
	if err != nil {
		as.logger.Error().Err(err).Msg("failed to get refresh token")
		return nil, errors.Wrap(err, "invalid refresh token")
	}

	// Check if refresh token is expired
	if tokenRow.ExpiresAt.Time.Before(time.Now()) {
		as.logger.Error().Msg("refresh token expired")

		return nil, errors.New("refresh token expired")
	}

	// Generate a new access token
	newToken, newRefreshToken, err := as.tokenGeneration(ctx, tokenRow.UserID)

	as.logger.Info().Str("user_id", tokenRow.UserID.String()).Msg("refresh token successfully rotated")

	return &pb.RefreshTokenResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (as *Service) tokenGeneration(ctx context.Context, userID pgtype.UUID) (string, string, error) {
	userId := userID.String()

	tokenExp := time.Now().Add(TokenExpiration).Unix()

	token, err := utils.GenerateJWT(userId, as.cfg.JwtSecret, tokenExp)
	if err != nil {
		return "", "", errors.Wrap(err, "error generating token")
	}

	err = as.memStorage.Set(ctx, token, userId, TokenExpiration)
	if err != nil {
		as.logger.Error().Err(err).Msg("failed to set token")

		return "", "", errors.Wrap(err, "error setting token")
	}

	refreshTokenExp := time.Now().Add(RefreshTokenExpiration)

	refreshToken, err := utils.GenerateJWT(userId, as.cfg.JwtSecret, refreshTokenExp.Unix())
	if err != nil {
		return "", "", errors.Wrap(err, "error generating refresh token")
	}

	err = as.storage.StoreToken(ctx, userID, refreshToken, refreshTokenExp)
	if err != nil {
		return "", "", errors.Wrap(err, "error storing token")
	}

	return token, refreshToken, nil
}
