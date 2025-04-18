package storage_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/db"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func TestStoreBinary(t *testing.T) {
	t.Parallel()

	dbStorage, mock := testutils.SetupDBStorage(t)
	ctx := t.Context()

	userID := uuid.New()
	id := uuid.New().String()
	now := time.Now()

	params := db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: userID, Valid: true},
		FileName: "file.txt",
		FileSize: 1234,
		FileUrl:  "https://example.com/file.txt",
	}

	rows := pgxmock.NewRows([]string{
		"id", "user_id", "file_name", "file_size", "file_url", "created_at", "updated_at",
	}).AddRow(id, userID.String(), params.FileName, params.FileSize, params.FileUrl, now, now)

	mock.ExpectQuery(`INSERT INTO binary_entries`).
		WithArgs(params.UserID, params.FileName, params.FileUrl, params.FileSize).
		WillReturnRows(rows)

	entry, err := dbStorage.StoreBinary(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, params.FileName, entry.FileName)
	assert.Equal(t, params.FileUrl, entry.FileUrl)
	assert.Equal(t, params.FileSize, entry.FileSize)
	assert.Equal(t, userID.String(), entry.UserID.String())
}

func TestGetBinary(t *testing.T) {
	t.Parallel()

	dbStorage, mock := testutils.SetupDBStorage(t)
	ctx := t.Context()

	binaryID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	params := db.GetBinaryEntryByIDParams{
		ID:     pgtype.UUID{Bytes: binaryID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	}

	rows := pgxmock.NewRows([]string{
		"id", "user_id", "file_name", "file_size", "file_url", "created_at", "updated_at",
	}).AddRow(binaryID.String(), userID.String(), "img.png", int64(1024), "https://example.com/img.png", now, now)

	mock.ExpectQuery(`SELECT`).
		WithArgs(params.ID, params.UserID).
		WillReturnRows(rows)

	entry, err := dbStorage.GetBinary(ctx, binaryID.String(), params.UserID)
	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.Equal(t, "img.png", entry.FileName)
}

func TestGetBinaries(t *testing.T) {
	t.Parallel()

	dbStorage, mock := testutils.SetupDBStorage(t)
	ctx := t.Context()

	userID := uuid.New()
	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"id", "user_id", "file_name", "file_size", "file_url", "created_at", "updated_at",
	}).AddRow(uuid.New().String(), userID.String(), "a.txt", int64(200), "url1", now, now).
		AddRow(uuid.New().String(), userID.String(), "b.txt", int64(300), "url2", now, now)

	mock.ExpectQuery(`SELECT`).
		WithArgs(pgtype.UUID{Bytes: userID, Valid: true}).
		WillReturnRows(rows)

	entries, err := dbStorage.GetBinaries(ctx, userID.String())
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "a.txt", entries[0].FileName)
	assert.Equal(t, "b.txt", entries[1].FileName)
}

func TestDeleteBinary(t *testing.T) {
	t.Parallel()

	dbStorage, mock := testutils.SetupDBStorage(t)
	ctx := t.Context()

	param := db.DeleteBinaryEntryParams{
		ID:     pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
	}

	mock.ExpectExec(`DELETE FROM binary_entries WHERE id = \$1 and user_id = \$2`).
		WithArgs(param.ID, param.UserID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := dbStorage.DeleteBinary(ctx, param)
	require.NoError(t, err)
}

func TestStoreBinary_Error(t *testing.T) {
	t.Parallel()

	dbStorage, mock := testutils.SetupDBStorage(t)
	ctx := t.Context()

	userID := uuid.New()
	params := db.StoreBinaryEntryParams{
		UserID:   pgtype.UUID{Bytes: userID, Valid: true},
		FileName: "file.txt",
		FileSize: 1234,
		FileUrl:  "https://example.com/file.txt",
	}

	expectedErr := errors.New("database error")
	mock.ExpectQuery(`INSERT INTO binary_entries`).
		WithArgs(params.UserID, params.FileName, params.FileUrl, params.FileSize).
		WillReturnError(expectedErr)

	entry, err := dbStorage.StoreBinary(ctx, params)
	require.Error(t, err)
	require.Nil(t, entry)
	assert.Contains(t, err.Error(), "failed to store binary entry")
}

func TestGetBinary_NotFound(t *testing.T) {
	t.Parallel()

	dbStorage, mock := testutils.SetupDBStorage(t)
	ctx := t.Context()

	binaryID := uuid.New()
	userID := uuid.New()

	params := db.GetBinaryEntryByIDParams{
		ID:     pgtype.UUID{Bytes: binaryID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	}

	mock.ExpectQuery(`SELECT`).
		WithArgs(params.ID, params.UserID).
		WillReturnError(pgx.ErrNoRows)

	entry, err := dbStorage.GetBinary(ctx, binaryID.String(), params.UserID)
	require.Error(t, err)
	require.Nil(t, entry)
	assert.Contains(t, err.Error(), "failed to find binary")
}

func TestGetBinary_InvalidID(t *testing.T) {
	t.Parallel()

	dbStorage, mock := testutils.SetupDBStorage(t)
	ctx := t.Context()

	binaryID := uuid.New()
	userID := uuid.New()

	params := db.GetBinaryEntryByIDParams{
		ID:     pgtype.UUID{Bytes: binaryID, Valid: false},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	}

	mock.ExpectQuery(`SELECT`).
		WithArgs(params.ID, params.UserID).
		WillReturnError(pgx.ErrNoRows)

	entry, err := dbStorage.GetBinary(ctx, params.ID.String(), params.UserID)
	require.Error(t, err)
	require.Nil(t, entry)
	assert.Contains(t, err.Error(), "failed to find binary")
}

func TestGetBinaries_EmptyResult(t *testing.T) {
	t.Parallel()

	dbStorage, mock := testutils.SetupDBStorage(t)
	ctx := t.Context()

	userID := uuid.New()

	rows := pgxmock.NewRows([]string{
		"id", "user_id", "file_name", "file_size", "file_url", "created_at", "updated_at",
	})

	mock.ExpectQuery(`SELECT`).
		WithArgs(pgtype.UUID{Bytes: userID, Valid: true}).
		WillReturnRows(rows)

	entries, err := dbStorage.GetBinaries(ctx, userID.String())
	require.NoError(t, err)
	require.Empty(t, entries)
}

func TestGetBinaries_Error(t *testing.T) {
	t.Parallel()

	dbStorage, mock := testutils.SetupDBStorage(t)
	ctx := t.Context()

	userID := uuid.New()
	expectedErr := errors.New("database error")

	mock.ExpectQuery(`SELECT`).
		WithArgs(pgtype.UUID{Bytes: userID, Valid: true}).
		WillReturnError(expectedErr)

	entries, err := dbStorage.GetBinaries(ctx, userID.String())
	require.Error(t, err)
	require.Nil(t, entries)
	assert.Contains(t, err.Error(), "failed to find cards")
}

func TestDeleteBinary_NotFound(t *testing.T) {
	t.Parallel()

	dbStorage, mock := testutils.SetupDBStorage(t)
	ctx := t.Context()

	param := db.DeleteBinaryEntryParams{
		ID:     pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
	}

	mock.ExpectExec(`DELETE FROM binary_entries WHERE id = \$1 and user_id = \$2`).
		WithArgs(param.ID, param.UserID).
		WillReturnError(pgx.ErrNoRows)

	err := dbStorage.DeleteBinary(ctx, param)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete binary entry")
}

func TestDeleteBinary_Error(t *testing.T) {
	t.Parallel()

	dbStorage, mock := testutils.SetupDBStorage(t)
	ctx := t.Context()

	param := db.DeleteBinaryEntryParams{
		ID:     pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
	}

	expectedErr := errors.New("database error")
	mock.ExpectExec(`DELETE FROM binary_entries WHERE id = \$1 and user_id = \$2`).
		WithArgs(param.ID, param.UserID).
		WillReturnError(expectedErr)

	err := dbStorage.DeleteBinary(ctx, param)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete binary entry")
}
