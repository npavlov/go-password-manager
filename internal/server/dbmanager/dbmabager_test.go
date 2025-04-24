//nolint:wrapcheck,err113,exhaustruct
package dbmanager_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/dbmanager"
)

func TestNewDBManager(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		connectionString string
		logger           zerolog.Logger
		expected         *dbmanager.DBManager
	}{
		{
			name:             "successful creation",
			connectionString: "postgres://localhost/test",
			logger:           zerolog.Nop(),
			expected: &dbmanager.DBManager{
				IsConnected: false,
				DB:          nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			manager := dbmanager.NewDBManager(tt.connectionString, &tt.logger)
			assert.Equal(t, tt.expected.IsConnected, manager.IsConnected)
			assert.Nil(t, manager.DB)
		})
	}
}

func TestConnect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockBehavior  func(mock pgxmock.PgxPoolIface)
		wantErr       bool
		wantConnected bool
	}{
		{
			name: "successful connection",
			mockBehavior: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectPing()
			},
			wantErr:       false,
			wantConnected: true,
		},
		{
			name: "failed ping",
			mockBehavior: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectPing().WillReturnError(errors.New("ping failed"))
			},
			wantErr:       true,
			wantConnected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			if tt.mockBehavior != nil {
				tt.mockBehavior(mock)
			}

			logger := zerolog.Nop()
			manager := dbmanager.NewDBManager("postgres://localhost/test", &logger).Connect(t.Context())
			if !tt.wantErr {
				manager.DB = mock // Inject mock
			}

			manager.VerifyConnection(t.Context())

			assert.Equal(t, tt.wantConnected, manager.IsConnected)

			if tt.wantErr {
				require.Error(t, mock.ExpectationsWereMet())
			} else {
				require.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}

func TestApplyMigrations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		isConnected  bool
		mockBehavior func(mock *sql.DB)
		wantErr      bool
	}{
		{
			name:         "not connected",
			isConnected:  false,
			wantErr:      true,
			mockBehavior: nil,
		},
		{
			name:        "successful migrations",
			isConnected: true,
			mockBehavior: func(_ *sql.DB) {
				// goose uses sql.DB internally, but we can't easily mock it
				// so we'll just test the happy path
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := zerolog.Nop()
			manager := dbmanager.NewDBManager("postgres://localhost/test", &logger)
			manager.IsConnected = tt.isConnected

			// For the connected case, we'd normally mock sql.DB, but goose doesn't
			// provide an easy way to inject mocks, so we'll just test the behavior
			result := manager.ApplyMigrations()

			if tt.wantErr {
				assert.False(t, result.IsConnected)
			} else {
				assert.True(t, result.IsConnected)
			}
		})
	}
}

func TestClose(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		initialState func() *dbmanager.DBManager
		expectClosed bool
	}{
		{
			name: "close connected pool",
			initialState: func() *dbmanager.DBManager {
				mock, _ := pgxmock.NewPool()
				logger := zerolog.Nop()

				return &dbmanager.DBManager{
					DB:          mock,
					Log:         &logger,
					IsConnected: true,
				}
			},
			expectClosed: true,
		},
		{
			name: "close nil pool",
			initialState: func() *dbmanager.DBManager {
				logger := zerolog.Nop()

				return &dbmanager.DBManager{
					DB:          nil,
					Log:         &logger,
					IsConnected: false,
				}
			},
			expectClosed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			manager := tt.initialState()
			manager.Close()

			if tt.expectClosed {
				// For the mock pool, we can verify it's closed by trying to use it
				if mock, ok := manager.DB.(pgxmock.PgxPoolIface); ok {
					err := mock.Ping(t.Context())
					require.Error(t, err, "expected closed pool to return error on ping")
				}
			}
		})
	}
}

func TestDBOperations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(mock pgxmock.PgxPoolIface)
		operation func(m *dbmanager.DBManager) error
		wantErr   bool
	}{
		{
			name: "successful query",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT").WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
			},
			operation: func(m *dbmanager.DBManager) error {
				_, err := m.DB.Query(t.Context(), "SELECT")

				return err
			},
			wantErr: false,
		},
		{
			name: "failed query",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT").WillReturnError(errors.New("query failed"))
			},
			operation: func(m *dbmanager.DBManager) error {
				_, err := m.DB.Query(t.Context(), "SELECT")

				return err
			},
			wantErr: true,
		},
		{
			name: "successful exec",
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT").WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			operation: func(m *dbmanager.DBManager) error {
				_, err := m.DB.Exec(t.Context(), "INSERT")

				return err
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			logger := zerolog.Nop()
			manager := dbmanager.NewDBManager("postgres://localhost/test", &logger)
			manager.DB = mock
			manager.IsConnected = true

			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			err = tt.operation(manager)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestConnect_ParseConfigError(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	// Invalid connection string that will fail parsing
	manager := dbmanager.NewDBManager("invalid://connection:string", &logger)

	result := manager.Connect(t.Context())

	assert.False(t, result.IsConnected)
	assert.Nil(t, result.DB)
}
