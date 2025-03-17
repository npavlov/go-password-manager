package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/npavlov/go-password-manager/internal/server/db"
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
func (ds *DBStorage) GetPassword(ctx context.Context, passwordId string) (*db.Password, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(passwordId); err != nil {
		ds.log.Error().Err(err).Msg("failed to scan uuid")

		return nil, errors.Wrap(err, "failed to parse uuid")
	}

	password, err := ds.Queries.GetPasswordEntryByID(ctx, uuid)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create user")

		return nil, errors.Wrap(err, "failed to create user")
	}

	return &password, nil
}

// GetPasswords retrieves password record
func (ds *DBStorage) GetPasswords(ctx context.Context, userId string) ([]db.Password, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(userId); err != nil {
		ds.log.Error().Err(err).Msg("failed to scan uuid")

		return nil, errors.Wrap(err, "failed to parse uuid")
	}

	passwords, err := ds.Queries.GetPasswordEntriesByUserID(ctx, uuid)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create user")

		return nil, errors.Wrap(err, "failed to create user")
	}

	return passwords, nil
}

func (ds *DBStorage) DeletePassword(ctx context.Context, passwordId string) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(passwordId); err != nil {
		ds.log.Error().Err(err).Msg("failed to scan uuid")

		return errors.Wrap(err, "failed to parse uuid")
	}

	err := ds.Queries.DeletePasswordEntry(ctx, uuid)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to delete user")

		return errors.Wrap(err, "failed to delete user")
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
