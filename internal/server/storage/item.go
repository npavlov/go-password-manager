package storage

import (
	"context"

	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/pkg/errors"
)

// GetItems get all type of docs
func (ds *DBStorage) GetItems(ctx context.Context, getItems db.GetItemsByUserIDParams) ([]db.GetItemsByUserIDRow, error) {
	items, err := ds.Queries.GetItemsByUserID(ctx, getItems)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to create user")

		return nil, errors.Wrap(err, "failed to create user")
	}

	return items, nil
}
