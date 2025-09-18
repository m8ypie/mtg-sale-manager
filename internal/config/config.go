package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once        sync.Once
	PostGresUri string
)

func initConfig() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		PostGresUri = os.Getenv("POSTGRES_URI")
	})
}

func init() {
	initConfig()
}
