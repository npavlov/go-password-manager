package items_test

import (
	"context"
	"errors"
	"testing"

	"github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/internal/client/grpc/items"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// Mocks

type MockItemServiceClient struct {
	mock.Mock
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockItemServiceClient) GetItems(ctx context.Context, in *item.GetItemsRequest, opts ...grpc.CallOption) (*item.GetItemsResponse, error) {
	args := m.Called(ctx, in)

	// Safely handle nil to avoid type assertion panic
	arg, ok := args.Get(0).(*item.GetItemsResponse)
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

	mockClient.On("GetItems", mock.Anything, &item.GetItemsRequest{
		Page:     1,
		PageSize: 10,
	}).Return(&item.GetItemsResponse{
		Items:      expectedItems,
		TotalCount: expectedTotal,
	}, nil)

	client := &items.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	items, total, err := client.GetItems(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
	assert.Equal(t, expectedTotal, total)
}

func TestGetItems_EmptyResult(t *testing.T) {
	mockClient := new(MockItemServiceClient)
	logger := zerolog.Nop()

	mockClient.On("GetItems", mock.Anything, &item.GetItemsRequest{
		Page:     2,
		PageSize: 20,
	}).Return(&item.GetItemsResponse{
		Items:      []*item.ItemData{},
		TotalCount: 0,
	}, nil)

	client := &items.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	items, total, err := client.GetItems(context.Background(), 2, 20)
	assert.NoError(t, err)
	assert.Empty(t, items)
	assert.Equal(t, int32(0), total)
}

func TestGetItems_Error(t *testing.T) {
	mockClient := new(MockItemServiceClient)
	logger := zerolog.Nop()

	mockClient.On("GetItems", mock.Anything, &item.GetItemsRequest{
		Page:     1,
		PageSize: 10,
	}).Return(nil, errors.New("grpc error"))

	client := &items.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	items, total, err := client.GetItems(context.Background(), 1, 10)
	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Equal(t, int32(0), total)
	assert.Contains(t, err.Error(), "GetItems failed")
	assert.Contains(t, err.Error(), "page=1")
	assert.Contains(t, err.Error(), "pageSize=10")
}

func TestGetItems_InvalidPageParams(t *testing.T) {
	mockClient := new(MockItemServiceClient)
	logger := zerolog.Nop()

	// Test with zero page size (should still make the call)
	mockClient.On("GetItems", mock.Anything, &item.GetItemsRequest{
		Page:     1,
		PageSize: 0,
	}).Return(&item.GetItemsResponse{
		Items:      []*item.ItemData{},
		TotalCount: 0,
	}, nil)

	client := &items.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, _, err := client.GetItems(context.Background(), 1, 0)
	assert.NoError(t, err) // The client doesn't validate parameters

	// Test with negative page number (should still make the call)
	mockClient.On("GetItems", mock.Anything, &item.GetItemsRequest{
		Page:     -1,
		PageSize: 10,
	}).Return(&item.GetItemsResponse{
		Items:      []*item.ItemData{},
		TotalCount: 0,
	}, nil)

	_, _, err = client.GetItems(context.Background(), -1, 10)
	assert.NoError(t, err) // The client doesn't validate parameters
}

func TestNewItemsClient(t *testing.T) {

	tm := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	client := items.NewItemsClient(nil, tm, &logger)

	assert.NotNil(t, client)
}
