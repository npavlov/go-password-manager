package auth_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/service/auth"
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

	req := &pb.RegisterRequest{
		Username: "testuser",
		Password: "securePass123!",
		Email:    "test@example.com",
	}

	// Register
	resp, err := service.Register(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, resp.GetToken())
	require.NotEmpty(t, resp.GetRefreshToken())
	require.NotEmpty(t, resp.GetUserKey())

	// Login
	loginResp, err := service.Login(ctx, &pb.LoginRequest{
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

	req := &pb.RegisterRequest{
		Username: "duplicateuser",
		Password: "pass1234",
		Email:    "dup@example.com",
	}

	_, err := service.Register(ctx, req)
	require.NoError(t, err)

	_, err = service.Register(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "username already exists")
}

func TestLoginInvalidPassword(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)

	registerReq := &pb.RegisterRequest{
		Username: "wrongpass",
		Password: "correctPass",
		Email:    "wrong@example.com",
	}

	_, err := service.Register(ctx, registerReq)
	require.NoError(t, err)

	_, err = service.Login(ctx, &pb.LoginRequest{
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

	registerResp, err := service.Register(ctx, &pb.RegisterRequest{
		Username: "refreshuser",
		Password: "refreshpass",
		Email:    "refresh@example.com",
	})
	require.NoError(t, err)

	refreshResp, err := service.RefreshToken(ctx, &pb.RefreshTokenRequest{
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

	_, err = service.RefreshToken(ctx, &pb.RefreshTokenRequest{
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
	_, err := service.Register(ctx, &pb.RegisterRequest{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestLoginValidationError(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)

	_, err := service.Login(ctx, &pb.LoginRequest{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestRefreshTokenValidationError(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)

	_, err := service.RefreshToken(ctx, &pb.RefreshTokenRequest{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestTokenGenerationFailsToStoreToken(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	service := newTestService(t)
	mockStorage := serviceStorage(service)

	// Simulate failure
	mockStorage.RegisterError = errors.New("error storing token")

	registerResp, err := service.Register(ctx, &pb.RegisterRequest{
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
