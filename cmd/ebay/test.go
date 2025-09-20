package main

import (
	"context"
	"log"
	"os"

	"github.com/m8ypie/mtg-sale-manager/internal/config"
	"github.com/m8ypie/mtg-sale-manager/internal/db"
	"github.com/m8ypie/mtg-sale-manager/internal/models"
	"github.com/m8ypie/mtg-sale-manager/internal/repositories"

	"github.com/m8ypie/mtg-sale-manager/internal/services"
)

func main() {
	config.Init()

	// res, err := clients.NewEbayListingsFetcher(clients.NewEbayListingsFetcherHttpClientProdDefault()).GetListings("MTG The One Ring 246")
	// if err != nil {
	// 	log.Fatalf("Error fetching eBay listings: %v", err)
	// }

	// fmt.Printf("Content Length: %d\n", res.HTTPResponse.ContentLength)
	// reqBodyString := string(res.Body)
	// fmt.Printf("resp.JSON200: %v\n", reqBodyString)
	// writeStringToFile("ebay_response.json", reqBodyString)
	// ctx := context.Background()
	// cardListing, error := repositories.NewGlobalRepository(db.GetDb()).EbayListing.FindByID(ctx, 1)
	// if error != nil {
	// 	log.Fatalf("Error fetching ebay listing: %v", error)
	// }
	// card, errorscry := clients.NewScryfallClient().GetCardDetails(&cardListing)
	// if errorscry != nil {
	// 	log.Fatalf("Error fetching card details: %v", errorscry)
	// }
	// log.Printf("Card Details: %+v\n", card)
	// popDb()
	newMe()
}

func newMe() {
	ctx := context.Background()
	cardListing, error := repositories.NewGlobalRepository(db.GetDb()).EbayListing.FindByID(ctx, 5)
	if error != nil {
		log.Fatalf("Error fetching ebay listing: %v", error)
	}
	comp := services.GetCompositeEbayListing(&cardListing)
	log.Printf("comp %v", comp)
}

func popDb() {
	sku := "d5806e68-1054-458e-866d-1f2470f682b2"
	offerId := "70191510011"
	listingId := "187575984289"
	emid := "150075"
	scryFallId := "d5806e68-1054-458e-866d-1f2470f682b2"
	ebayListing := models.EbayListing{
		Sku:                sku,
		EbayOfferId:        offerId,
		EbayListingId:      listingId,
		EchoMtgInventoryId: emid,
		ScryFallId:         scryFallId,
	}
	ctx := context.Background()
	error := repositories.NewGlobalRepository(db.GetDb()).EbayListing.Insert(ctx, &ebayListing)
	if error != nil {
		log.Fatalf("Error inserting ebay listing: %v", error)
	}
	listings, err := repositories.NewGlobalRepository(db.GetDb()).EbayListing.FindAll(ctx)
	if err != nil {
		log.Fatalf("Error fetching ebay listings: %v", err)
	}
	for _, listing := range listings {
		log.Printf("Ebay Listing: %+v\n", listing)
	}
}

func writeStringToFile(filename, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}
