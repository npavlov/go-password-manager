//nolint:exhaustruct
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
	GetPassword(ctx context.Context, passwordID string, userID pgtype.UUID) (*db.Password, error)
	UpdatePassword(ctx context.Context, updatePassword db.UpdatePasswordEntryParams) (*db.Password, error)
	DeletePassword(ctx context.Context, passwordID string, userID pgtype.UUID) error
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

func (ps *Service) StorePasswordV1(ctx context.Context,
	req *pb.StorePasswordV1Request,
) (*pb.StorePasswordV1Response, error) {
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

	return &pb.StorePasswordV1Response{
		PasswordId: password.ID.String(),
	}, nil
}

func (ps *Service) GetPasswordV1(ctx context.Context, req *pb.GetPasswordV1Request) (*pb.GetPasswordV1Response, error) {
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

	return &pb.GetPasswordV1Response{
		Password: &pb.PasswordData{
			Login:    password.Login,
			Password: decryptedPassword,
		},
		LastUpdate: timestamppb.New(password.UpdatedAt.Time),
	}, nil
}

func (ps *Service) GetPasswordsV1(_ context.Context,
	req *pb.GetPasswordsV1Request,
) (*pb.GetPasswordsV1Response, error) {
	if err := ps.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	return &pb.GetPasswordsV1Response{
		Passwords: []*pb.PasswordData{},
	}, nil
}

func (ps *Service) UpdatePasswordV1(
	ctx context.Context,
	req *pb.UpdatePasswordV1Request,
) (*pb.UpdatePasswordV1Response, error) {
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

	return &pb.UpdatePasswordV1Response{
		PasswordId: password.ID.String(),
	}, nil
}

func (ps *Service) DeletePasswordV1(
	ctx context.Context,
	req *pb.DeletePasswordV1Request,
) (*pb.DeletePasswordV1Response, error) {
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

	return &pb.DeletePasswordV1Response{
		Ok: true,
	}, nil
}
