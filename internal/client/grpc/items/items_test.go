//nolint:wrapcheck,err113
package items_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/internal/client/grpc/items"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

// Mocks

type MockItemServiceClient struct {
	mock.Mock
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockItemServiceClient) GetItemsV1(ctx context.Context,
	in *item.GetItemsV1Request,
	_ ...grpc.CallOption,
) (*item.GetItemsV1Response, error) {
	args := m.Called(ctx, in)

	// Safely handle nil to avoid type assertion panic
	arg, ok := args.Get(0).(*item.GetItemsV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockTokenManager) GetToken() (string, error) {
	args := m.Called()

	return args.String(0), args.Error(1)
}

func TestGetItems_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockItemServiceClient)
	logger := zerolog.Nop()

	expectedItems := []*item.ItemData{
		{
			Id:   "item1",
			Type: item.ItemType_ITEM_TYPE_PASSWORD,
		},
		{
			Id:   "item2",
			Type: item.ItemType_ITEM_TYPE_PASSWORD,
		},
	}
	expectedTotal := int32(2)

	mockClient.On("GetItemsV1", mock.Anything, &item.GetItemsV1Request{
		Page:     1,
		PageSize: 10,
	}).Return(&item.GetItemsV1Response{
		Items:      expectedItems,
		TotalCount: expectedTotal,
	}, nil)

	client := &items.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	items, total, err := client.GetItems(t.Context(), 1, 10)
	require.NoError(t, err)
	assert.Equal(t, expectedItems, items)
	assert.Equal(t, expectedTotal, total)
}

func TestGetItems_EmptyResult(t *testing.T) {
	t.Parallel()

	mockClient := new(MockItemServiceClient)
	logger := zerolog.Nop()

	mockClient.On("GetItemsV1", mock.Anything, &item.GetItemsV1Request{
		Page:     2,
		PageSize: 20,
	}).Return(&item.GetItemsV1Response{
		Items:      []*item.ItemData{},
		TotalCount: 0,
	}, nil)

	client := &items.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	items, total, err := client.GetItems(t.Context(), 2, 20)
	require.NoError(t, err)
	assert.Empty(t, items)
	assert.Equal(t, int32(0), total)
}

func TestGetItems_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockItemServiceClient)
	logger := zerolog.Nop()

	mockClient.On("GetItemsV1", mock.Anything, &item.GetItemsV1Request{
		Page:     1,
		PageSize: 10,
	}).Return(nil, errors.New("grpc error"))

	client := &items.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	items, total, err := client.GetItems(t.Context(), 1, 10)
	require.Error(t, err)
	assert.Nil(t, items)
	assert.Equal(t, int32(0), total)
	assert.Contains(t, err.Error(), "GetItems failed")
	assert.Contains(t, err.Error(), "page=1")
	assert.Contains(t, err.Error(), "pageSize=10")
}

func TestGetItems_InvalidPageParams(t *testing.T) {
	t.Parallel()

	mockClient := new(MockItemServiceClient)
	logger := zerolog.Nop()

	// Test with zero page size (should still make the call)
	mockClient.On("GetItemsV1", mock.Anything, &item.GetItemsV1Request{
		Page:     1,
		PageSize: 0,
	}).Return(&item.GetItemsV1Response{
		Items:      []*item.ItemData{},
		TotalCount: 0,
	}, nil)

	client := &items.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, _, err := client.GetItems(t.Context(), 1, 0)
	require.NoError(t, err) // The client doesn't validate parameters

	// Test with negative page number (should still make the call)
	mockClient.On("GetItemsV1", mock.Anything, &item.GetItemsV1Request{
		Page:     -1,
		PageSize: 10,
	}).Return(&item.GetItemsV1Response{
		Items:      []*item.ItemData{},
		TotalCount: 0,
	}, nil)

	_, _, err = client.GetItems(t.Context(), -1, 10)
	require.NoError(t, err) // The client doesn't validate parameters
}

func TestNewItemsClient(t *testing.T) {
	t.Parallel()

	tm := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	client := items.NewItemsClient(nil, tm, &logger)

	assert.NotNil(t, client)
}
