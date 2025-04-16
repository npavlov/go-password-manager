package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"

	"github.com/npavlov/go-password-manager/internal/server/db"
)

// RegisterUser creates new user record.
func (ds *DBStorage) RegisterUser(ctx context.Context, createUser db.CreateUserParams) (*db.User, error) {
	user, err := ds.Queries.CreateUser(ctx, createUser)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create user")

		return nil, errors.Wrap(err, "failed to create user")
	}

	return &user, nil
}

// GetUser retrieves user record.
func (ds *DBStorage) GetUser(ctx context.Context, username string) (*db.User, error) {
	user, err := ds.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create user")

		return nil, errors.Wrap(err, "failed to create user")
	}

	return &user, nil
}

// GetUserById retrieves user record.
func (ds *DBStorage) GetUserByID(ctx context.Context, userId pgtype.UUID) (*db.User, error) {
	user, err := ds.Queries.GetUserByID(ctx, userId)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create user")

		return nil, errors.Wrap(err, "failed to create user")
	}

	return &user, nil
}
