//nolint:lll,exhaustruct
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
	ctx := testutils.InjectUserToContext(t.Context(), userID.String())

	return password.NewPasswordService(&logger, mockStorage, cfg), mockStorage, ctx, encryptionKey
}

func TestStorePassword_Success(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupPasswordService(t)

	testLogin := "test@example.com"
	testPassword := "securepassword123"

	req := &pb.StorePasswordV1Request{
		Password: &pb.PasswordData{
			Login:    testLogin,
			Password: testPassword,
		},
	}

	resp, err := svc.StorePasswordV1(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, resp.GetPasswordId())

	// Verify password was stored
	storedPass, err := storage.GetPassword(ctx, resp.GetPasswordId(), pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true})
	require.NoError(t, err)
	require.Equal(t, testLogin, storedPass.Login)
}

func TestStorePassword_ValidationError(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupPasswordService(t)

	// Missing password data
	req := &pb.StorePasswordV1Request{
		Password: &pb.PasswordData{
			Login: "", // Invalid empty login
		},
	}

	_, err := svc.StorePasswordV1(ctx, req)
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

	req := &pb.StorePasswordV1Request{
		Password: &pb.PasswordData{
			Login:    "test@example.com",
			Password: "password123",
		},
	}

	_, err := svc.StorePasswordV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error getting decryption key")
}

func TestStorePassword_StorageError(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupPasswordService(t)

	// Make storage return error
	storage.CallError = errors.New("database failure")

	req := &pb.StorePasswordV1Request{
		Password: &pb.PasswordData{
			Login:    "test@example.com",
			Password: "password123",
		},
	}

	_, err := svc.StorePasswordV1(ctx, req)
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

	req := &pb.GetPasswordV1Request{
		PasswordId: storedPass.ID.String(),
	}

	resp, err := svc.GetPasswordV1(ctx, req)
	require.NoError(t, err)
	require.Equal(t, testLogin, resp.GetPassword().GetLogin())
	require.Equal(t, testPassword, resp.GetPassword().GetPassword())
	require.True(t, timestamppb.New(storedPass.UpdatedAt.Time).AsTime().Equal(resp.GetLastUpdate().AsTime()))
}

func TestGetPassword_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupPasswordService(t)

	req := &pb.GetPasswordV1Request{
		PasswordId: uuid.NewString(),
	}

	_, err := svc.GetPasswordV1(ctx, req)
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

	req := &pb.GetPasswordV1Request{
		PasswordId: storedPass.ID.String(),
	}

	_, err = svc.GetPasswordV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error decrypting password")
}

func TestGetPassword_Unauthorized(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupPasswordService(t)

	// Store a password with different user
	otherUserID := uuid.New()
	otherPass, err := storage.StorePassword(t.Context(), db.CreatePasswordEntryParams{
		UserID:   pgtype.UUID{Bytes: otherUserID, Valid: true},
		Login:    "other@example.com",
		Password: "encrypted-password",
	})
	require.NoError(t, err)

	req := &pb.GetPasswordV1Request{
		PasswordId: otherPass.ID.String(),
	}

	_, err = svc.GetPasswordV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized access to password")
}

func TestGetPasswords_NotImplemented(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupPasswordService(t)

	req := &pb.GetPasswordsV1Request{}
	resp, err := svc.GetPasswordsV1(ctx, req)

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

	req := &pb.UpdatePasswordV1Request{
		PasswordId: initialPass.ID.String(),
		Data: &pb.PasswordData{
			Login:    newLogin,
			Password: newPassword,
		},
	}

	resp, err := svc.UpdatePasswordV1(ctx, req)
	require.NoError(t, err)
	require.Equal(t, initialPass.ID.String(), resp.GetPasswordId())

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

	req := &pb.UpdatePasswordV1Request{
		PasswordId: uuid.NewString(),
		Data:       &pb.PasswordData{}, // Missing required fields
	}

	_, err := svc.UpdatePasswordV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestUpdatePassword_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupPasswordService(t)

	req := &pb.UpdatePasswordV1Request{
		PasswordId: uuid.NewString(),
		Data: &pb.PasswordData{
			Login:    "test@example.com",
			Password: "password123",
		},
	}

	_, err := svc.UpdatePasswordV1(ctx, req)
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

	req := &pb.DeletePasswordV1Request{
		PasswordId: testPass.ID.String(),
	}

	resp, err := svc.DeletePasswordV1(ctx, req)
	require.NoError(t, err)
	require.True(t, resp.GetOk())

	// Verify password was deleted
	_, err = storage.GetPassword(ctx, testPass.ID.String(), pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true})
	require.Error(t, err)
}

func TestDeletePassword_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupPasswordService(t)

	req := &pb.DeletePasswordV1Request{
		PasswordId: uuid.NewString(),
	}

	_, err := svc.DeletePasswordV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error deleting password")
}

func TestDeletePassword_Unauthorized(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupPasswordService(t)

	// Store a password with different user
	otherUserID := uuid.New()
	otherPass, err := storage.StorePassword(t.Context(), db.CreatePasswordEntryParams{
		UserID:   pgtype.UUID{Bytes: otherUserID, Valid: true},
		Login:    "other@example.com",
		Password: "encrypted-password",
	})
	require.NoError(t, err)

	req := &pb.DeletePasswordV1Request{
		PasswordId: otherPass.ID.String(),
	}

	_, err = svc.DeletePasswordV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized access to password")
}

func TestStorePassword_Validation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		req         *pb.StorePasswordV1Request
		expectedErr string
	}{
		{
			name: "nil password data",
			req: &pb.StorePasswordV1Request{
				Password: nil,
			},
			expectedErr: "error validating input",
		},
		{
			name: "empty login",
			req: &pb.StorePasswordV1Request{
				Password: &pb.PasswordData{
					Login:    "",
					Password: "validpassword",
				},
			},
			expectedErr: "error validating input",
		},
		{
			name: "empty password",
			req: &pb.StorePasswordV1Request{
				Password: &pb.PasswordData{
					Login:    "valid@example.com",
					Password: "",
				},
			},
			expectedErr: "error validating input",
		},
	}

	svc, _, ctx, _ := setupPasswordService(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.StorePasswordV1(ctx, tc.req)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestGetPassword_Validation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		req         *pb.GetPasswordV1Request
		expectedErr string
	}{
		{
			name: "empty password id",
			req: &pb.GetPasswordV1Request{
				PasswordId: "",
			},
			expectedErr: "error validating input",
		},
		{
			name: "invalid password id format",
			req: &pb.GetPasswordV1Request{
				PasswordId: "not-a-uuid",
			},
			expectedErr: "error validating input",
		},
	}

	svc, _, ctx, _ := setupPasswordService(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.GetPasswordV1(ctx, tc.req)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestUpdatePassword_Validation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		req         *pb.UpdatePasswordV1Request
		expectedErr string
	}{
		{
			name: "empty password id",
			req: &pb.UpdatePasswordV1Request{
				PasswordId: "",
				Data: &pb.PasswordData{
					Login:    "valid@example.com",
					Password: "validpassword",
				},
			},
			expectedErr: "error validating input",
		},
		{
			name: "invalid password id format",
			req: &pb.UpdatePasswordV1Request{
				PasswordId: "not-a-uuid",
				Data: &pb.PasswordData{
					Login:    "valid@example.com",
					Password: "validpassword",
				},
			},
			expectedErr: "error validating input",
		},
		{
			name: "nil password data",
			req: &pb.UpdatePasswordV1Request{
				PasswordId: uuid.NewString(),
				Data:       nil,
			},
			expectedErr: "error validating input",
		},
		{
			name: "empty login in data",
			req: &pb.UpdatePasswordV1Request{
				PasswordId: uuid.NewString(),
				Data: &pb.PasswordData{
					Login:    "",
					Password: "validpassword",
				},
			},
			expectedErr: "error validating input",
		},
		{
			name: "empty password in data",
			req: &pb.UpdatePasswordV1Request{
				PasswordId: uuid.NewString(),
				Data: &pb.PasswordData{
					Login:    "valid@example.com",
					Password: "",
				},
			},
			expectedErr: "error validating input",
		},
	}

	svc, _, ctx, _ := setupPasswordService(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := svc.UpdatePasswordV1(ctx, tc.req)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestDeletePassword_Validation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		req         *pb.DeletePasswordV1Request
		expectedErr string
	}{
		{
			name: "empty password id",
			req: &pb.DeletePasswordV1Request{
				PasswordId: "",
			},
			expectedErr: "error validating input",
		},
		{
			name: "invalid password id format",
			req: &pb.DeletePasswordV1Request{
				PasswordId: "not-a-uuid",
			},
			expectedErr: "error validating input",
		},
	}

	svc, _, ctx, _ := setupPasswordService(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := svc.DeletePasswordV1(ctx, tc.req)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestGetPasswords_Validation(t *testing.T) {
	t.Parallel()

	// Currently empty since GetPasswordsV1Request has no fields to validate
	// This test is included for completeness and future-proofing
	svc, _, ctx, _ := setupPasswordService(t)

	req := &pb.GetPasswordsV1Request{}
	_, err := svc.GetPasswordsV1(ctx, req)
	require.NoError(t, err)
}
