package interceptors_test

import (
	"context"
	"testing"

	"github.com/npavlov/go-password-manager/internal/server/service/interceptors"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock handler.
func mockHandler(_ context.Context, _ interface{}) (interface{}, error) {
	return "mockResponse", nil
}

// Test LoggingServerInterceptor.
func TestLoggingServerInterceptor(t *testing.T) {
	t.Parallel()

	logger := testutils.GetTLogger()

	interceptor := interceptors.LoggingServerInterceptor(logger)

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Method"}

	resp, err := interceptor(context.Background(), "testRequest", info, mockHandler)

	require.NoError(t, err)
	assert.Equal(t, "mockResponse", resp)
}

// Test LoggingServerInterceptor with Error.
func TestLoggingServerInterceptorWithError(t *testing.T) {
	t.Parallel()

	logger := testutils.GetTLogger()

	interceptor := interceptors.LoggingServerInterceptor(logger)

	info := &grpc.UnaryServerInfo{FullMethod: "/test.Method"}

	mockErrorHandler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, status.Error(codes.Internal, "internal error")
	}

	resp, err := interceptor(context.Background(), "testRequest", info, mockErrorHandler)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Internal, status.Code(err))
}
