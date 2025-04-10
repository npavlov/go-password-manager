package meta

import (
	"context"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/metadata"
	"github.com/npavlov/go-password-manager/internal/server/config"
)

type Storage interface {
	AddMeta(ctx context.Context, recordId string, key string, value string) (*db.Metainfo, error)
	GetMetaInfo(ctx context.Context, recordId string) ([]db.GetMetaInfoByItemIDRow, error)
	DeleteMetaInfo(ctx context.Context, key string, itemId string) error
}

type Service struct {
	pb.UnimplementedMetadataServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   Storage
	cfg       *config.Config
}

func NewMetadataService(log *zerolog.Logger, storage Storage, cfg *config.Config) *Service {
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

// AddMetaInfo Add metadata to an item.
func (ms *Service) AddMetaInfo(ctx context.Context, req *pb.AddMetaInfoRequest) (*pb.AddMetaInfoResponse, error) {
	if err := ms.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	for key, value := range req.GetMetadata() {
		_, err := ms.storage.AddMeta(ctx, req.GetItemId(), key, value)
		if err != nil {
			return nil, errors.Wrap(err, "failed to add meta info")
		}
	}

	return &pb.AddMetaInfoResponse{Success: true}, nil
}

// RemoveMetaInfo Remove a metadata key from an item.
func (ms *Service) RemoveMetaInfo(ctx context.Context, req *pb.RemoveMetaInfoRequest) (*pb.RemoveMetaInfoResponse, error) {
	if err := ms.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	err := ms.storage.DeleteMetaInfo(ctx, req.GetKey(), req.GetItemId())
	if err != nil {
		ms.logger.Error().Err(err).Msg("failed to remove meta info")

		return nil, errors.Wrap(err, "failed to remove meta info")
	}

	return &pb.RemoveMetaInfoResponse{Success: true}, nil
}

// GetMetaInfo Get all metadata for an item.
func (ms *Service) GetMetaInfo(ctx context.Context, req *pb.GetMetaInfoRequest) (*pb.GetMetaInfoResponse, error) {
	if err := ms.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	rows, err := ms.storage.GetMetaInfo(ctx, req.GetItemId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}

	metadata := make(map[string]string)
	for _, row := range rows {
		metadata[row.Key] = row.Value
	}

	return &pb.GetMetaInfoResponse{Metadata: metadata}, nil
}
