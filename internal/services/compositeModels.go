package services

import (
	"github.com/BlueMonday/go-scryfall"
	"github.com/m8ypie/mtg-sale-manager/internal/models"
)

type CompositeEbayListing struct {
	EbayListing  models.EbayListing
	ScryfallInfo scryfall.Card
	PriceInfo    models.PriceInfo
}
