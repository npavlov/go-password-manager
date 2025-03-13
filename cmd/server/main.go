package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/joho/godotenv"
	"github.com/npavlov/go-password-manager/internal/pkg/logger"
	"github.com/npavlov/go-password-manager/internal/server/buildinfo"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/grpc"
	"github.com/npavlov/go-password-manager/internal/utils"
	"github.com/rs/zerolog"
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

	cfg := config.NewConfigBuilder(log).
		FromEnv().
		FromFlags().
		Build()

	log.Info().Interface("config", cfg).Msg("Configuration loaded")

	return cfg
}

func starServer(ctx context.Context, cfg *config.Config, log *zerolog.Logger, wg *sync.WaitGroup) {
	grpcServer := grpc.NewGRPCServer(cfg, log)
	grpcServer.Start(ctx, wg)
}
