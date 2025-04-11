package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"

	"github.com/npavlov/go-password-manager/internal/server/db"
)

// StoreToken stores refresh token.
func (ds *DBStorage) StoreToken(ctx context.Context, userID pgtype.UUID, refreshToken string, expiresAt time.Time) error {
	var pgExpiresAt pgtype.Timestamp
	err := pgExpiresAt.Scan(expiresAt)
	if err != nil {
		return errors.Wrap(err, "failed to scan expires at")
	}

	err = ds.Queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: pgExpiresAt,
	})

	return errors.Wrap(err, "error creating refresh token")
}

// GetToken gets refresh token.
func (ds *DBStorage) GetToken(ctx context.Context, token string) (db.GetRefreshTokenRow, error) {
	tokenDb, err := ds.Queries.GetRefreshToken(ctx, token)

	return tokenDb, errors.Wrap(err, "error getting refresh token")
}
