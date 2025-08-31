package main

import (
	"context"
	"ebay-sale-manager/backend/db"
	"fmt"
	"log"
)

func main() {
	cards, err := db.GetCards(context.Background())
	if err != nil {
		log.Fatalf("Error fetching cards: %v", err)
	}

	for _, card := range cards {
		fmt.Printf("Card ID: %d, Name: %s, UUID: %s\n", card.ID, card.Name, card.Uuid)
	}
}
