package meta

import (
	"context"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	pb "github.com/npavlov/go-password-manager/gen/proto/metadata"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Service struct {
	pb.UnimplementedMetadataServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   *storage.DBStorage
	cfg       *config.Config
}

func NewMetadataService(log *zerolog.Logger, storage *storage.DBStorage, cfg *config.Config) *Service {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create validator")
	}

	return &Service{
		logger:    log,
		validator: validator,
		storage:   storage,
		cfg:       cfg,
	}
}

func (ms *Service) RegisterService(grpcServer *grpc.Server) {
	pb.RegisterMetadataServiceServer(grpcServer, ms)
}

// AddMetaInfo Add metadata to an item
func (ms *Service) AddMetaInfo(ctx context.Context, req *pb.AddMetaInfoRequest) (*pb.AddMetaInfoResponse, error) {
	for key, value := range req.Metadata {

		_, err := ms.storage.AddMeta(ctx, req.ItemId, key, value)
		if err != nil {
			return nil, errors.Wrap(err, "failed to add meta info")
		}
	}

	return &pb.AddMetaInfoResponse{Success: true}, nil
}

// RemoveMetaInfo Remove a metadata key from an item
func (ms *Service) RemoveMetaInfo(ctx context.Context, req *pb.RemoveMetaInfoRequest) (*pb.RemoveMetaInfoResponse, error) {
	err := ms.storage.DeleteMetaInfo(ctx, req.Key, req.ItemId)
	if err != nil {
		ms.logger.Error().Err(err).Msg("failed to remove meta info")

		return nil, errors.Wrap(err, "failed to remove meta info")
	}

	return &pb.RemoveMetaInfoResponse{Success: true}, nil
}

// GetMetaInfo Get all metadata for an item
func (ms *Service) GetMetaInfo(ctx context.Context, req *pb.GetMetaInfoRequest) (*pb.GetMetaInfoResponse, error) {
	rows, err := ms.storage.GetMetaInfo(ctx, req.ItemId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %v", err)
	}

	metadata := make(map[string]string)
	for _, row := range rows {
		metadata[row.Key] = row.Value
	}

	return &pb.GetMetaInfoResponse{Metadata: metadata}, nil
}
