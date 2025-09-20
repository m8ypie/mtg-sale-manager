package services

import (
	"log"
	"math"
	"strconv"

	"github.com/m8ypie/mtg-sale-manager/internal/clients"
	generatedClients "github.com/m8ypie/mtg-sale-manager/internal/clients/ebayBrowse"
	"github.com/m8ypie/mtg-sale-manager/internal/models"
)

func GetCompositeEbayListing(ebayListing *models.EbayListing) CompositeEbayListing {
	scryfallCard, err := clients.NewScryfallClient().GetCardDetails(ebayListing)
	if err != nil {
		// Handle error
		log.Fatalf("Error fetching Scryfall card details: %v", err)
	}

	ebayListingsRes := clients.NewEbayListingsFetcher(clients.NewEbayListingsFetcherHttpClientProdDefault()).GetListingsForCard(scryfallCard)

	return CompositeEbayListing{
		EbayListing:  *ebayListing,
		ScryfallInfo: *scryfallCard,
		PriceInfo:    getPriceInfo(ebayListingsRes),
	}
}

func getPriceInfo(ebayListingsRes *generatedClients.SearchResponse) models.PriceInfo {
	min := math.MaxFloat64
	max := 0.
	totalPrice := 0.
	count := float64(*ebayListingsRes.JSON200.Total)
	for _, priceSum := range *ebayListingsRes.JSON200.ItemSummaries {
		min = math.Min(min, convertStringToFloat(priceSum.Price.Value))
		max = math.Max(max, convertStringToFloat(priceSum.Price.Value))
		totalPrice += convertStringToFloat(priceSum.Price.Value)
	}
	avg := totalPrice / count
	return models.PriceInfo{
		Min:     min,
		Max:     max,
		Average: avg,
	}
}

func convertStringToFloat(priceStr *string) float64 {
	price, err := strconv.ParseFloat(*priceStr, 64)
	if err != nil {
		log.Fatalf("Error converting string to float: %v", err)
	}
	return price
}
