package storage

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/npavlov/go-password-manager/internal/client/grpc/facade"
	"github.com/npavlov/go-password-manager/internal/client/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	pageLimit int32 = 10
)

// StManager manages client-side storage and background syncing
type StManager struct {
	facade     *facade.Facade
	Password   map[string]model.PasswordItem `json:"passwords"`
	Notes      map[string]model.NoteItem     `json:"notes"`
	lastSyncAt time.Time
	mutex      sync.Mutex
	stopChan   chan struct{}
	logger     *zerolog.Logger
	tokenMgr   *auth.TokenManager
	syncing    int32
}

// NewStorageManager creates a new StorageManager with background sync
func NewStorageManager(facade *facade.Facade, tokenMgr *auth.TokenManager, logger *zerolog.Logger) *StManager {
	sm := &StManager{
		facade:   facade,
		stopChan: make(chan struct{}),
		logger:   logger,
		Password: make(map[string]model.PasswordItem),
		tokenMgr: tokenMgr,
	}
	return sm
}

// FetchItems retrieves all items from the server
func (sm *StManager) FetchItems(ctx context.Context) ([]*pb.ItemData, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	var total int = 0
	var i int32 = 1
	allItems := make([]*pb.ItemData, 0)
	for {
		items, count, err := sm.facade.GetItems(ctx, i, pageLimit)
		if err != nil {
			sm.logger.Error().Err(err).Msg("error getting items")

			return nil, errors.Wrap(err, "error getting items")
		}

		allItems = append(allItems, items...)

		total = total + int(count)

		//next page
		i++

		if count < pageLimit {
			break
		}
	}

	log.Println("Initial fetch complete:", total, "items retrieved.")

	return allItems, nil
}

func (sm *StManager) processItem(ctx context.Context, item *pb.ItemData) bool {
	if item.UpdatedAt.AsTime().Before(sm.lastSyncAt) {
		return false
	}

	itemId := item.GetId()

	meta, err := sm.facade.GetMetainfo(ctx, itemId)
	if err != nil {
		sm.logger.Error().Err(err).Msg("error getting metainfo")
	}

	switch item.Type {
	case pb.ItemType_ITEM_TYPE_PASSWORD:
		err := sm.ProcessPassword(ctx, itemId, meta)
		if err != nil {
			return false
		}
	}

	return true
}

func (sm *StManager) ProcessPassword(ctx context.Context, passwordId string, meta map[string]string) error {
	password, lastUpdate, err := sm.facade.GetPassword(ctx, passwordId)

	if err != nil {
		sm.logger.Error().Err(err).Msg("error getting password")

		return errors.Wrap(err, "error getting password")
	}

	sm.Password[passwordId] = model.PasswordItem{
		Login:    password.Login,
		Password: password.Password,
		StorageItem: model.StorageItem{
			ID:        passwordId,
			UpdatedAt: lastUpdate,
			Metadata:  meta,
		},
	}

	return nil
}

// StartBackgroundSync runs a goroutine that checks for updates every minute
func (sm *StManager) StartBackgroundSync(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
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

// SyncItems checks for updates since the last sync
func (sm *StManager) SyncItems(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&sm.syncing, 0, 1) {
		sm.logger.Info().Msg("Sync already in progress, skipping this round.")
		return nil
	}
	defer atomic.StoreInt32(&sm.syncing, 0)

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

		ok := sm.processItem(ctx, item)

		if ok {
			updatedCount++
		}
	}

	if updatedCount > 0 {
		sm.lastSyncAt = time.Now()
		log.Println("Sync complete:", updatedCount, "items updated.")
	}

	return nil
}

// StopSync stops the background sync goroutine
func (sm *StManager) StopSync() {
	close(sm.stopChan)
}
