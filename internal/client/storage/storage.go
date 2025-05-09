//nolint:exhaustruct
package storage

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/npavlov/go-password-manager/internal/client/grpc/facade"
	"github.com/npavlov/go-password-manager/internal/client/model"
)

const (
	syncTime        = 5 * time.Minute
	pageLimit int32 = 10
)

// StManager manages client-sIDe storage and background syncing.
type StManager struct {
	facade     facade.IFacade
	Password   map[string]model.PasswordItem `json:"passwords"`
	Notes      map[string]model.NoteItem     `json:"notes"`
	Cards      map[string]model.CardItem     `json:"cards"`
	Binaries   map[string]model.BinaryItem   `json:"binaries"`
	LastSyncAt time.Time
	mutex      sync.Mutex
	stopChan   chan struct{}
	logger     *zerolog.Logger
	tokenMgr   auth.ITokenManager
	Syncing    int32
}

// NewStorageManager creates a new StorageManager with background sync.
func NewStorageManager(facade facade.IFacade, tokenMgr auth.ITokenManager, logger *zerolog.Logger) *StManager {
	//nolint:exhaustruct
	sm := &StManager{
		facade:   facade,
		stopChan: make(chan struct{}),
		logger:   logger,
		Password: make(map[string]model.PasswordItem),
		Notes:    make(map[string]model.NoteItem),
		Cards:    make(map[string]model.CardItem),
		Binaries: make(map[string]model.BinaryItem),
		tokenMgr: tokenMgr,
	}

	return sm
}

// FetchItems retrieves all items from the server.
func (sm *StManager) FetchItems(ctx context.Context) ([]*pb.ItemData, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	total := 0
	var itr int32 = 1
	allItems := make([]*pb.ItemData, 0)
	for {
		items, count, err := sm.facade.GetItems(ctx, itr, pageLimit)
		if err != nil {
			sm.logger.Error().Err(err).Msg("error getting items")

			return nil, errors.Wrap(err, "error getting items")
		}

		allItems = append(allItems, items...)

		total += int(count)

		// next page
		itr++

		if count < pageLimit {
			break
		}
	}

	log.Println("Initial fetch complete:", total, "items retrieved.")

	return allItems, nil
}

func (sm *StManager) ProcessItem(ctx context.Context, item *pb.ItemData) bool {
	if item.GetUpdatedAt().AsTime().Before(sm.LastSyncAt) {
		return false
	}

	itemID := item.GetId()

	meta, err := sm.facade.GetMetainfo(ctx, itemID)
	if err != nil {
		sm.logger.Error().Err(err).Msg("error getting metainfo")

		return false
	}

	switch item.GetType() {
	case pb.ItemType_ITEM_TYPE_PASSWORD:
		err = sm.ProcessPassword(ctx, itemID, meta)
	case pb.ItemType_ITEM_TYPE_NOTE:
		err = sm.ProcessNote(ctx, itemID, meta)
	case pb.ItemType_ITEM_TYPE_CARD:
		err = sm.ProcessCard(ctx, itemID, meta)
	case pb.ItemType_ITEM_TYPE_BINARY:
		err = sm.ProcessBinary(ctx, itemID, meta)
	case pb.ItemType_ITEM_TYPE_UNSPECIFIED:
		return false
	}

	return err == nil
}

func (sm *StManager) ProcessPassword(ctx context.Context, passwordID string, meta map[string]string) error {
	password, lastUpdate, err := sm.facade.GetPassword(ctx, passwordID)
	if err != nil {
		sm.logger.Error().Err(err).Msg("error getting password")

		return errors.Wrap(err, "error getting password")
	}

	sm.Password[passwordID] = model.PasswordItem{
		Login:    password.GetLogin(),
		Password: password.GetPassword(),
		StorageItem: model.StorageItem{
			Type:      model.ItemTypePassword,
			ID:        passwordID,
			UpdatedAt: lastUpdate,
			Metadata:  meta,
		},
	}

	return nil
}

func (sm *StManager) ProcessNote(ctx context.Context, noteID string, meta map[string]string) error {
	note, lastUpdate, err := sm.facade.GetNote(ctx, noteID)
	if err != nil {
		sm.logger.Error().Err(err).Msg("error getting password")

		return errors.Wrap(err, "error getting password")
	}

	sm.Notes[noteID] = model.NoteItem{
		Content: note.GetContent(),
		StorageItem: model.StorageItem{
			Type:      model.ItemTypeNote,
			ID:        noteID,
			UpdatedAt: lastUpdate,
			Metadata:  meta,
		},
	}

	return nil
}

func (sm *StManager) ProcessCard(ctx context.Context, cardID string, meta map[string]string) error {
	card, lastUpdate, err := sm.facade.GetCard(ctx, cardID)
	if err != nil {
		sm.logger.Error().Err(err).Msg("error getting card")

		return errors.Wrap(err, "error getting card")
	}

	sm.Cards[cardID] = model.CardItem{
		CardNumber:     card.GetCardNumber(),
		CVV:            card.GetCvv(),
		ExpiryDate:     card.GetExpiryDate(),
		CardholderName: card.GetCardholderName(),
		StorageItem: model.StorageItem{
			Type:      model.ItemTypeCard,
			ID:        cardID,
			UpdatedAt: lastUpdate,
			Metadata:  meta,
		},
	}

	return nil
}

// ProcessBinary retrieves binary metadata and downloads the file.
func (sm *StManager) ProcessBinary(ctx context.Context, fileID string, meta map[string]string) error {
	file, err := sm.facade.GetFile(ctx, fileID)
	if err != nil {
		sm.logger.Error().Err(err).Str("fileID", fileID).Msg("error getting binary metadata")

		return errors.Wrap(err, "error getting binary metadata")
	}

	sm.Binaries[fileID] = model.BinaryItem{
		Filename: file.GetFileName(),
		Size:     file.GetFileSize(),
		StorageItem: model.StorageItem{
			Type:      model.ItemTypeBinary,
			ID:        fileID,
			UpdatedAt: time.Now(),
			Metadata:  meta,
		},
	}

	sm.logger.Info().Str("filename", file.GetFileName()).Msg("binary downloaded and stored in temp folder")

	return nil
}

// StartBackgroundSync runs a goroutine that checks for updates every minute.
func (sm *StManager) StartBackgroundSync(ctx context.Context) {
	ticker := time.NewTicker(syncTime)
	defer ticker.Stop()

	sm.logger.Info().Msg("starting background sync")
	err := sm.SyncItems(ctx)
	if err != nil {
		sm.logger.Error().Err(err).Msg("error syncing items")
	}

	for {
		select {
		case <-ticker.C:
			err := sm.SyncItems(ctx)
			if err != nil {
				sm.logger.Error().Err(err).Msg("error syncing items")
			}
		case <-sm.stopChan:
			sm.logger.Info().Msg("stopping background sync")

			return
		}
	}
}

// SyncItems checks for updates since the last sync.
func (sm *StManager) SyncItems(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&sm.Syncing, 0, 1) {
		sm.logger.Info().Msg("Sync already in progress, skipping this round.")

		return nil
	}
	defer atomic.StoreInt32(&sm.Syncing, 0)

	if !sm.tokenMgr.IsAuthorized() {
		return errors.New("not authorized")
	}

	sm.logger.Info().Msg("syncing items...")

	items, err := sm.FetchItems(ctx)
	if err != nil {
		return errors.Wrap(err, "error fetching items")
	}

	// Update only items that have changed
	updatedCount := 0
	for i, item := range items {
		sm.logger.Info().Int("Item num", i).Msg("got Processing item")

		ok := sm.ProcessItem(ctx, item)

		if ok {
			updatedCount++
		}
	}

	if updatedCount > 0 {
		sm.LastSyncAt = time.Now()
		log.Println("Sync complete:", updatedCount, "items updated.")
	}

	return nil
}

// StopSync stops the background sync goroutine.
func (sm *StManager) StopSync() {
	close(sm.stopChan)
}

func (sm *StManager) GetBinaries() map[string]model.BinaryItem {
	return sm.Binaries
}

func (sm *StManager) GetPasswords() map[string]model.PasswordItem {
	return sm.Password
}

func (sm *StManager) GetNotes() map[string]model.NoteItem {
	return sm.Notes
}

func (sm *StManager) GetCards() map[string]model.CardItem {
	return sm.Cards
}

func (sm *StManager) DeleteBinary(id string) {
	delete(sm.Binaries, id)
}

func (sm *StManager) DeleteCards(id string) {
	delete(sm.Cards, id)
}

func (sm *StManager) DeleteNotes(id string) {
	delete(sm.Notes, id)
}

func (sm *StManager) DeletePassword(id string) {
	delete(sm.Password, id)
}
