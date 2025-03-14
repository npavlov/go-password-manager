package auth

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/grpc/utils"
	"github.com/npavlov/go-password-manager/internal/server/redis"
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

func (au *Service) RegisterService(grpcServer *grpc.Server) {
	pb.RegisterAuthServiceServer(grpcServer, au)
}

// Register a new user
func (au *Service) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := au.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "error hashing password")
	}

	user, err := au.storage.RegisterUser(ctx, db.CreateUserParams{
		Username: req.Username,
		Password: string(hashedPassword),
	})

	if err != nil {
		au.logger.Error().Err(err).Msg("failed to register user")

		return nil, errors.Wrap(err, "error creating user")
	}

	au.logger.Info().Interface("user", user).Msg("user created")

	userId := user.ID.String()

	token, err := utils.GenerateJWT(userId, au.cfg.JwtSecret)
	if err != nil {
		au.logger.Error().Err(err).Msg("failed to generate token")

		return nil, errors.Wrap(err, "error generating token")
	}

	err = au.memStorage.Set(ctx, token, userId, utils.TokenExpiration)
	if err != nil {
		au.logger.Error().Err(err).Msg("failed to set token")

		return nil, errors.Wrap(err, "error setting token")
	}

	return &pb.RegisterResponse{Token: token}, nil
}

// Login user and return JWT token
func (au *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := au.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	user, err := au.storage.GetUser(ctx, req.Username)
	if err != nil {
		au.logger.Error().Err(err).Msg("failed to get user")

		return nil, errors.Wrap(err, "error getting user")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		au.logger.Error().Err(err).Msg("invalid password")

		return nil, errors.Wrap(err, "invalid password")
	}

	userId := user.ID.String()

	token, err := utils.GenerateJWT(userId, au.cfg.JwtSecret)
	if err != nil {
		au.logger.Error().Err(err).Msg("failed to generate token")

		return nil, errors.Wrap(err, "error generating token")
	}

	err = au.memStorage.Set(ctx, token, userId, utils.TokenExpiration)
	if err != nil {
		au.logger.Error().Err(err).Msg("failed to set token")

		return nil, errors.Wrap(err, "error setting token")
	}

	return &pb.LoginResponse{Token: token}, nil
}
