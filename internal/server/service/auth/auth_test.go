//nolint:exhaustruct,forcetypeassert
package auth_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/service/auth"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	generalutils "github.com/npavlov/go-password-manager/internal/utils"
)

func newTestService(t *testing.T) *auth.Service {
	t.Helper()

	masterKey, _ := utils.GenerateRandomKey()

	logger := testutils.GetTLogger()
	mockStorage := testutils.NewMockDBStorage(logger, masterKey)
	mockRedis := testutils.NewMockRedis()

	cfg := &config.Config{
		JwtSecret:        "test-jwt-secret",
		SecuredMasterKey: generalutils.NewString(masterKey),
	}

	return auth.NewAuthService(logger, mockStorage, cfg, mockRedis)
}

func TestRegisterLoginFlow(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)

	req := &pb.RegisterV1Request{
		Username: "testuser",
		Password: "securePass123!",
		Email:    "test@example.com",
	}

	// Register
	resp, err := service.RegisterV1(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, resp.GetToken())
	require.NotEmpty(t, resp.GetRefreshToken())
	require.NotEmpty(t, resp.GetUserKey())

	// Login
	loginResp, err := service.LoginV1(ctx, &pb.LoginV1Request{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.GetToken())
	require.NotEmpty(t, loginResp.GetRefreshToken())
}

func TestRegisterDuplicateUsername(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)

	req := &pb.RegisterV1Request{
		Username: "duplicateuser",
		Password: "pass1234",
		Email:    "dup@example.com",
	}

	_, err := service.RegisterV1(ctx, req)
	require.NoError(t, err)

	_, err = service.RegisterV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "username already exists")
}

func TestLoginInvalidPassword(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)

	registerReq := &pb.RegisterV1Request{
		Username: "wrongpass",
		Password: "correctPass",
		Email:    "wrong@example.com",
	}

	_, err := service.RegisterV1(ctx, registerReq)
	require.NoError(t, err)

	_, err = service.LoginV1(ctx, &pb.LoginV1Request{
		Username: registerReq.GetUsername(),
		Password: "wrongPass",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid password")
}

func TestRefreshTokenFlow(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)

	registerResp, err := service.RegisterV1(ctx, &pb.RegisterV1Request{
		Username: "refreshuser",
		Password: "refreshpass",
		Email:    "refresh@example.com",
	})
	require.NoError(t, err)

	refreshResp, err := service.RefreshTokenV1(ctx, &pb.RefreshTokenV1Request{
		RefreshToken: registerResp.GetRefreshToken(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, refreshResp.GetToken())
	require.NotEmpty(t, refreshResp.GetRefreshToken())
}

func TestExpiredRefreshToken(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)

	mockStorage := serviceStorage(service)
	userID := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	expiredTime := time.Now().Add(-1 * time.Hour)

	token := "expiredToken"
	err := mockStorage.StoreToken(ctx, userID, token, expiredTime)
	require.NoError(t, err)

	_, err = service.RefreshTokenV1(ctx, &pb.RefreshTokenV1Request{
		RefreshToken: token,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "refresh token expired")
}

func TestRegisterValidationError(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)

	// Missing required fields
	_, err := service.RegisterV1(ctx, &pb.RegisterV1Request{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestLoginValidationError(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)

	_, err := service.LoginV1(ctx, &pb.LoginV1Request{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestRefreshTokenValidationError(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)

	_, err := service.RefreshTokenV1(ctx, &pb.RefreshTokenV1Request{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestTokenGenerationFailsToStoreToken(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)
	mockStorage := serviceStorage(service)

	// Simulate failure
	mockStorage.CallError = errors.New("error storing token")

	registerResp, err := service.RegisterV1(ctx, &pb.RegisterV1Request{
		Username: "failstore",
		Password: "storepass",
		Email:    "fail@example.com",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "error storing token")
	require.Nil(t, registerResp)
}

// helper to extract storage from service (type assertion).
func serviceStorage(s *auth.Service) *testutils.MockDBStorage {
	return s.Storage.(*testutils.MockDBStorage)
}

func TestRefreshTokenValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		req         *pb.RefreshTokenV1Request
		expectedErr string
	}{
		{
			name: "empty refresh token",
			req: &pb.RefreshTokenV1Request{
				RefreshToken: "",
			},
		},
		{
			name: "invalid refresh token format",
			req: &pb.RefreshTokenV1Request{
				RefreshToken: "invalid.token.format",
			},
		},
	}

	ctx := t.Context()
	service := newTestService(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := service.RefreshTokenV1(ctx, tc.req)
			require.Error(t, err)
		})
	}
}

func TestRegisterValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		req         *pb.RegisterV1Request
		expectedErr string
	}{
		{
			name: "valid request",
			req: &pb.RegisterV1Request{
				Username: "validuser",
				Password: "ValidPass123!",
				Email:    "valid@example.com",
			},
			expectedErr: "",
		},
		{
			name: "empty username",
			req: &pb.RegisterV1Request{
				Username: "",
				Password: "ValidPass123!",
				Email:    "valid@example.com",
			},
			expectedErr: "value length must be at least 3 characters",
		},
		{
			name: "short username",
			req: &pb.RegisterV1Request{
				Username: "ab",
				Password: "ValidPass123!",
				Email:    "valid@example.com",
			},
			expectedErr: "value length must be at least 3 characters",
		},
		{
			name: "empty password",
			req: &pb.RegisterV1Request{
				Username: "validuser",
				Password: "",
				Email:    "valid@example.com",
			},
			expectedErr: "value length must be at least 8 characters",
		},
		{
			name: "weak password",
			req: &pb.RegisterV1Request{
				Username: "validuser",
				Password: "weak",
				Email:    "valid@example.com",
			},
			expectedErr: "value length must be at least 8 characters",
		},
		{
			name: "invalid email",
			req: &pb.RegisterV1Request{
				Username: "validuser",
				Password: "ValidPass123!",
				Email:    "invalid-email",
			},
			expectedErr: "value must be a valid email address",
		},
	}

	ctx := t.Context()
	service := newTestService(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := service.RegisterV1(ctx, tc.req)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			}
		})
	}
}

func TestLoginValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		req         *pb.LoginV1Request
		expectedErr bool
	}{
		{
			name: "valid request",
			req: &pb.LoginV1Request{
				Username: "validuser",
				Password: "ValidPass123!",
			},
			expectedErr: false,
		},
		{
			name: "empty username",
			req: &pb.LoginV1Request{
				Username: "",
				Password: "ValidPass123!",
			},
			expectedErr: true,
		},
		{
			name: "empty password",
			req: &pb.LoginV1Request{
				Username: "validuser",
				Password: "",
			},
			expectedErr: true,
		},
	}

	ctx := t.Context()
	service := newTestService(t)

	// First register a valid user for login tests
	_, err := service.RegisterV1(ctx, &pb.RegisterV1Request{
		Username: "validuser",
		Password: "ValidPass123!",
		Email:    "valid@example.com",
	})
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := service.LoginV1(ctx, tc.req)
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
