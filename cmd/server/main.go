package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/npavlov/go-password-manager/internal/pkg/logger"
	"github.com/npavlov/go-password-manager/internal/server/buildinfo"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/dbmanager"
	"github.com/npavlov/go-password-manager/internal/server/grpc/auth"
	"github.com/npavlov/go-password-manager/internal/server/storage"
	"github.com/npavlov/go-password-manager/internal/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	ErrDatabaseNotConnected = errors.New("database is not connected")
	ErrJWTisNotPorvided     = errors.New("JWT token is not provided")
)

func main() {
	log := logger.NewLogger(zerolog.DebugLevel).Get()

	log.Info().Str("buildVersion", buildinfo.Version).
		Str("buildCommit", buildinfo.Commit).
		Str("buildDate", buildinfo.Date).Msg("Starting agent")

	cfg := loadConfig(&log)
	var wg sync.WaitGroup

	ctx, cancel := utils.WithSignalCancel(context.Background(), &log)
	defer cancel()

	wg.Add(1)
	starServer(ctx, cfg, &log, &wg)

	fmt.Println("Hello World")

	utils.WaitForShutdown(&wg)
}

func loadConfig(log *zerolog.Logger) *config.Config {
	if err := godotenv.Load("server.env"); err != nil {
		log.Error().Err(err).Msg("Error loading server.env file")
	}

	cfg := config.NewConfigBuilder(log).FromEnv().FromFlags().Build()
	log.Info().Interface("config", cfg).Msg("Configuration loaded")

	if cfg.JwtSecret == "" {
		panic(ErrJWTisNotPorvided)
	}

	return cfg
}

func starServer(ctx context.Context, cfg *config.Config, log *zerolog.Logger, wg *sync.WaitGroup) {

	dbManager := setupDatabase(ctx, cfg, log)
	defer dbManager.Close()

	dbStorage := storage.NewDBStorage(dbManager.DB, log)

	// Create gRPC server
	creds, err := credentials.NewServerTLSFromFile(cfg.Certificate, cfg.PrivateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate credentials")
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		LoggingServerInterceptor(log), // Logs all requests/responses
	), grpc.Creds(creds))
	authService := auth.NewAuthService(log, dbStorage, cfg)
	authService.RegisterService(grpcServer)

	reflection.Register(grpcServer)

	// Start listening
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen on port")
	}

	log.Info().Msgf("gRPC server listening on %s", cfg.Address)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("failed to serve gRPC server")
	}

	//grpcServer := grpc.NewGRPCServer(cfg, log)
	//grpcCon := *grpcServer.GetServer()
	//
	//log.Info().Interface("grpcCon", grpcCon).Msg("Connecting to database")
	//
	//authService := auth.NewAuthService(log, dbStorage, cfg, &grpcCon)
	//
	//authService.RegisterService()
	//
	//grpcServer.Start(ctx, wg)
}

func setupDatabase(ctx context.Context, cfg *config.Config, log *zerolog.Logger) *dbmanager.DBManager {
	dbManager := dbmanager.NewDBManager(cfg.Database, log).Connect(ctx).ApplyMigrations()
	if dbManager.DB == nil {
		log.Fatal().Err(ErrDatabaseNotConnected).Msg("Database is not connected")
	}

	return dbManager
}

// LoggingServerInterceptor logs incoming requests and responses.
func LoggingServerInterceptor(logger *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Start time
		start := time.Now()

		// Log the request details before handling
		logger.Info().
			Str("method", info.FullMethod).
			Interface("request", req).
			Msg("gRPC Request received")

		// Call the actual handler
		resp, err := handler(ctx, req)

		// Calculate the duration
		duration := time.Since(start)

		// Log the response details
		logEvent := logger.Info().
			Str("method", info.FullMethod).
			Dur("duration", duration)

		// Add status code and error details if there's an error
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				logEvent = logEvent.
					Int("status", int(st.Code())).
					Str("error", st.Message())
			} else {
				logEvent = logEvent.
					Int("status", int(codes.Unknown)).
					Str("error", err.Error())
			}
		} else {
			logEvent = logEvent.
				Int("status", int(codes.OK))
		}

		// Log the final message
		logEvent.Msg("gRPC Request completed")

		return resp, err
	}
}
