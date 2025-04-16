package testutils

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/internal/client/model"
	"github.com/npavlov/go-password-manager/internal/client/storage"
)

// MockStorageManager is a mock implementation of IStorageManager for testing.
type MockStorageManager struct {
	mock.Mock
	FetchItemsFunc          func(ctx context.Context) ([]*pb.ItemData, error)
	ProcessItemFunc         func(ctx context.Context, item *pb.ItemData) bool
	ProcessPasswordFunc     func(ctx context.Context, passwordId string, meta map[string]string) error
	ProcessNoteFunc         func(ctx context.Context, noteId string, meta map[string]string) error
	ProcessCardFunc         func(ctx context.Context, cardId string, meta map[string]string) error
	ProcessBinaryFunc       func(ctx context.Context, fileID string, meta map[string]string) error
	StartBackgroundSyncFunc func(ctx context.Context)
	SyncItemsFunc           func(ctx context.Context) error
	StopSyncFunc            func()
	GetBinariesFunc         func() map[string]model.BinaryItem
	GetCardsFunc            func() map[string]model.CardItem
	GetPasswordsFunc        func() map[string]model.PasswordItem
	GetNotesFunc            func() map[string]model.NoteItem
	DeleteBinaryFunc        func(Id string)
	DeleteCardsFunc         func(Id string)
	DeleteNotesFunc         func(Id string)
	DeletePasswordFunc      func(Id string)

	// Internal state for testing
	Binaries  map[string]model.BinaryItem
	Passwords map[string]model.PasswordItem
	Notes     map[string]model.NoteItem
	Cards     map[string]model.CardItem
	mutex     sync.Mutex
}

// Verify MockStorageManager implements IStorageManager.
var _ storage.IStorageManager = (*MockStorageManager)(nil)

// NewMockStorageManager creates a new mock storage manager with initialized maps.
func NewMockStorageManager() *MockStorageManager {
	//nolint:exhaustruct
	return &MockStorageManager{
		Binaries:  make(map[string]model.BinaryItem),
		Passwords: make(map[string]model.PasswordItem),
		Notes:     make(map[string]model.NoteItem),
		Cards:     make(map[string]model.CardItem),
	}
}

func (m *MockStorageManager) FetchItems(ctx context.Context) ([]*pb.ItemData, error) {
	if m.FetchItemsFunc != nil {
		return m.FetchItemsFunc(ctx)
	}

	return nil, errors.New("FetchItemsFunc not implemented")
}

func (m *MockStorageManager) ProcessItem(ctx context.Context, item *pb.ItemData) bool {
	if m.ProcessItemFunc != nil {
		return m.ProcessItemFunc(ctx, item)
	}

	return false
}

func (m *MockStorageManager) ProcessPassword(ctx context.Context, passwordId string, meta map[string]string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.ProcessPasswordFunc != nil {
		return m.ProcessPasswordFunc(ctx, passwordId, meta)
	}

	return errors.New("ProcessPasswordFunc not implemented")
}

func (m *MockStorageManager) ProcessNote(ctx context.Context, noteId string, meta map[string]string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.ProcessNoteFunc != nil {
		return m.ProcessNoteFunc(ctx, noteId, meta)
	}

	return errors.New("ProcessNoteFunc not implemented")
}

func (m *MockStorageManager) ProcessCard(ctx context.Context, cardId string, meta map[string]string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.ProcessCardFunc != nil {
		return m.ProcessCardFunc(ctx, cardId, meta)
	}

	return errors.New("ProcessCardFunc not implemented")
}

func (m *MockStorageManager) ProcessBinary(ctx context.Context, fileID string, meta map[string]string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.ProcessBinaryFunc != nil {
		return m.ProcessBinaryFunc(ctx, fileID, meta)
	}

	return errors.New("ProcessBinaryFunc not implemented")
}

func (m *MockStorageManager) StartBackgroundSync(ctx context.Context) {
	if m.StartBackgroundSyncFunc != nil {
		m.Called()
		m.StartBackgroundSyncFunc(ctx)
	}
}

func (m *MockStorageManager) SyncItems(ctx context.Context) error {
	if m.SyncItemsFunc != nil {
		m.Called(ctx)

		return m.SyncItemsFunc(ctx)
	}

	return errors.New("SyncItemsFunc not implemented")
}

func (m *MockStorageManager) StopSync() {
	if m.StopSyncFunc != nil {
		m.StopSyncFunc()
	}
}

func (m *MockStorageManager) GetBinaries() map[string]model.BinaryItem {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.GetBinariesFunc != nil {
		return m.GetBinariesFunc()
	}

	return m.Binaries
}

func (m *MockStorageManager) GetPasswords() map[string]model.PasswordItem {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.GetPasswordsFunc != nil {
		return m.GetPasswordsFunc()
	}

	return m.Passwords
}

func (m *MockStorageManager) GetNotes() map[string]model.NoteItem {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.GetNotesFunc != nil {
		return m.GetNotesFunc()
	}

	return m.Notes
}

func (m *MockStorageManager) GetCards() map[string]model.CardItem {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.GetCardsFunc != nil {
		return m.GetCardsFunc()
	}

	return m.Cards
}

func (m *MockStorageManager) DeleteBinary(Id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.DeleteBinaryFunc != nil {
		m.DeleteBinaryFunc(Id)

		return
	}
	delete(m.Binaries, Id)
}

func (m *MockStorageManager) DeleteCards(Id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.DeleteCardsFunc != nil {
		m.DeleteCardsFunc(Id)

		return
	}
	delete(m.Cards, Id)
}

func (m *MockStorageManager) DeleteNotes(Id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.DeleteNotesFunc != nil {
		m.DeleteNotesFunc(Id)

		return
	}
	delete(m.Notes, Id)
}

func (m *MockStorageManager) DeletePassword(Id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.DeletePasswordFunc != nil {
		m.DeletePasswordFunc(Id)

		return
	}
	delete(m.Passwords, Id)
}
