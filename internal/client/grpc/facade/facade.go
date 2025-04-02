package facade

import (
	"context"
	"time"

	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	pb_password "github.com/npavlov/go-password-manager/gen/proto/password"
	tokenMgr "github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/npavlov/go-password-manager/internal/client/grpc/auth"
	"github.com/npavlov/go-password-manager/internal/client/grpc/items"
	"github.com/npavlov/go-password-manager/internal/client/grpc/metainfo"
	"github.com/npavlov/go-password-manager/internal/client/grpc/passwords"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Facade struct {
	authClient     *auth.Client
	itemsClient    *items.Client
	passwordClient *passwords.Client
	metaClient     *metainfo.Client
}

func NewFacade(conn *grpc.ClientConn, tokenManager *tokenMgr.TokenManager, log *zerolog.Logger) *Facade {
	return &Facade{
		authClient:     auth.NewAuthClient(conn, tokenManager, log),
		itemsClient:    items.NewItemsClient(conn, tokenManager, log),
		passwordClient: passwords.NewItemsClient(conn, tokenManager, log),
		metaClient:     metainfo.NewMetainfoClient(conn, tokenManager, log),
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
	storePassword, err := fa.passwordClient.StorePassword(ctx, login, password)

	return storePassword, errors.Wrap(err, "error storing password")
}

func (fa *Facade) GetPassword(ctx context.Context, id string) (*pb_password.PasswordData, time.Time, error) {
	password, lastUpdate, err := fa.passwordClient.GetPassword(ctx, id)

	return password, lastUpdate, errors.Wrap(err, "error getting password")
}

func (fa *Facade) UpdatePassword(ctx context.Context, id, login, password string) error {
	err := fa.passwordClient.UpdatePassword(ctx, id, login, password)

	return errors.Wrap(err, "error updating password")
}

func (fa *Facade) GetMetainfo(ctx context.Context, id string) (map[string]string, error) {
	meta, err := fa.metaClient.GetMetainfo(ctx, id)

	return meta, errors.Wrap(err, "error getting metainfo")
}

func (fa *Facade) SetMetainfo(ctx context.Context, id string, meta map[string]string) (bool, error) {
	result, err := fa.metaClient.SetMetainfo(ctx, id, meta)

	return result, errors.Wrap(err, "error setting metainfo")
}

func (fa *Facade) DeleteMetainfo(ctx context.Context, id, key string) (bool, error) {
	result, err := fa.metaClient.DeleteMetainfo(ctx, id, key)

	return result, errors.Wrap(err, "error deleting metainfo")
}
