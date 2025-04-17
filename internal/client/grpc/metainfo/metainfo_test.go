//nolint:wrapcheck,err113,lll
package metainfo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/npavlov/go-password-manager/gen/proto/metadata"
	"github.com/npavlov/go-password-manager/internal/client/grpc/metainfo"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

// Mocks

type MockMetadataServiceClient struct {
	mock.Mock
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockMetadataServiceClient) GetMetaInfo(ctx context.Context, in *metadata.GetMetaInfoRequest, _ ...grpc.CallOption) (*metadata.GetMetaInfoResponse, error) {
	args := m.Called(ctx, in)

	// Safely handle nil to avoid type assertion panic
	arg, ok := args.Get(0).(*metadata.GetMetaInfoResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockMetadataServiceClient) AddMetaInfo(ctx context.Context, in *metadata.AddMetaInfoRequest, _ ...grpc.CallOption) (*metadata.AddMetaInfoResponse, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*metadata.AddMetaInfoResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockMetadataServiceClient) RemoveMetaInfo(ctx context.Context, in *metadata.RemoveMetaInfoRequest, _ ...grpc.CallOption) (*metadata.RemoveMetaInfoResponse, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*metadata.RemoveMetaInfoResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockTokenManager) GetToken() (string, error) {
	args := m.Called()

	return args.String(0), args.Error(1)
}

func TestGetMetainfo_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockMetadataServiceClient)
	logger := zerolog.Nop()

	expectedMeta := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	mockClient.On("GetMetaInfo", mock.Anything, &metadata.GetMetaInfoRequest{
		ItemId: "item123",
	}).Return(&metadata.GetMetaInfoResponse{
		Metadata: expectedMeta,
	}, nil)

	client := &metainfo.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	meta, err := client.GetMetainfo(t.Context(), "item123")
	require.NoError(t, err)
	assert.Equal(t, expectedMeta, meta)
}

func TestGetMetainfo_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockMetadataServiceClient)
	logger := zerolog.Nop()

	mockClient.On("GetMetaInfo", mock.Anything, &metadata.GetMetaInfoRequest{
		ItemId: "item123",
	}).Return(nil, errors.New("get meta failed"))

	client := &metainfo.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.GetMetainfo(t.Context(), "item123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error getting metainfo")
}

func TestSetMetainfo_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockMetadataServiceClient)
	logger := zerolog.Nop()

	meta := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	mockClient.On("AddMetaInfo", mock.Anything, &metadata.AddMetaInfoRequest{
		ItemId:   "item123",
		Metadata: meta,
	}).Return(&metadata.AddMetaInfoResponse{
		Success: true,
	}, nil)

	client := &metainfo.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	success, err := client.SetMetainfo(t.Context(), "item123", meta)
	require.NoError(t, err)
	assert.True(t, success)
}

func TestSetMetainfo_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockMetadataServiceClient)
	logger := zerolog.Nop()

	meta := map[string]string{
		"key1": "value1",
	}

	mockClient.On("AddMetaInfo", mock.Anything, &metadata.AddMetaInfoRequest{
		ItemId:   "item123",
		Metadata: meta,
	}).Return(nil, errors.New("set meta failed"))

	client := &metainfo.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.SetMetainfo(t.Context(), "item123", meta)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error setting metainfo")
}

func TestSetMetainfo_NotSuccessful(t *testing.T) {
	t.Parallel()

	mockClient := new(MockMetadataServiceClient)
	logger := zerolog.Nop()

	meta := map[string]string{
		"key1": "value1",
	}

	mockClient.On("AddMetaInfo", mock.Anything, &metadata.AddMetaInfoRequest{
		ItemId:   "item123",
		Metadata: meta,
	}).Return(&metadata.AddMetaInfoResponse{
		Success: false,
	}, nil)

	client := &metainfo.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	success, err := client.SetMetainfo(t.Context(), "item123", meta)
	require.NoError(t, err)
	assert.False(t, success)
}

func TestDeleteMetainfo_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockMetadataServiceClient)
	logger := zerolog.Nop()

	mockClient.On("RemoveMetaInfo", mock.Anything, &metadata.RemoveMetaInfoRequest{
		ItemId: "item123",
		Key:    "key1",
	}).Return(&metadata.RemoveMetaInfoResponse{
		Success: true,
	}, nil)

	client := &metainfo.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	success, err := client.DeleteMetainfo(t.Context(), "item123", "key1")
	require.NoError(t, err)
	assert.True(t, success)
}

func TestDeleteMetainfo_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockMetadataServiceClient)
	logger := zerolog.Nop()

	mockClient.On("RemoveMetaInfo", mock.Anything, &metadata.RemoveMetaInfoRequest{
		ItemId: "item123",
		Key:    "key1",
	}).Return(nil, errors.New("delete meta failed"))

	client := &metainfo.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.DeleteMetainfo(t.Context(), "item123", "key1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting metainfo")
}

func TestDeleteMetainfo_NotSuccessful(t *testing.T) {
	t.Parallel()

	mockClient := new(MockMetadataServiceClient)
	logger := zerolog.Nop()

	mockClient.On("RemoveMetaInfo", mock.Anything, &metadata.RemoveMetaInfoRequest{
		ItemId: "item123",
		Key:    "key1",
	}).Return(&metadata.RemoveMetaInfoResponse{
		Success: false,
	}, nil)

	client := &metainfo.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	success, err := client.DeleteMetainfo(t.Context(), "item123", "key1")
	require.NoError(t, err)
	assert.False(t, success)
}

func TestNewMetaClient(t *testing.T) {
	t.Parallel()

	tm := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	client := metainfo.NewMetainfoClient(nil, tm, &logger)

	assert.NotNil(t, client)
}
