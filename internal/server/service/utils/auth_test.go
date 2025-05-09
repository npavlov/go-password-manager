//nolint:exhaustruct,revive,staticcheck,forcetypeassert,wrapcheck
package utils_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
)

// ---- MOCK STORAGE ----

type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) GetUserByID(ctx context.Context, id pgtype.UUID) (*db.User, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*db.User), args.Error(1)
}

// ---- TESTS ----

func TestGenerateAndValidateJWT(t *testing.T) {
	t.Parallel()

	userID := "test-user-id"
	secret := "my-secret"
	exp := time.Now().Add(time.Hour).Unix()

	token, err := utils.GenerateJWT(userID, secret, exp)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the token
	validatedUserID, err := utils.ValidateJWT(token, secret)
	require.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	t.Parallel()

	_, err := utils.ValidateJWT("bad-token", "my-secret")
	require.Error(t, err)
}

func TestValidateJWT_InvalidSecret(t *testing.T) {
	t.Parallel()

	secret := "correct-secret"
	wrongSecret := "wrong-secret"

	token, err := utils.GenerateJWT("user-id", secret, time.Now().Add(time.Hour).Unix())
	require.NoError(t, err)

	_, err = utils.ValidateJWT(token, wrongSecret)
	require.Error(t, err)
}

func TestGetUserId_Success(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(t.Context(), "user_id", "550e8400-e29b-41d4-a716-446655440000")

	uuid, err := utils.GetUserID(ctx)
	require.NoError(t, err)

	// Just check if status is valid (UUID in pgtype is not nil)
	assert.True(t, uuid.Valid)
}

func TestGetUserId_Missing(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	_, err := utils.GetUserID(ctx)
	require.Error(t, err)
}

func TestGetUserId_Invalid(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(t.Context(), "user_id", "not-a-uuid")
	_, err := utils.GetUserID(ctx)
	require.Error(t, err)
}

func TestGetUserKey_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := t.Context()
	userUUID := pgtype.UUID{}
	err := userUUID.Scan("550e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)

	mockedStorage := new(mockStorage)
	masterKey := "yIdg5TUnZsQH8anXm18Uss18Q7CAt3lvPQ2wg9JDXZY="

	// Create an encrypted key using real Encrypt function
	plainUserKey := "user-private-key"
	encryptedKey, err := utils.Encrypt(plainUserKey, masterKey)
	require.NoError(t, err)

	mockedUser := &db.User{EncryptionKey: encryptedKey}
	mockedStorage.On("GetUserByID", mock.Anything, userUUID).Return(mockedUser, nil)

	// Act
	decryptedKey, err := utils.GetUserKey(ctx, mockedStorage, userUUID, masterKey)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, plainUserKey, decryptedKey)
	mockedStorage.AssertCalled(t, "GetUserByID", ctx, userUUID)
}

func TestGetUserKey_Failure(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	userUUID := pgtype.UUID{}
	err := userUUID.Scan("550e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)

	mockedStorage := new(mockStorage)
	mockedUser := &db.User{}
	mockedStorage.On("GetUserByID", mock.Anything, userUUID).Return(mockedUser, assert.AnError)

	_, err = utils.GetUserKey(ctx, mockedStorage, userUUID, "some-master-key")
	require.Error(t, err)
}
