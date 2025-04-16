package storage

import (
	"context"

	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/internal/client/model"
)

//nolint:interfacebloat
type IStorageManager interface {
	FetchItems(ctx context.Context) ([]*pb.ItemData, error)
	ProcessItem(ctx context.Context, item *pb.ItemData) bool
	ProcessPassword(ctx context.Context, passwordID string, meta map[string]string) error
	ProcessNote(ctx context.Context, noteID string, meta map[string]string) error
	ProcessCard(ctx context.Context, cardID string, meta map[string]string) error
	ProcessBinary(ctx context.Context, fileID string, meta map[string]string) error
	StartBackgroundSync(ctx context.Context)
	SyncItems(ctx context.Context) error
	StopSync()
	GetBinaries() map[string]model.BinaryItem
	GetCards() map[string]model.CardItem
	GetPasswords() map[string]model.PasswordItem
	GetNotes() map[string]model.NoteItem
	DeleteBinary(ID string)
	DeleteCards(ID string)
	DeleteNotes(ID string)
	DeletePassword(ID string)
}
