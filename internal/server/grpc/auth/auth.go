package services

import (
	"context"
	"fmt"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

// JWT secret key (should be in env variables)
var jwtSecret = []byte("supersecretkey")

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
}

func NewAuthService(log *zerolog.Logger) *AuthService {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create validator")
	}

	return &AuthService{
		logger: log,
		validator: validator
	}
}

// Register a new user
func (au *AuthService) Register(_ context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := au.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "error hashing password")
	}

	s.

	return &pb.RegisterResponse{Message: "User registered successfully"}, nil
}

// Login user and return JWT token
func (au *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var userID, hashedPassword string

	query := `SELECT id, password_hash FROM users WHERE username=$1`
	err := s.db.Conn.QueryRow(ctx, query, req.Username).Scan(&userID, &hashedPassword)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &pb.LoginResponse{Token: tokenString}, nil
}
