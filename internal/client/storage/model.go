package storage

import (
	"context"

	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/internal/client/model"
)

type IStorageManager interface {
	FetchItems(ctx context.Context) ([]*pb.ItemData, error)
	ProcessItem(ctx context.Context, item *pb.ItemData) bool
	ProcessPassword(ctx context.Context, passwordId string, meta map[string]string) error
	ProcessNote(ctx context.Context, noteId string, meta map[string]string) error
	ProcessCard(ctx context.Context, cardId string, meta map[string]string) error
	ProcessBinary(ctx context.Context, fileID string, meta map[string]string) error
	StartBackgroundSync(ctx context.Context)
	SyncItems(ctx context.Context) error
	StopSync()
	GetBinaries() map[string]model.BinaryItem
	GetCards() map[string]model.CardItem
	GetPasswords() map[string]model.PasswordItem
	GetNotes() map[string]model.NoteItem
	DeleteBinary(Id string)
	DeleteCards(Id string)
	DeleteNotes(Id string)
	DeletePassword(Id string)
}
