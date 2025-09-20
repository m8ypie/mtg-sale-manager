package db

import (
	"log"
	"sync"

	"github.com/m8ypie/mtg-sale-manager/internal/config"
	"github.com/m8ypie/mtg-sale-manager/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
	// GetCCards   func(ctx context.Context) ([]Card, error)
)

func GetDb() *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(postgres.Open(config.PostGresUri), &gorm.Config{})
		db.AutoMigrate(&models.EbayListingGorm{})
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

	})
	return db
}
