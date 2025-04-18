//nolint:wrapcheck,exhaustruct,forcetypeassert
package facade_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	pb_card "github.com/npavlov/go-password-manager/gen/proto/card"
	pb_file "github.com/npavlov/go-password-manager/gen/proto/file"
	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	pb_note "github.com/npavlov/go-password-manager/gen/proto/note"
	pb_password "github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/client/grpc/facade"
)

// Mock implementations for all client interfaces

type MockAuthClient struct{ mock.Mock }

func (m *MockAuthClient) Login(username, password string) error {
	return m.Called(username, password).Error(0)
}

func (m *MockAuthClient) Register(username, password, email string) (string, error) {
	args := m.Called(username, password, email)

	return args.String(0), args.Error(1)
}

type MockItemsClient struct{ mock.Mock }

func (m *MockItemsClient) GetItems(ctx context.Context, page, pageSize int32) ([]*pb.ItemData, int32, error) {
	args := m.Called(ctx, page, pageSize)

	return args.Get(0).([]*pb.ItemData), args.Get(1).(int32), args.Error(2)
}

type MockPasswordClient struct{ mock.Mock }

func (m *MockPasswordClient) StorePassword(ctx context.Context, login, password string) (string, error) {
	args := m.Called(ctx, login, password)

	return args.String(0), args.Error(1)
}

func (m *MockPasswordClient) GetPassword(ctx context.Context, id string) (*pb_password.PasswordData, time.Time, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*pb_password.PasswordData), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockPasswordClient) UpdatePassword(ctx context.Context, id, login, password string) error {
	return m.Called(ctx, id, login, password).Error(0)
}

func (m *MockPasswordClient) DeletePassword(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)

	return args.Bool(0), args.Error(1)
}

type MockMetaClient struct{ mock.Mock }

func (m *MockMetaClient) GetMetainfo(ctx context.Context, id string) (map[string]string, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockMetaClient) SetMetainfo(ctx context.Context, id string, meta map[string]string) (bool, error) {
	args := m.Called(ctx, id, meta)

	return args.Bool(0), args.Error(1)
}

func (m *MockMetaClient) DeleteMetainfo(ctx context.Context, id, key string) (bool, error) {
	args := m.Called(ctx, id, key)

	return args.Bool(0), args.Error(1)
}

type MockNoteClient struct{ mock.Mock }

func (m *MockNoteClient) StoreNote(ctx context.Context, content string) (string, error) {
	args := m.Called(ctx, content)

	return args.String(0), args.Error(1)
}

func (m *MockNoteClient) GetNote(ctx context.Context, id string) (*pb_note.NoteData, time.Time, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*pb_note.NoteData), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockNoteClient) DeleteNote(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)

	return args.Bool(0), args.Error(1)
}

type MockCardClient struct{ mock.Mock }

func (m *MockCardClient) StoreCard(ctx context.Context, cardNum, expDate, cvv, cardHolder string) (string, error) {
	args := m.Called(ctx, cardNum, expDate, cvv, cardHolder)

	return args.String(0), args.Error(1)
}

func (m *MockCardClient) UpdateCard(ctx context.Context, id, cardNum, expDate, cvv, cardHolder string) error {
	return m.Called(ctx, id, cardNum, expDate, cvv, cardHolder).Error(0)
}

func (m *MockCardClient) GetCard(ctx context.Context, id string) (*pb_card.CardData, time.Time, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*pb_card.CardData), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockCardClient) DeleteCard(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)

	return args.Bool(0), args.Error(1)
}

type MockBinaryClient struct{ mock.Mock }

func (m *MockBinaryClient) UploadFile(ctx context.Context, filename string, reader io.Reader) (string, error) {
	args := m.Called(ctx, filename, reader)

	return args.String(0), args.Error(1)
}

func (m *MockBinaryClient) DownloadFile(ctx context.Context, fileID string, writer io.Writer) error {
	return m.Called(ctx, fileID, writer).Error(0)
}

func (m *MockBinaryClient) GetFile(ctx context.Context, fileID string) (*pb_file.FileMeta, error) {
	args := m.Called(ctx, fileID)

	return args.Get(0).(*pb_file.FileMeta), args.Error(1)
}

func (m *MockBinaryClient) DeleteFile(ctx context.Context, fileID string) (bool, error) {
	args := m.Called(ctx, fileID)

	return args.Bool(0), args.Error(1)
}

// Test setup helper.
func setupFacadeTest() (*facade.Facade,
	*MockAuthClient,
	*MockItemsClient,
	*MockPasswordClient,
	*MockMetaClient,
	*MockNoteClient,
	*MockCardClient,
	*MockBinaryClient,
) {
	authMock := &MockAuthClient{}
	itemsMock := &MockItemsClient{}
	passMock := &MockPasswordClient{}
	metaMock := &MockMetaClient{}
	noteMock := &MockNoteClient{}
	cardMock := &MockCardClient{}
	binaryMock := &MockBinaryClient{}

	facade := facade.NewFacadeWithOptions(facade.Options{
		AuthClient:     authMock,
		ItemsClient:    itemsMock,
		PasswordClient: passMock,
		MetaClient:     metaMock,
		NoteClient:     noteMock,
		CardClient:     cardMock,
		BinaryClient:   binaryMock,
	})

	return facade, authMock, itemsMock, passMock, metaMock, noteMock, cardMock, binaryMock
}

// Test cases

func TestFacade_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "successful login",
			username: "user1",
			password: "pass1",
			wantErr:  false,
		},
		{
			name:     "failed login",
			username: "user1",
			password: "wrongpass",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fClient, authMock, _, _, _, _, _, _ := setupFacadeTest()
			if tt.wantErr {
				authMock.On("Login", tt.username, tt.password).Return(errors.New("error")).Once()
			} else {
				authMock.On("Login", tt.username, tt.password).Return(nil)
			}

			err := fClient.Login(tt.username, tt.password)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			authMock.AssertExpectations(t)
		})
	}
}

func TestFacade_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		username  string
		password  string
		email     string
		mockKey   string
		mockErr   error
		wantKey   string
		wantErr   bool
		wantError string
	}{
		{
			name:     "successful registration",
			username: "newuser",
			password: "newpass",
			email:    "new@email.com",
			mockKey:  "master-key-123",
			mockErr:  nil,
			wantKey:  "master-key-123",
			wantErr:  false,
		},
		{
			name:      "failed registration",
			username:  "existinguser",
			password:  "pass",
			email:     "existing@email.com",
			mockKey:   "",
			mockErr:   errors.New("username exists"),
			wantKey:   "",
			wantErr:   true,
			wantError: "failed to register user existinguser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fClient, authMock, _, _, _, _, _, _ := setupFacadeTest()

			authMock.On("Register", tt.username, tt.password, tt.email).
				Return(tt.mockKey, tt.mockErr).Once()

			key, err := fClient.Register(tt.username, tt.password, tt.email)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantError)
				assert.Equal(t, tt.wantKey, key)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantKey, key)
			}

			authMock.AssertExpectations(t)
		})
	}
}

func TestFacade_GetItems(t *testing.T) {
	t.Parallel()

	page := int32(1)
	pageSize := int32(10)
	items := []*pb.ItemData{
		{Id: "item1", Type: pb.ItemType_ITEM_TYPE_PASSWORD},
		{Id: "item2", Type: pb.ItemType_ITEM_TYPE_NOTE},
	}
	total := int32(2)

	tests := []struct {
		name      string
		mockItems []*pb.ItemData
		mockTotal int32
		mockErr   error
		wantItems []*pb.ItemData
		wantTotal int32
		wantErr   bool
	}{
		{
			name:      "successful get items",
			mockItems: items,
			mockTotal: total,
			mockErr:   nil,
			wantItems: items,
			wantTotal: total,
			wantErr:   false,
		},
		{
			name:      "failed get items",
			mockItems: nil,
			mockTotal: 0,
			mockErr:   errors.New("database error"),
			wantItems: nil,
			wantTotal: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fClient, _, itemsMock, _, _, _, _, _ := setupFacadeTest()

			ctx := t.Context()

			itemsMock.On("GetItems", ctx, page, pageSize).
				Return(tt.mockItems, tt.mockTotal, tt.mockErr).Once()

			resultItems, resultTotal, err := fClient.GetItems(ctx, page, pageSize)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "error getting items")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantItems, resultItems)
				assert.Equal(t, tt.wantTotal, resultTotal)
			}

			itemsMock.AssertExpectations(t)
		})
	}
}

func TestFacade_PasswordOperations(t *testing.T) {
	t.Parallel()

	testPass := &pb_password.PasswordData{
		Login:    "user1",
		Password: "pass1",
	}

	// StorePassword
	t.Run("StorePassword success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, passMock, _, _, _, _ := setupFacadeTest()
		ctx := t.Context()

		passMock.On("StorePassword", ctx, "user1", "pass1").
			Return("pass-123", nil).Once()

		id, err := fClient.StorePassword(ctx, "user1", "pass1")
		require.NoError(t, err)
		assert.Equal(t, "pass-123", id)
		passMock.AssertExpectations(t)
	})

	// GetPassword
	t.Run("GetPassword success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, passMock, _, _, _, _ := setupFacadeTest()
		ctx := t.Context()
		now := time.Now()

		passMock.On("GetPassword", ctx, "pass-123").
			Return(testPass, now, nil).Once()

		pass, updated, err := fClient.GetPassword(ctx, "pass-123")
		require.NoError(t, err)
		assert.Equal(t, testPass, pass)
		assert.Equal(t, now, updated)
		passMock.AssertExpectations(t)
	})

	// UpdatePassword
	t.Run("UpdatePassword success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, passMock, _, _, _, _ := setupFacadeTest()
		ctx := t.Context()

		passMock.On("UpdatePassword", ctx, "pass-123", "user1", "newpass").
			Return(nil).Once()

		err := fClient.UpdatePassword(ctx, "pass-123", "user1", "newpass")
		require.NoError(t, err)
		passMock.AssertExpectations(t)
	})

	// DeletePassword
	t.Run("DeletePassword success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, passMock, _, _, _, _ := setupFacadeTest()
		ctx := t.Context()

		passMock.On("DeletePassword", ctx, "pass-123").
			Return(true, nil).Once()

		deleted, err := fClient.DeletePassword(ctx, "pass-123")
		require.NoError(t, err)
		assert.True(t, deleted)
		passMock.AssertExpectations(t)
	})
}

func TestFacade_MetaOperations(t *testing.T) {
	t.Parallel()

	// GetMetainfo
	t.Run("GetMetainfo success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, metaMock, _, _, _ := setupFacadeTest()
		ctx := t.Context()
		meta := map[string]string{"key1": "value1", "key2": "value2"}

		metaMock.On("GetMetainfo", ctx, "item-123").
			Return(meta, nil).Once()

		result, err := fClient.GetMetainfo(ctx, "item-123")
		require.NoError(t, err)
		assert.Equal(t, meta, result)
		metaMock.AssertExpectations(t)
	})

	// SetMetainfo
	t.Run("SetMetainfo success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, metaMock, _, _, _ := setupFacadeTest()
		ctx := t.Context()
		meta := map[string]string{"key1": "value1", "key2": "value2"}

		metaMock.On("SetMetainfo", ctx, "item-123", meta).
			Return(true, nil).Once()

		success, err := fClient.SetMetainfo(ctx, "item-123", meta)
		require.NoError(t, err)
		assert.True(t, success)
		metaMock.AssertExpectations(t)
	})

	// DeleteMetainfo
	t.Run("DeleteMetainfo success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, metaMock, _, _, _ := setupFacadeTest()
		ctx := t.Context()

		metaMock.On("DeleteMetainfo", ctx, "item-123", "key1").
			Return(true, nil).Once()

		success, err := fClient.DeleteMetainfo(ctx, "item-123", "key1")
		require.NoError(t, err)
		assert.True(t, success)
		metaMock.AssertExpectations(t)
	})
}

func TestFacade_NoteOperations(t *testing.T) {
	t.Parallel()

	testNote := &pb_note.NoteData{Content: "test content"}

	// StoreNote
	t.Run("StoreNote success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, _, noteMock, _, _ := setupFacadeTest()
		ctx := t.Context()

		noteMock.On("StoreNote", ctx, "test content").
			Return("note-123", nil).Once()

		id, err := fClient.StoreNote(ctx, "test content")
		require.NoError(t, err)
		assert.Equal(t, "note-123", id)
		noteMock.AssertExpectations(t)
	})

	// GetNote
	t.Run("GetNote success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, _, noteMock, _, _ := setupFacadeTest()
		ctx := t.Context()
		now := time.Now()

		noteMock.On("GetNote", ctx, "note-123").
			Return(testNote, now, nil).Once()

		note, updated, err := fClient.GetNote(ctx, "note-123")
		require.NoError(t, err)
		assert.Equal(t, testNote, note)
		assert.Equal(t, now, updated)
		noteMock.AssertExpectations(t)
	})

	// DeleteNote
	t.Run("DeleteNote success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, _, noteMock, _, _ := setupFacadeTest()
		ctx := t.Context()

		noteMock.On("DeleteNote", ctx, "note-123").
			Return(true, nil).Once()

		deleted, err := fClient.DeleteNote(ctx, "note-123")
		require.NoError(t, err)
		assert.True(t, deleted)
		noteMock.AssertExpectations(t)
	})
}

func TestFacade_CardOperations(t *testing.T) {
	t.Parallel()

	testCard := &pb_card.CardData{
		CardNumber:     "4111111111111111",
		ExpiryDate:     "12/25",
		Cvv:            "123",
		CardholderName: "John Doe",
	}

	// StoreCard
	t.Run("StoreCard success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, _, _, cardMock, _ := setupFacadeTest()
		ctx := t.Context()

		cardMock.On("StoreCard", ctx, "4111111111111111", "12/25", "123", "John Doe").
			Return("card-123", nil).Once()

		id, err := fClient.StoreCard(ctx, "4111111111111111", "12/25", "123", "John Doe")
		require.NoError(t, err)
		assert.Equal(t, "card-123", id)
		cardMock.AssertExpectations(t)
	})

	// GetCard
	t.Run("GetCard success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, _, _, cardMock, _ := setupFacadeTest()
		ctx := t.Context()
		now := time.Now()

		cardMock.On("GetCard", ctx, "card-123").
			Return(testCard, now, nil).Once()

		card, updated, err := fClient.GetCard(ctx, "card-123")
		require.NoError(t, err)
		assert.Equal(t, testCard, card)
		assert.Equal(t, now, updated)
		cardMock.AssertExpectations(t)
	})

	// UpdateCard
	t.Run("UpdateCard success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, _, _, cardMock, _ := setupFacadeTest()
		ctx := t.Context()

		cardMock.On("UpdateCard", ctx, "card-123", "4111111111111111", "12/26", "456", "John Doe").
			Return(nil).Once()

		err := fClient.UpdateCard(ctx, "card-123", "4111111111111111", "12/26", "456", "John Doe")
		require.NoError(t, err)
		cardMock.AssertExpectations(t)
	})

	// DeleteCard
	t.Run("DeleteCard success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, _, _, cardMock, _ := setupFacadeTest()
		ctx := t.Context()

		cardMock.On("DeleteCard", ctx, "card-123").
			Return(true, nil).Once()

		deleted, err := fClient.DeleteCard(ctx, "card-123")
		require.NoError(t, err)
		assert.True(t, deleted)
		cardMock.AssertExpectations(t)
	})
}

func TestFacade_BinaryOperations(t *testing.T) {
	t.Parallel()

	testFile := &pb_file.FileMeta{
		Id:       "file-123",
		FileName: "test.txt",
		FileSize: 1024,
	}

	// UploadBinary
	t.Run("UploadBinary success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, _, _, _, binaryMock := setupFacadeTest()
		ctx := t.Context()

		reader := bytes.NewBufferString("test data")
		binaryMock.On("UploadFile", ctx, "test.txt", reader).
			Return("file-123", nil).Once()

		id, err := fClient.UploadBinary(ctx, "test.txt", reader)
		require.NoError(t, err)
		assert.Equal(t, "file-123", id)
		binaryMock.AssertExpectations(t)
	})

	// DownloadBinary
	t.Run("DownloadBinary success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, _, _, _, binaryMock := setupFacadeTest()
		ctx := t.Context()

		writer := &bytes.Buffer{}
		binaryMock.On("DownloadFile", ctx, "file-123", writer).
			Return(nil).Once()

		err := fClient.DownloadBinary(ctx, "file-123", writer)
		require.NoError(t, err)
		binaryMock.AssertExpectations(t)
	})

	// GetFile
	t.Run("GetFile success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, _, _, _, binaryMock := setupFacadeTest()
		ctx := t.Context()

		binaryMock.On("GetFile", ctx, "file-123").
			Return(testFile, nil).Once()

		file, err := fClient.GetFile(ctx, "file-123")
		require.NoError(t, err)
		assert.Equal(t, testFile, file)
		binaryMock.AssertExpectations(t)
	})

	// DeleteBinary
	t.Run("DeleteBinary success", func(t *testing.T) {
		t.Parallel()

		fClient, _, _, _, _, _, _, binaryMock := setupFacadeTest()
		ctx := t.Context()

		binaryMock.On("DeleteFile", ctx, "file-123").
			Return(true, nil).Once()

		deleted, err := fClient.DeleteBinary(ctx, "file-123")
		require.NoError(t, err)
		assert.True(t, deleted)
		binaryMock.AssertExpectations(t)
	})
}
