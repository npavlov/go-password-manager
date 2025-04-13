package note

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/npavlov/go-password-manager/gen/proto/note"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
)

type Storage interface {
	StoreNote(ctx context.Context, createNote db.CreateNoteEntryParams) (*db.Note, error)
	GetNote(ctx context.Context, noteId string, userId pgtype.UUID) (*db.Note, error)
	GetUserById(ctx context.Context, id pgtype.UUID) (*db.User, error)
	DeleteNote(ctx context.Context, noteID string, userID pgtype.UUID) error
}

type Service struct {
	pb.UnimplementedNoteServiceServer
	validator protovalidate.Validator
	logger    *zerolog.Logger
	storage   Storage
	cfg       *config.Config
}

func NewNoteService(log *zerolog.Logger, storage Storage, cfg *config.Config) *Service {
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

	userUUID, decryptedUserKey, err := utils.GetDecryptionKey(ctx, ns.storage, ns.cfg.SecuredMasterKey.Get())
	if err != nil {
		return nil, err
	}

	encryptedNote, err := utils.Encrypt(req.GetNote().GetContent(), decryptedUserKey)
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

	userUUID, decryptedUserKey, err := utils.GetDecryptionKey(ctx, ns.storage, ns.cfg.SecuredMasterKey.Get())

	note, err := ns.storage.GetNote(ctx, req.GetNoteId(), userUUID)
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
			Content: content,
		},
		LastUpdate: timestamppb.New(note.UpdatedAt.Time),
	}, nil
}

func (ns *Service) GetNotes(ctx context.Context, req *pb.GetNotesRequest) (*pb.GetNotesResponse, error) {
	if err := ns.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	return &pb.GetNotesResponse{}, nil
}

func (ns *Service) DeleteNote(ctx context.Context, req *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error) {
	if err := ns.validator.Validate(req); err != nil {
		return nil, errors.Wrap(err, "error validating input")
	}

	userUUID, err := utils.GetUserId(ctx)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error getting user id")

		return nil, errors.Wrap(err, "error getting user id")
	}

	err = ns.storage.DeleteNote(ctx, req.GetNoteId(), userUUID)
	if err != nil {
		ns.logger.Error().Err(err).Msg("error deleting note")

		return nil, errors.Wrap(err, "error deleting note")
	}

	return &pb.DeleteNoteResponse{
		Ok: true,
	}, nil
}
