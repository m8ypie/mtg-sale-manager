package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once             sync.Once
	PostGresUri      string
	EbayAppId        string
	EbayCertId       string
	EbayEnvironment  string
	EbayRedirectUri  string
	EbayRefreshToken string
)

func initConfig() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		PostGresUri = os.Getenv("POSTGRES_URI")
		EbayAppId = os.Getenv("EBAY_APP_ID")
		EbayCertId = os.Getenv("EBAY_CERT_ID")
		EbayEnvironment = os.Getenv("EBAY_ENVIRONMENT")
		EbayRedirectUri = os.Getenv("EBAY_REDIRECT_URI")
		EbayRefreshToken = os.Getenv("EBAY_REFRESH_TOKEN")
	})
}

func Init() {
	initConfig()
}
