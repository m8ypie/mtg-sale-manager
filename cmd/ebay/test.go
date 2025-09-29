package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/m8ypie/mtg-sale-manager/internal/config"

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
	re()
	services.GetAllListingsWithUnderCut()

}

func re() {
	listings := services.GetEbayListingsWithOffer()
	for _, listing := range *listings {
		services.RecordNewEbayListingWithOffer(&listing)
		if listing.EbayListing.Product != nil && listing.EbayListing.Product.Title != nil {
			writeJsonToFile(*listing.EbayListing.Product.Title+".json", listing)
		} else {
			log.Printf("Warning: listing.Product.Title is nil for listing: %+v\n", listing)
		}
	}
	//writeStringToFile("ebay_response.json", string(listings.))
}

func writeStringToFile(filename, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

func writeJsonToFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
