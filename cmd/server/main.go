package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/joho/godotenv"
	"github.com/npavlov/go-password-manager/internal/pkg/logger"
	"github.com/npavlov/go-password-manager/internal/server/buildinfo"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/dbmanager"
	"github.com/npavlov/go-password-manager/internal/server/grpc"
	"github.com/npavlov/go-password-manager/internal/server/grpc/auth"
	"github.com/npavlov/go-password-manager/internal/server/storage"
	"github.com/npavlov/go-password-manager/internal/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
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
	grpcServer := grpc.NewGRPCServer(cfg, log)

	dbManager := setupDatabase(ctx, cfg, log)
	defer dbManager.Close()

	dbStorage := storage.NewDBStorage(dbManager.DB, log)

	auth.NewAuthService(log, dbStorage, cfg, grpcServer.GetServer())

	grpcServer.Start(ctx, wg)
}

func setupDatabase(ctx context.Context, cfg *config.Config, log *zerolog.Logger) *dbmanager.DBManager {
	dbManager := dbmanager.NewDBManager(cfg.Database, log).Connect(ctx).ApplyMigrations()
	if dbManager.DB == nil {
		log.Fatal().Err(ErrDatabaseNotConnected).Msg("Database is not connected")
	}

	return dbManager
}
