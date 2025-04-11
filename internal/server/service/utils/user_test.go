package utils_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	"github.com/stretchr/testify/require"
)

func TestGetDecryptionKey_Success(t *testing.T) {
	t.Parallel()

	masterKey, _ := utils.GenerateRandomKey()

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

	mockStorage := testutils.SetupMockUserStorage(masterKey)
	mockStorage.AddTestUser(testUser)
	// Inject user ID and encryption key into context
	ctx := testutils.InjectUserToContext(context.Background(), testUser.ID.String())

	userUUID, decryptedKey, err := utils.GetDecryptionKey(ctx, mockStorage, masterKey)

	require.NoError(t, err)
	require.True(t, userUUID.Valid)
	require.NotEmpty(t, decryptedKey)
}

func TestGetDecryptionKey_MissingUserID(t *testing.T) {
	t.Parallel()

	masterKey, _ := utils.GenerateRandomKey()
	mockStorage := testutils.SetupMockUserStorage(masterKey)

	// Use empty context – no user injected
	ctx := context.Background()

	_, _, err := utils.GetDecryptionKey(ctx, mockStorage, masterKey)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error getting user id")
}

func TestGetDecryptionKey_InvalidMasterKey(t *testing.T) {
	t.Parallel()

	masterKey, _ := utils.GenerateRandomKey()
	encryptionKey, _ := utils.GenerateRandomKey()
	encryptionKeyEncrypted, _ := utils.Encrypt(encryptionKey, masterKey)

	testUser := db.User{
		ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Username:      "tester",
		Email:         "test@example.com",
		Password:      "hashed-password",
		EncryptionKey: encryptionKeyEncrypted,
	}

	mockStorage := testutils.SetupMockUserStorage(masterKey)
	mockStorage.AddTestUser(testUser)

	ctx := testutils.InjectUserToContext(context.Background(), testUser.ID.String())

	// Use wrong master key for decryption
	badMasterKey := "wrong-master-key"

	_, _, err := utils.GetDecryptionKey(ctx, mockStorage, badMasterKey)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error getting decryption key")
}
