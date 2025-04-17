//nolint:exhaustruct
package storage_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/db"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func TestStoreNote(t *testing.T) {
	t.Parallel()
	noteID := uuid.New()
	userID := uuid.New()
	userUUID := pgtype.UUID{Bytes: userID, Valid: true}

	tests := []struct {
		name       string
		createNote db.CreateNoteEntryParams
		mock       func(mock pgxmock.PgxPoolIface)
		want       *db.Note
		wantErr    bool
	}{
		{
			name: "successful note creation",
			createNote: db.CreateNoteEntryParams{
				UserID:           userUUID,
				EncryptedContent: "encrypted_content",
			},
			mock: func(mock pgxmock.PgxPoolIface) {
				now := time.Now()
				rows := pgxmock.NewRows([]string{"id", "user_id", "encrypted_content", "created_at", "updated_at"}).
					AddRow(noteID.String(), userID.String(), "encrypted_content", now, now)
				mock.ExpectQuery("INSERT INTO notes").
					WithArgs(userUUID, "encrypted_content").
					WillReturnRows(rows)
			},
			want: &db.Note{
				EncryptedContent: "encrypted_content",
			},
			wantErr: false,
		},
		{
			name: "database error",
			createNote: db.CreateNoteEntryParams{
				UserID:           userUUID,
				EncryptedContent: "encrypted_content",
			},
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("INSERT INTO notes").
					WithArgs(userUUID, "encrypted_content").
					WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			result, err := storage.StoreNote(t.Context(), tt.createNote)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to store note")
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.EncryptedContent, result.EncryptedContent)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetNote(t *testing.T) {
	t.Parallel()

	noteID := uuid.New()
	noteUUID := pgtype.UUID{Bytes: noteID, Valid: true}
	userID := uuid.New()
	userUUID := pgtype.UUID{Bytes: userID, Valid: true}

	tests := []struct {
		name    string
		noteID  string
		userID  pgtype.UUID
		mock    func(mock pgxmock.PgxPoolIface)
		want    *db.Note
		wantErr bool
	}{
		{
			name:   "successful note retrieval",
			noteID: noteID.String(),
			userID: userUUID,
			mock: func(mock pgxmock.PgxPoolIface) {
				now := time.Now()
				rows := pgxmock.NewRows([]string{"id", "user_id", "encrypted_content", "created_at", "updated_at"}).
					AddRow(noteID.String(), userID.String(), "encrypted_content", now, now)
				mock.ExpectQuery("SELECT").
					WithArgs(noteUUID, userUUID).
					WillReturnRows(rows)
			},
			want: &db.Note{
				EncryptedContent: "encrypted_content",
			},
			wantErr: false,
		},
		{
			name:   "note not found",
			noteID: noteID.String(),
			userID: userUUID,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT").
					WithArgs(noteUUID, userUUID).
					WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			result, err := storage.GetNote(t.Context(), tt.noteID, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to create user")
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.EncryptedContent, result.EncryptedContent)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetNotes(t *testing.T) {
	t.Parallel()
	noteID := uuid.New()
	userID := uuid.New()
	userUUID := pgtype.UUID{Bytes: userID, Valid: true}

	tests := []struct {
		name    string
		userID  string
		mock    func(mock pgxmock.PgxPoolIface)
		want    []db.Note
		wantErr bool
	}{
		{
			name:   "successful notes retrieval",
			userID: userID.String(),
			mock: func(mock pgxmock.PgxPoolIface) {
				now := time.Now()
				rows := pgxmock.NewRows([]string{"id", "user_id", "encrypted_content", "created_at", "updated_at"}).
					AddRow(noteID.String(), userID.String(), "note1", now, now).
					AddRow(noteID.String(), userID.String(), "note2", now.Add(-time.Hour), now.Add(-time.Hour))
				mock.ExpectQuery("SELECT id, user_id, encrypted_content, created_at, updated_at FROM notes").
					WithArgs(userUUID).
					WillReturnRows(rows)
			},
			want: []db.Note{
				{EncryptedContent: "note1"},
				{EncryptedContent: "note2"},
			},
			wantErr: false,
		},
		{
			name:   "no notes found",
			userID: userID.String(),
			mock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "user_id", "encrypted_content", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, user_id, encrypted_content, created_at, updated_at FROM notes").
					WithArgs(userUUID).
					WillReturnRows(rows)
			},
			want:    []db.Note{},
			wantErr: false,
		},
		{
			name:   "database error",
			userID: userID.String(),
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT id, user_id, encrypted_content, created_at, updated_at FROM notes").
					WithArgs(userUUID).
					WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			result, err := storage.GetNotes(t.Context(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to create note")
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tt.want), len(result))
				for i, note := range tt.want {
					require.Equal(t, note.EncryptedContent, result[i].EncryptedContent)
				}
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteNote(t *testing.T) {
	t.Parallel()

	noteID := uuid.New()
	noteUUID := pgtype.UUID{Bytes: noteID, Valid: true}
	userID := uuid.New()
	userUUID := pgtype.UUID{Bytes: userID, Valid: true}

	tests := []struct {
		name    string
		noteID  string
		userID  pgtype.UUID
		mock    func(mock pgxmock.PgxPoolIface)
		wantErr bool
	}{
		{
			name:   "successful note deletion",
			noteID: noteID.String(),
			userID: userUUID,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM notes").
					WithArgs(noteUUID, userUUID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: false,
		},
		{
			name:   "database error",
			noteID: noteID.String(),
			userID: userUUID,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM notes").
					WithArgs(noteUUID, userUUID).
					WillReturnError(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			err := storage.DeleteNote(t.Context(), tt.noteID, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to delete note")
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
