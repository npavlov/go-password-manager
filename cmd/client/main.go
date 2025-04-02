package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/npavlov/go-password-manager/internal/client/buildinfo"
	"github.com/npavlov/go-password-manager/internal/client/config"
	"github.com/npavlov/go-password-manager/internal/client/grpc/facade"
	"github.com/npavlov/go-password-manager/internal/client/interceptors"
	"github.com/npavlov/go-password-manager/internal/client/storage"
	"github.com/npavlov/go-password-manager/internal/client/tui"
	"github.com/npavlov/go-password-manager/internal/pkg/logger"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	log := logger.NewLogger(zerolog.DebugLevel).Get()

	log.Info().Str("buildVersion", buildinfo.Version).
		Str("buildCommit", buildinfo.Commit).
		Str("buildDate", buildinfo.Date).Msg("Starting agent")

	cfg := loadConfig(&log)

	tokenManager := auth.NewTokenManager(&log, cfg)
	err := tokenManager.LoadTokens()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load tokens")
	}

	// Create the authInterceptor
	authInterceptor := interceptors.NewAuthInterceptor(*cfg, tokenManager)

	// Initialize gRPC conn
	conn, err := makeConnection(*cfg, authInterceptor)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to make connection")
	}
	defer conn.Close()

	facadeClient := facade.NewFacade(conn, tokenManager, &log)

	storageManager := storage.NewStorageManager(facadeClient, tokenManager, &log)

	ctx := context.Background()

	go storageManager.StartBackgroundSync(ctx)

	// Initialize TUI
	app := tview.NewApplication()
	tuiView := tui.NewTUI(app, facadeClient, storageManager, tokenManager, &log)
	tokenManager.SetAuthFailCallback(tuiView.ResetToLoginScreen)
	// Start TUI
	if err := tuiView.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run tui")
	}

	fmt.Println("Shutting down...")

	storageManager.StopSync()

}

func loadConfig(log *zerolog.Logger) *config.Config {
	err := godotenv.Load("client.env")
	if err != nil {
		log.Error().Err(err).Msg("Error loading client.env file")
	}

	cfg := config.NewConfigBuilder(log).
		FromEnv().
		FromFlags().
		Build()

	log.Info().Interface("config", cfg).Msg("Configuration loaded")

	return cfg
}

func makeConnection(cfg config.Config, interceptor *interceptors.AuthInterceptor) (*grpc.ClientConn, error) {
	creds, err := credentials.NewClientTLSFromFile(cfg.Certificate, "")
	if err != nil {
		return nil, errors.Wrap(err, "could not load TLS keys")
	}
	// Dial the gRPC server with the TLS credentials.
	conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(creds), grpc.WithUnaryInterceptor(interceptor.UnaryInterceptor))
	if err != nil {
		return nil, errors.Wrap(err, "could not dial gRPC server")
	}

	interceptor.SetAuthClient(conn)

	return conn, nil
}
