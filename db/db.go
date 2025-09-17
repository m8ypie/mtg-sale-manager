package db

import (
	"database/sql"
	"log"
	"sync"

	"github.com/m8ypie/mtg-sale-manager/config"
	"github.com/xo/dburl"
)

var (
	dbInstance *sql.DB
	once       sync.Once
	// GetC÷ards   func(ctx context.Context) ([]Card, error)
)

func GetDb() *sql.DB {
	once.Do(func() {
		var err error

		dbInstance, err = dburl.Open(config.PostGresUri)

		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

	})
	return dbInstance
}
