package main

import (
	"context"
	"fmt"
	"log"

	"github.com/m8ypie/mtg-sale-manager/config"
	"github.com/m8ypie/mtg-sale-manager/db"
	"github.com/m8ypie/mtg-sale-manager/models"
)

func main() {
	fmt.Printf("config.PostGresUri: %v\n", config.PostGresUri)
	cards, err := models.CardsByUUID(context.Background(), db.GetDb(), "some-uuid")
	if err != nil {
		log.Fatalf("Error fetching cards: %v", err)
	}

	for _, card := range cards {
		fmt.Printf("Card ID: %d, Name: %s, UUID: %s\n", card.ID, card.Name, card.UUID)
	}
}
