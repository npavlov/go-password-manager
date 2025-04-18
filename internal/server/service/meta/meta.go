//nolint:exhaustruct
package meta

import (
	"context"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/metadata"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
)

type Storage interface {
	AddMeta(ctx context.Context, recordID string, key string, value string) (*db.Metainfo, error)
	GetMetaInfo(ctx context.Context, recordID string) ([]db.GetMetaInfoByItemIDRow, error)
	DeleteMetaInfo(ctx context.Context, key string, itemID string) error
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

// AddMetaInfoV1 Add metadata to an item.
func (ms *Service) AddMetaInfoV1(ctx context.Context, req *pb.AddMetaInfoV1Request) (*pb.AddMetaInfoV1Response, error) {
	if err := ms.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	for key, value := range req.GetMetadata() {
		_, err := ms.storage.AddMeta(ctx, req.GetItemId(), key, value)
		if err != nil {
			return nil, errors.Wrap(err, "failed to add meta info")
		}
	}

	return &pb.AddMetaInfoV1Response{Success: true}, nil
}

// RemoveMetaInfoV1 Remove a metadata key from an item.
func (ms *Service) RemoveMetaInfoV1(ctx context.Context,
	req *pb.RemoveMetaInfoV1Request,
) (*pb.RemoveMetaInfoV1Response, error) {
	if err := ms.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	err := ms.storage.DeleteMetaInfo(ctx, req.GetKey(), req.GetItemId())
	if err != nil {
		ms.logger.Error().Err(err).Msg("failed to remove meta info")

		return nil, errors.Wrap(err, "failed to remove meta info")
	}

	return &pb.RemoveMetaInfoV1Response{Success: true}, nil
}

// GetMetaInfoV1 Get all metadata for an item.
func (ms *Service) GetMetaInfoV1(ctx context.Context, req *pb.GetMetaInfoV1Request) (*pb.GetMetaInfoV1Response, error) {
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

	return &pb.GetMetaInfoV1Response{Metadata: metadata}, nil
}
