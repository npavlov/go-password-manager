package main

import (
	"context"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/npavlov/go-password-manager/internal/pkg/logger"
	"github.com/npavlov/go-password-manager/internal/server/buildinfo"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/dbmanager"
	"github.com/npavlov/go-password-manager/internal/server/redis"
	"github.com/npavlov/go-password-manager/internal/server/service"
	"github.com/npavlov/go-password-manager/internal/server/service/auth"
	"github.com/npavlov/go-password-manager/internal/server/service/item"
	"github.com/npavlov/go-password-manager/internal/server/service/note"
	"github.com/npavlov/go-password-manager/internal/server/service/password"
	"github.com/npavlov/go-password-manager/internal/server/storage"
	"github.com/npavlov/go-password-manager/internal/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
		Str("buildDate", buildinfo.Date).Msg("Starting server")

	cfg := loadConfig(&log)
	var wg sync.WaitGroup

	ctx, cancel := utils.WithSignalCancel(context.Background(), &log)
	defer cancel()

	dbManager := setupDatabase(ctx, cfg, &log)
	defer dbManager.Close()

	wg.Add(1)
	starServer(ctx, cfg, &log, &wg, dbManager)

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

func starServer(ctx context.Context, cfg *config.Config, log *zerolog.Logger, wg *sync.WaitGroup, dbM *dbmanager.DBManager) {

	dbStorage, memStorage := setupStorage(ctx, cfg, dbM, log)

	grpcManager := service.NewGRPCManager(cfg, log, memStorage)
	grpcServer := grpcManager.GetServer()

	authService := auth.NewAuthService(log, dbStorage, cfg, memStorage)
	authService.RegisterService(grpcServer)

	passwordService := password.NewPasswordService(log, dbStorage, cfg)
	passwordService.RegisterService(grpcServer)

	noteService := note.NewNoteService(log, dbStorage, cfg)
	noteService.RegisterService(grpcServer)

	itemService := item.NewItemService(log, dbStorage, cfg)
	itemService.RegisterService(grpcServer)

	grpcManager.Start(ctx, wg)
}

func setupDatabase(ctx context.Context, cfg *config.Config, log *zerolog.Logger) *dbmanager.DBManager {
	dbManager := dbmanager.NewDBManager(cfg.Database, log).Connect(ctx).ApplyMigrations()
	if dbManager.DB == nil {
		log.Fatal().Err(ErrDatabaseNotConnected).Msg("Database is not connected")
	}

	return dbManager
}

func setupStorage(
	ctx context.Context,
	cfg *config.Config,
	dbManager *dbmanager.DBManager,
	log *zerolog.Logger,
) (*storage.DBStorage, *redis.RStorage) {
	st := storage.NewDBStorage(dbManager.DB, log)
	memStorage := redis.NewRStorage(*cfg, log)

	if err := memStorage.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("Error connecting to redis")
	}

	return st, memStorage
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
