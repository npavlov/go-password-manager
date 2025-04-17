//nolint:lll,exhaustruct
package storage_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/db"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func TestStoreCard(t *testing.T) {
	t.Parallel()

	storage, mock := testutils.SetupDBStorage(t)

	userID := uuid.New()
	id := uuid.New().String()

	card := db.StoreCardParams{
		UserID:              pgtype.UUID{Bytes: userID, Valid: true},
		EncryptedCardNumber: "123456",
		EncryptedExpiryDate: "12/12",
		EncryptedCvv:        "123",
		CardholderName:      "Ivan Ivanov",
		HashedCardNumber: pgtype.Text{
			String: "123456",
			Valid:  true,
		},
	}

	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"id", "user_id", "encrypted_card_number", "encrypted_expiry_date", "encrypted_cvv",
		"cardholder_name", "created_at", "updated_at", "hashed_card_number",
	}).AddRow(id, userID.String(), card.EncryptedCardNumber, card.EncryptedExpiryDate,
		card.EncryptedCvv, card.CardholderName, now, now, card.HashedCardNumber)

	mock.ExpectQuery("INSERT INTO cards").
		WithArgs(card.UserID, card.HashedCardNumber, card.EncryptedCardNumber,
			card.EncryptedExpiryDate, card.EncryptedCvv, card.CardholderName).
		WillReturnRows(rows)

	result, err := storage.StoreCard(t.Context(), card)
	require.NoError(t, err)
	require.Equal(t, result.UserID.String(), userID.String())
	require.Equal(t, result.EncryptedCardNumber, card.EncryptedCardNumber)
	require.Equal(t, result.EncryptedExpiryDate, card.EncryptedExpiryDate)
	require.Equal(t, result.EncryptedCvv, card.EncryptedCvv)
	require.Equal(t, result.CardholderName, card.CardholderName)
	require.Equal(t, result.HashedCardNumber, card.HashedCardNumber)
}

func TestGetCard(t *testing.T) {
	t.Parallel()

	storage, mock := testutils.SetupDBStorage(t)

	userID := uuid.New()
	cardID := uuid.New()

	expectedCard := db.Card{
		ID:                  pgtype.UUID{Bytes: cardID, Valid: true},
		UserID:              pgtype.UUID{Bytes: userID, Valid: true},
		EncryptedCardNumber: "123456",
		EncryptedExpiryDate: "12/12",
		EncryptedCvv:        "123",
		CardholderName:      "Ivan Ivanov",
		HashedCardNumber: pgtype.Text{
			String: "123456",
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	rows := pgxmock.NewRows([]string{
		"id", "user_id", "encrypted_card_number", "encrypted_expiry_date", "encrypted_cvv",
		"cardholder_name", "created_at", "updated_at", "hashed_card_number",
	}).AddRow(
		expectedCard.ID.String(), expectedCard.UserID.String(), expectedCard.EncryptedCardNumber,
		expectedCard.EncryptedExpiryDate, expectedCard.EncryptedCvv,
		expectedCard.CardholderName, expectedCard.CreatedAt, expectedCard.UpdatedAt,
		expectedCard.HashedCardNumber,
	)

	mock.ExpectQuery("SELECT (.+) FROM cards").
		WithArgs(expectedCard.ID, expectedCard.UserID).
		WillReturnRows(rows)

	result, err := storage.GetCard(t.Context(), cardID.String(), pgtype.UUID{Bytes: userID, Valid: true})
	require.NoError(t, err)
	require.Equal(t, expectedCard.ID.String(), result.ID.String())
	require.Equal(t, expectedCard.UserID, result.UserID)
	require.Equal(t, expectedCard.EncryptedCardNumber, result.EncryptedCardNumber)
	require.Equal(t, expectedCard.EncryptedExpiryDate, result.EncryptedExpiryDate)
	require.Equal(t, expectedCard.EncryptedCvv, result.EncryptedCvv)
	require.Equal(t, expectedCard.CardholderName, result.CardholderName)
}

func TestGetCards(t *testing.T) {
	t.Parallel()

	storage, mock := testutils.SetupDBStorage(t)

	userID := uuid.New()

	expectedCards := []db.Card{
		{
			ID:                  pgtype.UUID{Bytes: uuid.New(), Valid: true},
			UserID:              pgtype.UUID{Bytes: userID, Valid: true},
			EncryptedCardNumber: "123456",
			EncryptedExpiryDate: "12/12",
			EncryptedCvv:        "123",
			CardholderName:      "Ivan Ivanov",
			HashedCardNumber: pgtype.Text{
				String: "123456",
				Valid:  true,
			},
			CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		},
		{
			ID:                  pgtype.UUID{Bytes: uuid.New(), Valid: true},
			UserID:              pgtype.UUID{Bytes: userID, Valid: true},
			EncryptedCardNumber: "654321",
			EncryptedExpiryDate: "01/25",
			EncryptedCvv:        "456",
			CardholderName:      "Petr Petrov",
			HashedCardNumber: pgtype.Text{
				String: "123456",
				Valid:  true,
			},
			CreatedAt: pgtype.Timestamp{Time: time.Now().Add(-time.Hour), Valid: true},
			UpdatedAt: pgtype.Timestamp{Time: time.Now().Add(-time.Hour), Valid: true},
		},
	}

	rows := pgxmock.NewRows([]string{
		"id", "user_id", "encrypted_card_number", "encrypted_expiry_date", "encrypted_cvv",
		"cardholder_name", "created_at", "updated_at", "hashed_card_number",
	})
	for _, card := range expectedCards {
		rows.AddRow(
			card.ID, card.UserID, card.EncryptedCardNumber,
			card.EncryptedExpiryDate, card.EncryptedCvv, card.CardholderName,
			card.CreatedAt, card.UpdatedAt, card.HashedCardNumber,
		)
	}

	mock.ExpectQuery("SELECT (.+) FROM cards").
		WithArgs(expectedCards[0].UserID).
		WillReturnRows(rows)

	result, err := storage.GetCards(t.Context(), userID.String())
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, expectedCards[0].ID, result[0].ID)
	require.Equal(t, expectedCards[1].ID, result[1].ID)
	require.Equal(t, expectedCards[0].CardholderName, result[0].CardholderName)
	require.Equal(t, expectedCards[1].CardholderName, result[1].CardholderName)
}

func TestDeleteCard(t *testing.T) {
	t.Parallel()

	storage, mock := testutils.SetupDBStorage(t)

	userID := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	cardID := pgtype.UUID{Bytes: uuid.New(), Valid: true}

	mock.ExpectExec("DELETE FROM cards").
		WithArgs(cardID, userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := storage.DeleteCard(t.Context(), cardID.String(), userID)
	require.NoError(t, err)
}

func TestUpdateCard(t *testing.T) {
	t.Parallel()

	storage, mock := testutils.SetupDBStorage(t)

	userID := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	cardID := pgtype.UUID{Bytes: uuid.New(), Valid: true}

	updateParams := db.UpdateCardParams{
		ID:                  cardID,
		EncryptedCardNumber: "new_encrypted_number",
		EncryptedExpiryDate: "new_expiry",
		EncryptedCvv:        "new_cvv",
		CardholderName:      "New Name",
		HashedCardNumber: pgtype.Text{
			String: "123456",
			Valid:  true,
		},
	}

	expectedCard := db.Card{
		ID:                  cardID,
		UserID:              userID,
		EncryptedCardNumber: updateParams.EncryptedCardNumber,
		EncryptedExpiryDate: updateParams.EncryptedExpiryDate,
		EncryptedCvv:        updateParams.EncryptedCvv,
		CardholderName:      updateParams.CardholderName,
		HashedCardNumber:    updateParams.HashedCardNumber,
		CreatedAt:           pgtype.Timestamp{Time: time.Now().Add(-time.Hour), Valid: true},
		UpdatedAt:           pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	rows := pgxmock.NewRows([]string{
		"id", "user_id", "encrypted_card_number", "encrypted_expiry_date", "encrypted_cvv",
		"cardholder_name", "created_at", "updated_at", "hashed_card_number",
	}).AddRow(
		expectedCard.ID.String(), expectedCard.UserID.String(), expectedCard.EncryptedCardNumber,
		expectedCard.EncryptedExpiryDate, expectedCard.EncryptedCvv,
		expectedCard.CardholderName, expectedCard.CreatedAt, expectedCard.UpdatedAt,
		expectedCard.HashedCardNumber,
	)

	mock.ExpectQuery("UPDATE cards").
		WithArgs(
			updateParams.EncryptedCardNumber,
			updateParams.EncryptedExpiryDate,
			updateParams.EncryptedCvv,
			updateParams.CardholderName,
			updateParams.HashedCardNumber,
			updateParams.ID,
		).
		WillReturnRows(rows)

	result, err := storage.UpdateCard(t.Context(), updateParams)
	require.NoError(t, err)
	require.Equal(t, expectedCard.ID, result.ID)
	require.Equal(t, expectedCard.EncryptedCardNumber, result.EncryptedCardNumber)
	require.Equal(t, expectedCard.EncryptedExpiryDate, result.EncryptedExpiryDate)
	require.Equal(t, expectedCard.EncryptedCvv, result.EncryptedCvv)
	require.Equal(t, expectedCard.CardholderName, result.CardholderName)
	require.Equal(t, expectedCard.HashedCardNumber, result.HashedCardNumber)
}

func TestGetCard_NotFound(t *testing.T) {
	t.Parallel()

	storage, mock := testutils.SetupDBStorage(t)

	userID := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	cardID := pgtype.UUID{Bytes: uuid.New(), Valid: true}

	mock.ExpectQuery("SELECT (.+) FROM cards").
		WithArgs(cardID, userID).
		WillReturnError(pgx.ErrNoRows)

	result, err := storage.GetCard(t.Context(), cardID.String(), userID)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "failed to find card")
}

func TestDeleteCard_NotFound(t *testing.T) {
	t.Parallel()

	storage, mock := testutils.SetupDBStorage(t)

	userID := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	cardID := pgtype.UUID{Bytes: uuid.New(), Valid: true}

	mock.ExpectExec("DELETE FROM cards").
		WithArgs(cardID, userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err := storage.DeleteCard(t.Context(), cardID.String(), userID)
	require.NoError(t, err) // Note: This depends on your business logic - you might want to return an error if no rows were deleted
}
