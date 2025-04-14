package storage_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/npavlov/go-password-manager/gen/proto/card"
	"github.com/npavlov/go-password-manager/gen/proto/file"
	"github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/gen/proto/note"
	"github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/client/storage"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mocks

type MockFacade struct {
	mock.Mock
}

func (m *MockFacade) Login(username, password string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) Register(username, password, email string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) StorePassword(ctx context.Context, login string, password string) (string, error) {
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

func (m *MockFacade) DeleteBinary(ctx context.Context, fileID string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFacade) GetItems(ctx context.Context, page, pageSize int32) ([]*item.ItemData, int32, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*item.ItemData), args.Get(1).(int32), args.Error(2)
}

func (m *MockFacade) GetMetainfo(ctx context.Context, id string) (map[string]string, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockFacade) GetPassword(ctx context.Context, id string) (*password.PasswordData, time.Time, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, time.Now(), args.Error(2)
	}

	return args.Get(0).(*password.PasswordData), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockFacade) GetNote(ctx context.Context, id string) (*note.NoteData, time.Time, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*note.NoteData), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockFacade) GetCard(ctx context.Context, id string) (*card.CardData, time.Time, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*card.CardData), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockFacade) GetFile(ctx context.Context, id string) (*file.FileMeta, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*file.FileMeta), args.Error(1)
}

func TestNewStorageManager(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	sm := storage.NewStorageManager(f, tm, &logger)

	assert.NotNil(t, sm)
	assert.Empty(t, sm.Password)
	assert.Empty(t, sm.Notes)
	assert.Empty(t, sm.Cards)
	assert.Empty(t, sm.Binaries)
}

func TestFetchItems_Success(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	// Mock two pages of results
	itemsPage1 := []*item.ItemData{
		{Id: "1", Type: item.ItemType_ITEM_TYPE_PASSWORD, UpdatedAt: timestamppb.Now()},
		{Id: "2", Type: item.ItemType_ITEM_TYPE_NOTE, UpdatedAt: timestamppb.Now()},
	}

	f.On("GetItems", mock.Anything, int32(1), int32(10)).Return(itemsPage1, int32(2), nil)

	sm := storage.NewStorageManager(f, tm, &logger)
	items, err := sm.FetchItems(context.Background())
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	f.AssertExpectations(t)
}

func TestFetchItems_Error(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	f.On("GetItems", mock.Anything, int32(1), int32(10)).Return([]*item.ItemData{}, int32(0), errors.New("fetch error"))

	sm := storage.NewStorageManager(f, tm, &logger)
	_, err := sm.FetchItems(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting items")
	f.AssertExpectations(t)
}

func TestProcessItem_NotUpdated(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	sm := storage.NewStorageManager(f, tm, &logger)
	sm.LastSyncAt = time.Now().Add(time.Hour) // Set future time

	item := &item.ItemData{
		Id:        "1",
		Type:      item.ItemType_ITEM_TYPE_PASSWORD,
		UpdatedAt: timestamppb.Now(),
	}

	processed := sm.ProcessItem(context.Background(), item)
	assert.False(t, processed)
}

func TestProcessPassword_Success(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	passwordID := "pass123"
	meta := map[string]string{"key": "value"}
	passwordData := &password.PasswordData{
		Login:    "user",
		Password: "pass",
	}
	lastUpdate := time.Now()

	f.On("GetPassword", mock.Anything, passwordID).Return(passwordData, lastUpdate, nil)

	sm := storage.NewStorageManager(f, tm, &logger)
	err := sm.ProcessPassword(context.Background(), passwordID, meta)

	assert.NoError(t, err)
	assert.Contains(t, sm.Password, passwordID)
	assert.Equal(t, "user", sm.Password[passwordID].Login)
	assert.Equal(t, "pass", sm.Password[passwordID].Password)
	f.AssertExpectations(t)
}

func TestProcessNote_Success(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	noteID := "note123"
	meta := map[string]string{"key": "value"}
	noteData := &note.NoteData{
		Content: "test content",
	}
	lastUpdate := time.Now()

	f.On("GetNote", mock.Anything, noteID).Return(noteData, lastUpdate, nil)

	sm := storage.NewStorageManager(f, tm, &logger)
	err := sm.ProcessNote(context.Background(), noteID, meta)

	assert.NoError(t, err)
	assert.Contains(t, sm.Notes, noteID)
	assert.Equal(t, "test content", sm.Notes[noteID].Content)
	f.AssertExpectations(t)
}

func TestProcessCard_Success(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	cardID := "card123"
	meta := map[string]string{"key": "value"}
	cardData := &card.CardData{
		CardNumber:     "4111111111111111",
		Cvv:            "123",
		ExpiryDate:     "12/25",
		CardholderName: "John Doe",
	}
	lastUpdate := time.Now()

	f.On("GetCard", mock.Anything, cardID).Return(cardData, lastUpdate, nil)

	sm := storage.NewStorageManager(f, tm, &logger)
	err := sm.ProcessCard(context.Background(), cardID, meta)

	assert.NoError(t, err)
	assert.Contains(t, sm.Cards, cardID)
	assert.Equal(t, "4111111111111111", sm.Cards[cardID].CardNumber)
	f.AssertExpectations(t)
}

func TestProcessBinary_Success(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	fileID := "file123"
	meta := map[string]string{"key": "value"}
	fileMeta := &file.FileMeta{
		FileName: "test.txt",
		FileSize: 1024,
	}

	f.On("GetFile", mock.Anything, fileID).Return(fileMeta, nil)

	sm := storage.NewStorageManager(f, tm, &logger)
	err := sm.ProcessBinary(context.Background(), fileID, meta)

	assert.NoError(t, err)
	assert.Contains(t, sm.Binaries, fileID)
	assert.Equal(t, "test.txt", sm.Binaries[fileID].Filename)
	f.AssertExpectations(t)
}

func TestSyncItems_NotAuthorized(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)
	tm.On("IsAuthorized").Return(false)

	sm := storage.NewStorageManager(f, tm, &logger)
	err := sm.SyncItems(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not authorized")
}

func TestSyncItems_AlreadySyncing(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	tm.On("IsAuthorized").Return(true)

	sm := storage.NewStorageManager(f, tm, &logger)
	sm.Syncing = 1 // Simulate ongoing sync

	err := sm.SyncItems(context.Background())
	assert.NoError(t, err) // Should return nil when already syncing
}

func TestStartBackgroundSync(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	// Setup mocks for initial sync
	tm.On("IsAuthorized").Return(true)
	f.On("GetItems", mock.Anything, int32(1), int32(10)).Return([]*item.ItemData{}, int32(0), nil)

	sm := storage.NewStorageManager(f, tm, &logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start sync in separate goroutine
	go sm.StartBackgroundSync(ctx)

	// Wait a bit to ensure sync starts
	time.Sleep(100 * time.Millisecond)

	// Stop the sync
	sm.StopSync()

}

func TestProcessItem_Success(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	sm := storage.NewStorageManager(f, tm, &logger)
	sm.LastSyncAt = time.Now().Add(-time.Hour) // Set past time

	testCases := []struct {
		name     string
		itemType item.ItemType
		mockFunc func()
	}{
		{
			name:     "Password",
			itemType: item.ItemType_ITEM_TYPE_PASSWORD,
			mockFunc: func() {
				f.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
				f.On("GetPassword", mock.Anything, "item1").Return(&password.PasswordData{}, time.Now(), nil)
			},
		},
		{
			name:     "Note",
			itemType: item.ItemType_ITEM_TYPE_NOTE,
			mockFunc: func() {
				f.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
				f.On("GetNote", mock.Anything, "item1").Return(&note.NoteData{}, time.Now(), nil)
			},
		},
		{
			name:     "Card",
			itemType: item.ItemType_ITEM_TYPE_CARD,
			mockFunc: func() {
				f.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
				f.On("GetCard", mock.Anything, "item1").Return(&card.CardData{}, time.Now(), nil)
			},
		},
		{
			name:     "Binary",
			itemType: item.ItemType_ITEM_TYPE_BINARY,
			mockFunc: func() {
				f.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
				f.On("GetFile", mock.Anything, "item1").Return(&file.FileMeta{}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			item := &item.ItemData{
				Id:        "item1",
				Type:      tc.itemType,
				UpdatedAt: timestamppb.Now(),
			}

			processed := sm.ProcessItem(context.Background(), item)
			assert.True(t, processed)
			f.AssertExpectations(t)
		})
	}
}

func TestProcessItem_MetaError(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	sm := storage.NewStorageManager(f, tm, &logger)
	sm.LastSyncAt = time.Now().Add(-time.Hour)

	item := &item.ItemData{
		Id:        "item1",
		Type:      item.ItemType_ITEM_TYPE_PASSWORD,
		UpdatedAt: timestamppb.Now(),
	}

	f.On("GetMetainfo", mock.Anything, "item1").Return(nil, errors.New("meta error"))

	processed := sm.ProcessItem(context.Background(), item)
	assert.False(t, processed)
	f.AssertExpectations(t)
}

func TestProcessItem_ProcessingError(t *testing.T) {
	logger := zerolog.Nop()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)

	sm := storage.NewStorageManager(f, tm, &logger)
	sm.LastSyncAt = time.Now().Add(-time.Hour)

	item := &item.ItemData{
		Id:        "item1",
		Type:      item.ItemType_ITEM_TYPE_PASSWORD,
		UpdatedAt: timestamppb.Now(),
	}

	f.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
	f.On("GetPassword", mock.Anything, "item1").Return(nil, time.Time{}, errors.New("password error"))

	processed := sm.ProcessItem(context.Background(), item)
	assert.False(t, processed)
	f.AssertExpectations(t)
}

func TestSyncItems_Logging(t *testing.T) {
	logger := testutils.GetTLogger()
	f := new(MockFacade)
	tm := new(testutils.MockTokenManager)
	tm.Authorized = true

	// Setup test items
	items := []*item.ItemData{
		{
			Id:        "item1",
			Type:      item.ItemType_ITEM_TYPE_PASSWORD,
			UpdatedAt: timestamppb.Now(),
		},
		{
			Id:        "item2",
			Type:      item.ItemType_ITEM_TYPE_NOTE,
			UpdatedAt: timestamppb.Now(),
		},
	}

	// Setup mocks
	tm.On("IsAuthorized").Return(true)
	f.On("GetItems", mock.Anything, int32(1), int32(10)).Return(items, int32(2), nil)
	f.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
	f.On("GetPassword", mock.Anything, "item1").Return(&password.PasswordData{}, time.Now(), nil)
	f.On("GetMetainfo", mock.Anything, "item2").Return(map[string]string{}, nil)
	f.On("GetNote", mock.Anything, "item2").Return(&note.NoteData{}, time.Now(), nil)

	sm := storage.NewStorageManager(f, tm, logger)
	err := sm.SyncItems(context.Background())

	assert.NoError(t, err)

}
