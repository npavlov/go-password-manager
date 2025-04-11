package facade

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb_card "github.com/npavlov/go-password-manager/gen/proto/card"
	pb_file "github.com/npavlov/go-password-manager/gen/proto/file"
	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	pb_note "github.com/npavlov/go-password-manager/gen/proto/note"
	pb_password "github.com/npavlov/go-password-manager/gen/proto/password"
	tokenMgr "github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/npavlov/go-password-manager/internal/client/grpc/auth"
	binary "github.com/npavlov/go-password-manager/internal/client/grpc/binaries"
	"github.com/npavlov/go-password-manager/internal/client/grpc/cards"
	"github.com/npavlov/go-password-manager/internal/client/grpc/items"
	"github.com/npavlov/go-password-manager/internal/client/grpc/metainfo"
	"github.com/npavlov/go-password-manager/internal/client/grpc/notes"
	"github.com/npavlov/go-password-manager/internal/client/grpc/passwords"
)

type Facade struct {
	authClient     *auth.Client
	itemsClient    *items.Client
	passwordClient *passwords.Client
	metaClient     *metainfo.Client
	noteClient     *notes.Client
	cardsClient    *cards.Client
	binariesClient *binary.Client
}

func NewFacade(conn *grpc.ClientConn, tokenManager *tokenMgr.TokenManager, log *zerolog.Logger) *Facade {
	return &Facade{
		authClient:     auth.NewAuthClient(conn, tokenManager, log),
		itemsClient:    items.NewItemsClient(conn, tokenManager, log),
		passwordClient: passwords.NewPasswordClient(conn, tokenManager, log),
		metaClient:     metainfo.NewMetainfoClient(conn, tokenManager, log),
		noteClient:     notes.NewNoteClient(conn, tokenManager, log),
		cardsClient:    cards.NewCardClient(conn, tokenManager, log),
		binariesClient: binary.NewBinaryClient(conn, tokenManager, log),
	}
}

func (fa *Facade) Login(username, password string) error {
	return fa.authClient.Login(username, password)
}

func (fa *Facade) Register(username, password, email string) (string, error) {
	masterKey, err := fa.authClient.Register(username, password, email)

	return masterKey, errors.WithMessagef(err, "failed to register user %s", username)
}

func (fa *Facade) GetItems(ctx context.Context, page, pageSize int32) ([]*pb.ItemData, int32, error) {
	getItems, i, err := fa.itemsClient.GetItems(ctx, page, pageSize)

	return getItems, i, errors.Wrap(err, "error getting items")
}

func (fa *Facade) StorePassword(ctx context.Context, login string, password string) (string, error) {
	passwordId, err := fa.passwordClient.StorePassword(ctx, login, password)

	return passwordId, errors.Wrap(err, "error storing password")
}

func (fa *Facade) GetPassword(ctx context.Context, id string) (*pb_password.PasswordData, time.Time, error) {
	password, lastUpdate, err := fa.passwordClient.GetPassword(ctx, id)

	return password, lastUpdate, errors.Wrap(err, "error getting password")
}

func (fa *Facade) UpdatePassword(ctx context.Context, id, login, password string) error {
	err := fa.passwordClient.UpdatePassword(ctx, id, login, password)

	return errors.Wrap(err, "error updating password")
}

func (fa *Facade) DeletePassword(ctx context.Context, id string) (bool, error) {
	result, err := fa.passwordClient.DeletePassword(ctx, id)

	return result, errors.Wrap(err, "error deleting password")
}

func (fa *Facade) GetMetainfo(ctx context.Context, id string) (map[string]string, error) {
	meta, err := fa.metaClient.GetMetainfo(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "error getting metainfo")
	}

	if meta == nil {
		return make(map[string]string), nil
	}

	return meta, nil
}

func (fa *Facade) SetMetainfo(ctx context.Context, id string, meta map[string]string) (bool, error) {
	result, err := fa.metaClient.SetMetainfo(ctx, id, meta)

	return result, errors.Wrap(err, "error setting metainfo")
}

func (fa *Facade) DeleteMetainfo(ctx context.Context, id, key string) (bool, error) {
	result, err := fa.metaClient.DeleteMetainfo(ctx, id, key)

	return result, errors.Wrap(err, "error deleting metainfo")
}

func (fa *Facade) StoreNote(ctx context.Context, content string) (string, error) {
	noteId, err := fa.noteClient.StoreNote(ctx, content)

	return noteId, errors.Wrap(err, "error storing note")
}

func (fa *Facade) GetNote(ctx context.Context, id string) (*pb_note.NoteData, time.Time, error) {
	note, lastUpdate, err := fa.noteClient.GetNote(ctx, id)

	return note, lastUpdate, errors.Wrap(err, "error getting note")
}

func (fa *Facade) DeleteNote(ctx context.Context, id string) (bool, error) {
	result, err := fa.noteClient.DeleteNote(ctx, id)

	return result, errors.Wrap(err, "error deleting note")
}

func (fa *Facade) StoreCard(ctx context.Context, cardNum, expDate, Cvv, cardHolder string) (string, error) {
	cardId, err := fa.cardsClient.StoreCard(ctx, cardNum, expDate, Cvv, cardHolder)

	return cardId, errors.Wrap(err, "error storing card")
}

func (fa *Facade) UpdateCard(ctx context.Context, id, cardNum, expDate, Cvv, cardHolder string) error {
	err := fa.cardsClient.UpdateCard(ctx, id, cardNum, expDate, Cvv, cardHolder)

	return errors.Wrap(err, "error updating card")
}

func (fa *Facade) GetCard(ctx context.Context, id string) (*pb_card.CardData, time.Time, error) {
	card, lastUpdate, err := fa.cardsClient.GetCard(ctx, id)

	return card, lastUpdate, errors.Wrap(err, "error getting note")
}

func (fa *Facade) DeleteCard(ctx context.Context, id string) (bool, error) {
	result, err := fa.cardsClient.DeleteCard(ctx, id)

	return result, errors.Wrap(err, "error deleting card")
}

// UploadBinary streams and stores a binary file.
func (fa *Facade) UploadBinary(ctx context.Context, filename string, reader io.Reader) (string, error) {
	fileID, err := fa.binariesClient.UploadFile(ctx, filename, reader)

	return fileID, errors.Wrap(err, "failed to upload binary")
}

// DownloadBinary retrieves and writes a binary file.
func (fa *Facade) DownloadBinary(ctx context.Context, fileID string, writer io.Writer) error {
	err := fa.binariesClient.DownloadFile(ctx, fileID, writer)

	return errors.Wrap(err, "failed to download binary")
}

// GetFile retrieves and writes a binary file.
func (fa *Facade) GetFile(ctx context.Context, fileID string) (*pb_file.FileMeta, error) {
	file, err := fa.binariesClient.GetFile(ctx, fileID)

	return file, errors.Wrap(err, "error getting file")
}

// DeleteBinary removes a binary file by ID.
func (fa *Facade) DeleteBinary(ctx context.Context, fileID string) (bool, error) {
	ok, err := fa.binariesClient.DeleteFile(ctx, fileID)

	return ok, errors.Wrap(err, "failed to delete binary")
}
