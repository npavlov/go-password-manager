//nolint:wrapcheck,exhaustruct
package interceptors_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/npavlov/go-password-manager/gen/proto/auth"
	"github.com/npavlov/go-password-manager/internal/server/service/interceptors"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
)

type MockMemStorage struct {
	mock.Mock
}

func (m *MockMemStorage) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)

	return args.String(0), args.Error(1)
}

func (m *MockMemStorage) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)

	return args.Error(0)
}

func TestTokenInterceptor(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(nil)
	jwtSecret := "test-secret1"
	validUserID := "user-123"
	validToken, _ := utils.GenerateJWT("user-123", jwtSecret, time.Now().Add(time.Hour).Unix())
	invalidToken := "invalid.token"

	tests := []struct {
		name          string
		method        string
		token         string
		mockSetup     func(*MockMemStorage)
		expectedError bool
		expectedCode  codes.Code
	}{
		{
			name:   "successful authentication",
			method: "/service.privateMethod",
			token:  validToken,
			mockSetup: func(m *MockMemStorage) {
				m.On("Get", mock.Anything, validToken).Return(validUserID, nil)
			},
			expectedError: false,
		},
		{
			name:          "skip public method",
			method:        pb.AuthService_Register_FullMethodName,
			token:         "",
			mockSetup:     func(_ *MockMemStorage) {},
			expectedError: false,
		},
		{
			name:          "missing metadata",
			method:        "/service.privateMethod",
			token:         "",
			mockSetup:     func(_ *MockMemStorage) {},
			expectedError: true,
			expectedCode:  codes.Unauthenticated,
		},
		{
			name:   "invalid token",
			method: "/service.privateMethod",
			token:  invalidToken,
			mockSetup: func(m *MockMemStorage) {
				m.On("Get", mock.Anything, invalidToken).Return("", errors.New("invalid"))
			},
			expectedError: true,
			expectedCode:  codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockMemStorage{}
			tt.mockSetup(mockStorage)

			interceptor := interceptors.TokenInterceptor(&logger, jwtSecret, mockStorage)

			// Create test context with or without token
			ctx := t.Context()
			if tt.token != "" {
				md := metadata.New(map[string]string{"authorization": tt.token})
				ctx = metadata.NewIncomingContext(ctx, md)
			}

			// Mock handler that checks for user_id in context
			handler := func(ctx context.Context, _ interface{}) (interface{}, error) {
				if tt.expectedError {
					t.Error("Handler should not be called for error cases")

					//nolint:nilnil
					return nil, nil
				}

				if tt.token != "" {
					// For successful cases, verify user_id is in context
					userID := ctx.Value("user_id")
					require.Equal(t, validUserID, userID)
				}

				return "success", nil
			}

			// Create mock server info
			info := &grpc.UnaryServerInfo{
				FullMethod: tt.method,
			}

			// Call the interceptor
			resp, err := interceptor(ctx, "request", info, handler)

			if tt.expectedError {
				require.Error(t, err)
				if st, ok := status.FromError(err); ok {
					require.Equal(t, tt.expectedCode, st.Code())
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, "success", resp)
				mockStorage.AssertExpectations(t)
			}
		})
	}
}

func TestStreamTokenInterceptor(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(nil)
	jwtSecret := "test-secret"
	validToken, _ := utils.GenerateJWT("user-123", jwtSecret, time.Now().Add(time.Hour).Unix())
	validUserID := "user-123"

	tests := []struct {
		name          string
		method        string
		token         string
		mockSetup     func(*MockMemStorage)
		expectedError bool
		expectedCode  codes.Code
	}{
		{
			name:   "successful authentication",
			method: "/service.StreamMethod",
			token:  validToken,
			mockSetup: func(m *MockMemStorage) {
				m.On("Get", mock.Anything, validToken).Return(validUserID, nil)
			},
			expectedError: false,
		},
		{
			name:          "skip reflection",
			method:        "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
			token:         "",
			mockSetup:     func(_ *MockMemStorage) {},
			expectedError: false,
		},
		{
			name:          "missing metadata",
			method:        "/service.StreamMethod",
			token:         "",
			mockSetup:     func(_ *MockMemStorage) {},
			expectedError: true,
			expectedCode:  codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockMemStorage{}
			tt.mockSetup(mockStorage)

			interceptor := interceptors.StreamTokenInterceptor(&logger, jwtSecret, mockStorage)

			// Create test context with or without token
			ctx := t.Context()
			if tt.token != "" {
				md := metadata.New(map[string]string{"authorization": tt.token})
				ctx = metadata.NewIncomingContext(ctx, md)
			}

			// Mock stream
			mockStream := &mockServerStream{ctx: ctx}

			// Mock handler that checks for user_id in context
			handler := func(_ interface{}, stream grpc.ServerStream) error {
				if tt.expectedError {
					t.Error("Handler should not be called for error cases")

					return nil
				}

				if tt.token != "" {
					// For successful cases, verify user_id is in context
					userID := stream.Context().Value("user_id")
					require.Equal(t, validUserID, userID)
				}

				return nil
			}

			// Create mock server info
			info := &grpc.StreamServerInfo{
				FullMethod: tt.method,
			}

			// Call the interceptor
			err := interceptor(nil, mockStream, info, handler)

			if tt.expectedError {
				require.Error(t, err)
				if st, ok := status.FromError(err); ok {
					require.Equal(t, tt.expectedCode, st.Code())
				}
			} else {
				require.NoError(t, err)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

//nolint:containedctx
type mockServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}

func TestAuthenticateToken(t *testing.T) {
	t.Parallel()

	jwtSecret := "test-secret"
	validToken, _ := utils.GenerateJWT("user-123", jwtSecret, time.Now().Add(time.Hour).Unix())

	tests := []struct {
		name          string
		token         string
		mockSetup     func(*MockMemStorage)
		expectedError bool
		expectedCode  codes.Code
	}{
		{
			name:  "valid token",
			token: validToken,
			mockSetup: func(m *MockMemStorage) {
				m.On("Get", mock.Anything, validToken).Return("user-123", nil)
			},
			expectedError: false,
		},
		{
			name:          "missing metadata",
			token:         "",
			mockSetup:     func(_ *MockMemStorage) {},
			expectedError: true,
			expectedCode:  codes.Unauthenticated,
		},
		{
			name:          "invalid token format",
			token:         "badformat",
			mockSetup:     func(_ *MockMemStorage) {},
			expectedError: true,
			expectedCode:  codes.Unauthenticated,
		},
		{
			name:  "token not in storage",
			token: validToken,
			mockSetup: func(m *MockMemStorage) {
				m.On("Get", mock.Anything, validToken).Return("", errors.New("not found"))
			},
			expectedError: true,
			expectedCode:  codes.Unauthenticated,
		},
		{
			name:  "token user mismatch",
			token: validToken,
			mockSetup: func(m *MockMemStorage) {
				m.On("Get", mock.Anything, validToken).Return("different-user", nil)
			},
			expectedError: true,
			expectedCode:  codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockMemStorage{}

			mockStorage.ExpectedCalls = nil
			mockStorage.Calls = nil
			tt.mockSetup(mockStorage)

			// Create context with or without token
			ctx := t.Context()
			if tt.token != "" {
				md := metadata.New(map[string]string{"authorization": tt.token})
				ctx = metadata.NewIncomingContext(ctx, md)
			}

			// Call authenticateToken
			_, err := interceptors.AuthenticateToken(ctx, jwtSecret, mockStorage)

			if tt.expectedError {
				require.Error(t, err)
				if st, ok := status.FromError(err); ok {
					require.Equal(t, tt.expectedCode, st.Code())
				}
			} else {
				require.NoError(t, err)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}
