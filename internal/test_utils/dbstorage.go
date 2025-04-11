package testutils

import (
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/storage"
)

// SetupDBStorage creates a DBStorage with a pgxmock pool and returns the DBStorate and mock pool.
//
//nolint:ireturn
func SetupDBStorage(t *testing.T) (*storage.DBStorage, pgxmock.PgxPoolIface) {
	t.Helper()

	mockDB, err := pgxmock.NewPool()
	require.NoError(t, err)

	log := GetTLogger()
	dbStorage := storage.NewDBStorage(mockDB, log)

	return dbStorage, mockDB
}
