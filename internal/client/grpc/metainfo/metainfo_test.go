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

func (m *MockMetadataServiceClient) GetMetaInfoV1(ctx context.Context, in *metadata.GetMetaInfoV1Request, _ ...grpc.CallOption) (*metadata.GetMetaInfoV1Response, error) {
	args := m.Called(ctx, in)

	// Safely handle nil to avoid type assertion panic
	arg, ok := args.Get(0).(*metadata.GetMetaInfoV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockMetadataServiceClient) AddMetaInfoV1(ctx context.Context, in *metadata.AddMetaInfoV1Request, _ ...grpc.CallOption) (*metadata.AddMetaInfoV1Response, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*metadata.AddMetaInfoV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockMetadataServiceClient) RemoveMetaInfoV1(ctx context.Context, in *metadata.RemoveMetaInfoV1Request, _ ...grpc.CallOption) (*metadata.RemoveMetaInfoV1Response, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*metadata.RemoveMetaInfoV1Response)
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

	mockClient.On("GetMetaInfoV1", mock.Anything, &metadata.GetMetaInfoV1Request{
		ItemId: "item123",
	}).Return(&metadata.GetMetaInfoV1Response{
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

	mockClient.On("GetMetaInfoV1", mock.Anything, &metadata.GetMetaInfoV1Request{
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

	mockClient.On("AddMetaInfoV1", mock.Anything, &metadata.AddMetaInfoV1Request{
		ItemId:   "item123",
		Metadata: meta,
	}).Return(&metadata.AddMetaInfoV1Response{
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

	mockClient.On("AddMetaInfoV1", mock.Anything, &metadata.AddMetaInfoV1Request{
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

	mockClient.On("AddMetaInfoV1", mock.Anything, &metadata.AddMetaInfoV1Request{
		ItemId:   "item123",
		Metadata: meta,
	}).Return(&metadata.AddMetaInfoV1Response{
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

	mockClient.On("RemoveMetaInfoV1", mock.Anything, &metadata.RemoveMetaInfoV1Request{
		ItemId: "item123",
		Key:    "key1",
	}).Return(&metadata.RemoveMetaInfoV1Response{
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

	mockClient.On("RemoveMetaInfoV1", mock.Anything, &metadata.RemoveMetaInfoV1Request{
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

	mockClient.On("RemoveMetaInfoV1", mock.Anything, &metadata.RemoveMetaInfoV1Request{
		ItemId: "item123",
		Key:    "key1",
	}).Return(&metadata.RemoveMetaInfoV1Response{
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
