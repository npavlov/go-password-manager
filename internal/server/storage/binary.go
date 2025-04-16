package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"

	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/utils"
)

// StoreBinary creates new binary record.
func (ds *DBStorage) StoreBinary(ctx context.Context, createBinary db.StoreBinaryEntryParams) (*db.BinaryEntry, error) {
	binary, err := ds.Queries.StoreBinaryEntry(ctx, createBinary)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to store binary entry")

		return nil, errors.Wrap(err, "failed to store binary entry")
	}

	return &binary, nil
}

// GetBinary retrieves binary record.
func (ds *DBStorage) GetBinary(ctx context.Context, binaryID string, userId pgtype.UUID) (*db.BinaryEntry, error) {
	uuid := utils.GetIDFromString(binaryID)

	binary, err := ds.Queries.GetBinaryEntryByID(ctx, db.GetBinaryEntryByIDParams{
		ID:     uuid,
		UserID: userId,
	})
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to find binary")

		return nil, errors.Wrap(err, "failed to find binary")
	}

	return &binary, nil
}

// GetBinaries retrieves binary record.
func (ds *DBStorage) GetBinaries(ctx context.Context, userID string) ([]db.BinaryEntry, error) {
	uuid := utils.GetIDFromString(userID)

	cards, err := ds.Queries.GetBinaryEntriesByUserID(ctx, uuid)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to find cards")

		return nil, errors.Wrap(err, "failed to find cards")
	}

	return cards, nil
}

// DeleteBinary removes binary record.
func (ds *DBStorage) DeleteBinary(ctx context.Context, arg db.DeleteBinaryEntryParams) error {
	err := ds.Queries.DeleteBinaryEntry(ctx, arg)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to delete binary entry")

		return errors.Wrap(err, "failed to delete binary entry")
	}

	return nil
}
