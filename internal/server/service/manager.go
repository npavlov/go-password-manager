package service

import (
	"context"
	"net"
	"sync"

	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/redis"
	"github.com/npavlov/go-password-manager/internal/server/service/interceptors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type GManager struct {
	logger     *zerolog.Logger // Logger for logging errors and info.
	cfg        *config.Config
	grpcServer *grpc.Server
}

func NewGRPCManager(cfg *config.Config, logger *zerolog.Logger, memStorage redis.MemStorage) *GManager {
	// Create gRPC server
	creds, err := credentials.NewServerTLSFromFile(cfg.Certificate, cfg.PrivateKey)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to generate credentials")
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.LoggingServerInterceptor(logger), // Logs all requests/responses
		interceptors.TokenInterceptor(logger, cfg.JwtSecret, memStorage),
	), grpc.Creds(creds))
	reflection.Register(grpcServer)

	//nolint:exhaustruct
	return &GManager{
		logger:     logger,
		cfg:        cfg,
		grpcServer: grpcServer,
	}
}

func (gs *GManager) GetServer() *grpc.Server {
	return gs.grpcServer
}

func (gs *GManager) Start(ctx context.Context, wg *sync.WaitGroup) {
	// Start gRPC-server in goroutine
	go func() {
		gs.logger.Info().Str("address", gs.cfg.Address).Msg("starting gRPC server")

		// Start listening
		listener, err := net.Listen("tcp", gs.cfg.Address)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to listen on port")
		}

		log.Info().Msgf("gRPC server listening on %s", gs.cfg.Address)

		if err := gs.grpcServer.Serve(listener); err != nil {
			log.Fatal().Err(err).Msg("failed to serve gRPC server")
		}

	}()

	go func() {
		<-ctx.Done()
		wg.Done()
		gs.logger.Info().Msg("shutting down gRPC server")
		gs.grpcServer.GracefulStop()
	}()
}
