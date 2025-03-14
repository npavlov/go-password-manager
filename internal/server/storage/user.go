package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/pkg/errors"
)

// RegisterUser creates new user record
func (ds *DBStorage) RegisterUser(ctx context.Context, createUser db.CreateUserParams) (*db.User, error) {
	user, err := ds.Queries.CreateUser(ctx, createUser)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create user")

		return nil, errors.Wrap(err, "failed to create user")
	}

	return &user, nil
}

// GetUser retrieves user record
func (ds *DBStorage) GetUser(ctx context.Context, id string) (*db.User, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		ds.log.Error().Err(err).Msg("failed to scan uuid")

		return nil, errors.Wrap(err, "failed to parse uuid")
	}

	user, err := ds.Queries.GetUserByID(ctx, uuid)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create user")

		return nil, errors.Wrap(err, "failed to create user")
	}

	return &user, nil
}
