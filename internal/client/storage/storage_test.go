//nolint:err113,exhaustruct
package storage_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/npavlov/go-password-manager/gen/proto/card"
	"github.com/npavlov/go-password-manager/gen/proto/file"
	"github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/gen/proto/note"
	"github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/client/storage"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

// Mocks

func TestNewStorageManager(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	f := new(testutils.MockFacade)
	tm := new(testutils.MockTokenManager)

	sm := storage.NewStorageManager(f, tm, &logger)

	assert.NotNil(t, sm)
	assert.Empty(t, sm.Password)
	assert.Empty(t, sm.Notes)
	assert.Empty(t, sm.Cards)
	assert.Empty(t, sm.Binaries)
}

func TestFetchItems_Success(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	// Mock two pages of results
	itemsPage1 := []*item.ItemData{
		{Id: "1", Type: item.ItemType_ITEM_TYPE_PASSWORD, UpdatedAt: timestamppb.Now()},
		{Id: "2", Type: item.ItemType_ITEM_TYPE_NOTE, UpdatedAt: timestamppb.Now()},
	}
	facade := new(testutils.MockFacade)
	facade.GetItemsFunc = func(_ context.Context, _, _ int32) ([]*item.ItemData, int32, error) {
		return itemsPage1, 0, nil
	}
	tm := new(testutils.MockTokenManager)

	facade.On("GetItems", mock.Anything, int32(1), int32(10)).Return(itemsPage1, int32(2), nil)

	sm := storage.NewStorageManager(facade, tm, &logger)
	items, err := sm.FetchItems(t.Context())
	require.NoError(t, err)
	assert.Len(t, items, 2)
	facade.AssertExpectations(t)
}

func TestFetchItems_Error(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	facade := new(testutils.MockFacade)
	tm := new(testutils.MockTokenManager)

	facade.GetItemsFunc = func(_ context.Context, _, _ int32) ([]*item.ItemData, int32, error) {
		return []*item.ItemData{}, 0, errors.New("error")
	}
	facade.On("GetItems", mock.Anything, int32(1), int32(10)).
		Return([]*item.ItemData{}, int32(0), errors.New("fetch error"))

	sm := storage.NewStorageManager(facade, tm, &logger)
	_, err := sm.FetchItems(t.Context())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error getting items")
	facade.AssertExpectations(t)
}

func TestProcessItem_NotUpdated(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	f := new(testutils.MockFacade)
	tm := new(testutils.MockTokenManager)

	sm := storage.NewStorageManager(f, tm, &logger)
	sm.LastSyncAt = time.Now().Add(time.Hour) // Set future time

	itemData := &item.ItemData{
		Id:        "1",
		Type:      item.ItemType_ITEM_TYPE_PASSWORD,
		UpdatedAt: timestamppb.Now(),
	}

	processed := sm.ProcessItem(t.Context(), itemData)
	assert.False(t, processed)
}

func TestProcessPassword_Success(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	facade := new(testutils.MockFacade)
	passwordData := &password.PasswordData{
		Login:    "user",
		Password: "pass",
	}
	facade.GetPasswordFunc = func(_ context.Context, _ string) (*password.PasswordData, time.Time, error) {
		return passwordData, time.Time{}, nil
	}
	passwordID := "pass123"
	tm := new(testutils.MockTokenManager)

	meta := map[string]string{"key": "value"}
	lastUpdate := time.Now()

	facade.On("GetPassword", mock.Anything, passwordID).Return(passwordData, lastUpdate, nil)

	sm := storage.NewStorageManager(facade, tm, &logger)
	err := sm.ProcessPassword(t.Context(), passwordID, meta)

	require.NoError(t, err)
	assert.Contains(t, sm.Password, passwordID)
	assert.Equal(t, "user", sm.Password[passwordID].Login)
	assert.Equal(t, "pass", sm.Password[passwordID].Password)
	facade.AssertExpectations(t)
}

func TestProcessNote_Success(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	fClient := new(testutils.MockFacade)
	tm := new(testutils.MockTokenManager)

	noteID := "note123"
	meta := map[string]string{"key": "value"}
	noteData := &note.NoteData{
		Content: "test content",
	}
	lastUpdate := time.Now()
	fClient.GetNoteFunc = func(_ context.Context, _ string) (*note.NoteData, time.Time, error) {
		return noteData, time.Time{}, nil
	}
	fClient.On("GetNote", mock.Anything, noteID).Return(noteData, lastUpdate, nil)

	sm := storage.NewStorageManager(fClient, tm, &logger)
	err := sm.ProcessNote(t.Context(), noteID, meta)

	require.NoError(t, err)
	assert.Contains(t, sm.Notes, noteID)
	assert.Equal(t, "test content", sm.Notes[noteID].Content)
	fClient.AssertExpectations(t)
}

func TestProcessCard_Success(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	facade := new(testutils.MockFacade)
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
	facade.GetCardFunc = func(_ context.Context, _ string) (*card.CardData, time.Time, error) {
		return cardData, time.Time{}, nil
	}
	facade.On("GetCard", mock.Anything, cardID).Return(cardData, lastUpdate, nil)

	sm := storage.NewStorageManager(facade, tm, &logger)
	err := sm.ProcessCard(t.Context(), cardID, meta)

	require.NoError(t, err)
	assert.Contains(t, sm.Cards, cardID)
	assert.Equal(t, "4111111111111111", sm.Cards[cardID].CardNumber)
	facade.AssertExpectations(t)
}

func TestProcessBinary_Success(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	facade := new(testutils.MockFacade)
	tm := new(testutils.MockTokenManager)

	fileID := "file123"
	meta := map[string]string{"key": "value"}
	fileMeta := &file.FileMeta{
		FileName: "test.txt",
		FileSize: 1024,
	}

	facade.GetFileFunc = func(_ context.Context, _ string) (*file.FileMeta, error) {
		return fileMeta, nil
	}
	facade.On("GetFile", mock.Anything, fileID).Return(fileMeta, nil)

	sm := storage.NewStorageManager(facade, tm, &logger)
	err := sm.ProcessBinary(t.Context(), fileID, meta)

	require.NoError(t, err)
	assert.Contains(t, sm.Binaries, fileID)
	assert.Equal(t, "test.txt", sm.Binaries[fileID].Filename)
	facade.AssertExpectations(t)
}

func TestSyncItems_NotAuthorized(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	f := new(testutils.MockFacade)
	tm := new(testutils.MockTokenManager)
	tm.On("IsAuthorized").Return(false)

	sm := storage.NewStorageManager(f, tm, &logger)
	err := sm.SyncItems(t.Context())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not authorized")
}

func TestSyncItems_AlreadySyncing(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	f := new(testutils.MockFacade)
	tm := new(testutils.MockTokenManager)

	tm.On("IsAuthorized").Return(true)

	sm := storage.NewStorageManager(f, tm, &logger)
	sm.Syncing = 1 // Simulate ongoing sync

	err := sm.SyncItems(t.Context())
	require.NoError(t, err) // Should return nil when already syncing
}

func TestStartBackgroundSync(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	facade := new(testutils.MockFacade)
	tm := new(testutils.MockTokenManager)

	// Setup mocks for initial sync
	tm.On("IsAuthorized").Return(true)
	facade.On("GetItems", mock.Anything, int32(1), int32(10)).Return([]*item.ItemData{}, int32(0), nil)

	sm := storage.NewStorageManager(facade, tm, &logger)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	// Start sync in separate goroutine
	go sm.StartBackgroundSync(ctx)

	// Wait a bit to ensure sync starts
	time.Sleep(100 * time.Millisecond)

	// Stop the sync
	sm.StopSync()
}

func TestProcessItem_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		itemType item.ItemType
		mockFunc func(facade *testutils.MockFacade)
	}{
		{
			name:     "Password",
			itemType: item.ItemType_ITEM_TYPE_PASSWORD,
			mockFunc: func(facade *testutils.MockFacade) {
				facade.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
				facade.On("GetPassword", mock.Anything, "item1").Return(&password.PasswordData{}, time.Now(), nil)
			},
		},
		{
			name:     "Note",
			itemType: item.ItemType_ITEM_TYPE_NOTE,
			mockFunc: func(facade *testutils.MockFacade) {
				facade.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
				facade.On("GetNote", mock.Anything, "item1").Return(&note.NoteData{}, time.Now(), nil)
			},
		},
		{
			name:     "Card",
			itemType: item.ItemType_ITEM_TYPE_CARD,
			mockFunc: func(facade *testutils.MockFacade) {
				facade.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
				facade.On("GetCard", mock.Anything, "item1").Return(&card.CardData{}, time.Now(), nil)
			},
		},
		{
			name:     "Binary",
			itemType: item.ItemType_ITEM_TYPE_BINARY,
			mockFunc: func(facade *testutils.MockFacade) {
				facade.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
				facade.On("GetFile", mock.Anything, "item1").Return(&file.FileMeta{}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logger := zerolog.Nop()
			facade := new(testutils.MockFacade)
			facade.GetMetainfoFunc = func(_ context.Context, _ string) (map[string]string, error) {
				return map[string]string{}, nil
			}
			facade.GetNoteFunc = func(_ context.Context, _ string) (*note.NoteData, time.Time, error) {
				return &note.NoteData{}, time.Time{}, nil
			}
			facade.GetPasswordFunc = func(_ context.Context, _ string) (*password.PasswordData, time.Time, error) {
				return &password.PasswordData{}, time.Time{}, nil
			}
			facade.GetCardFunc = func(_ context.Context, _ string) (*card.CardData, time.Time, error) {
				return &card.CardData{}, time.Time{}, nil
			}
			facade.GetFileFunc = func(_ context.Context, _ string) (*file.FileMeta, error) {
				return &file.FileMeta{}, nil
			}
			tm := new(testutils.MockTokenManager)

			sm := storage.NewStorageManager(facade, tm, &logger)
			sm.LastSyncAt = time.Now().Add(-time.Hour) // Set past time

			tc.mockFunc(facade)
			item := &item.ItemData{
				Id:        "item1",
				Type:      tc.itemType,
				UpdatedAt: timestamppb.Now(),
			}

			processed := sm.ProcessItem(t.Context(), item)
			assert.True(t, processed)
			facade.AssertExpectations(t)
		})
	}
}

func TestProcessItem_MetaError(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	fClient := new(testutils.MockFacade)
	tm := new(testutils.MockTokenManager)

	sm := storage.NewStorageManager(fClient, tm, &logger)
	sm.LastSyncAt = time.Now().Add(-time.Hour)

	item := &item.ItemData{
		Id:        "item1",
		Type:      item.ItemType_ITEM_TYPE_PASSWORD,
		UpdatedAt: timestamppb.Now(),
	}

	fClient.GetMetainfoFunc = func(_ context.Context, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	}
	fClient.On("GetMetainfo", mock.Anything, "item1").Return(nil, errors.New("meta error"))

	processed := sm.ProcessItem(t.Context(), item)
	assert.False(t, processed)
	fClient.AssertExpectations(t)
}

func TestProcessItem_ProcessingError(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	facade := new(testutils.MockFacade)
	tm := new(testutils.MockTokenManager)

	sm := storage.NewStorageManager(facade, tm, &logger)
	sm.LastSyncAt = time.Now().Add(-time.Hour)

	newItem := &item.ItemData{
		Id:        "item1",
		Type:      item.ItemType_ITEM_TYPE_PASSWORD,
		UpdatedAt: timestamppb.Now(),
	}

	facade.GetMetainfoFunc = func(_ context.Context, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	}
	facade.GetPasswordFunc = func(_ context.Context, _ string) (*password.PasswordData, time.Time, error) {
		return &password.PasswordData{}, time.Time{}, errors.New("password error")
	}

	facade.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
	facade.On("GetPassword", mock.Anything, "item1").Return(nil, time.Time{}, errors.New("password error"))

	processed := sm.ProcessItem(t.Context(), newItem)
	assert.False(t, processed)
	facade.AssertExpectations(t)
}

func TestSyncItems_Logging(t *testing.T) {
	t.Parallel()

	logger := testutils.GetTLogger()
	facade := new(testutils.MockFacade)
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
	facade.On("GetItems", mock.Anything, int32(1), int32(10)).Return(items, int32(2), 2)
	facade.On("GetMetainfo", mock.Anything, "item1").Return(map[string]string{}, nil)
	facade.On("GetPassword", mock.Anything, "item1").Return(&password.PasswordData{}, time.Now(), nil)
	facade.On("GetMetainfo", mock.Anything, "item2").Return(map[string]string{}, nil)
	facade.On("GetNote", mock.Anything, "item2").Return(&note.NoteData{}, time.Now(), nil)
	facade.GetItemsFunc = func(_ context.Context, _, _ int32) ([]*item.ItemData, int32, error) {
		return items, 2, nil
	}

	sm := storage.NewStorageManager(facade, tm, logger)
	err := sm.SyncItems(t.Context())

	require.NoError(t, err)
}
