//nolint:gochecknoglobals,gosec
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

// Static IDs and values for testing.
const (
	testToken   = "test_refresh_token"
	testTokenID = "223e4567-e89b-12d3-a456-426614174000"
)

var (
	tokenUUID   = pgtype.UUID{Bytes: uuid.MustParse(testTokenID), Valid: true}
	fixedTime   = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	pgFixedTime = pgtype.Timestamp{Time: fixedTime, Valid: true}
)

func TestStoreToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		userID        pgtype.UUID
		refreshToken  string
		expiresAt     time.Time
		mock          func(mock pgxmock.PgxPoolIface)
		wantErr       bool
		expectedError string
	}{
		{
			name:         "successful token storage",
			userID:       userUUID,
			refreshToken: testToken,
			expiresAt:    fixedTime,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO refresh_tokens").
					WithArgs(userUUID, testToken, pgFixedTime).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: false,
		},
		{
			name:         "database error",
			userID:       userUUID,
			refreshToken: testToken,
			expiresAt:    fixedTime,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO refresh_tokens").
					WithArgs(userUUID, testToken, pgFixedTime).
					WillReturnError(errors.New("failed to scan expires at"))
			},
			wantErr:       true,
			expectedError: "failed to scan expires at",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			err := storage.StoreToken(t.Context(), tt.userID, tt.refreshToken, tt.expiresAt)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		token         string
		mock          func(mock pgxmock.PgxPoolIface)
		want          db.GetRefreshTokenRow
		wantErr       bool
		expectedError string
	}{
		{
			name:  "successful token retrieval",
			token: testToken,
			mock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "user_id", "token", "expires_at"}).
					AddRow(tokenUUID, userUUID, testToken, pgFixedTime)
				mock.ExpectQuery("SELECT id, user_id, token, expires_at FROM refresh_tokens").
					WithArgs(testToken).
					WillReturnRows(rows)
			},
			want: db.GetRefreshTokenRow{
				ID:        tokenUUID,
				UserID:    userUUID,
				Token:     testToken,
				ExpiresAt: pgFixedTime,
			},
			wantErr: false,
		},
		{
			name:  "token not found",
			token: testToken,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT id, user_id, token, expires_at FROM refresh_tokens").
					WithArgs(testToken).
					WillReturnError(errors.New("db error"))
			},
			want:          db.GetRefreshTokenRow{},
			wantErr:       true,
			expectedError: "",
		},
		{
			name:  "database error",
			token: testToken,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT id, user_id, token, expires_at FROM refresh_tokens").
					WithArgs(testToken).
					WillReturnError(errors.New("db error"))
			},
			want:          db.GetRefreshTokenRow{},
			wantErr:       true,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			result, err := storage.GetToken(t.Context(), tt.token)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, result.ID)
				require.Equal(t, tt.want.UserID, result.UserID)
				require.Equal(t, tt.want.Token, result.Token)
				require.Equal(t, tt.want.ExpiresAt, result.ExpiresAt)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
