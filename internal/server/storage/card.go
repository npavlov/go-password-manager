//nolint:dupl
package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"

	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/utils"
)

// StoreCard creates new card record.
func (ds *DBStorage) StoreCard(ctx context.Context, createCard db.StoreCardParams) (*db.Card, error) {
	card, err := ds.Queries.StoreCard(ctx, createCard)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to store card")

		return nil, errors.Wrap(err, "failed to store card")
	}

	return &card, nil
}

// GetCard retrieves card record.
func (ds *DBStorage) GetCard(ctx context.Context, cardID string, userID pgtype.UUID) (*db.Card, error) {
	uuid := utils.GetIDFromString(cardID)

	card, err := ds.Queries.GetCardByID(ctx, db.GetCardByIDParams{
		ID:     uuid,
		UserID: userID,
	})
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to find card")

		return nil, errors.Wrap(err, "failed to find card")
	}

	return &card, nil
}

// GetCards retrieves password record.
func (ds *DBStorage) GetCards(ctx context.Context, userID string) ([]db.Card, error) {
	uuid := utils.GetIDFromString(userID)

	cards, err := ds.Queries.GetCardsByUserID(ctx, uuid)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to find cards")

		return nil, errors.Wrap(err, "failed to find cards")
	}

	return cards, nil
}

func (ds *DBStorage) DeleteCard(ctx context.Context, cardID string, userID pgtype.UUID) error {
	uuid := utils.GetIDFromString(cardID)

	err := ds.Queries.DeleteCard(ctx, db.DeleteCardParams{
		ID:     uuid,
		UserID: userID,
	})
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to delete card")

		return errors.Wrap(err, "failed to delete card")
	}

	return nil
}

// UpdateCard updates card record.
func (ds *DBStorage) UpdateCard(ctx context.Context, updateCard db.UpdateCardParams) (*db.Card, error) {
	card, err := ds.Queries.UpdateCard(ctx, updateCard)
	if err != nil {
		ds.log.Error().Err(err).Msg("failed to update card")

		return nil, errors.Wrap(err, "failed to update card")
	}

	return &card, nil
}
