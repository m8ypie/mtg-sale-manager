package main

import (
	"fmt"

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
	// re()
	for _, re := range services.GetAllListingsWithUnderCutNoDb() {
		fmt.Println(*re.ToString())
	}

}
