package clients

import (
	"context"

	"log"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/m8ypie/mtg-sale-manager/internal/models"
)

type ScryfallClient struct {
	Scryfall *scryfall.Client
	ctx      context.Context
}

func NewScryfallClient() *ScryfallClient {
	ctx := context.Background()
	client, err := scryfall.NewClient()
	if err != nil {
		panic(err)
	}
	return &ScryfallClient{
		Scryfall: client,
		ctx:      ctx,
	}
}

func (s ScryfallClient) GetCardDetails(ebayListing *models.EbayListing) (*scryfall.Card, error) {
	card, err := s.Scryfall.GetCard(s.ctx, ebayListing.ScryFallId)
	if err != nil {
		// Handle error
		log.Fatalf("Error fetching Scryfall card details: %v", err)
	}

	return &card, err
}
