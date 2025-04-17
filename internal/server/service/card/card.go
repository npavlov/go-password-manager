//nolint:exhaustruct
package card

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/npavlov/go-password-manager/gen/proto/card"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	gu "github.com/npavlov/go-password-manager/internal/utils"
)

type Service struct {
	pb.UnimplementedCardServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   Storage
	cfg       *config.Config
}

type Storage interface {
	UpdateCard(ctx context.Context, updateCard db.UpdateCardParams) (*db.Card, error)
	DeleteCard(ctx context.Context, cardID string, userID pgtype.UUID) error
	GetCards(ctx context.Context, userID string) ([]db.Card, error)
	GetCard(ctx context.Context, cardID string, userID pgtype.UUID) (*db.Card, error)
	StoreCard(ctx context.Context, createCard db.StoreCardParams) (*db.Card, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (*db.User, error)
}

func NewCardService(log *zerolog.Logger, storage Storage, cfg *config.Config) *Service {
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

func (ns *Service) RegisterService(grpcServer *grpc.Server) {
	pb.RegisterCardServiceServer(grpcServer, ns)
}

func (ns *Service) StoreCard(ctx context.Context, req *pb.StoreCardRequest) (*pb.StoreCardResponse, error) {
	if err := ns.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, decryptedUserKey, err := utils.GetDecryptionKey(ctx, ns.storage, ns.cfg.SecuredMasterKey.Get())
	if err != nil {
		return nil, errors.Wrap(err, "error getting decrypted user UUID")
	}

	data := req.GetCard()

	encryptedCardNumber, encryptedCVV, encryptedExpiryDate, err := ns.EncryptCard(
		decryptedUserKey,
		data.GetCardNumber(),
		data.GetCvv(),
		data.GetExpiryDate(),
	)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error encrypting card")

		return nil, errors.Wrap(err, "error encrypting card")
	}

	// Hash card number for uniqueness check
	hashedCardNumber := utils.HashCardNumber(req.GetCard().GetCardNumber())

	Card, err := ns.storage.StoreCard(ctx, db.StoreCardParams{
		UserID:              userUUID,
		EncryptedCardNumber: encryptedCardNumber,
		HashedCardNumber:    hashedCardNumber,
		EncryptedCvv:        encryptedCVV,
		EncryptedExpiryDate: encryptedExpiryDate,
		CardholderName:      req.GetCard().GetCardholderName(),
	})
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to store card")

		return nil, errors.Wrap(err, "failed to store card")
	}

	return &pb.StoreCardResponse{
		CardId: Card.ID.String(),
	}, nil
}

func (ns *Service) UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*pb.UpdateCardResponse, error) {
	if err := ns.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	data := req.GetData()

	_, decryptedUserKey, err := utils.GetDecryptionKey(ctx, ns.storage, ns.cfg.SecuredMasterKey.Get())
	if err != nil {
		return nil, errors.Wrap(err, "error getting decrypted user UUID")
	}

	encryptedCardNumber, encryptedCVV, encryptedExpiryDate, err := ns.EncryptCard(
		decryptedUserKey,
		data.GetCardNumber(),
		data.GetCvv(),
		data.GetExpiryDate(),
	)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error encrypting card")

		return nil, errors.Wrap(err, "error encrypting card")
	}

	// Hash card number for uniqueness check
	hashedCardNumber := utils.HashCardNumber(req.GetData().GetCardNumber())

	card, err := ns.storage.UpdateCard(ctx, db.UpdateCardParams{
		ID:                  gu.GetIDFromString(req.GetCardId()),
		EncryptedCardNumber: encryptedCardNumber,
		HashedCardNumber:    hashedCardNumber,
		EncryptedCvv:        encryptedCVV,
		EncryptedExpiryDate: encryptedExpiryDate,
		CardholderName:      req.GetData().GetCardholderName(),
	})
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to store password")

		return nil, errors.Wrap(err, "failed to store password")
	}

	return &pb.UpdateCardResponse{
		CardId: card.ID.String(),
	}, nil
}

func (ns *Service) EncryptCard(decryptedUserKey, cardNum, cvv, expiryDate string) (string, string, string, error) {
	encryptedCardNumber, err := utils.Encrypt(cardNum, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to encrypt card number")

		return "", "", "", errors.Wrap(err, "failed to encrypt card number")
	}

	encryptedCVV, err := utils.Encrypt(cvv, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to encrypt card cvv")

		return "", "", "", errors.Wrap(err, "failed to encrypt card cvv")
	}

	encryptedExpiryDate, err := utils.Encrypt(expiryDate, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to encrypt card Expiry Date")

		return "", "", "", errors.Wrap(err, "failed to encrypt card Expiry Date")
	}

	return encryptedCardNumber, encryptedCVV, encryptedExpiryDate, nil
}

func (ns *Service) GetCard(ctx context.Context, req *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	if err := ns.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, decryptedUserKey, err := utils.GetDecryptionKey(ctx, ns.storage, ns.cfg.SecuredMasterKey.Get())
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting decrypted user UUID")

		return nil, errors.Wrap(err, "error getting decrypted user UUID")
	}

	Card, err := ns.storage.GetCard(ctx, req.GetCardId(), userUUID)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	cardNumber, err := utils.Decrypt(Card.EncryptedCardNumber, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error decrypting Card Number")

		return nil, errors.Wrap(err, "error decrypting Card Number")
	}

	cardCvv, err := utils.Decrypt(Card.EncryptedCvv, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error decrypting Card CVV")

		return nil, errors.Wrap(err, "error decrypting Card CVV")
	}

	expiryDate, err := utils.Decrypt(Card.EncryptedExpiryDate, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error decrypting Expiry Date")

		return nil, errors.Wrap(err, "error decrypting Expiry Date")
	}

	return &pb.GetCardResponse{
		Card: &pb.CardData{
			CardNumber:     cardNumber,
			Cvv:            cardCvv,
			ExpiryDate:     expiryDate,
			CardholderName: Card.CardholderName,
		},
		LastUpdate: timestamppb.New(Card.UpdatedAt.Time),
	}, nil
}

func (ns *Service) GetCards(_ context.Context, req *pb.GetCardsRequest) (*pb.GetCardsResponse, error) {
	if err := ns.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	return &pb.GetCardsResponse{
		Cards: make([]*pb.CardData, 0),
	}, nil
}

func (ns *Service) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	if err := ns.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, _, err := utils.GetDecryptionKey(ctx, ns.storage, ns.cfg.SecuredMasterKey.Get())
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	err = ns.storage.DeleteCard(ctx, req.GetCardId(), userUUID)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error deleting Card")

		return nil, errors.Wrap(err, "error deleting Card")
	}

	return &pb.DeleteCardResponse{
		Ok: true,
	}, nil
}
