package password

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	gu "github.com/npavlov/go-password-manager/internal/utils"
)

type Storage interface {
	StorePassword(ctx context.Context, createPassword db.CreatePasswordEntryParams) (*db.Password, error)
	GetPassword(ctx context.Context, passwordId string, userId pgtype.UUID) (*db.Password, error)
	UpdatePassword(ctx context.Context, updatePassword db.UpdatePasswordEntryParams) (*db.Password, error)
	DeletePassword(ctx context.Context, passwordId string, userId pgtype.UUID) error
	GetUserByID(ctx context.Context, id pgtype.UUID) (*db.User, error)
}

type Service struct {
	pb.UnimplementedPasswordServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   Storage
	cfg       *config.Config
}

func NewPasswordService(log *zerolog.Logger, storage Storage, cfg *config.Config) *Service {
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

func (ps *Service) RegisterService(grpcServer *grpc.Server) {
	pb.RegisterPasswordServiceServer(grpcServer, ps)
}

func (ps *Service) StorePassword(ctx context.Context, req *pb.StorePasswordRequest) (*pb.StorePasswordResponse, error) {
	if err := ps.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, decryptedUserKey, err := utils.GetDecryptionKey(ctx, ps.storage, ps.cfg.SecuredMasterKey.Get())
	if err != nil {
		return nil, errors.Wrap(err, "error getting decrypted user UUID")
	}

	encryptedPassword, err := utils.Encrypt(req.GetPassword().GetPassword(), decryptedUserKey)
	if err != nil {
		ps.logger.Error().Err(err).Msg("failed to encrypt password")

		return nil, errors.Wrap(err, "failed to encrypt password")
	}

	password, err := ps.storage.StorePassword(ctx, db.CreatePasswordEntryParams{
		UserID:   userUUID,
		Login:    req.GetPassword().GetLogin(),
		Password: encryptedPassword,
	})
	if err != nil {
		ps.logger.Error().Err(err).Msg("failed to store password")

		return nil, errors.Wrap(err, "failed to store password")
	}

	return &pb.StorePasswordResponse{
		PasswordId: password.ID.String(),
	}, nil
}

func (ps *Service) GetPassword(ctx context.Context, req *pb.GetPasswordRequest) (*pb.GetPasswordResponse, error) {
	if err := ps.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, decryptedUserKey, err := utils.GetDecryptionKey(ctx, ps.storage, ps.cfg.SecuredMasterKey.Get())
	if err != nil {
		return nil, errors.Wrap(err, "error getting decrypted user UUID")
	}

	password, err := ps.storage.GetPassword(ctx, req.GetPasswordId(), userUUID)
	if err != nil {
		ps.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	decryptedPassword, err := utils.Decrypt(password.Password, decryptedUserKey)
	if err != nil {
		ps.logger.Error().Err(err).Msg("error decrypting password")

		return nil, errors.Wrap(err, "error decrypting password")
	}

	return &pb.GetPasswordResponse{
		Password: &pb.PasswordData{
			Login:    password.Login,
			Password: decryptedPassword,
		},
		LastUpdate: timestamppb.New(password.UpdatedAt.Time),
	}, nil
}

func (ps *Service) GetPasswords(ctx context.Context, req *pb.GetPasswordsRequest) (*pb.GetPasswordsResponse, error) {
	if err := ps.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	return &pb.GetPasswordsResponse{}, nil
}

func (ps *Service) UpdatePassword(
	ctx context.Context,
	req *pb.UpdatePasswordRequest,
) (*pb.UpdatePasswordResponse, error) {
	if err := ps.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	_, decryptedUserKey, err := utils.GetDecryptionKey(ctx, ps.storage, ps.cfg.SecuredMasterKey.Get())
	if err != nil {
		return nil, errors.Wrap(err, "error getting decrypted user UUID")
	}

	encryptedPassword, err := utils.Encrypt(req.GetData().GetPassword(), decryptedUserKey)
	if err != nil {
		ps.logger.Error().Err(err).Msg("failed to encrypt password")

		return nil, errors.Wrap(err, "failed to encrypt password")
	}

	password, err := ps.storage.UpdatePassword(ctx, db.UpdatePasswordEntryParams{
		ID:       gu.GetIDFromString(req.GetPasswordId()),
		Login:    req.GetData().GetLogin(),
		Password: encryptedPassword,
	})
	if err != nil {
		ps.logger.Error().Err(err).Msg("failed to store password")

		return nil, errors.Wrap(err, "failed to store password")
	}

	return &pb.UpdatePasswordResponse{
		PasswordId: password.ID.String(),
	}, nil
}

func (ps *Service) DeletePassword(
	ctx context.Context,
	req *pb.DeletePasswordRequest,
) (*pb.DeletePasswordResponse, error) {
	if err := ps.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserID(ctx)
	if err != nil {
		ps.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	err = ps.storage.DeletePassword(ctx, req.GetPasswordId(), userUUID)
	if err != nil {
		ps.logger.Error().Err(err).Msg("error deleting password")

		return nil, errors.Wrap(err, "error deleting password")
	}

	return &pb.DeletePasswordResponse{
		Ok: true,
	}, nil
}
