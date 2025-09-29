package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/m8ypie/mtg-sale-manager/internal/clients"
	generatedClients "github.com/m8ypie/mtg-sale-manager/internal/clients/ebayBrowse"
	generatedEbayListingClient "github.com/m8ypie/mtg-sale-manager/internal/clients/ebayListing"
	"github.com/m8ypie/mtg-sale-manager/internal/db"
	"github.com/m8ypie/mtg-sale-manager/internal/models"
	"github.com/m8ypie/mtg-sale-manager/internal/repositories"
)

func GetCompositeEbayListing(ebayListing *models.EbayListing) CompositeEbayListing {
	scryfallCard, err := clients.NewScryfallClient().GetCardDetails(ebayListing)
	if err != nil {
		// Handle error
		log.Fatalf("Error fetching Scryfall card details: %v", err)
	}

	ebayListingsRes := clients.NewEbayCardClient(clients.NewEbayCardClientHttpClientProdDefault()).GetPublicListingsForCard(ebayListing)

	return CompositeEbayListing{
		EbayListing:  *ebayListing,
		ScryfallInfo: *scryfallCard,
		PriceInfo:    getPriceInfo(ebayListingsRes),
	}
}

func RecordNewEbayListingWithOffer(ebayListingWithOffer *EbayListingWithOffer) {
	ctx := context.Background()
	listing := ebayListingWithOffer.ToListing()
	repositories.NewGlobalRepository(db.GetDb()).EbayListing.Upsert(ctx, listing)
}

func GetAllListingsWithUnderCut() {
	ctx := context.Background()
	allListings := repositories.NewGlobalRepository(db.GetDb()).EbayListing.GetAll(ctx)
	for _, listing := range allListings {
		println("running ", listing.EbayTitle)
		comp := GetCompositeEbayListing(listing)
		offer := GetLowestOffer(&listing.Sku)
		if offer != nil && offer.PricingSummary != nil && offer.PricingSummary.Price != nil {
			offerPrice := convertStringToFloat(offer.PricingSummary.Price.Value)
			if comp.PriceInfo.Min < offerPrice {
				println("Listing undercut detected", "Title:", comp.ScryfallInfo.Name, "Lowest Price:", fmt.Sprintf("%.2f", comp.PriceInfo.Min), "My Price:", fmt.Sprintf("%.2f", offerPrice))
			} else {
				println("No undercut", "Title:", comp.ScryfallInfo.Name, "Lowest Price:", fmt.Sprintf("%.2f", comp.PriceInfo.Min), "My Price:", fmt.Sprintf("%.2f", offerPrice))
			}

		}
	}
}

func GetAllListingsWithUnderCutNoDb() []*CompositeEbayListing {
	allListings := GetEbayListingsWithOffer()
	var compositeEbayListing []*CompositeEbayListing
	for _, listing := range *allListings {
		println("running ", listing.EbayListing.Product.Title)
		comp := GetCompositeEbayListing(listing.ToListing())
		offer := GetLowestOffer(listing.EbayListing.Sku)
		if offer != nil && offer.PricingSummary != nil && offer.PricingSummary.Price != nil {
			offerPrice := convertStringToFloat(offer.PricingSummary.Price.Value)
			comp.MyPrice = &offerPrice
			comp.CompetitionsPrice = &comp.PriceInfo.Min
			compositeEbayListing = append(compositeEbayListing, &comp)
		}
	}
	return compositeEbayListing
}

func GetEbayListingsWithOffer() *[]EbayListingWithOffer {

	ebayRes, err := clients.NewEbayCardClient(clients.NewEbayCardClientHttpClientProdDefault()).GetActiveListings()

	if err != nil {
		// Handle error
		log.Fatalf("Error fetching Scryfall card details: %v", err)
	}
	var listingsWithOffers []EbayListingWithOffer
	for _, item := range *ebayRes.JSON200.InventoryItems {
		offers, err := clients.NewEbayCardClient(clients.NewEbayCardClientHttpClientProdDefault()).GetLiveOffers(item.Sku)
		if err != nil {
			log.Printf("Error fetching offer for SKU %s: %v", item.Sku, err)
			continue
		}
		if len(*offers) == 0 {
			log.Printf("No offers found for SKU %s", item.Sku)
			continue
		}
		listingsWithOffers = append(listingsWithOffers, EbayListingWithOffer{
			EbayListing: item,
			EbayOffer:   (*offers)[0],
		})
	}
	return &listingsWithOffers
}

func GetLowestOffer(sku *string) *generatedEbayListingClient.EbayOfferDetailsWithAll {
	offers, err := clients.NewEbayCardClient(clients.NewEbayCardClientHttpClientProdDefault()).GetLiveOffers(sku)
	if err != nil {
		log.Printf("Error fetching offer for SKU %s: %v", sku, err)
		return nil
	}
	if len(*offers) == 0 {
		log.Printf("No offers found for SKU %s", sku)
		return nil
	}
	return &(*offers)[0]
}

func getPriceInfo(ebayListingsRes *generatedClients.SearchResponse) models.PriceInfo {
	min := math.MaxFloat64
	max := 0.
	totalPrice := 0.
	count := float64(*ebayListingsRes.JSON200.Total)
	if count > 0 {
		for _, priceSum := range *ebayListingsRes.JSON200.ItemSummaries {
			min = math.Min(min, convertStringToFloat(priceSum.Price.Value))
			max = math.Max(max, convertStringToFloat(priceSum.Price.Value))
			totalPrice += convertStringToFloat(priceSum.Price.Value)
		}
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
