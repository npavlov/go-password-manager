package storage

import (
	"context"

	"github.com/pkg/errors"

	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/utils"
)

// AddMeta add meta data.
func (ds *DBStorage) AddMeta(ctx context.Context, recordId, key, value string) (*db.Metainfo, error) {
	uuid := utils.GetIDFromString(recordId)

	meta, err := ds.Queries.AddMetaInfo(ctx, db.AddMetaInfoParams{
		ItemID: uuid,
		Key:    key,
		Value:  value,
	})
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to get items")

		return nil, errors.Wrap(err, "failed to get items")
	}

	return &meta, nil
}

func (ds *DBStorage) GetMetaInfo(ctx context.Context, recordId string) ([]db.GetMetaInfoByItemIDRow, error) {
	uuid := utils.GetIDFromString(recordId)

	meta, err := ds.Queries.GetMetaInfoByItemID(ctx, uuid)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to get items")

		return nil, errors.Wrap(err, "failed to get items")
	}

	return meta, nil
}

func (ds *DBStorage) DeleteMetaInfo(ctx context.Context, key, itemId string) error {
	uuid := utils.GetIDFromString(itemId)

	err := ds.Queries.DeleteMetaInfo(ctx, db.DeleteMetaInfoParams{
		Key:    key,
		ItemID: uuid,
	})
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to delete items")

		return errors.Wrap(err, "failed to delete items")
	}

	return nil
}
