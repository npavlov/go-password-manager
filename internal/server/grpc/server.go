package grpc

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/bufbuild/protovalidate-go"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	logger    *zerolog.Logger // Logger for logging errors and info.
	cfg       *config.Config
	gServer   *grpc.Server
	validator protovalidate.Validator
}

func NewGRPCServer(cfg *config.Config, logger *zerolog.Logger) *Server {
	validator, err := protovalidate.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create validator")
	}

	//TODO: move to config
	creds, err := credentials.NewServerTLSFromFile("certs/cert.pem", "certs/key.pem")
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
	}

	//nolint:exhaustruct
	return &Server{
		logger: logger,
		cfg:    cfg,
		gServer: grpc.NewServer(grpc.ChainUnaryInterceptor(
			LoggingServerInterceptor(logger), // Logs all requests/responses
		), grpc.Creds(creds)),
		validator: validator,
	}
}

func (s *Server) GetServer() *grpc.Server {
	return s.gServer
}

func (gs *Server) Start(ctx context.Context, wg *sync.WaitGroup) {
	// Start gRPC-server in goroutine
	go func() {
		gs.logger.Info().Str("address", gs.cfg.Address).Msg("starting gRPC server")

		tcpListen, err := net.Listen("tcp", gs.cfg.Address)
		if err != nil {
			gs.logger.Fatal().Err(err).Str("address", gs.cfg.Address).Msg("failed to listen")
		}

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
