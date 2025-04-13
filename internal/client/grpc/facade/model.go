package facade

import (
	"context"
	"io"
	"time"

	pb_card "github.com/npavlov/go-password-manager/gen/proto/card"
	pb_file "github.com/npavlov/go-password-manager/gen/proto/file"
	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	pb_note "github.com/npavlov/go-password-manager/gen/proto/note"
	pb_password "github.com/npavlov/go-password-manager/gen/proto/password"
)

// IFacade defines the interface for the Facade package
type IFacade interface {
	Login(username, password string) error
	Register(username, password, email string) (string, error)
	GetItems(ctx context.Context, page, pageSize int32) ([]*pb.ItemData, int32, error)
	StorePassword(ctx context.Context, login string, password string) (string, error)
	GetPassword(ctx context.Context, id string) (*pb_password.PasswordData, time.Time, error)
	UpdatePassword(ctx context.Context, id, login, password string) error
	DeletePassword(ctx context.Context, id string) (bool, error)
	GetMetainfo(ctx context.Context, id string) (map[string]string, error)
	SetMetainfo(ctx context.Context, id string, meta map[string]string) (bool, error)
	DeleteMetainfo(ctx context.Context, id, key string) (bool, error)
	StoreNote(ctx context.Context, content string) (string, error)
	GetNote(ctx context.Context, id string) (*pb_note.NoteData, time.Time, error)
	DeleteNote(ctx context.Context, id string) (bool, error)
	StoreCard(ctx context.Context, cardNum, expDate, Cvv, cardHolder string) (string, error)
	UpdateCard(ctx context.Context, id, cardNum, expDate, Cvv, cardHolder string) error
	GetCard(ctx context.Context, id string) (*pb_card.CardData, time.Time, error)
	DeleteCard(ctx context.Context, id string) (bool, error)
	UploadBinary(ctx context.Context, filename string, reader io.Reader) (string, error)
	DownloadBinary(ctx context.Context, fileID string, writer io.Writer) error
	GetFile(ctx context.Context, fileID string) (*pb_file.FileMeta, error)
	DeleteBinary(ctx context.Context, fileID string) (bool, error)
}
