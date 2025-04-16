//nolint:dupl
package storage_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/db"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

// Static test data.
const (
	testUsername = "testuser"
	testEmail    = "test@example.com"
	testPassword = "testpassword"
	testEncKey   = "testenckey"
)

func TestRegisterUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		createParams  db.CreateUserParams
		mock          func(mock pgxmock.PgxPoolIface)
		want          *db.User
		wantErr       bool
		expectedError string
	}{
		{
			name: "successful user registration",
			createParams: db.CreateUserParams{
				Username:      testUsername,
				Password:      testPassword,
				EncryptionKey: testEncKey,
				Email:         testEmail,
			},
			mock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "username", "email", "password", "encryption_key"}).
					AddRow(userUUID, testUsername, testEmail, testPassword, testEncKey)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(testUsername, testPassword, testEncKey, testEmail).
					WillReturnRows(rows)
			},
			want: &db.User{
				ID:            userUUID,
				Username:      testUsername,
				Email:         testEmail,
				Password:      testPassword,
				EncryptionKey: testEncKey,
			},
			wantErr: false,
		},
		{
			name: "duplicate username",
			createParams: db.CreateUserParams{
				Username:      testUsername,
				Password:      testPassword,
				EncryptionKey: testEncKey,
				Email:         testEmail,
			},
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(testUsername, testPassword, testEncKey, testEmail).
					WillReturnError(errors.New("db error"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "failed to create user",
		},
		{
			name: "database error",
			createParams: db.CreateUserParams{
				Username:      testUsername,
				Password:      testPassword,
				EncryptionKey: testEncKey,
				Email:         testEmail,
			},
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(testUsername, testPassword, testEncKey, testEmail).
					WillReturnError(errors.New("db error"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			result, err := storage.RegisterUser(t.Context(), tt.createParams)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, result.ID)
				require.Equal(t, tt.want.Username, result.Username)
				require.Equal(t, tt.want.Email, result.Email)
				require.Equal(t, tt.want.Password, result.Password)
				require.Equal(t, tt.want.EncryptionKey, result.EncryptionKey)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name          string
		username      string
		mock          func(mock pgxmock.PgxPoolIface)
		want          *db.User
		wantErr       bool
		expectedError string
	}{
		{
			name:     "successful user retrieval by username",
			username: testUsername,
			mock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "username", "email", "password", "encryption_key"}).
					AddRow(userUUID, testUsername, testEmail, testPassword, testEncKey)
				mock.ExpectQuery("SELECT").
					WithArgs(testUsername).
					WillReturnRows(rows)
			},
			want: &db.User{
				ID:            userUUID,
				Username:      testUsername,
				Email:         testEmail,
				Password:      testPassword,
				EncryptionKey: testEncKey,
			},
			wantErr: false,
		},
		{
			name:     "user not found",
			username: "nonexistent",
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT").
					WithArgs("nonexistent").
					WillReturnError(errors.New("failed to create user"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "failed to create user",
		},
		{
			name:     "database error",
			username: testUsername,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT").
					WithArgs(testUsername).
					WillReturnError(errors.New("failed to create user"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			result, err := storage.GetUser(t.Context(), tt.username)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, result.ID)
				require.Equal(t, tt.want.Username, result.Username)
				require.Equal(t, tt.want.Email, result.Email)
				require.Equal(t, tt.want.Password, result.Password)
				require.Equal(t, tt.want.EncryptionKey, result.EncryptionKey)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUserById(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		userId        pgtype.UUID
		mock          func(mock pgxmock.PgxPoolIface)
		want          *db.User
		wantErr       bool
		expectedError string
	}{
		{
			name:   "successful user retrieval by ID",
			userId: userUUID,
			mock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "username", "email", "password", "encryption_key"}).
					AddRow(userUUID, testUsername, testEmail, testPassword, testEncKey)
				mock.ExpectQuery("SELECT").
					WithArgs(userUUID).
					WillReturnRows(rows)
			},
			want: &db.User{
				ID:            userUUID,
				Username:      testUsername,
				Email:         testEmail,
				Password:      testPassword,
				EncryptionKey: testEncKey,
			},
			wantErr: false,
		},
		{
			name:   "user not found by ID",
			userId: userUUID,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT").
					WithArgs(userUUID).
					WillReturnError(errors.New("db error"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "failed to create user",
		},
		{
			name:   "database error",
			userId: userUUID,
			mock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT").
					WithArgs(userUUID).
					WillReturnError(errors.New("db error"))
			},
			want:          nil,
			wantErr:       true,
			expectedError: "failed to create user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage, mock := testutils.SetupDBStorage(t)
			tt.mock(mock)

			result, err := storage.GetUserByID(t.Context(), tt.userId)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, result.ID)
				require.Equal(t, tt.want.Username, result.Username)
				require.Equal(t, tt.want.Email, result.Email)
				require.Equal(t, tt.want.Password, result.Password)
				require.Equal(t, tt.want.EncryptionKey, result.EncryptionKey)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
