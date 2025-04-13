package notes_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/npavlov/go-password-manager/gen/proto/note"
	"github.com/npavlov/go-password-manager/internal/client/grpc/notes"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mocks

type MockNoteServiceClient struct {
	mock.Mock
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockNoteServiceClient) GetNote(ctx context.Context, in *note.GetNoteRequest, opts ...grpc.CallOption) (*note.GetNoteResponse, error) {
	args := m.Called(ctx, in)

	// Safely handle nil to avoid type assertion panic
	arg, ok := args.Get(0).(*note.GetNoteResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockNoteServiceClient) GetNotes(ctx context.Context, in *note.GetNotesRequest, opts ...grpc.CallOption) (*note.GetNotesResponse, error) {
	args := m.Called(ctx, in)

	return args.Get(0).(*note.GetNotesResponse), args.Error(1)
}

func (m *MockNoteServiceClient) StoreNote(ctx context.Context, in *note.StoreNoteRequest, opts ...grpc.CallOption) (*note.StoreNoteResponse, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*note.StoreNoteResponse)
	if !ok && args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return arg, args.Error(1)
}

func (m *MockNoteServiceClient) DeleteNote(ctx context.Context, in *note.DeleteNoteRequest, opts ...grpc.CallOption) (*note.DeleteNoteResponse, error) {
	args := m.Called(ctx, in)

	arg, ok := args.Get(0).(*note.DeleteNoteResponse)
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
	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	expectedNote := &note.NoteData{
		Content: "This is a test note",
	}
	expectedTime := time.Now()

	mockClient.On("GetNote", mock.Anything, &note.GetNoteRequest{
		NoteId: "note123",
	}).Return(&note.GetNoteResponse{
		Note:       expectedNote,
		LastUpdate: timestamppb.New(expectedTime),
	}, nil)

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	noteData, _, err := client.GetNote(context.Background(), "note123")
	assert.NoError(t, err)
	assert.Equal(t, expectedNote, noteData)
}

func TestGetNote_Error(t *testing.T) {
	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	mockClient.On("GetNote", mock.Anything, &note.GetNoteRequest{
		NoteId: "note123",
	}).Return(nil, errors.New("get note failed"))

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, _, err := client.GetNote(context.Background(), "note123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting password")
}

func TestStoreNote_Success(t *testing.T) {
	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	content := "This is a test note content"

	mockClient.On("StoreNote", mock.Anything, &note.StoreNoteRequest{
		Note: &note.NoteData{
			Content: content,
		},
	}).Return(&note.StoreNoteResponse{
		NoteId: "new-note-123",
	}, nil)

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	noteID, err := client.StoreNote(context.Background(), content)
	assert.NoError(t, err)
	assert.Equal(t, "new-note-123", noteID)
}

func TestStoreNote_Error(t *testing.T) {
	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	content := "This is a test note content"

	mockClient.On("StoreNote", mock.Anything, &note.StoreNoteRequest{
		Note: &note.NoteData{
			Content: content,
		},
	}).Return(nil, errors.New("store note failed"))

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.StoreNote(context.Background(), content)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error storing note")
}

func TestDeleteNote_Success(t *testing.T) {
	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteNote", mock.Anything, &note.DeleteNoteRequest{
		NoteId: "note123",
	}).Return(&note.DeleteNoteResponse{
		Ok: true,
	}, nil)

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeleteNote(context.Background(), "note123")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestDeleteNote_Error(t *testing.T) {
	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteNote", mock.Anything, &note.DeleteNoteRequest{
		NoteId: "note123",
	}).Return(nil, errors.New("delete note failed"))

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	_, err := client.DeleteNote(context.Background(), "note123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting note")
}

func TestDeleteNote_NotSuccessful(t *testing.T) {
	mockClient := new(MockNoteServiceClient)
	logger := zerolog.Nop()

	mockClient.On("DeleteNote", mock.Anything, &note.DeleteNoteRequest{
		NoteId: "note123",
	}).Return(&note.DeleteNoteResponse{
		Ok: false,
	}, nil)

	client := &notes.Client{
		Client:       mockClient,
		TokenManager: new(testutils.MockTokenManager),
		Log:          &logger,
	}

	ok, err := client.DeleteNote(context.Background(), "note123")
	assert.NoError(t, err)
	assert.False(t, ok)
}
