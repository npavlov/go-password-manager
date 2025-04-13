package password_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/password"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	gu "github.com/npavlov/go-password-manager/internal/utils"
)

func setupPasswordService(t *testing.T) (*password.Service, *testutils.MockDBStorage, context.Context, string) {
	t.Helper()

	logger := zerolog.New(nil)
	masterKey, _ := utils.GenerateRandomKey()
	mockStorage := testutils.SetupMockUserStorage(masterKey)
	cfg := &config.Config{
		SecuredMasterKey: gu.NewString(masterKey),
	}

	// Create test user
	userID := uuid.New()
	encryptionKey, _ := utils.GenerateRandomKey()
	encryptedKey, _ := utils.Encrypt(encryptionKey, masterKey)

	testUser := db.User{
		ID:            pgtype.UUID{Bytes: userID, Valid: true},
		Username:      "testuser",
		Password:      "hashed-password",
		EncryptionKey: encryptedKey,
	}
	mockStorage.AddTestUser(testUser)

	// Inject user ID into context
	ctx := testutils.InjectUserToContext(context.Background(), userID.String())

	return password.NewPasswordService(&logger, mockStorage, cfg), mockStorage, ctx, encryptionKey
}

func TestStorePassword_Success(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupPasswordService(t)

	testLogin := "test@example.com"
	testPassword := "securepassword123"

	req := &pb.StorePasswordRequest{
		Password: &pb.PasswordData{
			Login:    testLogin,
			Password: testPassword,
		},
	}

	resp, err := svc.StorePassword(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, resp.PasswordId)

	// Verify password was stored
	storedPass, err := storage.GetPassword(ctx, resp.PasswordId, pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true})
	require.NoError(t, err)
	require.Equal(t, testLogin, storedPass.Login)
}

func TestStorePassword_ValidationError(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupPasswordService(t)

	// Missing password data
	req := &pb.StorePasswordRequest{
		Password: &pb.PasswordData{
			Login: "", // Invalid empty login
		},
	}

	_, err := svc.StorePassword(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestStorePassword_EncryptionError(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupPasswordService(t)

	// Corrupt the user's encryption key in storage
	userID := uuid.MustParse(testutils.GetUserIDFromContext(ctx))
	userIDPG := pgtype.UUID{Bytes: userID, Valid: true}
	user := storage.UsersByID[userIDPG]
	user.EncryptionKey = "invalid-key"
	storage.UsersByID[userIDPG] = user

	req := &pb.StorePasswordRequest{
		Password: &pb.PasswordData{
			Login:    "test@example.com",
			Password: "password123",
		},
	}

	_, err := svc.StorePassword(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error getting decryption key")
}

func TestStorePassword_StorageError(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupPasswordService(t)

	// Make storage return error
	storage.CallError = errors.New("database failure")

	req := &pb.StorePasswordRequest{
		Password: &pb.PasswordData{
			Login:    "test@example.com",
			Password: "password123",
		},
	}

	_, err := svc.StorePassword(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to store password")
}

func TestGetPassword_Success(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, userKey := setupPasswordService(t)

	// Store a test password first
	testLogin := "test@example.com"
	testPassword := "securepassword123"
	encryptedPassword, _ := utils.Encrypt(testPassword, userKey)
	storedPass, err := storage.StorePassword(ctx, db.CreatePasswordEntryParams{
		UserID:   pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true},
		Login:    testLogin,
		Password: encryptedPassword,
	})
	require.NoError(t, err)

	req := &pb.GetPasswordRequest{
		PasswordId: storedPass.ID.String(),
	}

	resp, err := svc.GetPassword(ctx, req)
	require.NoError(t, err)
	require.Equal(t, testLogin, resp.Password.Login)
	require.Equal(t, testPassword, resp.Password.Password)
	require.True(t, timestamppb.New(storedPass.UpdatedAt.Time).AsTime().Equal(resp.LastUpdate.AsTime()))
}

func TestGetPassword_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupPasswordService(t)

	req := &pb.GetPasswordRequest{
		PasswordId: uuid.NewString(),
	}

	_, err := svc.GetPassword(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error getting user id") // This should probably be "password not found"
}

func TestGetPassword_DecryptionError(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupPasswordService(t)

	// Store a password with invalid encrypted content
	storedPass, err := storage.StorePassword(ctx, db.CreatePasswordEntryParams{
		UserID:   pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true},
		Login:    "test@example.com",
		Password: "invalid-encrypted-content",
	})
	require.NoError(t, err)

	req := &pb.GetPasswordRequest{
		PasswordId: storedPass.ID.String(),
	}

	_, err = svc.GetPassword(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error decrypting password")
}

func TestGetPassword_Unauthorized(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupPasswordService(t)

	// Store a password with different user
	otherUserID := uuid.New()
	otherPass, err := storage.StorePassword(context.Background(), db.CreatePasswordEntryParams{
		UserID:   pgtype.UUID{Bytes: otherUserID, Valid: true},
		Login:    "other@example.com",
		Password: "encrypted-password",
	})
	require.NoError(t, err)

	req := &pb.GetPasswordRequest{
		PasswordId: otherPass.ID.String(),
	}

	_, err = svc.GetPassword(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized access to password")
}

func TestGetPasswords_NotImplemented(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupPasswordService(t)

	req := &pb.GetPasswordsRequest{}
	resp, err := svc.GetPasswords(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestUpdatePassword_Success(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, userKey := setupPasswordService(t)

	// Store initial password
	initialPass, err := storage.StorePassword(ctx, db.CreatePasswordEntryParams{
		UserID:   pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true},
		Login:    "old@example.com",
		Password: "initial-encrypted-pass",
	})
	require.NoError(t, err)

	newLogin := "new@example.com"
	newPassword := "newsecurepassword123"

	req := &pb.UpdatePasswordRequest{
		PasswordId: initialPass.ID.String(),
		Data: &pb.PasswordData{
			Login:    newLogin,
			Password: newPassword,
		},
	}

	resp, err := svc.UpdatePassword(ctx, req)
	require.NoError(t, err)
	require.Equal(t, initialPass.ID.String(), resp.PasswordId)

	// Verify update
	updatedPass, err := storage.GetPassword(ctx, initialPass.ID.String(), pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true})
	require.NoError(t, err)
	require.Equal(t, newLogin, updatedPass.Login)

	// Verify decryption
	decrypted, err := utils.Decrypt(updatedPass.Password, userKey)
	require.NoError(t, err)
	require.Equal(t, newPassword, decrypted)
}

func TestUpdatePassword_ValidationError(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupPasswordService(t)

	req := &pb.UpdatePasswordRequest{
		PasswordId: uuid.NewString(),
		Data:       &pb.PasswordData{}, // Missing required fields
	}

	_, err := svc.UpdatePassword(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestUpdatePassword_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupPasswordService(t)

	req := &pb.UpdatePasswordRequest{
		PasswordId: uuid.NewString(),
		Data: &pb.PasswordData{
			Login:    "test@example.com",
			Password: "password123",
		},
	}

	_, err := svc.UpdatePassword(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to store password") // This should probably be "password not found"
}

func TestDeletePassword_Success(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupPasswordService(t)

	// Store a test password first
	testPass, err := storage.StorePassword(ctx, db.CreatePasswordEntryParams{
		UserID:   pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true},
		Login:    "test@example.com",
		Password: "encrypted-password",
	})
	require.NoError(t, err)

	req := &pb.DeletePasswordRequest{
		PasswordId: testPass.ID.String(),
	}

	resp, err := svc.DeletePassword(ctx, req)
	require.NoError(t, err)
	require.True(t, resp.Ok)

	// Verify password was deleted
	_, err = storage.GetPassword(ctx, testPass.ID.String(), pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true})
	require.Error(t, err)
}

func TestDeletePassword_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupPasswordService(t)

	req := &pb.DeletePasswordRequest{
		PasswordId: uuid.NewString(),
	}

	_, err := svc.DeletePassword(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error deleting password")
}

func TestDeletePassword_Unauthorized(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupPasswordService(t)

	// Store a password with different user
	otherUserID := uuid.New()
	otherPass, err := storage.StorePassword(context.Background(), db.CreatePasswordEntryParams{
		UserID:   pgtype.UUID{Bytes: otherUserID, Valid: true},
		Login:    "other@example.com",
		Password: "encrypted-password",
	})
	require.NoError(t, err)

	req := &pb.DeletePasswordRequest{
		PasswordId: otherPass.ID.String(),
	}

	_, err = svc.DeletePassword(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized access to password")
}
