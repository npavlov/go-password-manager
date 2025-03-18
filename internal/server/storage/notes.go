package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/pkg/errors"
)

// StoreNote creates new note record
func (ds *DBStorage) StoreNote(ctx context.Context, createNote db.CreateNoteEntryParams) (*db.Note, error) {
	note, err := ds.Queries.CreateNoteEntry(ctx, createNote)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to store note")

		return nil, errors.Wrap(err, "failed to store note")
	}

	return &note, nil
}

// GetNote retrieves note record
func (ds *DBStorage) GetNote(ctx context.Context, noteId string) (*db.Note, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(noteId); err != nil {
		ds.log.Error().Err(err).Msg("failed to scan uuid")

		return nil, errors.Wrap(err, "failed to parse uuid")
	}

	note, err := ds.Queries.GetNoteByID(ctx, uuid)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create user")

		return nil, errors.Wrap(err, "failed to create user")
	}

	return &note, nil
}

// GetNotes retrieves note records
func (ds *DBStorage) GetNotes(ctx context.Context, userId string) ([]db.Note, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(userId); err != nil {
		ds.log.Error().Err(err).Msg("failed to scan uuid")

		return nil, errors.Wrap(err, "failed to parse uuid")
	}

	notes, err := ds.Queries.GetNotesByUserID(ctx, uuid)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create note")

		return nil, errors.Wrap(err, "failed to create note")
	}

	return notes, nil
}

func (ds *DBStorage) DeleteNote(ctx context.Context, noteId string) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(noteId); err != nil {
		ds.log.Error().Err(err).Msg("failed to scan uuid")

		return errors.Wrap(err, "failed to parse uuid")
	}

	err := ds.Queries.DeleteNoteEntry(ctx, uuid)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to delete note")

		return errors.Wrap(err, "failed to delete note")
	}

	return nil
}
