package item

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	"github.com/npavlov/go-password-manager/internal/server/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	pb.UnimplementedItemServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   *storage.DBStorage
	cfg       *config.Config
}

func NewItemService(log *zerolog.Logger, storage *storage.DBStorage, cfg *config.Config) *Service {
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

func (is *Service) GetItems(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	if err := is.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		is.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	data, err := is.storage.GetItems(ctx, db.GetItemsByUserIDParams{
		UserID: userUUID,
		Offset: (req.Page - 1) * req.PageSize,
		Limit:  req.PageSize,
	})

	if err != nil {
		is.logger.Error().Err(err).Msg("error getting items")

		return nil, errors.Wrap(err, "error getting items")
	}

	totalCount := int32(len(data))

	items := make([]*pb.ItemData, totalCount)
	for i, item := range data {
		items[i] = &pb.ItemData{
			Id:        item.IDResource.String(),
			UpdatedAt: timestamppb.New(item.UpdatedAt.Time),
			CreatedAt: timestamppb.New(item.CreatedAt.Time),
		}

		switch item.Type {
		case db.ItemTypeCard:
			items[i].Type = pb.ItemType_ITEM_TYPE_CARD
		case db.ItemTypeBinary:
			items[i].Type = pb.ItemType_ITEM_TYPE_BINARY
		case db.ItemTypePassword:
			items[i].Type = pb.ItemType_ITEM_TYPE_PASSWORD
		case db.ItemTypeText:
			items[i].Type = pb.ItemType_ITEM_TYPE_NOTE
		}

	}

	return &pb.GetItemsResponse{
		Items:      items,
		TotalCount: totalCount,
	}, nil
}
