//nolint:lll,dogsled,exhaustruct
package note_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/npavlov/go-password-manager/gen/proto/note"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/note"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	generalutils "github.com/npavlov/go-password-manager/internal/utils"
)

func setupNoteService(t *testing.T) (*note.Service, *testutils.MockDBStorage, context.Context, string) {
	t.Helper()

	logger := zerolog.New(nil)
	masterKey, _ := utils.GenerateRandomKey()
	storage := testutils.SetupMockUserStorage(masterKey)
	cfg := &config.Config{
		SecuredMasterKey: generalutils.NewString(masterKey),
	}

	// Create test user
	userID := uuid.New()
	encryptionKey, _ := utils.GenerateRandomKey()
	encryptedKey, _ := utils.Encrypt(encryptionKey, masterKey)

	testUser := db.User{
		ID:            pgtype.UUID{Bytes: userID, Valid: true},
		Username:      "testuser",
		Password:      "hashed-password",
		EncryptionKey: encryptedKey,
	}
	storage.AddTestUser(testUser)

	// Inject user ID and encryption key into context
	ctx := testutils.InjectUserToContext(t.Context(), userID.String())

	return note.NewNoteService(&logger, storage, cfg), storage, ctx, encryptionKey
}

func TestStoreNote_Success(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupNoteService(t)

	testContent := "This is a secret note"

	req := &pb.StoreNoteV1Request{
		Note: &pb.NoteData{
			Content: testContent,
		},
	}

	resp, err := svc.StoreNoteV1(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, resp.GetNoteId())

	// Verify note was stored
	notes, err := storage.GetNotes(ctx, pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true})
	require.NoError(t, err)
	require.Len(t, notes, 1)

	getNote, err := svc.GetNoteV1(ctx, &pb.GetNoteV1Request{NoteId: resp.GetNoteId()})
	require.NoError(t, err)

	require.Equal(t, getNote.GetNote().GetContent(), testContent)
}

func TestStoreNote_InvalidInput(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupNoteService(t)

	req := &pb.StoreNoteV1Request{
		Note: &pb.NoteData{}, // Missing content
	}

	_, err := svc.StoreNoteV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestGetNote_Success(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, userKey := setupNoteService(t)

	// Store a test note first
	testContent := "Secret note content"
	encryptedContent, _ := utils.Encrypt(testContent, userKey)
	note, err := storage.StoreNote(ctx, db.CreateNoteEntryParams{
		UserID:           pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true},
		EncryptedContent: encryptedContent,
	})
	require.NoError(t, err)

	req := &pb.GetNoteV1Request{
		NoteId: note.ID.String(),
	}

	resp, err := svc.GetNoteV1(ctx, req)
	require.NoError(t, err)
	require.Equal(t, testContent, resp.GetNote().GetContent())
	require.True(t, timestamppb.New(note.UpdatedAt.Time).AsTime().Equal(resp.GetLastUpdate().AsTime()))
}

func TestGetNote_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupNoteService(t)

	req := &pb.GetNoteV1Request{
		NoteId: uuid.NewString(),
	}

	_, err := svc.GetNoteV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "note not found")
}

func TestGetNote_Unauthorized(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupNoteService(t)

	// Store a note with different user
	otherUserID := uuid.New()
	otherNote, err := storage.StoreNote(t.Context(), db.CreateNoteEntryParams{
		UserID:           pgtype.UUID{Bytes: otherUserID, Valid: true},
		EncryptedContent: "encrypted-content",
	})
	require.NoError(t, err)

	req := &pb.GetNoteV1Request{
		NoteId: otherNote.ID.String(),
	}

	_, err = svc.GetNoteV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized access to note")
}

func TestDeleteNote_Success(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, userKey := setupNoteService(t)

	// Store a test note first
	testContent := "Note to be deleted"
	encryptedContent, _ := utils.Encrypt(testContent, userKey)
	note, err := storage.StoreNote(ctx, db.CreateNoteEntryParams{
		UserID:           pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true},
		EncryptedContent: encryptedContent,
	})
	require.NoError(t, err)

	req := &pb.DeleteNoteV1Request{
		NoteId: note.ID.String(),
	}

	resp, err := svc.DeleteNoteV1(ctx, req)
	require.NoError(t, err)
	require.True(t, resp.GetOk())

	// Verify note was deleted
	_, err = storage.GetNote(ctx, note.ID.String(), pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true})
	require.Error(t, err)
}

func TestGetNotes_NotImplemented(t *testing.T) {
	t.Parallel()

	svc, _, ctx, _ := setupNoteService(t)

	req := &pb.GetNotesV1Request{}
	resp, err := svc.GetNotesV1(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	// Add more assertions when method is implemented
}

func TestStoreNote_EncryptionFailure(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupNoteService(t)

	// Corrupt the user's encryption key in storage
	userID := uuid.MustParse(testutils.GetUserIDFromContext(ctx))
	userIDPG := pgtype.UUID{Bytes: userID, Valid: true}
	user := storage.UsersByID[userIDPG]

	user.EncryptionKey = "invalid-key"
	storage.UsersByID[userIDPG] = user

	req := &pb.StoreNoteV1Request{
		Note: &pb.NoteData{
			Content: "test content",
		},
	}

	_, err := svc.StoreNoteV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Error decrypting user id")
}

func TestStoreNote_DatabaseFailure(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupNoteService(t)

	// Make storage return error
	storage.CallError = errors.New("database failure")

	req := &pb.StoreNoteV1Request{
		Note: &pb.NoteData{
			Content: "test content",
		},
	}

	_, err := svc.StoreNoteV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to store password")
}

func TestGetNote_DecryptionFailure(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, _ := setupNoteService(t)

	// Store storeNote with invalid encrypted content
	storeNote, err := storage.StoreNote(ctx, db.CreateNoteEntryParams{
		UserID:           pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true},
		EncryptedContent: "invalid-encrypted-content",
	})
	require.NoError(t, err)

	req := &pb.GetNoteV1Request{
		NoteId: storeNote.ID.String(),
	}

	_, err = svc.GetNoteV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error decrypting password")
}

func TestDeleteNote_DatabaseFailure(t *testing.T) {
	t.Parallel()

	svc, storage, ctx, userKey := setupNoteService(t)

	// Store a test note first
	encryptedContent, _ := utils.Encrypt("test content", userKey)
	note, err := storage.StoreNote(ctx, db.CreateNoteEntryParams{
		UserID:           pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true},
		EncryptedContent: encryptedContent,
	})
	require.NoError(t, err)

	// Make storage return error
	storage.CallError = errors.New("database failure")

	req := &pb.DeleteNoteV1Request{
		NoteId: note.ID.String(),
	}

	_, err = svc.DeleteNoteV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error deleting note")
}

func TestStoreNote_MissingUserContext(t *testing.T) {
	t.Parallel()

	svc, _, _, _ := setupNoteService(t)

	req := &pb.StoreNoteV1Request{
		Note: &pb.NoteData{
			Content: "test content",
		},
	}

	_, err := svc.StoreNoteV1(t.Context(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error getting user id")
}

func TestGetNote_MissingUserContext(t *testing.T) {
	t.Parallel()

	svc, storage, _, userKey := setupNoteService(t)

	// Store a test note first with valid context
	ctx := testutils.InjectUserToContext(t.Context(), uuid.NewString())
	encryptedContent, _ := utils.Encrypt("test content", userKey)
	note, err := storage.StoreNote(ctx, db.CreateNoteEntryParams{
		UserID:           pgtype.UUID{Bytes: uuid.MustParse(testutils.GetUserIDFromContext(ctx)), Valid: true},
		EncryptedContent: encryptedContent,
	})
	require.NoError(t, err)

	// Try to get note without user context
	req := &pb.GetNoteV1Request{
		NoteId: note.ID.String(),
	}

	_, err = svc.GetNoteV1(t.Context(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error getting user id")
}
