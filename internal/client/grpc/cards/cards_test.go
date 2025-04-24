//nolint:wrapcheck,err113,forcetypeassert
package cards_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/npavlov/go-password-manager/gen/proto/card"
	"github.com/npavlov/go-password-manager/internal/client/grpc/cards"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

// Mocks

type MockCardServiceClient struct {
	mock.Mock
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockCardServiceClient) GetCardsV1(ctx context.Context,
	in *card.GetCardsV1Request,
	_ ...grpc.CallOption,
) (*card.GetCardsV1Response, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*card.GetCardsV1Response), args.Error(1)
}

func (m *MockCardServiceClient) GetCardV1(ctx context.Context,
	in *card.GetCardV1Request,
	_ ...grpc.CallOption,
) (*card.GetCardV1Response, error) {
	args := m.Called(ctx, in)

	// Safely handle nil to avoid type assertion panic
	arg, ok := args.Get(0).(*card.GetCardV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockCardServiceClient) UpdateCardV1(ctx context.Context,
	in *card.UpdateCardV1Request,
	_ ...grpc.CallOption,
) (*card.UpdateCardV1Response, error) {
	args := m.Called(ctx, in)

	// Safely handle nil to avoid type assertion panic
	arg, ok := args.Get(0).(*card.UpdateCardV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockCardServiceClient) StoreCardV1(ctx context.Context,
	in *card.StoreCardV1Request,
	_ ...grpc.CallOption,
) (*card.StoreCardV1Response, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*card.StoreCardV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockCardServiceClient) DeleteCardV1(ctx context.Context,
	in *card.DeleteCardV1Request,
	_ ...grpc.CallOption,
) (*card.DeleteCardV1Response, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*card.DeleteCardV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockTokenManager) GetToken() (string, error) {
	args := m.Called()

	return args.String(0), args.Error(1)
}

func TestGetCard_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	expectedCard := &card.CardData{
		CardNumber:     "4111111111111111",
		ExpiryDate:     "12/25",
		Cvv:            "123",
		CardholderName: "John Doe",
	}
	expectedTime := time.Now()

	mockClient.On("GetCardV1", mock.Anything, &card.GetCardV1Request{CardId: "card123"}).
		Return(&card.GetCardV1Response{
			Card:       expectedCard,
			LastUpdate: timestamppb.New(expectedTime),
		}, nil)

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	cardData, _, err := client.GetCard(t.Context(), "card123")
	require.NoError(t, err)
	assert.Equal(t, expectedCard, cardData)
}

func TestGetCard_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	mockClient.On("GetCardV1", mock.Anything, &card.GetCardV1Request{CardId: "card123"}).
		Return(nil, errors.New("get card failed"))

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, _, err := client.GetCard(t.Context(), "card123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error getting password")
}

func TestUpdateCard_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	updateReq := &card.UpdateCardV1Request{
		CardId: "card123",
		Data: &card.CardData{
			CardNumber:     "4111111111111111",
			ExpiryDate:     "12/25",
			Cvv:            "123",
			CardholderName: "John Doe",
		},
	}

	mockClient.On("UpdateCardV1", mock.Anything, updateReq).
		Return(&card.UpdateCardV1Response{
			CardId: "card123",
		}, nil)

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	err := client.UpdateCard(t.Context(), "card123", "4111111111111111", "12/25", "123", "John Doe")
	require.NoError(t, err)
}

func TestUpdateCard_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	updateReq := &card.UpdateCardV1Request{
		CardId: "card123",
		Data: &card.CardData{
			CardNumber:     "4111111111111111",
			ExpiryDate:     "12/25",
			Cvv:            "123",
			CardholderName: "John Doe",
		},
	}

	mockClient.On("UpdateCardV1", mock.Anything, updateReq).
		Return(nil, errors.New("update failed"))

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	err := client.UpdateCard(t.Context(), "card123", "4111111111111111", "12/25", "123", "John Doe")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error updating card")
}

func TestStoreCard_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	storeReq := &card.StoreCardV1Request{
		Card: &card.CardData{
			CardNumber:     "4111111111111111",
			ExpiryDate:     "12/25",
			Cvv:            "123",
			CardholderName: "John Doe",
		},
	}

	mockClient.On("StoreCardV1", mock.Anything, storeReq).
		Return(&card.StoreCardV1Response{CardId: "new-card-123"}, nil)

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	cardID, err := client.StoreCard(t.Context(), "4111111111111111", "12/25", "123", "John Doe")
	require.NoError(t, err)
	assert.Equal(t, "new-card-123", cardID)
}

func TestStoreCard_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	storeReq := &card.StoreCardV1Request{
		Card: &card.CardData{
			CardNumber:     "4111111111111111",
			ExpiryDate:     "12/25",
			Cvv:            "123",
			CardholderName: "John Doe",
		},
	}

	mockClient.On("StoreCardV1", mock.Anything, storeReq).
		Return(nil, errors.New("store failed"))

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.StoreCard(t.Context(), "4111111111111111", "12/25", "123", "John Doe")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error storing card")
}

func TestDeleteCard_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteCardV1", mock.Anything, &card.DeleteCardV1Request{CardId: "card123"}).
		Return(&card.DeleteCardV1Response{Ok: true}, nil)

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeleteCard(t.Context(), "card123")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestDeleteCard_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteCardV1", mock.Anything, &card.DeleteCardV1Request{CardId: "card123"}).
		Return(nil, errors.New("delete failed"))

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.DeleteCard(t.Context(), "card123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting card")
}

func TestDeleteCard_NotOk(t *testing.T) {
	t.Parallel()

	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteCardV1", mock.Anything, &card.DeleteCardV1Request{CardId: "card123"}).
		Return(&card.DeleteCardV1Response{Ok: false}, nil)

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeleteCard(t.Context(), "card123")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestNewCardsClient(t *testing.T) {
	t.Parallel()

	tm := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	client := cards.NewCardClient(nil, tm, &logger)

	assert.NotNil(t, client)
}
