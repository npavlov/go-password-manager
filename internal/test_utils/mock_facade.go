//nolint:forcetypeassert,wrapcheck
package testutils

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	pb_card "github.com/npavlov/go-password-manager/gen/proto/card"
	pb_file "github.com/npavlov/go-password-manager/gen/proto/file"
	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	pb_note "github.com/npavlov/go-password-manager/gen/proto/note"
	pb_password "github.com/npavlov/go-password-manager/gen/proto/password"
	"github.com/npavlov/go-password-manager/internal/client/grpc/facade"
)

// MockFacade is a mock implementation of IFacade for testing.
type MockFacade struct {
	mock.Mock
	LoginFunc          func(username, password string) error
	RegisterFunc       func(username, password, email string) (string, error)
	GetItemsFunc       func(ctx context.Context, page, pageSize int32) ([]*pb.ItemData, int32, error)
	StorePasswordFunc  func(ctx context.Context, login string, password string) (string, error)
	GetPasswordFunc    func(ctx context.Context, id string) (*pb_password.PasswordData, time.Time, error)
	UpdatePasswordFunc func(ctx context.Context, id, login, password string) error
	DeletePasswordFunc func(ctx context.Context, id string) (bool, error)
	GetMetainfoFunc    func(ctx context.Context, id string) (map[string]string, error)
	SetMetainfoFunc    func(ctx context.Context, id string, meta map[string]string) (bool, error)
	DeleteMetainfoFunc func(ctx context.Context, id, key string) (bool, error)
	StoreNoteFunc      func(ctx context.Context, content string) (string, error)
	GetNoteFunc        func(ctx context.Context, id string) (*pb_note.NoteData, time.Time, error)
	DeleteNoteFunc     func(ctx context.Context, id string) (bool, error)
	StoreCardFunc      func(ctx context.Context, cardNum, expDate, Cvv, cardHolder string) (string, error)
	UpdateCardFunc     func(ctx context.Context, id, cardNum, expDate, Cvv, cardHolder string) error
	GetCardFunc        func(ctx context.Context, id string) (*pb_card.CardData, time.Time, error)
	DeleteCardFunc     func(ctx context.Context, id string) (bool, error)
	UploadBinaryFunc   func(ctx context.Context, filename string, reader io.Reader) (string, error)
	DownloadBinaryFunc func(ctx context.Context, fileID string, writer io.Writer) error
	GetFileFunc        func(ctx context.Context, fileID string) (*pb_file.FileMeta, error)
	DeleteBinaryFunc   func(ctx context.Context, fileID string) (bool, error)
}

// Verify MockFacade implements IFacade.
var _ facade.IFacade = (*MockFacade)(nil)

func (m *MockFacade) Login(username, password string) error {
	if m.LoginFunc != nil {
		m.Called(username, password)

		return m.LoginFunc(username, password)
	}

	return errors.New("LoginFunc not implemented")
}

func (m *MockFacade) Register(username, password, email string) (string, error) {
	if m.RegisterFunc != nil {
		args := m.Called(username, password, email)

		return args.Get(0).(string), args.Error(1)
	}

	return "", errors.New("RegisterFunc not implemented")
}

func (m *MockFacade) GetItems(ctx context.Context, page, pageSize int32) ([]*pb.ItemData, int32, error) {
	if m.GetItemsFunc != nil {
		m.Called(ctx, page, pageSize)

		return m.GetItemsFunc(ctx, page, pageSize)
	}

	return nil, 0, errors.New("GetItemsFunc not implemented")
}

func (m *MockFacade) StorePassword(ctx context.Context, login string, password string) (string, error) {
	if m.StorePasswordFunc != nil {
		return m.StorePasswordFunc(ctx, login, password)
	}

	return "", errors.New("StorePasswordFunc not implemented")
}

func (m *MockFacade) GetPassword(ctx context.Context, id string) (*pb_password.PasswordData, time.Time, error) {
	if m.GetPasswordFunc != nil {
		m.Called(ctx, id)

		return m.GetPasswordFunc(ctx, id)
	}

	return nil, time.Time{}, errors.New("GetPasswordFunc not implemented")
}

func (m *MockFacade) UpdatePassword(ctx context.Context, id, login, password string) error {
	if m.UpdatePasswordFunc != nil {
		return m.UpdatePasswordFunc(ctx, id, login, password)
	}

	return errors.New("UpdatePasswordFunc not implemented")
}

func (m *MockFacade) DeletePassword(ctx context.Context, id string) (bool, error) {
	if m.DeletePasswordFunc != nil {
		return m.DeletePasswordFunc(ctx, id)
	}

	return false, errors.New("DeletePasswordFunc not implemented")
}

func (m *MockFacade) GetMetainfo(ctx context.Context, id string) (map[string]string, error) {
	if m.GetMetainfoFunc != nil {
		m.Called(ctx, id)

		return m.GetMetainfoFunc(ctx, id)
	}

	return nil, errors.New("GetMetainfoFunc not implemented")
}

func (m *MockFacade) SetMetainfo(ctx context.Context, id string, meta map[string]string) (bool, error) {
	if m.SetMetainfoFunc != nil {
		return m.SetMetainfoFunc(ctx, id, meta)
	}

	return false, errors.New("SetMetainfoFunc not implemented")
}

func (m *MockFacade) DeleteMetainfo(ctx context.Context, id, key string) (bool, error) {
	if m.DeleteMetainfoFunc != nil {
		return m.DeleteMetainfoFunc(ctx, id, key)
	}

	return false, errors.New("DeleteMetainfoFunc not implemented")
}

func (m *MockFacade) StoreNote(ctx context.Context, content string) (string, error) {
	if m.StoreNoteFunc != nil {
		return m.StoreNoteFunc(ctx, content)
	}

	return "", errors.New("StoreNoteFunc not implemented")
}

func (m *MockFacade) GetNote(ctx context.Context, id string) (*pb_note.NoteData, time.Time, error) {
	if m.GetNoteFunc != nil {
		m.Called(ctx, id)

		return m.GetNoteFunc(ctx, id)
	}

	return nil, time.Time{}, errors.New("GetNoteFunc not implemented")
}

func (m *MockFacade) DeleteNote(ctx context.Context, id string) (bool, error) {
	if m.DeleteNoteFunc != nil {
		return m.DeleteNoteFunc(ctx, id)
	}

	return false, errors.New("DeleteNoteFunc not implemented")
}

func (m *MockFacade) StoreCard(ctx context.Context, cardNum, expDate, cvv, cardHolder string) (string, error) {
	if m.StoreCardFunc != nil {
		return m.StoreCardFunc(ctx, cardNum, expDate, cvv, cardHolder)
	}

	return "", errors.New("StoreCardFunc not implemented")
}

func (m *MockFacade) UpdateCard(ctx context.Context, id, cardNum, expDate, cvv, cardHolder string) error {
	if m.UpdateCardFunc != nil {
		return m.UpdateCardFunc(ctx, id, cardNum, expDate, cvv, cardHolder)
	}

	return errors.New("UpdateCardFunc not implemented")
}

func (m *MockFacade) GetCard(ctx context.Context, id string) (*pb_card.CardData, time.Time, error) {
	if m.GetCardFunc != nil {
		m.Called(ctx, id)

		return m.GetCardFunc(ctx, id)
	}

	return nil, time.Time{}, errors.New("GetCardFunc not implemented")
}

func (m *MockFacade) DeleteCard(ctx context.Context, id string) (bool, error) {
	if m.DeleteCardFunc != nil {
		return m.DeleteCardFunc(ctx, id)
	}

	return false, errors.New("DeleteCardFunc not implemented")
}

func (m *MockFacade) UploadBinary(ctx context.Context, filename string, reader io.Reader) (string, error) {
	if m.UploadBinaryFunc != nil {
		return m.UploadBinaryFunc(ctx, filename, reader)
	}

	return "", errors.New("UploadBinaryFunc not implemented")
}

func (m *MockFacade) DownloadBinary(ctx context.Context, fileID string, writer io.Writer) error {
	if m.DownloadBinaryFunc != nil {
		return m.DownloadBinaryFunc(ctx, fileID, writer)
	}

	return errors.New("DownloadBinaryFunc not implemented")
}

func (m *MockFacade) GetFile(ctx context.Context, fileID string) (*pb_file.FileMeta, error) {
	if m.GetFileFunc != nil {
		m.Called(ctx, fileID)

		return m.GetFileFunc(ctx, fileID)
	}

	return nil, errors.New("GetFileFunc not implemented")
}

func (m *MockFacade) DeleteBinary(ctx context.Context, fileID string) (bool, error) {
	if m.DeleteBinaryFunc != nil {
		return m.DeleteBinaryFunc(ctx, fileID)
	}

	return false, errors.New("DeleteBinaryFunc not implemented")
}
