package auth

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
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
	encryptedKey, err := utils.Encrypt(userKey, as.cfg.MasterKey)
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

	userId := user.ID.String()

	token, err := utils.GenerateJWT(userId, as.cfg.JwtSecret)
	if err != nil {
		as.logger.Error().Err(err).Msg("failed to generate token")

		return nil, errors.Wrap(err, "error generating token")
	}

	err = as.memStorage.Set(ctx, token, userId, utils.TokenExpiration)
	if err != nil {
		as.logger.Error().Err(err).Msg("failed to set token")

		return nil, errors.Wrap(err, "error setting token")
	}

	return &pb.RegisterResponse{Token: token, UserKey: userKey}, nil
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

	userId := user.ID.String()

	token, err := utils.GenerateJWT(userId, as.cfg.JwtSecret)
	if err != nil {
		as.logger.Error().Err(err).Msg("failed to generate token")

		return nil, errors.Wrap(err, "error generating token")
	}

	err = as.memStorage.Set(ctx, token, userId, utils.TokenExpiration)
	if err != nil {
		as.logger.Error().Err(err).Msg("failed to set token")

		return nil, errors.Wrap(err, "error setting token")
	}

	return &pb.LoginResponse{Token: token}, nil
}
