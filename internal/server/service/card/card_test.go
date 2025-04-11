package card_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	pb "github.com/npavlov/go-password-manager/gen/proto/card"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/card"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	generalutils "github.com/npavlov/go-password-manager/internal/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
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
	ctx := testutils.InjectUserToContext(context.Background(), testUser.ID.String())

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
	require.NotEmpty(t, resp.CardId)
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
		CardId: card.CardId,
	}

	resp, err := svc.GetCard(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "4111111111111111", resp.Card.CardNumber)
	require.Equal(t, "123", resp.Card.Cvv)
	require.Equal(t, "12/30", resp.Card.ExpiryDate)
	require.Equal(t, "John Doe", resp.Card.CardholderName)
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

	req := &pb.DeleteCardRequest{CardId: created.CardId}
	resp, err := svc.DeleteCard(ctx, req)
	require.NoError(t, err)
	require.True(t, resp.Ok)
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
		CardId: created.CardId,
		Data: &pb.CardData{
			CardNumber:     "4000000000000002",
			Cvv:            "999",
			ExpiryDate:     "01/31",
			CardholderName: "Jane Smith",
		},
	}

	resp, err := svc.UpdateCard(ctx, updateReq)
	require.NoError(t, err)
	require.Equal(t, created.CardId, resp.CardId)

	// Optionally verify updated data
	card, err := svc.GetCard(ctx, &pb.GetCardRequest{CardId: created.CardId})
	require.NoError(t, err)
	require.Equal(t, "4000000000000002", card.Card.CardNumber)
	require.Equal(t, "999", card.Card.Cvv)
	require.Equal(t, "01/31", card.Card.ExpiryDate)
	require.Equal(t, "Jane Smith", card.Card.CardholderName)
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

	_, err := svc.StoreCard(context.Background(), req) // no user context
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

	_, err := svc.UpdateCard(context.Background(), req) // no user context
	require.Error(t, err)
	require.Contains(t, err.Error(), "error encrypting card")
}

func TestGetCards_ValidEmpty(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	resp, err := svc.GetCards(ctx, &pb.GetCardsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Cards, 0)
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
