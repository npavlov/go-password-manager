//nolint:exhaustruct
package card_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	pb "github.com/npavlov/go-password-manager/gen/proto/card"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/card"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	generalutils "github.com/npavlov/go-password-manager/internal/utils"
)

func setupCardService(t *testing.T) (*card.Service, *testutils.MockDBStorage, context.Context) {
	t.Helper()

	logger := zerolog.New(nil)
	masterKey, _ := utils.GenerateRandomKey()

	cfg := &config.Config{
		SecuredMasterKey: generalutils.NewString(masterKey),
	}

	storage := testutils.SetupMockUserStorage(masterKey)
	svc := card.NewCardService(&logger, storage, cfg)
	encryptionKey, _ := utils.GenerateRandomKey()

	encryptionKeyEncrypted, _ := utils.Encrypt(encryptionKey, masterKey)

	// Create test user
	testUser := db.User{
		ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Username:      "tester",
		Email:         "test@example.com",
		Password:      "hashed-password",
		EncryptionKey: encryptionKeyEncrypted,
	}
	storage.AddTestUser(testUser)

	// Inject user ID and encryption key into context
	ctx := testutils.InjectUserToContext(t.Context(), testUser.ID.String())

	return svc, storage, ctx
}

func TestStoreCard_Success(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.StoreCardRequest{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	}

	resp, err := svc.StoreCard(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, resp.GetCardId())
}

func TestStoreCard_InvalidInput(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.StoreCardRequest{
		Card: &pb.CardData{}, // Missing required fields
	}

	_, err := svc.StoreCard(ctx, req)
	require.Error(t, err)
}

func TestGetCard_Success(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	card, err := svc.StoreCard(ctx, &pb.StoreCardRequest{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	})
	require.NoError(t, err)

	req := &pb.GetCardRequest{
		CardId: card.GetCardId(),
	}

	resp, err := svc.GetCard(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "4111111111111111", resp.GetCard().GetCardNumber())
	require.Equal(t, "123", resp.GetCard().GetCvv())
	require.Equal(t, "12/30", resp.GetCard().GetExpiryDate())
	require.Equal(t, "John Doe", resp.GetCard().GetCardholderName())
}

func TestDeleteCard_Success(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	created, err := svc.StoreCard(ctx, &pb.StoreCardRequest{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	})
	require.NoError(t, err)

	req := &pb.DeleteCardRequest{CardId: created.GetCardId()}
	resp, err := svc.DeleteCard(ctx, req)
	require.NoError(t, err)
	require.True(t, resp.GetOk())
}

func TestUpdateCard_Success(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	created, err := svc.StoreCard(ctx, &pb.StoreCardRequest{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	})
	require.NoError(t, err)

	updateReq := &pb.UpdateCardRequest{
		CardId: created.GetCardId(),
		Data: &pb.CardData{
			CardNumber:     "4000000000000002",
			Cvv:            "999",
			ExpiryDate:     "01/31",
			CardholderName: "Jane Smith",
		},
	}

	resp, err := svc.UpdateCard(ctx, updateReq)
	require.NoError(t, err)
	require.Equal(t, created.GetCardId(), resp.GetCardId())

	// Optionally verify updated data
	card, err := svc.GetCard(ctx, &pb.GetCardRequest{CardId: created.GetCardId()})
	require.NoError(t, err)
	require.Equal(t, "4000000000000002", card.GetCard().GetCardNumber())
	require.Equal(t, "999", card.GetCard().GetCvv())
	require.Equal(t, "01/31", card.GetCard().GetExpiryDate())
	require.Equal(t, "Jane Smith", card.GetCard().GetCardholderName())
}

func TestUpdateCard_Invalid(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.UpdateCardRequest{
		CardId: "invalid-uuid",
		Data:   &pb.CardData{}, // invalid data
	}

	_, err := svc.UpdateCard(ctx, req)
	require.Error(t, err)
}

func TestGetCards_Empty(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.GetCardsRequest{}
	resp, err := svc.GetCards(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestStoreCard_MissingContext(t *testing.T) {
	t.Parallel()

	svc, _, _ := setupCardService(t)

	req := &pb.StoreCardRequest{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	}

	_, err := svc.StoreCard(t.Context(), req) // no user context
	require.Error(t, err)
}

func TestDeleteCard_InvalidId(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.DeleteCardRequest{CardId: "invalid-id"}
	_, err := svc.DeleteCard(ctx, req)
	require.Error(t, err)
}

func TestUpdateCard_GetUserIdFailure(t *testing.T) {
	t.Parallel()

	svc, _, _ := setupCardService(t)

	req := &pb.UpdateCardRequest{
		CardId: uuid.NewString(),
		Data: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	}

	_, err := svc.UpdateCard(t.Context(), req) // no user context
	require.Error(t, err)
	require.Contains(t, err.Error(), "error getting decrypted user UUID")
}

func TestGetCards_ValidEmpty(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	resp, err := svc.GetCards(ctx, &pb.GetCardsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Empty(t, resp.GetCards())
}

func TestDeleteCard_InvalidUUID(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	_, err := svc.DeleteCard(ctx, &pb.DeleteCardRequest{
		CardId: "not-a-uuid",
	})
	require.Error(t, err)
}

func TestEncryptCard_Success(t *testing.T) {
	t.Parallel()

	svc, _, _ := setupCardService(t)

	encryptionKey, err := utils.GenerateRandomKey()
	require.NoError(t, err)

	encryptedCardNumber, encryptedCVV, encryptedExpiryDate, err := svc.EncryptCard(
		encryptionKey,
		"4111111111111111",
		"123",
		"12/30",
	)

	require.NoError(t, err)
	require.NotEmpty(t, encryptedCardNumber)
	require.NotEmpty(t, encryptedCVV)
	require.NotEmpty(t, encryptedExpiryDate)
}

func TestStoreCard_DuplicateCardNumber(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.StoreCardRequest{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	}

	// First store should succeed
	_, err := svc.StoreCard(ctx, req)
	require.NoError(t, err)

	// Second store with same card number should fail
	_, err = svc.StoreCard(ctx, req)
	require.Error(t, err)
}

func TestGetCard_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.GetCardRequest{
		CardId: uuid.NewString(), // Non-existent ID
	}

	_, err := svc.GetCard(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestGetCard_WrongUser(t *testing.T) {
	t.Parallel()

	svc, storage, ctx := setupCardService(t)

	// Create card for first user
	card, err := svc.StoreCard(ctx, &pb.StoreCardRequest{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	})
	require.NoError(t, err)

	// Create second user
	user2 := db.User{
		ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Username:      "tester2",
		Email:         "test2@example.com",
		Password:      "hashed-password2",
		EncryptionKey: "enc-key-2",
	}
	storage.AddTestUser(user2)
	ctx2 := testutils.InjectUserToContext(t.Context(), user2.ID.String())

	// Try to access first user's card with second user
	req := &pb.GetCardRequest{CardId: card.GetCardId()}
	_, err = svc.GetCard(ctx2, req)
	require.Error(t, err)
}

func TestUpdateCard_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.UpdateCardRequest{
		CardId: uuid.NewString(), // Non-existent ID
		Data: &pb.CardData{
			CardNumber:     "4000000000000002",
			Cvv:            "999",
			ExpiryDate:     "01/31",
			CardholderName: "Jane Smith",
		},
	}

	_, err := svc.UpdateCard(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestDeleteCard_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.DeleteCardRequest{
		CardId: uuid.NewString(), // Non-existent ID
	}

	_, err := svc.DeleteCard(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestEncryptCard_InvalidKey(t *testing.T) {
	t.Parallel()

	svc, _, _ := setupCardService(t)

	//nolint:dogsled
	_, _, _, err := svc.EncryptCard(
		"invalid-key", // Invalid encryption key
		"4111111111111111",
		"123",
		"12/30",
	)
	require.Error(t, err)
}

func TestStoreCard_EncryptionFailure(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	// Force encryption failure by providing invalid user key
	//nolint:revive,staticcheck
	ctx = context.WithValue(ctx, "user_id", uuid.NewString())

	req := &pb.StoreCardRequest{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	}

	_, err := svc.StoreCard(ctx, req)
	require.Error(t, err)
}

func TestGetCard_DecryptionFailure(t *testing.T) {
	t.Parallel()

	svc, storage, ctx := setupCardService(t)

	// Store card normally
	card, err := svc.StoreCard(ctx, &pb.StoreCardRequest{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	})
	require.NoError(t, err)

	// Force decryption failure by changing user's encryption key
	userID := testutils.GetUserIDFromContext(ctx)
	userGUID := pgtype.UUID{
		Bytes: uuid.MustParse(userID),
		Valid: true,
	}
	newKey, _ := utils.GenerateRandomKey()

	newMasterKey, _ := utils.GenerateRandomKey()
	user := storage.UsersByID[userGUID]
	user.EncryptionKey, _ = utils.Encrypt(newKey, newMasterKey)
	storage.UsersByID[userGUID] = user

	req := &pb.GetCardRequest{CardId: card.GetCardId()}
	_, err = svc.GetCard(ctx, req)
	require.Error(t, err)
}
