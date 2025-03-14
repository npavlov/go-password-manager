package storage

import (
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/dbmanager"
	"github.com/rs/zerolog"
)

type DBStorage struct {
	Queries *db.Queries
	log     *zerolog.Logger
	dbCon   dbmanager.PgxPool
}

// NewDBStorage initializes a new DBStorage instance.
func NewDBStorage(dbCon dbmanager.PgxPool, log *zerolog.Logger) *DBStorage {
	return &DBStorage{
		dbCon:   dbCon,
		Queries: db.New(dbCon),
		log:     log,
	}
}
