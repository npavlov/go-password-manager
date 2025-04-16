//nolint:wrapcheck,lll,err113
package interceptors_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/npavlov/go-password-manager/internal/client/config"
	"github.com/npavlov/go-password-manager/internal/client/interceptors"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

// Mocks

type MockAuthServiceClient struct {
	mock.Mock
}

func (m *MockAuthServiceClient) RefreshToken(ctx context.Context, in *auth.RefreshTokenRequest, opts ...grpc.CallOption) (*auth.RefreshTokenResponse, error) {
	args := m.Called(ctx, in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*auth.RefreshTokenResponse), args.Error(1)
}

type MockUnaryInvoker struct {
	mock.Mock
}

type MockStreamer struct {
	mock.Mock
}

func (m *MockUnaryInvoker) Invoke(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
	args := m.Called(ctx, method, req, reply, cc, opts)

	return args.Error(0)
}

func (m *MockStreamer) Stream(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	args := m.Called(ctx, desc, cc, method, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(grpc.ClientStream), args.Error(1)
}

func TestNewAuthInterceptor(t *testing.T) {
	t.Parallel()

	cfg := config.Config{}
	tm := new(testutils.MockTokenManager)

	interceptor := interceptors.NewAuthInterceptor(cfg, tm)

	assert.NotNil(t, interceptor)
}

func TestSetAuthClient(t *testing.T) {
	t.Parallel()

	cfg := config.Config{}
	tm := new(testutils.MockTokenManager)
	interceptor := interceptors.NewAuthInterceptor(cfg, tm)
	mockAuthClient := new(MockAuthServiceClient)

	interceptor.SetAuthClient(mockAuthClient)

	assert.NotNil(t, interceptor.AuthClient)
}

func TestUnaryInterceptor_SkipAuthMethods(t *testing.T) {
	tests := []struct {
		method string
	}{
		{method: auth.AuthService_Register_FullMethodName},
		{method: auth.AuthService_Login_FullMethodName},
		{method: auth.AuthService_RefreshToken_FullMethodName},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			cfg := config.Config{}
			tm := new(testutils.MockTokenManager)
			interceptor := interceptors.NewAuthInterceptor(cfg, tm)
			invoker := new(MockUnaryInvoker)
			invoker.On("Invoke", mock.Anything, tt.method, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(nil)

			err := interceptor.UnaryInterceptor(t.Context(), tt.method, nil, nil, nil, invoker.Invoke)

			require.NoError(t, err)
			invoker.AssertExpectations(t)
		})
	}
}

func TestUnaryInterceptor_Success(t *testing.T) {
	cfg := config.Config{}
	tm := new(testutils.MockTokenManager)
	tm.On("GetAccessToken").Return("valid_token")
	interceptor := interceptors.NewAuthInterceptor(cfg, tm)
	invoker := new(MockUnaryInvoker)
	invoker.On("Invoke", mock.Anything, "some.method", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	err := interceptor.UnaryInterceptor(t.Context(), "some.method", nil, nil, nil, invoker.Invoke)

	require.NoError(t, err)
	tm.AssertExpectations(t)
	invoker.AssertExpectations(t)
}

func TestUnaryInterceptor_UnauthenticatedError(t *testing.T) {
	t.Parallel()

	cfg := config.Config{}
	tm := new(testutils.MockTokenManager)
	tm.On("GetRefreshToken").Return("valid_refresh").Once()
	tm.On("UpdateTokens", "new_token", "new_refresh").Return(nil)
	tm.On("GetAccessToken").Return("new_token").Once()

	interceptor := interceptors.NewAuthInterceptor(cfg, tm)
	mockAuthClient := new(MockAuthServiceClient)
	mockAuthClient.On("RefreshToken", mock.Anything, &auth.RefreshTokenRequest{
		RefreshToken: "valid_refresh",
	}).Return(&auth.RefreshTokenResponse{
		Token:        "new_token",
		RefreshToken: "new_refresh",
	}, nil)
	interceptor.SetAuthClient(mockAuthClient)

	invoker := new(MockUnaryInvoker)
	invoker.On("Invoke", mock.Anything, "some.method", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(status.Error(codes.Unauthenticated, "invalid token")).Once()
	invoker.On("Invoke", mock.Anything, "some.method", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	err := interceptor.UnaryInterceptor(t.Context(), "some.method", nil, nil, nil, invoker.Invoke)

	require.NoError(t, err)
	tm.AssertExpectations(t)
	invoker.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
}

func TestUnaryInterceptor_RefreshTokenFailure(t *testing.T) {
	t.Parallel()

	cfg := config.Config{}
	tm := new(testutils.MockTokenManager)
	tm.On("GetAccessToken").Return("expired_token")
	tm.On("GetRefreshToken").Return("invalid_refresh")
	tm.On("HandleAuthFailure").Once()

	interceptor := interceptors.NewAuthInterceptor(cfg, tm)
	mockAuthClient := new(MockAuthServiceClient)
	mockAuthClient.On("RefreshToken", mock.Anything, &auth.RefreshTokenRequest{
		RefreshToken: "invalid_refresh",
	}).Return(nil, errors.New("refresh failed"))
	interceptor.SetAuthClient(mockAuthClient)

	invoker := new(MockUnaryInvoker)
	invoker.On("Invoke", mock.Anything, "some.method", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(status.Error(codes.Unauthenticated, "invalid token"))

	err := interceptor.UnaryInterceptor(t.Context(), "some.method", nil, nil, nil, invoker.Invoke)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to refresh access token")
	tm.AssertExpectations(t)
	invoker.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
}

func TestUnaryInterceptor_OtherError(t *testing.T) {
	t.Parallel()

	cfg := config.Config{}
	tm := new(testutils.MockTokenManager)
	tm.On("GetAccessToken").Return("valid_token")

	interceptor := interceptors.NewAuthInterceptor(cfg, tm)
	invoker := new(MockUnaryInvoker)
	invoker.On("Invoke", mock.Anything, "some.method", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(status.Error(codes.Internal, "server error"))

	err := interceptor.UnaryInterceptor(t.Context(), "some.method", nil, nil, nil, invoker.Invoke)

	require.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
	tm.AssertExpectations(t)
	invoker.AssertExpectations(t)
}

func TestStreamInterceptor_NoTokenManager(t *testing.T) {
	t.Parallel()

	interceptor := interceptors.NewAuthInterceptor(config.Config{}, nil)
	streamer := new(MockStreamer)
	streamer.On("Stream", mock.Anything, mock.Anything, mock.Anything, "some.method", mock.Anything).
		Return(nil, nil)

	_, err := interceptor.StreamInterceptor(t.Context(), nil, nil, "some.method", streamer.Stream)

	require.NoError(t, err)
	streamer.AssertExpectations(t)
}

func TestStreamInterceptor_NoToken(t *testing.T) {
	t.Parallel()

	tm := new(testutils.MockTokenManager)
	tm.On("GetAccessToken").Return("")

	interceptor := interceptors.NewAuthInterceptor(config.Config{}, tm)
	_, err := interceptor.StreamInterceptor(t.Context(), nil, nil, "some.method", nil)

	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
	tm.AssertExpectations(t)
}

func TestStreamInterceptor_Success(t *testing.T) {
	t.Parallel()

	tm := new(testutils.MockTokenManager)
	tm.On("GetAccessToken").Return("valid_token")

	interceptor := interceptors.NewAuthInterceptor(config.Config{}, tm)
	streamer := new(MockStreamer)
	streamer.On("Stream", mock.Anything, mock.Anything, mock.Anything, "some.method", mock.Anything).
		Return(nil, nil)

	ctx := t.Context()
	_, err := interceptor.StreamInterceptor(ctx, nil, nil, "some.method", streamer.Stream)

	require.NoError(t, err)
	streamer.AssertExpectations(t)
	tm.AssertExpectations(t)
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	cfg := config.Config{}
	tm := new(testutils.MockTokenManager)
	tm.On("GetRefreshToken").Return("refresh").Once()
	tm.On("UpdateTokens", "new_token", "new_refresh").Return(nil)
	tm.On("GetAccessToken").Return("new_token").Times(2)

	interceptor := interceptors.NewAuthInterceptor(cfg, tm)
	mockAuthClient := new(MockAuthServiceClient)
	mockAuthClient.On("RefreshToken", mock.Anything, &auth.RefreshTokenRequest{
		RefreshToken: "refresh",
	}).Return(&auth.RefreshTokenResponse{
		Token:        "new_token",
		RefreshToken: "new_refresh",
	}, nil)
	interceptor.SetAuthClient(mockAuthClient)

	invoker := new(MockUnaryInvoker)
	invoker.On("Invoke", mock.Anything, "method1", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(status.Error(codes.Unauthenticated, "invalid token")).Once()
	invoker.On("Invoke", mock.Anything, "method1", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	invoker.On("Invoke", mock.Anything, "method2", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := interceptor.UnaryInterceptor(t.Context(), "method1", nil, nil, nil, invoker.Invoke)
		assert.NoError(t, err)
	}()

	go func() {
		defer wg.Done()
		err := interceptor.UnaryInterceptor(t.Context(), "method2", nil, nil, nil, invoker.Invoke)
		assert.NoError(t, err)
	}()

	wg.Wait()
	tm.AssertExpectations(t)
	invoker.AssertExpectations(t)
	mockAuthClient.AssertExpectations(t)
}
