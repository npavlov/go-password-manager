package testutils

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/npavlov/go-password-manager/internal/server/db"
)

// MockDBStorage implements a mock version of DBStorage with map-based storage.
type MockDBStorage struct {
	mu            sync.RWMutex
	usersByID     map[pgtype.UUID]db.User
	usersByName   map[string]db.User
	tokens        map[string]db.GetRefreshTokenRow
	cards         map[string]db.Card
	binaries      map[string]db.BinaryEntry
	items         map[string]db.Item
	metaInfo      map[string]map[string]string
	log           *zerolog.Logger
	RegisterError error
	masterKey     string
}

func NewMockDBStorage(logger *zerolog.Logger, masterKey string) *MockDBStorage {

	return &MockDBStorage{
		usersByID:   make(map[pgtype.UUID]db.User),
		usersByName: make(map[string]db.User),
		tokens:      make(map[string]db.GetRefreshTokenRow),
		cards:       make(map[string]db.Card),
		binaries:    make(map[string]db.BinaryEntry),
		items:       make(map[string]db.Item),
		metaInfo:    make(map[string]map[string]string),
		log:         logger,
		masterKey:   masterKey,
	}
}

// AddTestUser adds a user to the mock storage for testing.
func (m *MockDBStorage) AddTestUser(user db.User) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.usersByID[user.ID] = user
	m.usersByName[user.Username] = user
}

// RegisterUser mock implementation.
func (m *MockDBStorage) RegisterUser(_ context.Context, createUser db.CreateUserParams) (*db.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.RegisterError != nil {
		return nil, m.RegisterError
	}

	// Check if username already exists
	if _, exists := m.usersByName[createUser.Username]; exists {
		return nil, errors.New("username already exists")
	}

	user := db.User{
		ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Username:      createUser.Username,
		Email:         createUser.Email,
		Password:      createUser.Password,
		EncryptionKey: createUser.EncryptionKey,
	}

	// Add to both maps
	m.usersByID[user.ID] = user
	m.usersByName[user.Username] = user

	return &user, nil
}

// GetUser mock implementation.
func (m *MockDBStorage) GetUser(ctx context.Context, username string) (*db.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.usersByName[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// GetUserById mock implementation.
func (m *MockDBStorage) GetUserById(ctx context.Context, userId pgtype.UUID) (*db.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !userId.Valid {
		return nil, errors.New("invalid user ID")
	}

	user, exists := m.usersByID[userId]
	if !exists {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// SetupMockUserStorage is a helper function to configure mock storage with test data.
func SetupMockUserStorage(masterKey string, initialUsers ...db.User) *MockDBStorage {
	logger := GetTLogger()
	mockStorage := NewMockDBStorage(logger, masterKey)

	for _, user := range initialUsers {
		mockStorage.AddTestUser(user)
	}

	return mockStorage
}

func (m *MockDBStorage) StoreToken(_ context.Context, userID pgtype.UUID, refreshToken string, expiresAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !userID.Valid {
		return errors.New("invalid user ID")
	}

	pgExpiresAt := pgtype.Timestamp{}
	if err := pgExpiresAt.Scan(expiresAt); err != nil {
		return errors.Wrap(err, "failed to scan expires at")
	}

	m.tokens[refreshToken] = db.GetRefreshTokenRow{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: pgExpiresAt,
	}

	return nil
}

func (m *MockDBStorage) GetToken(_ context.Context, token string) (db.GetRefreshTokenRow, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	row, exists := m.tokens[token]
	if !exists {
		return db.GetRefreshTokenRow{}, errors.New("token not found")
	}

	return row, nil
}

func (m *MockDBStorage) StoreCard(ctx context.Context, createCard db.StoreCardParams) (*db.Card, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	card := db.Card{
		ID:                  pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID:              createCard.UserID,
		EncryptedCardNumber: createCard.EncryptedCardNumber,
		HashedCardNumber:    createCard.HashedCardNumber,
		EncryptedCvv:        createCard.EncryptedCvv,
		EncryptedExpiryDate: createCard.EncryptedExpiryDate,
		CardholderName:      createCard.CardholderName,
		UpdatedAt:           pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	// Store the card in a map (initialize it first if needed)
	if m.cards == nil {
		m.cards = make(map[string]db.Card)
	}
	m.cards[card.ID.String()] = card

	return &card, nil
}

func (m *MockDBStorage) UpdateCard(ctx context.Context, updateCard db.UpdateCardParams) (*db.Card, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cards == nil {
		return nil, errors.New("no cards stored")
	}

	id := updateCard.ID.String()
	card, exists := m.cards[id]
	if !exists {
		return nil, errors.New("card not found")
	}

	card.EncryptedCardNumber = updateCard.EncryptedCardNumber
	card.HashedCardNumber = updateCard.HashedCardNumber
	card.EncryptedCvv = updateCard.EncryptedCvv
	card.EncryptedExpiryDate = updateCard.EncryptedExpiryDate
	card.CardholderName = updateCard.CardholderName
	card.UpdatedAt = pgtype.Timestamp{Time: time.Now(), Valid: true}

	m.cards[id] = card
	return &card, nil
}

func (m *MockDBStorage) GetCard(ctx context.Context, cardId string, userId pgtype.UUID) (*db.Card, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.cards == nil {
		return nil, errors.New("no cards stored")
	}

	card, exists := m.cards[cardId]
	if !exists {
		return nil, errors.New("card not found")
	}

	if card.UserID != userId {
		return nil, errors.New("unauthorized access to card")
	}

	return &card, nil
}

func (m *MockDBStorage) GetCards(ctx context.Context, userId string) ([]db.Card, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []db.Card
	for _, card := range m.cards {
		if card.UserID.Valid && card.UserID.String() == userId {
			result = append(result, card)
		}
	}

	return result, nil
}

func (m *MockDBStorage) DeleteCard(_ context.Context, cardId string, _ pgtype.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.cards, cardId)

	return nil
}

// StoreBinary stores a binary entry in the mock storage
func (m *MockDBStorage) StoreBinary(ctx context.Context, createBinary db.StoreBinaryEntryParams) (*db.BinaryEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !createBinary.UserID.Valid {
		return nil, errors.New("invalid user ID")
	}

	binary := db.BinaryEntry{
		ID:       pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID:   createBinary.UserID,
		FileName: createBinary.FileName,
		FileSize: createBinary.FileSize,
		FileUrl:  createBinary.FileUrl,
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}

	m.binaries[binary.ID.String()] = binary
	return &binary, nil
}

// DeleteBinary removes a binary entry from the mock storage
func (m *MockDBStorage) DeleteBinary(ctx context.Context, arg db.DeleteBinaryEntryParams) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !arg.ID.Valid {
		return errors.New("invalid binary ID")
	}

	if !arg.UserID.Valid {
		return errors.New("invalid user ID")
	}

	binary, exists := m.binaries[arg.ID.String()]
	if !exists {
		return errors.New("binary not found")
	}

	if binary.UserID != arg.UserID {
		return errors.New("unauthorized access to binary")
	}

	delete(m.binaries, arg.ID.String())
	return nil
}

// GetBinaries returns all binaries for a user
func (m *MockDBStorage) GetBinaries(ctx context.Context, userId string) ([]db.BinaryEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []db.BinaryEntry
	for _, binary := range m.binaries {
		if binary.UserID.Valid && binary.UserID.String() == userId {
			result = append(result, binary)
		}
	}

	return result, nil
}

// GetBinary retrieves a specific binary entry
func (m *MockDBStorage) GetBinary(ctx context.Context, binaryId string, userId pgtype.UUID) (*db.BinaryEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	binary, exists := m.binaries[binaryId]
	if !exists {
		return nil, errors.New("binary not found")
	}

	if binary.UserID != userId {
		return nil, errors.New("unauthorized access to binary")
	}

	return &binary, nil
}

// GetItems Add these new methods to MockDBStorage
func (m *MockDBStorage) GetItems(ctx context.Context, params db.GetItemsByUserIDParams) ([]db.GetItemsByUserIDRow, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []db.GetItemsByUserIDRow
	for _, item := range m.items {
		if item.UserID.Valid && item.UserID.String() == params.UserID.String() {
			result = append(result, db.GetItemsByUserIDRow{
				IDResource: item.ID,
				Type:       item.Type,
				CreatedAt:  item.CreatedAt,
			})
		}
	}

	// Apply pagination
	start := int(params.Offset)
	end := start + int(params.Limit)
	if end > len(result) {
		end = len(result)
	}
	if start > len(result) {
		return []db.GetItemsByUserIDRow{}, nil
	}

	return result[start:end], nil
}

func (m *MockDBStorage) GetItem(ctx context.Context, itemID string, userID pgtype.UUID) (*db.Item, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[itemID]
	if !exists {
		return nil, errors.New("item not found")
	}

	if item.UserID != userID {
		return nil, errors.New("unauthorized access to item")
	}

	return &item, nil
}

func (m *MockDBStorage) StoreItem(ctx context.Context, userId pgtype.UUID, itemType db.ItemType) (*db.Item, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item := db.Item{
		ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID:    userId,
		Type:      itemType,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	m.items[item.ID.String()] = item
	return &item, nil
}

func (m *MockDBStorage) DeleteItem(ctx context.Context, itemID string, userID pgtype.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.items[itemID]
	if !exists {
		return errors.New("item not found")
	}

	if item.UserID != userID {
		return errors.New("unauthorized access to item")
	}

	delete(m.items, itemID)
	return nil
}

// AddMeta Add these new methods to MockDBStorage
func (m *MockDBStorage) AddMeta(ctx context.Context, itemID string, key string, value string) (*db.Metainfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.metaInfo[itemID]; !exists {
		m.metaInfo[itemID] = make(map[string]string)
	}

	m.metaInfo[itemID][key] = value
	return &db.Metainfo{
		ID:     pgtype.UUID{Bytes: uuid.New(), Valid: true},
		ItemID: pgtype.UUID{Bytes: uuid.MustParse(itemID), Valid: true},
		Key:    key,
		Value:  value,
	}, nil
}

func (m *MockDBStorage) DeleteMetaInfo(ctx context.Context, key string, itemID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.metaInfo[itemID]; !exists {
		return errors.New("item not found")
	}

	if _, exists := m.metaInfo[itemID][key]; !exists {
		return errors.New("metadata key not found")
	}

	delete(m.metaInfo[itemID], key)
	return nil
}

func (m *MockDBStorage) GetMetaInfo(ctx context.Context, itemID string) ([]db.GetMetaInfoByItemIDRow, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	itemMeta, exists := m.metaInfo[itemID]
	if !exists {
		return []db.GetMetaInfoByItemIDRow{}, nil
	}

	var result []db.GetMetaInfoByItemIDRow
	for key, value := range itemMeta {
		result = append(result, db.GetMetaInfoByItemIDRow{
			Key:   key,
			Value: value,
		})
	}

	return result, nil
}

func (m *MockDBStorage) ClearTestData() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.usersByID = make(map[pgtype.UUID]db.User)
	m.usersByName = make(map[string]db.User)
	m.tokens = make(map[string]db.GetRefreshTokenRow)
	m.cards = make(map[string]db.Card)
	m.binaries = make(map[string]db.BinaryEntry)
	m.items = make(map[string]db.Item)
	m.metaInfo = make(map[string]map[string]string)

	m.RegisterError = nil
}
