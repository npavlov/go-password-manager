//nolint:exhaustruct
package item

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
)

type Storage interface {
	GetItems(ctx context.Context, getItems db.GetItemsByUserIDParams) ([]db.GetItemsByUserIDRow, error)
}

type Service struct {
	pb.UnimplementedItemServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   Storage
	cfg       *config.Config
}

func NewItemService(log *zerolog.Logger, storage Storage, cfg *config.Config) *Service {
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

func (is *Service) RegisterService(grpcServer *grpc.Server) {
	pb.RegisterItemServiceServer(grpcServer, is)
}

func (is *Service) GetItemsV1(ctx context.Context, req *pb.GetItemsV1Request) (*pb.GetItemsV1Response, error) {
	if err := is.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserID(ctx)
	if err != nil {
		is.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	data, err := is.storage.GetItems(ctx, db.GetItemsByUserIDParams{
		UserID: userUUID,
		Offset: (req.GetPage() - 1) * req.GetPageSize(),
		Limit:  req.GetPageSize(),
	})
	if err != nil {
		is.logger.Error().Err(err).Msg("error getting items")

		return nil, errors.Wrap(err, "error getting items")
	}

	//nolint:gosec
	totalCount := int32(len(data))

	items := make([]*pb.ItemData, totalCount)
	for cursor, item := range data {
		items[cursor] = &pb.ItemData{
			Id:        item.IDResource.String(),
			UpdatedAt: timestamppb.New(item.UpdatedAt.Time),
			CreatedAt: timestamppb.New(item.CreatedAt.Time),
		}

		switch item.Type {
		case db.ItemTypeCard:
			items[cursor].Type = pb.ItemType_ITEM_TYPE_CARD
		case db.ItemTypeBinary:
			items[cursor].Type = pb.ItemType_ITEM_TYPE_BINARY
		case db.ItemTypePassword:
			items[cursor].Type = pb.ItemType_ITEM_TYPE_PASSWORD
		case db.ItemTypeText:
			items[cursor].Type = pb.ItemType_ITEM_TYPE_NOTE
		}
	}

	return &pb.GetItemsV1Response{
		Items:      items,
		TotalCount: totalCount,
	}, nil
}
