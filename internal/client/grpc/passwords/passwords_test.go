package passwords_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/client/grpc/passwords"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mocks

type MockPasswordServiceClient struct {
	mock.Mock
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockPasswordServiceClient) GetPasswords(ctx context.Context, in *password.GetPasswordsRequest, opts ...grpc.CallOption) (*password.GetPasswordsResponse, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*password.GetPasswordsResponse), args.Error(1)
}

func (m *MockPasswordServiceClient) GetPassword(ctx context.Context, in *password.GetPasswordRequest, opts ...grpc.CallOption) (*password.GetPasswordResponse, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*password.GetPasswordResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockPasswordServiceClient) UpdatePassword(ctx context.Context, in *password.UpdatePasswordRequest, opts ...grpc.CallOption) (*password.UpdatePasswordResponse, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*password.UpdatePasswordResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockPasswordServiceClient) StorePassword(ctx context.Context, in *password.StorePasswordRequest, opts ...grpc.CallOption) (*password.StorePasswordResponse, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*password.StorePasswordResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockPasswordServiceClient) DeletePassword(ctx context.Context, in *password.DeletePasswordRequest, opts ...grpc.CallOption) (*password.DeletePasswordResponse, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*password.DeletePasswordResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockTokenManager) GetToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestGetPassword_Success(t *testing.T) {
	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	expectedPassword := &password.PasswordData{
		Login:    "testuser",
		Password: "securepassword123",
	}
	expectedTime := time.Now()

	mockClient.On("GetPassword", mock.Anything, &password.GetPasswordRequest{
		PasswordId: "pass123",
	}).Return(&password.GetPasswordResponse{
		Password:   expectedPassword,
		LastUpdate: timestamppb.New(expectedTime),
	}, nil)

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	passData, _, err := client.GetPassword(context.Background(), "pass123")
	assert.NoError(t, err)
	assert.Equal(t, expectedPassword, passData)
}

func TestGetPassword_Error(t *testing.T) {
	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("GetPassword", mock.Anything, &password.GetPasswordRequest{
		PasswordId: "pass123",
	}).Return(nil, errors.New("get password failed"))

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, _, err := client.GetPassword(context.Background(), "pass123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting password")
}

func TestUpdatePassword_Success(t *testing.T) {
	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	updateReq := &password.UpdatePasswordRequest{
		PasswordId: "pass123",
		Data: &password.PasswordData{
			Login:    "updateduser",
			Password: "newpassword456",
		},
	}

	mockClient.On("UpdatePassword", mock.Anything, updateReq).
		Return(&password.UpdatePasswordResponse{}, nil)

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	err := client.UpdatePassword(context.Background(), "pass123", "updateduser", "newpassword456")
	assert.NoError(t, err)
}

func TestUpdatePassword_Error(t *testing.T) {
	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	updateReq := &password.UpdatePasswordRequest{
		PasswordId: "pass123",
		Data: &password.PasswordData{
			Login:    "updateduser",
			Password: "newpassword456",
		},
	}

	mockClient.On("UpdatePassword", mock.Anything, updateReq).
		Return(nil, errors.New("update failed"))

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	err := client.UpdatePassword(context.Background(), "pass123", "updateduser", "newpassword456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error updating password")
}

func TestStorePassword_Success(t *testing.T) {
	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	storeReq := &password.StorePasswordRequest{
		Password: &password.PasswordData{
			Login:    "newuser",
			Password: "securepass789",
		},
	}

	mockClient.On("StorePassword", mock.Anything, storeReq).
		Return(&password.StorePasswordResponse{
			PasswordId: "new-pass-456",
		}, nil)

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	passID, err := client.StorePassword(context.Background(), "newuser", "securepass789")
	assert.NoError(t, err)
	assert.Equal(t, "new-pass-456", passID)
}

func TestStorePassword_Error(t *testing.T) {
	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	storeReq := &password.StorePasswordRequest{
		Password: &password.PasswordData{
			Login:    "newuser",
			Password: "securepass789",
		},
	}

	mockClient.On("StorePassword", mock.Anything, storeReq).
		Return(nil, errors.New("store failed"))

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.StorePassword(context.Background(), "newuser", "securepass789")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error storing password")
}

func TestDeletePassword_Success(t *testing.T) {
	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeletePassword", mock.Anything, &password.DeletePasswordRequest{
		PasswordId: "pass123",
	}).Return(&password.DeletePasswordResponse{
		Ok: true,
	}, nil)

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeletePassword(context.Background(), "pass123")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestDeletePassword_Error(t *testing.T) {
	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeletePassword", mock.Anything, &password.DeletePasswordRequest{
		PasswordId: "pass123",
	}).Return(nil, errors.New("delete failed"))

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.DeletePassword(context.Background(), "pass123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting password")
}

func TestDeletePassword_NotSuccessful(t *testing.T) {
	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeletePassword", mock.Anything, &password.DeletePasswordRequest{
		PasswordId: "pass123",
	}).Return(&password.DeletePasswordResponse{
		Ok: false,
	}, nil)

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeletePassword(context.Background(), "pass123")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestNewPasswordClient(t *testing.T) {

	tm := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	client := passwords.NewPasswordClient(nil, tm, &logger)

	assert.NotNil(t, client)
}
