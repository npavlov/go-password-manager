//nolint:err113,wrapcheck,forcetypeassert
package notes_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/npavlov/go-password-manager/gen/proto/note"
	"github.com/npavlov/go-password-manager/internal/client/grpc/notes"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

// Mocks

type MockNoteServiceClient struct {
	mock.Mock
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockNoteServiceClient) GetNoteV1(ctx context.Context,
	in *note.GetNoteV1Request,
	_ ...grpc.CallOption,
) (*note.GetNoteV1Response, error) {
	args := m.Called(ctx, in)

	// Safely handle nil to avoid type assertion panic
	arg, ok := args.Get(0).(*note.GetNoteV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockNoteServiceClient) GetNotesV1(ctx context.Context,
	in *note.GetNotesV1Request,
	_ ...grpc.CallOption,
) (*note.GetNotesV1Response, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*note.GetNotesV1Response), args.Error(1)
}

func (m *MockNoteServiceClient) StoreNoteV1(ctx context.Context,
	in *note.StoreNoteV1Request,
	_ ...grpc.CallOption,
) (*note.StoreNoteV1Response, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*note.StoreNoteV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockNoteServiceClient) DeleteNoteV1(ctx context.Context,
	in *note.DeleteNoteV1Request,
	_ ...grpc.CallOption,
) (*note.DeleteNoteV1Response, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*note.DeleteNoteV1Response)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockTokenManager) GetToken() (string, error) {
	args := m.Called()

	return args.String(0), args.Error(1)
}

func TestGetNote_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	expectedNote := &note.NoteData{
		Content: "This is a test note",
	}
	expectedTime := time.Now()

	mockClient.On("GetNoteV1", mock.Anything, &note.GetNoteV1Request{
		NoteId: "note123",
	}).Return(&note.GetNoteV1Response{
		Note:       expectedNote,
		LastUpdate: timestamppb.New(expectedTime),
	}, nil)

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	noteData, _, err := client.GetNote(t.Context(), "note123")
	require.NoError(t, err)
	assert.Equal(t, expectedNote, noteData)
}

func TestGetNote_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	mockClient.On("GetNoteV1", mock.Anything, &note.GetNoteV1Request{
		NoteId: "note123",
	}).Return(nil, errors.New("get note failed"))

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, _, err := client.GetNote(t.Context(), "note123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error getting password")
}

func TestStoreNote_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	content := "This is a test note content"

	mockClient.On("StoreNoteV1", mock.Anything, &note.StoreNoteV1Request{
		Note: &note.NoteData{
			Content: content,
		},
	}).Return(&note.StoreNoteV1Response{
		NoteId: "new-note-123",
	}, nil)

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	noteID, err := client.StoreNote(t.Context(), content)
	require.NoError(t, err)
	assert.Equal(t, "new-note-123", noteID)
}

func TestStoreNote_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	content := "This is a test note content"

	mockClient.On("StoreNoteV1", mock.Anything, &note.StoreNoteV1Request{
		Note: &note.NoteData{
			Content: content,
		},
	}).Return(nil, errors.New("store note failed"))

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.StoreNote(t.Context(), content)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error storing note")
}

func TestDeleteNote_Success(t *testing.T) {
	t.Parallel()

	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteNoteV1", mock.Anything, &note.DeleteNoteV1Request{
		NoteId: "note123",
	}).Return(&note.DeleteNoteV1Response{
		Ok: true,
	}, nil)

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeleteNote(t.Context(), "note123")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestDeleteNote_Error(t *testing.T) {
	t.Parallel()

	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteNoteV1", mock.Anything, &note.DeleteNoteV1Request{
		NoteId: "note123",
	}).Return(nil, errors.New("delete note failed"))

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.DeleteNote(t.Context(), "note123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting note")
}

func TestDeleteNote_NotSuccessful(t *testing.T) {
	t.Parallel()

	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteNoteV1", mock.Anything, &note.DeleteNoteV1Request{
		NoteId: "note123",
	}).Return(&note.DeleteNoteV1Response{
		Ok: false,
	}, nil)

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeleteNote(t.Context(), "note123")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestNewNoteClient(t *testing.T) {
	t.Parallel()

	tm := new(testutils.MockTokenManager)
	logger := zerolog.Nop()

	client := notes.NewNoteClient(nil, tm, &logger)

	assert.NotNil(t, client)
}
