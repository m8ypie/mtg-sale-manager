package services

import (
	"fmt"
	"strconv"

	"github.com/BlueMonday/go-scryfall"
	generatedEbayListingClient "github.com/m8ypie/mtg-sale-manager/internal/clients/ebayListing"
	"github.com/m8ypie/mtg-sale-manager/internal/models"
)

type CompositeEbayListing struct {
	EbayListing       models.EbayListing
	ScryfallInfo      scryfall.Card
	PriceInfo         models.PriceInfo
	MyPrice           *float64
	CompetitionsPrice *float64
}

func countDigits(n float64) int {
	s := strconv.Itoa(int(n))
	if n < 0 {
		return len(s) - 1 // exclude minus sign
	}
	return len(s)
}

func (e CompositeEbayListing) ToString() *string {
	str := fmt.Sprintf("%s | is undercut: %v | my price: %.2f | their price: %.2f",
		e.EbayListing.EbayTitle,
		func() bool {
			if e.MyPrice != nil && e.CompetitionsPrice != nil {
				return *e.MyPrice > *e.CompetitionsPrice
			}
			return false
		}(),
		func() float64 {
			if e.MyPrice != nil {
				return *e.MyPrice
			} else {
				return 0
			}
		}(),
		func() float64 {
			if e.CompetitionsPrice != nil {
				if countDigits(*e.CompetitionsPrice) > 10 {
					return *e.MyPrice
				}
				return *e.CompetitionsPrice
			} else {
				return 0
			}
		}(),
	)
	return &str
}

type EbayListingWithOffer struct {
	EbayListing generatedEbayListingClient.InventoryItemWithSkuLocaleGroupid
	EbayOffer   generatedEbayListingClient.EbayOfferDetailsWithAll
}

func (e EbayListingWithOffer) ToListing() *models.EbayListing {
	return &models.EbayListing{
		EbayTitle:     *e.EbayListing.Product.Title,
		EbayListingId: *e.EbayOffer.Listing.ListingId,
		EbayOfferId:   *e.EbayOffer.OfferId,
		Sku:           *e.EbayListing.Sku,
		ScryFallId:    *e.EbayListing.Sku,
	}
}
