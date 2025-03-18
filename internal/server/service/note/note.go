package note

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	pb "github.com/npavlov/go-password-manager/gen/proto/note"
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
	pb.UnimplementedNoteServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   *storage.DBStorage
	cfg       *config.Config
}

func NewNoteService(log *zerolog.Logger, storage *storage.DBStorage, cfg *config.Config) *Service {
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
	pb.RegisterNoteServiceServer(grpcServer, ns)
}

func (ns *Service) StoreNote(ctx context.Context, req *pb.StoreNoteRequest) (*pb.StoreNoteResponse, error) {
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

	encryptedNote, err := utils.Encrypt(req.Content, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to encrypt password")

		return nil, errors.Wrap(err, "failed to encrypt password")
	}

	note, err := ns.storage.StoreNote(ctx, db.CreateNoteEntryParams{
		UserID:           userUUID,
		EncryptedContent: encryptedNote,
	})
	if err != nil {
		ns.logger.Error().Err(err).Msg("failed to store password")

		return nil, errors.Wrap(err, "failed to store password")
	}

	return &pb.StoreNoteResponse{
		NoteId: note.ID.String(),
	}, nil
}

func (ns *Service) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
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

	note, err := ns.storage.GetNote(ctx, req.NoteId)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	content, err := utils.Decrypt(note.EncryptedContent, decryptedUserKey)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error decrypting password")

		return nil, errors.Wrap(err, "error decrypting password")
	}

	return &pb.GetNoteResponse{
		Note: &pb.NoteData{
			Content:    content,
			LastUpdate: timestamppb.New(note.UpdatedAt.Time),
		},
	}, nil
}

func (ns *Service) GetNotes(ctx context.Context, req *pb.GetNotesRequest) (*pb.GetNotesResponse, error) {
	if err := ns.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	return &pb.GetNotesResponse{}, nil
}
