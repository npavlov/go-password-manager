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

func TestGetItems(t *testing.T) {
	tests := []struct {
		name     string
		params   db.GetItemsByUserIDParams
		mockRows func() *pgxmock.Rows
		want     []db.GetItemsByUserIDRow
		wantErr  bool
	}{
		{
			name: "successful retrieval with multiple items",
			params: db.GetItemsByUserIDParams{
				UserID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
				Limit:  10,
				Offset: 0,
			},
			mockRows: func() *pgxmock.Rows {
				now := pgtype.Timestamp{Time: time.Now(), Valid: true}

				return pgxmock.NewRows([]string{
					"id", "type", "id_resource", "created_at", "updated_at",
				}).
					AddRow(
						uuid.New(), "password", uuid.New(), now, now,
					).
					AddRow(
						uuid.New(), "card", uuid.New(), now, now,
					)
			},
			want: []db.GetItemsByUserIDRow{
				{
					ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
					Type:       "password",
					IDResource: pgtype.UUID{Bytes: uuid.New(), Valid: true},
					CreatedAt:  pgtype.Timestamp{Time: time.Now(), Valid: true},
					UpdatedAt:  pgtype.Timestamp{Time: time.Now(), Valid: true},
				},
				{
					ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
					Type:       "card",
					IDResource: pgtype.UUID{Bytes: uuid.New(), Valid: true},
					CreatedAt:  pgtype.Timestamp{Time: time.Now(), Valid: true},
					UpdatedAt:  pgtype.Timestamp{Time: time.Now(), Valid: true},
				},
			},
			wantErr: false,
		},
		{
			name: "successful retrieval with no items",
			params: db.GetItemsByUserIDParams{
				UserID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
				Limit:  10,
				Offset: 0,
			},
			mockRows: func() *pgxmock.Rows {
				return pgxmock.NewRows([]string{
					"id", "type", "id_resource", "created_at", "updated_at",
				})
			},
			want:    []db.GetItemsByUserIDRow{},
			wantErr: false,
		},
		{
			name: "database error",
			params: db.GetItemsByUserIDParams{
				UserID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
				Limit:  10,
				Offset: 0,
			},
			mockRows: func() *pgxmock.Rows {
				return nil // No rows expected
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, mock := testutils.SetupDBStorage(t)

			if tt.wantErr {
				mock.ExpectQuery("SELECT").
					WithArgs(tt.params.UserID, tt.params.Limit, tt.params.Offset).
					WillReturnError(errors.New("database error"))
			} else {
				rows := tt.mockRows()
				if len(tt.want) > 0 {
					// Update the mock rows with the expected values
					rows = pgxmock.NewRows([]string{
						"id", "type", "id_resource", "created_at", "updated_at",
					})
					for _, item := range tt.want {
						rows.AddRow(
							item.ID, item.Type, item.IDResource,
							item.CreatedAt, item.UpdatedAt,
						)
					}
				}
				mock.ExpectQuery("SELECT").
					WithArgs(tt.params.UserID, tt.params.Limit, tt.params.Offset).
					WillReturnRows(rows)
			}

			got, err := storage.GetItems(t.Context(), tt.params)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to get items")
			} else {
				require.NoError(t, err)
				if len(tt.want) > 0 {
					require.Len(t, got, len(tt.want))
					for i, item := range tt.want {
						require.Equal(t, item.ID, got[i].ID)
						require.Equal(t, item.Type, got[i].Type)
						require.Equal(t, item.IDResource, got[i].IDResource)
					}
				} else {
					require.Empty(t, got)
				}
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetItemsWithPagination(t *testing.T) {
	storage, mock := testutils.SetupDBStorage(t)

	userID := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	params := db.GetItemsByUserIDParams{
		UserID: userID,
		Limit:  5,
		Offset: 10,
	}

	now := pgtype.Timestamp{Time: time.Now(), Valid: true}
	rows := pgxmock.NewRows([]string{
		"id", "type", "id_resource", "created_at", "updated_at",
	}).
		AddRow(uuid.New().String(), "password", uuid.New().String(), now, now).
		AddRow(uuid.New().String(), "card", uuid.New().String(), now, now)

	mock.ExpectQuery("SELECT").
		WithArgs(params.UserID, params.Limit, params.Offset).
		WillReturnRows(rows)

	result, err := storage.GetItems(t.Context(), params)
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, db.ItemTypePassword, result[0].Type)
	require.Equal(t, db.ItemTypeCard, result[1].Type)
}

func TestGetItemsEmptyResult(t *testing.T) {
	t.Parallel()

	storage, mock := testutils.SetupDBStorage(t)

	userID := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	params := db.GetItemsByUserIDParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	rows := pgxmock.NewRows([]string{
		"id", "type", "id_resource", "created_at", "updated_at",
	})

	mock.ExpectQuery("SELECT").
		WithArgs(params.UserID, params.Limit, params.Offset).
		WillReturnRows(rows)

	result, err := storage.GetItems(t.Context(), params)
	require.NoError(t, err)
	require.Empty(t, result)
}

func TestGetItemsDatabaseError(t *testing.T) {
	storage, mock := testutils.SetupDBStorage(t)

	userID := pgtype.UUID{Bytes: uuid.New(), Valid: true}
	params := db.GetItemsByUserIDParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	mock.ExpectQuery("SELECT").
		WithArgs(params.UserID, params.Limit, params.Offset).
		WillReturnError(errors.New("database error"))

	result, err := storage.GetItems(t.Context(), params)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "failed to get items")
}
