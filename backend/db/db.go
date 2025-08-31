package db

import (
	"context"
	"ebay-sale-manager/backend/config"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Card struct {
	// Define fields for Card, for example:
	ID   int
	Name string
	Uuid string
	// Add other fields as needed
}

var (
	dbInstance *pgxpool.Pool
	once       sync.Once
	// GetC÷ards   func(ctx context.Context) ([]Card, error)
)

func GetCards(ctx context.Context) ([]Card, error) {
	rows, err := dbInstance.Query(ctx, "SELECT id, name, uuid FROM cards LIMIT 10")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var card Card
		if err := rows.Scan(&card.ID, &card.Name, &card.Uuid); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func initDatabase() *pgxpool.Pool {
	once.Do(func() {
		var err error
		dbInstance, err = pgxpool.New(context.Background(), config.PostGresUri)

		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

	})
	return dbInstance
}

func init() {
	initDatabase()
}
