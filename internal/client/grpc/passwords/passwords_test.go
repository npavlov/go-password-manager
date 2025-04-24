//nolint:err113,wrapcheck,lll,forcetypeassert,exhaustruct
package passwords_test

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

	"github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/client/grpc/passwords"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

// Mocks

type MockPasswordServiceClient struct {
	mock.Mock
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockPasswordServiceClient) GetPasswordsV1(ctx context.Context, in *password.GetPasswordsV1Request, _ ...grpc.CallOption) (*password.GetPasswordsV1Response, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*password.GetPasswordsV1Response), args.Error(1)
}

func (m *MockPasswordServiceClient) GetPasswordV1(ctx context.Context, in *password.GetPasswordV1Request, _ ...grpc.CallOption) (*password.GetPasswordV1Response, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*password.GetPasswordV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockPasswordServiceClient) UpdatePasswordV1(ctx context.Context, in *password.UpdatePasswordV1Request, _ ...grpc.CallOption) (*password.UpdatePasswordV1Response, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*password.UpdatePasswordV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockPasswordServiceClient) StorePasswordV1(ctx context.Context, in *password.StorePasswordV1Request, _ ...grpc.CallOption) (*password.StorePasswordV1Response, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*password.StorePasswordV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockPasswordServiceClient) DeletePasswordV1(ctx context.Context, in *password.DeletePasswordV1Request, _ ...grpc.CallOption) (*password.DeletePasswordV1Response, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*password.DeletePasswordV1Response)
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
	t.Parallel()

	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	expectedPassword := &password.PasswordData{
		Login:    "testuser",
		Password: "securepassword123",
	}
	expectedTime := time.Now()

	mockClient.On("GetPasswordV1", mock.Anything, &password.GetPasswordV1Request{
		PasswordId: "pass123",
	}).Return(&password.GetPasswordV1Response{
		Password:   expectedPassword,
		LastUpdate: timestamppb.New(expectedTime),
	}, nil)

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	passData, _, err := client.GetPassword(t.Context(), "pass123")
	require.NoError(t, err)
	assert.Equal(t, expectedPassword, passData)
}

func TestGetPassword_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("GetPasswordV1", mock.Anything, &password.GetPasswordV1Request{
		PasswordId: "pass123",
	}).Return(nil, errors.New("get password failed"))

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, _, err := client.GetPassword(t.Context(), "pass123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error getting password")
}

func TestUpdatePassword_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("UpdatePasswordV1", mock.Anything, mock.Anything).
		Return(&password.UpdatePasswordV1Response{}, nil)

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	err := client.UpdatePassword(t.Context(), "pass123", "updateduser", "newpassword456")
	require.NoError(t, err)
}

func TestUpdatePassword_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("UpdatePasswordV1", mock.Anything, mock.Anything).
		Return(nil, errors.New("update failed"))

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	err := client.UpdatePassword(t.Context(), "pass123", "updateduser", "newpassword456")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error updating password")
}

func TestStorePassword_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("StorePasswordV1", mock.Anything, mock.Anything).
		Return(&password.StorePasswordV1Response{
			PasswordId: "new-pass-456",
		}, nil)

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	passID, err := client.StorePassword(t.Context(), "newuser", "securepass789")
	require.NoError(t, err)
	assert.Equal(t, "new-pass-456", passID)
}

func TestStorePassword_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("StorePasswordV1", mock.Anything, mock.Anything).
		Return(nil, errors.New("store failed"))

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.StorePassword(t.Context(), "newuser", "securepass789")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error storing password")
}

func TestDeletePassword_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeletePasswordV1", mock.Anything, &password.DeletePasswordV1Request{
		PasswordId: "pass123",
	}).Return(&password.DeletePasswordV1Response{
		Ok: true,
	}, nil)

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeletePassword(t.Context(), "pass123")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestDeletePassword_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeletePasswordV1", mock.Anything, &password.DeletePasswordV1Request{
		PasswordId: "pass123",
	}).Return(nil, errors.New("delete failed"))

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.DeletePassword(t.Context(), "pass123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting password")
}

func TestDeletePassword_NotSuccessful(t *testing.T) {
	t.Parallel()

	mockClient := new(MockPasswordServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeletePasswordV1", mock.Anything, &password.DeletePasswordV1Request{
		PasswordId: "pass123",
	}).Return(&password.DeletePasswordV1Response{
		Ok: false,
	}, nil)

	client := &passwords.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeletePassword(t.Context(), "pass123")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestNewPasswordClient(t *testing.T) {
	t.Parallel()

	tm := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	client := passwords.NewPasswordClient(nil, tm, &logger)

	assert.NotNil(t, client)
}
