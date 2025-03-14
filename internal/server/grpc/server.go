package grpc

import (
	"context"
	"net"
	"sync"

	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type GServer struct {
	logger  *zerolog.Logger // Logger for logging errors and info.
	cfg     *config.Config
	gServer *grpc.Server
}

func NewGRPCServer(cfg *config.Config, logger *zerolog.Logger) *GServer {
	//TODO: move to config
	creds, err := credentials.NewServerTLSFromFile("certs/cert.pem", "certs/key.pem")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to generate credentials")
	}

	gServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		LoggingServerInterceptor(logger), // Logs all requests/responses
	), grpc.Creds(creds))

	logger.Info().Interface("gServer", gServer).Msg("Created gRPC server")

	//nolint:exhaustruct
	return &GServer{
		logger:  logger,
		cfg:     cfg,
		gServer: gServer,
	}
}

func (gs *GServer) GetServer() *grpc.Server {
	return gs.gServer
}

func (gs *GServer) Start(ctx context.Context, wg *sync.WaitGroup) {
	// Start gRPC-server in goroutine
	go func() {
		gs.logger.Info().Str("address", gs.cfg.Address).Msg("starting gRPC server")

		tcpListen, err := net.Listen("tcp", gs.cfg.Address)
		if err != nil {
			gs.logger.Fatal().Err(err).Str("address", gs.cfg.Address).Msg("failed to listen")
		}

		// Enable reflection
		reflection.Register(gs.gServer)

		if err := gs.gServer.Serve(tcpListen); err != nil {
			gs.logger.Fatal().Err(err).Msg("failed to start gRPC server")
		}
	}()

	go func() {
		<-ctx.Done()
		wg.Done()
		gs.logger.Info().Msg("shutting down gRPC server")
		gs.gServer.GracefulStop()
	}()
}
