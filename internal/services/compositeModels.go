package services

import (
	"github.com/BlueMonday/go-scryfall"
	generatedEbayListingClient "github.com/m8ypie/mtg-sale-manager/internal/clients/ebayListing"
	"github.com/m8ypie/mtg-sale-manager/internal/models"
)

type CompositeEbayListing struct {
	EbayListing  models.EbayListing
	ScryfallInfo scryfall.Card
	PriceInfo    models.PriceInfo
}

type EbayListingWithOffer struct {
	EbayListing generatedEbayListingClient.InventoryItemWithSkuLocaleGroupid
	EbayOffer   generatedEbayListingClient.EbayOfferDetailsWithAll
}
