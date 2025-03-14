package auth

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/grpc/utils"
	"github.com/npavlov/go-password-manager/internal/server/storage"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	validator  protovalidate.Validator
	logger     *zerolog.Logger
	storage    *storage.DBStorage
	cfg        *config.Config
	grpcServer *grpc.Server
}

func NewAuthService(log *zerolog.Logger, storage *storage.DBStorage, cfg *config.Config, grpcServer *grpc.Server) *AuthService {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create validator")
	}

	return &AuthService{
		logger:    log,
		validator: validator,
		storage:   storage,
		cfg:       cfg,
	}
}

func (au *AuthService) RegisterService() {
	pb.RegisterAuthServiceServer(au.grpcServer, au)
}

// Register a new user
func (au *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
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

	token, err := utils.GenerateJWT(user.ID.String(), au.cfg.JwtSecret)
	if err != nil {
		au.logger.Error().Err(err).Msg("failed to generate token")

		return nil, errors.Wrap(err, "error generating token")
	}

	return &pb.RegisterResponse{Token: token}, nil
}

// Login user and return JWT token
func (au *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
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

	token, err := utils.GenerateJWT(user.ID.String(), au.cfg.JwtSecret)
	if err != nil {
		au.logger.Error().Err(err).Msg("failed to generate token")

		return nil, errors.Wrap(err, "error generating token")
	}

	return &pb.LoginResponse{Token: token}, nil
}
