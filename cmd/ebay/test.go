package main

import (
	"fmt"
	"log"
	"os"

	"github.com/m8ypie/mtg-sale-manager/internal/clients"
	"github.com/m8ypie/mtg-sale-manager/internal/config"
)

func main() {
	config.Init()
	res, err := clients.NewEbayListingsFetcher(clients.NewEbayListingsFetcherHttpClientProdDefault()).GetListings("MTG The One Ring 246")
	if err != nil {
		log.Fatalf("Error fetching eBay listings: %v", err)
	}

	fmt.Printf("Content Length: %d\n", res.HTTPResponse.ContentLength)
	reqBodyString := string(res.Body)
	fmt.Printf("resp.JSON200: %v\n", reqBodyString)
	writeStringToFile("ebay_response.json", reqBodyString)

}

func writeStringToFile(filename, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}
