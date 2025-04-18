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

	req := &pb.StoreCardV1Request{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	}

	resp, err := svc.StoreCardV1(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, resp.GetCardId())
}

func TestStoreCard_InvalidInput(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.StoreCardV1Request{
		Card: &pb.CardData{}, // Missing required fields
	}

	_, err := svc.StoreCardV1(ctx, req)
	require.Error(t, err)
}

func TestGetCard_Success(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	card, err := svc.StoreCardV1(ctx, &pb.StoreCardV1Request{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	})
	require.NoError(t, err)

	req := &pb.GetCardV1Request{
		CardId: card.GetCardId(),
	}

	resp, err := svc.GetCardV1(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "4111111111111111", resp.GetCard().GetCardNumber())
	require.Equal(t, "123", resp.GetCard().GetCvv())
	require.Equal(t, "12/30", resp.GetCard().GetExpiryDate())
	require.Equal(t, "John Doe", resp.GetCard().GetCardholderName())
}

func TestDeleteCard_Success(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	created, err := svc.StoreCardV1(ctx, &pb.StoreCardV1Request{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	})
	require.NoError(t, err)

	req := &pb.DeleteCardV1Request{CardId: created.GetCardId()}
	resp, err := svc.DeleteCardV1(ctx, req)
	require.NoError(t, err)
	require.True(t, resp.GetOk())
}

func TestUpdateCard_Success(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	created, err := svc.StoreCardV1(ctx, &pb.StoreCardV1Request{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	})
	require.NoError(t, err)

	updateReq := &pb.UpdateCardV1Request{
		CardId: created.GetCardId(),
		Data: &pb.CardData{
			CardNumber:     "4000000000000002",
			Cvv:            "999",
			ExpiryDate:     "01/31",
			CardholderName: "Jane Smith",
		},
	}

	resp, err := svc.UpdateCardV1(ctx, updateReq)
	require.NoError(t, err)
	require.Equal(t, created.GetCardId(), resp.GetCardId())

	// Optionally verify updated data
	card, err := svc.GetCardV1(ctx, &pb.GetCardV1Request{CardId: created.GetCardId()})
	require.NoError(t, err)
	require.Equal(t, "4000000000000002", card.GetCard().GetCardNumber())
	require.Equal(t, "999", card.GetCard().GetCvv())
	require.Equal(t, "01/31", card.GetCard().GetExpiryDate())
	require.Equal(t, "Jane Smith", card.GetCard().GetCardholderName())
}

func TestUpdateCard_Invalid(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.UpdateCardV1Request{
		CardId: "invalid-uuid",
		Data:   &pb.CardData{}, // invalid data
	}

	_, err := svc.UpdateCardV1(ctx, req)
	require.Error(t, err)
}

func TestGetCards_Empty(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.GetCardsV1Request{}
	resp, err := svc.GetCardsV1(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestStoreCard_MissingContext(t *testing.T) {
	t.Parallel()

	svc, _, _ := setupCardService(t)

	req := &pb.StoreCardV1Request{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	}

	_, err := svc.StoreCardV1(t.Context(), req) // no user context
	require.Error(t, err)
}

func TestDeleteCard_InvalidId(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.DeleteCardV1Request{CardId: "invalid-id"}
	_, err := svc.DeleteCardV1(ctx, req)
	require.Error(t, err)
}

func TestUpdateCard_GetUserIdFailure(t *testing.T) {
	t.Parallel()

	svc, _, _ := setupCardService(t)

	req := &pb.UpdateCardV1Request{
		CardId: uuid.NewString(),
		Data: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	}

	_, err := svc.UpdateCardV1(t.Context(), req) // no user context
	require.Error(t, err)
	require.Contains(t, err.Error(), "error getting decrypted user UUID")
}

func TestGetCards_ValidEmpty(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	resp, err := svc.GetCardsV1(ctx, &pb.GetCardsV1Request{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Empty(t, resp.GetCards())
}

func TestDeleteCard_InvalidUUID(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	_, err := svc.DeleteCardV1(ctx, &pb.DeleteCardV1Request{
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

	req := &pb.StoreCardV1Request{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	}

	// First store should succeed
	_, err := svc.StoreCardV1(ctx, req)
	require.NoError(t, err)

	// Second store with same card number should fail
	_, err = svc.StoreCardV1(ctx, req)
	require.Error(t, err)
}

func TestGetCard_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.GetCardV1Request{
		CardId: uuid.NewString(), // Non-existent ID
	}

	_, err := svc.GetCardV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestGetCard_WrongUser(t *testing.T) {
	t.Parallel()

	svc, storage, ctx := setupCardService(t)

	// Create card for first user
	card, err := svc.StoreCardV1(ctx, &pb.StoreCardV1Request{
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
	req := &pb.GetCardV1Request{CardId: card.GetCardId()}
	_, err = svc.GetCardV1(ctx2, req)
	require.Error(t, err)
}

func TestUpdateCard_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.UpdateCardV1Request{
		CardId: uuid.NewString(), // Non-existent ID
		Data: &pb.CardData{
			CardNumber:     "4000000000000002",
			Cvv:            "999",
			ExpiryDate:     "01/31",
			CardholderName: "Jane Smith",
		},
	}

	_, err := svc.UpdateCardV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestDeleteCard_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupCardService(t)

	req := &pb.DeleteCardV1Request{
		CardId: uuid.NewString(), // Non-existent ID
	}

	_, err := svc.DeleteCardV1(ctx, req)
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

	req := &pb.StoreCardV1Request{
		Card: &pb.CardData{
			CardNumber:     "4111111111111111",
			Cvv:            "123",
			ExpiryDate:     "12/30",
			CardholderName: "John Doe",
		},
	}

	_, err := svc.StoreCardV1(ctx, req)
	require.Error(t, err)
}

func TestGetCard_DecryptionFailure(t *testing.T) {
	t.Parallel()

	svc, storage, ctx := setupCardService(t)

	// Store card normally
	card, err := svc.StoreCardV1(ctx, &pb.StoreCardV1Request{
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

	req := &pb.GetCardV1Request{CardId: card.GetCardId()}
	_, err = svc.GetCardV1(ctx, req)
	require.Error(t, err)
}

func TestCardValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		card        *pb.CardData
		shouldError bool
		errContains string
	}{
		{
			name: "valid card",
			card: &pb.CardData{
				CardNumber:     "4111111111111111",
				Cvv:            "123",
				ExpiryDate:     "12/30",
				CardholderName: "John Doe",
			},
			shouldError: false,
		},
		{
			name: "invalid card number - too short",
			card: &pb.CardData{
				CardNumber:     "411111",
				Cvv:            "123",
				ExpiryDate:     "12/30",
				CardholderName: "John Doe",
			},
			shouldError: true,
			errContains: "value does not match regex pattern",
		},
		{
			name: "invalid card number - non-numeric",
			card: &pb.CardData{
				CardNumber:     "4111-1111-1111-1111",
				Cvv:            "123",
				ExpiryDate:     "12/30",
				CardholderName: "John Doe",
			},
			shouldError: true,
			errContains: "value does not match regex pattern",
		},
		{
			name: "invalid expiry date - wrong format",
			card: &pb.CardData{
				CardNumber:     "4111111111111111",
				Cvv:            "123",
				ExpiryDate:     "12-30",
				CardholderName: "John Doe",
			},
			shouldError: true,
			errContains: "value does not match regex pattern",
		},
		{
			name: "invalid expiry date - invalid month",
			card: &pb.CardData{
				CardNumber:     "4111111111111111",
				Cvv:            "123",
				ExpiryDate:     "13/30",
				CardholderName: "John Doe",
			},
			shouldError: true,
			errContains: "value does not match regex pattern",
		},
		{
			name: "invalid CVV - too short",
			card: &pb.CardData{
				CardNumber:     "4111111111111111",
				Cvv:            "12",
				ExpiryDate:     "12/30",
				CardholderName: "John Doe",
			},
			shouldError: true,
			errContains: "value does not match regex pattern ",
		},
		{
			name: "invalid CVV - too long",
			card: &pb.CardData{
				CardNumber:     "4111111111111111",
				Cvv:            "12345",
				ExpiryDate:     "12/30",
				CardholderName: "John Doe",
			},
			shouldError: true,
			errContains: "value does not match regex pattern",
		},
		{
			name: "invalid CVV - non-numeric",
			card: &pb.CardData{
				CardNumber:     "4111111111111111",
				Cvv:            "12a",
				ExpiryDate:     "12/30",
				CardholderName: "John Doe",
			},
			shouldError: true,
			errContains: "value does not match regex pattern",
		},
		{
			name: "missing cardholder name",
			card: &pb.CardData{
				CardNumber:     "4111111111111111",
				Cvv:            "123",
				ExpiryDate:     "12/30",
				CardholderName: "",
			},
			shouldError: true,
			errContains: "value length must be at least 1 characters",
		},
		{
			name: "cardholder name too long",
			card: &pb.CardData{
				CardNumber: "4111111111111111",
				Cvv:        "123",
				ExpiryDate: "12/30",
				//nolint:lll
				CardholderName: "This name is way too long and exceeds the maximum allowed length of 100 characters which should trigger a validation error",
			},
			shouldError: true,
			errContains: "value length must be at most 100 characters",
		},
	}

	svc, _, ctx := setupCardService(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Test StoreCard validation
			storeReq := &pb.StoreCardV1Request{Card: tc.card}
			_, err := svc.StoreCardV1(ctx, storeReq)

			if tc.shouldError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errContains)
			} else {
				require.NoError(t, err)
			}

			// Test UpdateCard validation (only if we have a valid card to update)
			if !tc.shouldError {
				// First create a card to update
				storeResp, err := svc.StoreCardV1(ctx, &pb.StoreCardV1Request{
					Card: &pb.CardData{
						CardNumber:     "5555555555554444", // Different card number
						Cvv:            "321",
						ExpiryDate:     "01/25",
						CardholderName: "Initial Name",
					},
				})
				require.NoError(t, err)

				// Now try to update with test case data
				updateReq := &pb.UpdateCardV1Request{
					CardId: storeResp.GetCardId(),
					Data:   tc.card,
				}
				_, err = svc.UpdateCardV1(ctx, updateReq)

				if tc.shouldError {
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.errContains)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}
