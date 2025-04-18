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

func (m *MockAuthServiceClient) RefreshTokenV1(ctx context.Context,
	in *pb.RefreshTokenV1Request,
	_ ...grpc.CallOption,
) (*pb.RefreshTokenV1Response, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*pb.RefreshTokenV1Response), args.Error(1)
}

func (m *MockAuthServiceClient) RegisterV1(ctx context.Context,
	in *pb.RegisterV1Request,
	_ ...grpc.CallOption,
) (*pb.RegisterV1Response, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*pb.RegisterV1Response), args.Error(1)
}

func (m *MockAuthServiceClient) LoginV1(ctx context.Context,
	in *pb.LoginV1Request,
	_ ...grpc.CallOption,
) (*pb.LoginV1Response, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*pb.LoginV1Response), args.Error(1)
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

	mockClient.On("RegisterV1", mock.Anything, mock.AnythingOfType("*auth.RegisterV1Request")).
		Return(&pb.RegisterV1Response{
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

	mockClient.On("LoginV1", mock.Anything, mock.AnythingOfType("*auth.LoginV1Request")).
		Return(&pb.LoginV1Response{
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

	mockClient.On("RegisterV1", mock.Anything, mock.AnythingOfType("*auth.RegisterV1Request")).
		Return(&pb.RegisterV1Response{
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

	mockClient.On("LoginV1", mock.Anything, mock.AnythingOfType("*auth.LoginV1Request")).
		Return(&pb.LoginV1Response{
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
