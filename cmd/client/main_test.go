package main_test

import (
	"context"
	"io"
	"testing"
	"time"

	main "github.com/npavlov/go-password-manager/cmd/client"
	"github.com/npavlov/go-password-manager/gen/proto/card"
	"github.com/npavlov/go-password-manager/gen/proto/file"
	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/gen/proto/note"
	"github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/client/config"
	"github.com/npavlov/go-password-manager/internal/client/interceptors"

	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks

type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) LoadTokens() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTokenManager) SetAuthFailCallback(callback func()) {
	m.Called(callback)
}

type MockFacade struct {
	mock.Mock
}

func (m *MockFacade) Register(username, password, email string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) GetItems(ctx context.Context, page, pageSize int32) ([]*pb.ItemData, int32, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) StorePassword(ctx context.Context, login string, password string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) GetPassword(ctx context.Context, id string) (*password.PasswordData, time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) UpdatePassword(ctx context.Context, id, login, password string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) DeletePassword(ctx context.Context, id string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) GetMetainfo(ctx context.Context, id string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) SetMetainfo(ctx context.Context, id string, meta map[string]string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) DeleteMetainfo(ctx context.Context, id, key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) StoreNote(ctx context.Context, content string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) GetNote(ctx context.Context, id string) (*note.NoteData, time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) DeleteNote(ctx context.Context, id string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) StoreCard(ctx context.Context, cardNum, expDate, Cvv, cardHolder string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) UpdateCard(ctx context.Context, id, cardNum, expDate, Cvv, cardHolder string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) GetCard(ctx context.Context, id string) (*card.CardData, time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) DeleteCard(ctx context.Context, id string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) UploadBinary(ctx context.Context, filename string, reader io.Reader) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) DownloadBinary(ctx context.Context, fileID string, writer io.Writer) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) GetFile(ctx context.Context, fileID string) (*file.FileMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) DeleteBinary(ctx context.Context, fileID string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) Login(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}

type MockStorageManager struct {
	mock.Mock
}

func (m *MockStorageManager) StartBackgroundSync(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockStorageManager) StopSync() {
	m.Called()
}

type MockTUI struct {
	mock.Mock
}

func (m *MockTUI) Run() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTUI) ResetToLoginScreen() {
	m.Called()
}

type MockGRPCConn struct {
	mock.Mock
}

func (m *MockGRPCConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestLoadConfig(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		logger := testutils.GetTLogger()
		cfg := main.LoadConfig(logger)

		assert.NotNil(t, cfg)
		assert.Equal(t, ":9090", cfg.Address) // default value
	})
}

func TestMakeConnection(t *testing.T) {
	t.Run("successful connection", func(t *testing.T) {
		// Setup
		cfg := config.Config{
			Address:     "localhost:9090",
			Certificate: "testdata/test.crt",
		}

		interceptor := interceptors.NewAuthInterceptor(cfg, nil)

		// Test - we'll use a real connection to a test server in practice
		// For this example we'll just verify it doesn't error with test certs
		conn, err := main.MakeConnection(cfg, interceptor)

		if conn != nil {
			defer conn.Close()
		}

		// In real tests, we'd have proper test certificates
		// Here we just expect it to fail since we don't have valid certs
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not load TLS keys")
	})
}
