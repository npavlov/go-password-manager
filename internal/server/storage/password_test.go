//nolint:dupl,gochecknoglobals,lll,gosec
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

// Static IDs for testing.
const (
	testUserID     = "123e4567-e89b-12d3-a456-426614174000"
	testPasswordID = "223e4567-e89b-12d3-a456-426614174000"
)

var (
	userUUID     = pgtype.UUID{Bytes: uuid.MustParse(testUserID), Valid: true}
	passwordUUID = pgtype.UUID{Bytes: uuid.MustParse(testPasswordID), Valid: true}
)

func TestStorePassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		createParams  db.CreatePasswordEntryParams
		mock          func(mock pgxmock.PgxPoolIface)
		want          *db.Password
		wantErr       bool
		expectedError string
	}{
		{
			name: "successful password creation",
			createParams: db.CreatePasswordEntryParams{
				UserID:   userUUID,
				Login:    "test_login",
				Password: "test_password",
			},
			mock: func(mock pgxmock.PgxPoolIface) {
				now := time.Now()
				rows := pgxmock.NewRows([]string{"id", "user_id", "login", "password", "created_at", "updated_at"}).
					AddRow(passwordUUID, userUUID, "test_login", "test_password", now, now)
				mock.ExpectQuery("INSERT INTO passwords").
					WithArgs(userUUID, "test_login", "test_password").
					WillReturnRows(rows)
			},
			want: &db.Password{
				ID:       passwordUUID,
				UserID:   userUUID,
				Login:    "test_login",
				Password: "test_password",
			},
			wantErr: false,
		},
		{
			name: "database error",
			createParams: db.CreatePasswordEntryParams{
				UserID:   userUUID,
				Login:    "test_login",
				Password: "test_password",
			},
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("INSERT INTO passwords").
					WithArgs(userUUID, "test_login", "test_password").
					WillReturnError(errors.New("db error"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "failed to store password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			result, err := storage.StorePassword(t.Context(), tt.createParams)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, result.ID)
				require.Equal(t, tt.want.UserID, result.UserID)
				require.Equal(t, tt.want.Login, result.Login)
				require.Equal(t, tt.want.Password, result.Password)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		passwordID    string
		userID        pgtype.UUID
		mock          func(mock pgxmock.PgxPoolIface)
		want          *db.Password
		wantErr       bool
		expectedError string
	}{
		{
			name:       "successful password retrieval",
			passwordID: testPasswordID,
			userID:     userUUID,
			mock: func(mock pgxmock.PgxPoolIface) {
				now := time.Now()
				rows := pgxmock.NewRows([]string{"id", "user_id", "login", "password", "created_at", "updated_at"}).
					AddRow(passwordUUID, userUUID, "test_login", "test_password", now, now)
				mock.ExpectQuery("SELECT passwords.id, passwords.user_id, passwords.login, passwords.password, passwords.created_at, passwords.updated_at").
					WithArgs(passwordUUID, userUUID).
					WillReturnRows(rows)
			},
			want: &db.Password{
				ID:       passwordUUID,
				UserID:   userUUID,
				Login:    "test_login",
				Password: "test_password",
			},
			wantErr: false,
		},
		{
			name:       "password not found",
			passwordID: testPasswordID,
			userID:     userUUID,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT passwords.id, passwords.user_id, passwords.login, passwords.password, passwords.created_at, passwords.updated_at").
					WithArgs(passwordUUID, userUUID).
					WillReturnError(errors.New("db error"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "failed to create password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			result, err := storage.GetPassword(t.Context(), tt.passwordID, tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, result.ID)
				require.Equal(t, tt.want.UserID, result.UserID)
				require.Equal(t, tt.want.Login, result.Login)
				require.Equal(t, tt.want.Password, result.Password)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPasswords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		userID        string
		mock          func(mock pgxmock.PgxPoolIface)
		want          []db.Password
		wantErr       bool
		expectedError string
	}{
		{
			name:   "successful passwords retrieval",
			userID: testUserID,
			mock: func(mock pgxmock.PgxPoolIface) {
				now := time.Now()
				rows := pgxmock.NewRows([]string{"id", "user_id", "login", "password", "created_at", "updated_at"}).
					AddRow(passwordUUID, userUUID, "test_login1", "test_password1", now, now).
					AddRow(passwordUUID, userUUID, "test_login2", "test_password2", now.Add(-time.Hour), now.Add(-time.Hour))
				mock.ExpectQuery("SELECT id, user_id, login, password, created_at, updated_at FROM passwords").
					WithArgs(userUUID).
					WillReturnRows(rows)
			},
			want: []db.Password{
				{
					ID:       passwordUUID,
					UserID:   userUUID,
					Login:    "test_login1",
					Password: "test_password1",
				},
				{
					ID:       passwordUUID,
					UserID:   userUUID,
					Login:    "test_login2",
					Password: "test_password2",
				},
			},
			wantErr: false,
		},
		{
			name:   "no passwords found",
			userID: testUserID,
			mock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "user_id", "login", "password", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, user_id, login, password, created_at, updated_at FROM passwords").
					WithArgs(userUUID).
					WillReturnRows(rows)
			},
			want:    []db.Password{},
			wantErr: false,
		},
		{
			name:   "database error",
			userID: testUserID,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT id, user_id, login, password, created_at, updated_at FROM passwords").
					WithArgs(userUUID).
					WillReturnError(errors.New("db error"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "failed to create password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			result, err := storage.GetPasswords(t.Context(), tt.userID)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tt.want), len(result))
				for i, password := range tt.want {
					require.Equal(t, password.ID, result[i].ID)
					require.Equal(t, password.UserID, result[i].UserID)
					require.Equal(t, password.Login, result[i].Login)
					require.Equal(t, password.Password, result[i].Password)
				}
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeletePassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		passwordID    string
		userID        pgtype.UUID
		mock          func(mock pgxmock.PgxPoolIface)
		wantErr       bool
		expectedError string
	}{
		{
			name:       "successful password deletion",
			passwordID: testPasswordID,
			userID:     userUUID,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM passwords").
					WithArgs(passwordUUID, userUUID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: false,
		},
		{
			name:       "database error",
			passwordID: testPasswordID,
			userID:     userUUID,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("DELETE FROM passwords").
					WithArgs(passwordUUID, userUUID).
					WillReturnError(errors.New("db error"))
			},
			wantErr:       true,
			expectedError: "failed to delete password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			err := storage.DeletePassword(t.Context(), tt.passwordID, tt.userID)

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

func TestUpdatePassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		updateParams  db.UpdatePasswordEntryParams
		mock          func(mock pgxmock.PgxPoolIface)
		want          *db.Password
		wantErr       bool
		expectedError string
	}{
		{
			name: "successful password update",
			updateParams: db.UpdatePasswordEntryParams{
				Login:    "updated_login",
				Password: "updated_password",
				ID:       passwordUUID,
			},
			mock: func(mock pgxmock.PgxPoolIface) {
				now := time.Now()
				rows := pgxmock.NewRows([]string{"id", "user_id", "login", "password", "created_at", "updated_at"}).
					AddRow(passwordUUID, userUUID, "updated_login", "updated_password", now.Add(-time.Hour), now)
				mock.ExpectQuery("UPDATE passwords").
					WithArgs("updated_login", "updated_password", passwordUUID).
					WillReturnRows(rows)
			},
			want: &db.Password{
				ID:       passwordUUID,
				UserID:   userUUID,
				Login:    "updated_login",
				Password: "updated_password",
			},
			wantErr: false,
		},
		{
			name: "database error",
			updateParams: db.UpdatePasswordEntryParams{
				Login:    "updated_login",
				Password: "updated_password",
				ID:       passwordUUID,
			},
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("UPDATE passwords").
					WithArgs("updated_login", "updated_password", passwordUUID).
					WillReturnError(errors.New("db error"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "failed to update password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			result, err := storage.UpdatePassword(t.Context(), tt.updateParams)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, result.ID)
				require.Equal(t, tt.want.UserID, result.UserID)
				require.Equal(t, tt.want.Login, result.Login)
				require.Equal(t, tt.want.Password, result.Password)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
