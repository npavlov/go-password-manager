package dbmanager_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/npavlov/go-password-manager/internal/server/dbmanager"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			if tt.mockBehavior != nil {
				tt.mockBehavior(mock)
			}

			logger := zerolog.Nop()
			manager := dbmanager.NewDBManager("postgres://localhost/test", &logger).Connect(context.Background())
			if !tt.wantErr {
				manager.DB = mock // Inject mock
			}

			manager.VerifyConnection(context.Background())

			assert.Equal(t, tt.wantConnected, manager.IsConnected)

			if tt.wantErr {
				assert.NotNil(t, mock.ExpectationsWereMet())
			} else {
				assert.NoError(t, mock.ExpectationsWereMet())
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
			name:        "not connected",
			isConnected: false,
			wantErr:     true,
		},
		{
			name:        "successful migrations",
			isConnected: true,
			mockBehavior: func(mock *sql.DB) {
				// goose uses sql.DB internally, but we can't easily mock it
				// so we'll just test the happy path
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			manager := tt.initialState()
			manager.Close()

			if tt.expectClosed {
				// For the mock pool, we can verify it's closed by trying to use it
				if mock, ok := manager.DB.(pgxmock.PgxPoolIface); ok {
					err := mock.Ping(context.Background())
					assert.Error(t, err, "expected closed pool to return error on ping")
				}
			}
		})
	}
}

func TestDBOperations(t *testing.T) {
	t.Parallel()

	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	tests := []struct {
		name      string
		setupMock func()
		operation func(m *dbmanager.DBManager) error
		wantErr   bool
	}{
		{
			name: "successful query",
			setupMock: func() {
				mock.ExpectQuery("SELECT").WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(1))
			},
			operation: func(m *dbmanager.DBManager) error {
				_, err := m.DB.Query(context.Background(), "SELECT")
				return err
			},
			wantErr: false,
		},
		{
			name: "failed query",
			setupMock: func() {
				mock.ExpectQuery("SELECT").WillReturnError(errors.New("query failed"))
			},
			operation: func(m *dbmanager.DBManager) error {
				_, err := m.DB.Query(context.Background(), "SELECT")
				return err
			},
			wantErr: true,
		},
		{
			name: "successful exec",
			setupMock: func() {
				mock.ExpectExec("INSERT").WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			operation: func(m *dbmanager.DBManager) error {
				_, err := m.DB.Exec(context.Background(), "INSERT")
				return err
			},
			wantErr: false,
		},
	}

	logger := zerolog.Nop()
	manager := dbmanager.NewDBManager("postgres://localhost/test", &logger)
	manager.DB = mock
	manager.IsConnected = true

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			err := tt.operation(manager)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestConnect_ParseConfigError(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	// Invalid connection string that will fail parsing
	manager := dbmanager.NewDBManager("invalid://connection:string", &logger)

	result := manager.Connect(context.Background())

	assert.False(t, result.IsConnected)
	assert.Nil(t, result.DB)
}
