package password

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	pb "github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	"github.com/npavlov/go-password-manager/internal/server/storage"
	gu "github.com/npavlov/go-password-manager/internal/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	pb.UnimplementedPasswordServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   *storage.DBStorage
	cfg       *config.Config
}

func NewPasswordService(log *zerolog.Logger, storage *storage.DBStorage, cfg *config.Config) *Service {
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

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		ps.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	decryptedUserKey, err := utils.GetUserKey(ctx, ps.storage, userUUID, ps.cfg.SecuredMasterKey.Get())
	if err != nil {
		ps.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	encryptedPassword, err := utils.Encrypt(req.Password.Password, decryptedUserKey)
	if err != nil {
		ps.logger.Error().Err(err).Msg("failed to encrypt password")

		return nil, errors.Wrap(err, "failed to encrypt password")
	}

	password, err := ps.storage.StorePassword(ctx, db.CreatePasswordEntryParams{
		UserID:   userUUID,
		Name:     req.Password.Name,
		Login:    req.Password.Login,
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

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		ps.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	decryptedUserKey, err := utils.GetUserKey(ctx, ps.storage, userUUID, ps.cfg.SecuredMasterKey.Get())
	if err != nil {
		ps.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	password, err := ps.storage.GetPassword(ctx, req.PasswordId, userUUID)
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
			Name:     password.Name,
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

func (ps *Service) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	if err := ps.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		ps.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	decryptedUserKey, err := utils.GetUserKey(ctx, ps.storage, userUUID, ps.cfg.SecuredMasterKey.Get())
	if err != nil {
		ps.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	encryptedPassword, err := utils.Encrypt(req.Data.Password, decryptedUserKey)
	if err != nil {
		ps.logger.Error().Err(err).Msg("failed to encrypt password")

		return nil, errors.Wrap(err, "failed to encrypt password")
	}

	password, err := ps.storage.UpdatePassword(ctx, db.UpdatePasswordEntryParams{
		ID:       gu.GetIdFromString(req.PasswordId),
		Name:     req.Data.Name,
		Login:    req.Data.Login,
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

func (ps *Service) DeletePassword(ctx context.Context, req *pb.DeletePasswordRequest) (*pb.DeletePasswordResponse, error) {
	if err := ps.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		ps.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	err = ps.storage.DeletePassword(ctx, req.PasswordId, userUUID)
	if err != nil {
		ps.logger.Error().Err(err).Msg("error deleting password")

		return nil, errors.Wrap(err, "error deleting password")
	}

	return &pb.DeletePasswordResponse{
		Ok: true,
	}, nil
}
