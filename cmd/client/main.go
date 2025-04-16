// main.go
package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/npavlov/go-password-manager/internal/client/buildinfo"
	"github.com/npavlov/go-password-manager/internal/client/config"
	"github.com/npavlov/go-password-manager/internal/client/grpc/facade"
	"github.com/npavlov/go-password-manager/internal/client/interceptors"
	"github.com/npavlov/go-password-manager/internal/client/storage"
	"github.com/npavlov/go-password-manager/internal/client/tui"
	"github.com/npavlov/go-password-manager/internal/pkg/logger"
	"github.com/npavlov/go-password-manager/internal/utils"
)

func main() {
	log := logger.NewLogger(zerolog.DebugLevel).Get()

	ctx, cancel := utils.WithSignalCancel(context.Background(), &log)
	defer cancel()

	tokenMgr, facadeClient, storageMgr, conn := GetApp(&log)
	defer conn.Close()

	go storageMgr.StartBackgroundSync(ctx)

	uiApp := GetTUI(&log, facadeClient, storageMgr, tokenMgr)

	if err := uiApp.Run(); err != nil {
		log.Error().Err(err).Msg("Failed to run tui")
	}

	log.Info().Msg("Shutting down...")
	storageMgr.StopSync()
}

func GetApp(log *zerolog.Logger) (auth.ITokenManager, facade.IFacade, storage.IStorageManager, *grpc.ClientConn) {
	log.Info().Str("buildVersion", buildinfo.Version).
		Str("buildCommit", buildinfo.Commit).
		Str("buildDate", buildinfo.Date).Msg("Starting agent")

	cfg := LoadConfig(log)

	tokenManager := auth.NewTokenManager(log, cfg)
	err := tokenManager.LoadTokens()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load tokens")
	}

	authInterceptor := interceptors.NewAuthInterceptor(*cfg, tokenManager)
	conn, err := MakeConnection(*cfg, authInterceptor)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to make connection")
	}

	facadeClient := facade.NewFacade(conn, tokenManager, log)
	storageManager := storage.NewStorageManager(facadeClient, tokenManager, log)

	return tokenManager, facadeClient, storageManager, conn
}

func GetTUI(
	log *zerolog.Logger,
	facadeClient facade.IFacade,
	storageManager storage.IStorageManager,
	tokenManager auth.ITokenManager,
) *tview.Application {
	app := tview.NewApplication()
	tuiView := tui.NewTUI(app, facadeClient, storageManager, tokenManager, log)

	tokenManager.SetAuthFailCallback(tuiView.ResetToLoginScreen)

	return tuiView.GetApp()
}

func LoadConfig(log *zerolog.Logger) *config.Config {
	_ = godotenv.Load("client.env") // Optional
	cfg := config.NewConfigBuilder(log).FromEnv().FromFlags().Build()
	log.Info().Interface("config", cfg).Msg("Configuration loaded")

	return cfg
}

func MakeConnection(cfg config.Config, interceptor *interceptors.AuthInterceptor) (*grpc.ClientConn, error) {
	creds, err := credentials.NewClientTLSFromFile(cfg.Certificate, "")
	if err != nil {
		return nil, errors.Wrap(err, "could not load TLS keys")
	}
	conn, err := grpc.NewClient(cfg.Address,
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(interceptor.UnaryInterceptor),
		grpc.WithStreamInterceptor(interceptor.StreamInterceptor))
	if err != nil {
		return nil, errors.Wrap(err, "could not dial gRPC server")
	}

	authClient := pb.NewAuthServiceClient(conn)
	interceptor.SetAuthClient(authClient)

	return conn, nil
}
