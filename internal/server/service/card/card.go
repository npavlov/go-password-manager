package card

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/bufbuild/protovalidate-go"
	"github.com/jackc/pgx/v5/pgtype"
	pb "github.com/npavlov/go-password-manager/gen/proto/card"
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
	pb.UnimplementedCardServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   *storage.DBStorage
	cfg       *config.Config
}

func NewCardService(log *zerolog.Logger, storage *storage.DBStorage, cfg *config.Config) *Service {
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

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	decryptedUserKey, err := utils.GetUserKey(ctx, ns.storage, userUUID, ns.cfg.SecuredMasterKey.Get())
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	encryptedCardNumber, err := utils.Encrypt(req.Card.CardNumber, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to encrypt card number")

		return nil, errors.Wrap(err, "failed to encrypt card number")
	}

	encryptedCVV, err := utils.Encrypt(req.Card.Cvv, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to encrypt card CVV")

		return nil, errors.Wrap(err, "failed to encrypt card CVV")
	}

	encryptedExpiryDate, err := utils.Encrypt(req.Card.ExpiryDate, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to encrypt card Expiry Date")

		return nil, errors.Wrap(err, "failed to encrypt card Expiry Date")
	}

	// Hash card number for uniqueness check
	hashedCardNumber := ns.HashCardNumber(req.Card.CardNumber)

	Card, err := ns.storage.StoreCard(ctx, db.StoreCardParams{
		UserID:              userUUID,
		EncryptedCardNumber: encryptedCardNumber,
		HashedCardNumber:    hashedCardNumber,
		EncryptedCvv:        encryptedCVV,
		EncryptedExpiryDate: encryptedExpiryDate,
		CardholderName:      req.Card.CardholderName,
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

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	decryptedUserKey, err := utils.GetUserKey(ctx, ns.storage, userUUID, ns.cfg.SecuredMasterKey.Get())
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	encryptedCardNumber, err := utils.Encrypt(req.Data.CardNumber, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to encrypt card number")

		return nil, errors.Wrap(err, "failed to encrypt card number")
	}

	encryptedCVV, err := utils.Encrypt(req.Data.Cvv, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to encrypt card CVV")

		return nil, errors.Wrap(err, "failed to encrypt card CVV")
	}

	encryptedExpiryDate, err := utils.Encrypt(req.Data.ExpiryDate, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to encrypt card Expiry Date")

		return nil, errors.Wrap(err, "failed to encrypt card Expiry Date")
	}

	// Hash card number for uniqueness check
	hashedCardNumber := ns.HashCardNumber(req.Data.CardNumber)

	card, err := ns.storage.UpdateCard(ctx, db.UpdateCardParams{
		ID:                  utils.GetIdFromString(req.CardId),
		EncryptedCardNumber: encryptedCardNumber,
		HashedCardNumber:    hashedCardNumber,
		EncryptedCvv:        encryptedCVV,
		EncryptedExpiryDate: encryptedExpiryDate,
		CardholderName:      req.Data.CardholderName,
	})
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to store password")

		return nil, errors.Wrap(err, "failed to store password")
	}

	return &pb.UpdateCardResponse{
		CardId: card.ID.String(),
	}, nil
}

// HashCardNumber hashes the card number to enforce uniqueness
func (ns *Service) HashCardNumber(cardNumber string) pgtype.Text {
	hash := sha256.Sum256([]byte(cardNumber))

	text := pgtype.Text{}

	hashString := hex.EncodeToString(hash[:]) // Convert to hex string for storage

	_ = text.Scan(hashString)

	return text
}

func (ns *Service) GetCard(ctx context.Context, req *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	if err := ns.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	decryptedUserKey, err := utils.GetUserKey(ctx, ns.storage, userUUID, ns.cfg.SecuredMasterKey.Get())
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	Card, err := ns.storage.GetCard(ctx, req.CardId)
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

func (ns *Service) GetCards(ctx context.Context, req *pb.GetCardsRequest) (*pb.GetCardsResponse, error) {
	if err := ns.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	return &pb.GetCardsResponse{}, nil
}

func (ns *Service) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	if err := ns.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	err = ns.storage.DeleteCard(ctx, db.DeleteCardParams{
		ID:     utils.GetIdFromString(req.CardId),
		UserID: userUUID,
	})
	if err != nil {
		ns.logger.Error().Err(err).Msg("error deleting Card")

		return nil, errors.Wrap(err, "error deleting Card")
	}

	return &pb.DeleteCardResponse{
		Ok: true,
	}, nil
}
