package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/utils"
	"github.com/pkg/errors"
)

// StorePassword creates new password record
func (ds *DBStorage) StorePassword(ctx context.Context, createPassword db.CreatePasswordEntryParams) (*db.Password, error) {
	password, err := ds.Queries.CreatePasswordEntry(ctx, createPassword)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to store password")

		return nil, errors.Wrap(err, "failed to store password")
	}

	return &password, nil
}

// GetPassword retrieves password record
func (ds *DBStorage) GetPassword(ctx context.Context, passwordId string, userId pgtype.UUID) (*db.Password, error) {
	uuid := utils.GetIdFromString(passwordId)

	password, err := ds.Queries.GetPasswordEntryByID(ctx, db.GetPasswordEntryByIDParams{
		ID:     uuid,
		UserID: userId,
	})
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create password")

		return nil, errors.Wrap(err, "failed to create password")
	}

	return &password, nil
}

// GetPasswords retrieves password record
func (ds *DBStorage) GetPasswords(ctx context.Context, userId string) ([]db.Password, error) {
	uuid := utils.GetIdFromString(userId)

	passwords, err := ds.Queries.GetPasswordEntriesByUserID(ctx, uuid)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create password")

		return nil, errors.Wrap(err, "failed to create password")
	}

	return passwords, nil
}

func (ds *DBStorage) DeletePassword(ctx context.Context, passwordId string, userId pgtype.UUID) error {
	uuid := utils.GetIdFromString(passwordId)

	err := ds.Queries.DeletePasswordEntry(ctx, db.DeletePasswordEntryParams{
		ID:     uuid,
		UserID: userId,
	})
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to delete password")

		return errors.Wrap(err, "failed to delete password")
	}

	return nil
}

// UpdatePassword updates password record
func (ds *DBStorage) UpdatePassword(ctx context.Context, updatePassword db.UpdatePasswordEntryParams) (*db.Password, error) {
	password, err := ds.Queries.UpdatePasswordEntry(ctx, updatePassword)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to update password")

		return nil, errors.Wrap(err, "failed to update password")
	}

	return &password, nil

}
