package cards_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/npavlov/go-password-manager/gen/proto/card"
	"github.com/npavlov/go-password-manager/internal/client/grpc/cards"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mocks

type MockCardServiceClient struct {
	mock.Mock
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockCardServiceClient) GetCards(ctx context.Context, in *card.GetCardsRequest, opts ...grpc.CallOption) (*card.GetCardsResponse, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*card.GetCardsResponse), args.Error(1)
}

func (m *MockCardServiceClient) GetCard(ctx context.Context, in *card.GetCardRequest, opts ...grpc.CallOption) (*card.GetCardResponse, error) {
	args := m.Called(ctx, in)

	// Safely handle nil to avoid type assertion panic
	arg, ok := args.Get(0).(*card.GetCardResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockCardServiceClient) UpdateCard(ctx context.Context, in *card.UpdateCardRequest, opts ...grpc.CallOption) (*card.UpdateCardResponse, error) {
	args := m.Called(ctx, in)

	// Safely handle nil to avoid type assertion panic
	arg, ok := args.Get(0).(*card.UpdateCardResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockCardServiceClient) StoreCard(ctx context.Context, in *card.StoreCardRequest, opts ...grpc.CallOption) (*card.StoreCardResponse, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*card.StoreCardResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockCardServiceClient) DeleteCard(ctx context.Context, in *card.DeleteCardRequest, opts ...grpc.CallOption) (*card.DeleteCardResponse, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*card.DeleteCardResponse)
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
	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	expectedCard := &card.CardData{
		CardNumber:     "4111111111111111",
		ExpiryDate:     "12/25",
		Cvv:            "123",
		CardholderName: "John Doe",
	}
	expectedTime := time.Now()

	mockClient.On("GetCard", mock.Anything, &card.GetCardRequest{CardId: "card123"}).
		Return(&card.GetCardResponse{
			Card:       expectedCard,
			LastUpdate: timestamppb.New(expectedTime),
		}, nil)

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	cardData, _, err := client.GetCard(context.Background(), "card123")
	assert.NoError(t, err)
	assert.Equal(t, expectedCard, cardData)
}

func TestGetCard_Error(t *testing.T) {
	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	mockClient.On("GetCard", mock.Anything, &card.GetCardRequest{CardId: "card123"}).
		Return(nil, errors.New("get card failed"))

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, _, err := client.GetCard(context.Background(), "card123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting password")
}

func TestUpdateCard_Success(t *testing.T) {
	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	updateReq := &card.UpdateCardRequest{
		CardId: "card123",
		Data: &card.CardData{
			CardNumber:     "4111111111111111",
			ExpiryDate:     "12/25",
			Cvv:            "123",
			CardholderName: "John Doe",
		},
	}

	mockClient.On("UpdateCard", mock.Anything, updateReq).
		Return(&card.UpdateCardResponse{}, nil)

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	err := client.UpdateCard(context.Background(), "card123", "4111111111111111", "12/25", "123", "John Doe")
	assert.NoError(t, err)
}

func TestUpdateCard_Error(t *testing.T) {
	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	updateReq := &card.UpdateCardRequest{
		CardId: "card123",
		Data: &card.CardData{
			CardNumber:     "4111111111111111",
			ExpiryDate:     "12/25",
			Cvv:            "123",
			CardholderName: "John Doe",
		},
	}

	mockClient.On("UpdateCard", mock.Anything, updateReq).
		Return(nil, errors.New("update failed"))

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	err := client.UpdateCard(context.Background(), "card123", "4111111111111111", "12/25", "123", "John Doe")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error updating card")
}

func TestStoreCard_Success(t *testing.T) {
	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	storeReq := &card.StoreCardRequest{
		Card: &card.CardData{
			CardNumber:     "4111111111111111",
			ExpiryDate:     "12/25",
			Cvv:            "123",
			CardholderName: "John Doe",
		},
	}

	mockClient.On("StoreCard", mock.Anything, storeReq).
		Return(&card.StoreCardResponse{CardId: "new-card-123"}, nil)

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	cardID, err := client.StoreCard(context.Background(), "4111111111111111", "12/25", "123", "John Doe")
	assert.NoError(t, err)
	assert.Equal(t, "new-card-123", cardID)
}

func TestStoreCard_Error(t *testing.T) {
	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	storeReq := &card.StoreCardRequest{
		Card: &card.CardData{
			CardNumber:     "4111111111111111",
			ExpiryDate:     "12/25",
			Cvv:            "123",
			CardholderName: "John Doe",
		},
	}

	mockClient.On("StoreCard", mock.Anything, storeReq).
		Return(nil, errors.New("store failed"))

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.StoreCard(context.Background(), "4111111111111111", "12/25", "123", "John Doe")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error storing card")
}

func TestDeleteCard_Success(t *testing.T) {
	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteCard", mock.Anything, &card.DeleteCardRequest{CardId: "card123"}).
		Return(&card.DeleteCardResponse{Ok: true}, nil)

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeleteCard(context.Background(), "card123")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestDeleteCard_Error(t *testing.T) {
	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteCard", mock.Anything, &card.DeleteCardRequest{CardId: "card123"}).
		Return(nil, errors.New("delete failed"))

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.DeleteCard(context.Background(), "card123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting card")
}

func TestDeleteCard_NotOk(t *testing.T) {
	mockClient := new(MockCardServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteCard", mock.Anything, &card.DeleteCardRequest{CardId: "card123"}).
		Return(&card.DeleteCardResponse{Ok: false}, nil)

	client := &cards.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeleteCard(context.Background(), "card123")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestNewCardsClient(t *testing.T) {

	tm := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	client := cards.NewCardClient(nil, tm, &logger)

	assert.NotNil(t, client)
}
