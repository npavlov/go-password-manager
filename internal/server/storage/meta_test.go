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

func TestAddMeta(t *testing.T) {
	t.Parallel()

	recordID := uuid.New()
	uuid := pgtype.UUID{Bytes: recordID, Valid: true}

	tests := []struct {
		name     string
		recordID string
		key      string
		value    string
		mock     func(mock pgxmock.PgxPoolIface)
		want     *db.Metainfo
		wantErr  bool
	}{
		{
			name:     "successful meta addition",
			recordID: recordID.String(),
			key:      "test_key",
			value:    "test_value",
			mock: func(mock pgxmock.PgxPoolIface) {
				now := pgtype.Timestamp{Time: time.Now(), Valid: true}
				rows := pgxmock.NewRows([]string{"id", "item_id", "key", "value", "created_at", "updated_at"}).
					AddRow(uuid, uuid, "test_key", "test_value", now, now)
				mock.ExpectQuery("INSERT INTO metainfo").
					WithArgs(uuid, "test_key", "test_value").
					WillReturnRows(rows)
			},
			want: &db.Metainfo{
				ID:    uuid,
				Key:   "test_key",
				Value: "test_value",
			},
			wantErr: false,
		},
		{
			name:     "database error",
			recordID: recordID.String(),
			key:      "test_key",
			value:    "test_value",
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("INSERT INTO metainfo").
					WithArgs(uuid, "test_key", "test_value").
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

			result, err := storage.AddMeta(t.Context(), tt.recordID, tt.key, tt.value)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to get items")
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.Key, result.Key)
				require.Equal(t, tt.want.Value, result.Value)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetMetaInfo(t *testing.T) {
	t.Parallel()

	recordId := uuid.New()
	uuid := pgtype.UUID{Bytes: recordId, Valid: true}

	tests := []struct {
		name     string
		recordId string
		mock     func(mock pgxmock.PgxPoolIface)
		want     []db.GetMetaInfoByItemIDRow
		wantErr  bool
	}{
		{
			name:     "successful meta retrieval",
			recordId: recordId.String(),
			mock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"key", "value"}).
					AddRow("key1", "value1").
					AddRow("key2", "value2")
				mock.ExpectQuery("SELECT key, value FROM metainfo").
					WithArgs(uuid).
					WillReturnRows(rows)
			},
			want: []db.GetMetaInfoByItemIDRow{
				{Key: "key1", Value: "value1"},
				{Key: "key2", Value: "value2"},
			},
			wantErr: false,
		},
		{
			name:     "no meta found",
			recordId: recordId.String(),
			mock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"key", "value"})
				mock.ExpectQuery("SELECT key, value FROM metainfo").
					WithArgs(uuid).
					WillReturnRows(rows)
			},
			want:    []db.GetMetaInfoByItemIDRow{},
			wantErr: false,
		},
		{
			name:     "database error",
			recordId: recordId.String(),
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT key, value FROM metainfo").
					WithArgs(uuid).
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

			result, err := storage.GetMetaInfo(t.Context(), tt.recordId)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to get items")
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tt.want), len(result))
				for i, item := range tt.want {
					require.Equal(t, item.Key, result[i].Key)
					require.Equal(t, item.Value, result[i].Value)
				}
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteMetaInfo(t *testing.T) {
	t.Parallel()

	recordId := uuid.New()
	uuid := pgtype.UUID{Bytes: recordId, Valid: true}

	tests := []struct {
		name    string
		key     string
		itemId  string
		mock    func(mock pgxmock.PgxPoolIface)
		wantErr bool
	}{
		{
			name:   "successful deletion",
			key:    "test_key",
			itemId: recordId.String(),
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM metainfo").
					WithArgs(uuid, "test_key").
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: false,
		},
		{
			name:   "database error",
			key:    "test_key",
			itemId: recordId.String(),
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM metainfo").
					WithArgs(uuid, "test_key").
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

			err := storage.DeleteMetaInfo(t.Context(), tt.key, tt.itemId)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to delete items")
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
