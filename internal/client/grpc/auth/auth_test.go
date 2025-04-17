//nolint:goconst,forcetypeassert,wrapcheck
package auth_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/npavlov/go-password-manager/internal/client/grpc/auth"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

type MockAuthServiceClient struct {
	mock.Mock
}

func (m *MockAuthServiceClient) RefreshToken(ctx context.Context,
	in *pb.RefreshTokenRequest,
	_ ...grpc.CallOption,
) (*pb.RefreshTokenResponse, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*pb.RefreshTokenResponse), args.Error(1)
}

func (m *MockAuthServiceClient) Register(ctx context.Context,
	in *pb.RegisterRequest,
	_ ...grpc.CallOption,
) (*pb.RegisterResponse, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*pb.RegisterResponse), args.Error(1)
}

func (m *MockAuthServiceClient) Login(ctx context.Context,
	in *pb.LoginRequest,
	_ ...grpc.CallOption,
) (*pb.LoginResponse, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*pb.LoginResponse), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockAuthServiceClient)
	mockTokenManager := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	authClient := &auth.Client{
		Client:       mockClient,
		TokenManager: mockTokenManager,
		Log:          &logger,
	}

	username, password, email := "testuser", "testpass", "test@example.com"
	userKey := "user-123"
	accessToken := "access-token"
	refreshToken := "refresh-token"

	mockClient.On("Register", mock.Anything, mock.AnythingOfType("*auth.RegisterRequest")).
		Return(&pb.RegisterResponse{
			UserKey:      userKey,
			Token:        accessToken,
			RefreshToken: refreshToken,
		}, nil)

	mockTokenManager.On("UpdateTokens", accessToken, refreshToken).Return(nil)

	result, err := authClient.Register(username, password, email)

	require.NoError(t, err)
	assert.Equal(t, userKey, result)

	mockClient.AssertExpectations(t)
	mockTokenManager.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockAuthServiceClient)
	mockTokenManager := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	authClient := &auth.Client{
		Client:       mockClient,
		TokenManager: mockTokenManager,
		Log:          &logger,
	}

	username, password := "testuser", "testpass"
	accessToken := "access-token"
	refreshToken := "refresh-token"

	mockClient.On("Login", mock.Anything, mock.AnythingOfType("*auth.LoginRequest")).
		Return(&pb.LoginResponse{
			Token:        accessToken,
			RefreshToken: refreshToken,
		}, nil)

	mockTokenManager.On("UpdateTokens", accessToken, refreshToken).Return(nil)

	err := authClient.Login(username, password)

	require.NoError(t, err)
	mockClient.AssertExpectations(t)
	mockTokenManager.AssertExpectations(t)
}

func TestRegister_TokenUpdateFails(t *testing.T) {
	t.Parallel()

	mockClient := new(MockAuthServiceClient)
	mockTokenManager := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	authClient := &auth.Client{
		Client:       mockClient,
		TokenManager: mockTokenManager,
		Log:          &logger,
	}

	username, password, email := "testuser", "testpass", "test@example.com"
	accessToken := "access-token"
	refreshToken := "refresh-token"

	mockClient.On("Register", mock.Anything, mock.AnythingOfType("*auth.RegisterRequest")).
		Return(&pb.RegisterResponse{
			UserKey:      "user-123",
			Token:        accessToken,
			RefreshToken: refreshToken,
		}, nil)

	mockTokenManager.On("UpdateTokens", accessToken, refreshToken).Return(assert.AnError)

	result, err := authClient.Register(username, password, email)

	require.Error(t, err)
	assert.Empty(t, result)
}

func TestLogin_TokenUpdateFails(t *testing.T) {
	t.Parallel()

	mockClient := new(MockAuthServiceClient)
	mockTokenManager := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	authClient := &auth.Client{
		Client:       mockClient,
		TokenManager: mockTokenManager,
		Log:          &logger,
	}

	username, password := "testuser", "testpass"
	accessToken := "access-token"
	refreshToken := "refresh-token"

	mockClient.On("Login", mock.Anything, mock.AnythingOfType("*auth.LoginRequest")).
		Return(&pb.LoginResponse{
			Token:        accessToken,
			RefreshToken: refreshToken,
		}, nil)

	mockTokenManager.On("UpdateTokens", accessToken, refreshToken).Return(assert.AnError)

	err := authClient.Login(username, password)

	require.Error(t, err)
}

func TestNewBinaryClient(t *testing.T) {
	t.Parallel()

	tm := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	client := auth.NewAuthClient(nil, tm, &logger)

	assert.NotNil(t, client)
}
