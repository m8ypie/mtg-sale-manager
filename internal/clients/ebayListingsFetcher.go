package clients

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	ebayBrowseClient "github.com/m8ypie/mtg-sale-manager/internal/clients/ebayBrowse"
	ebayListingClient "github.com/m8ypie/mtg-sale-manager/internal/clients/ebayListing"
	"github.com/m8ypie/mtg-sale-manager/internal/config"
	"github.com/m8ypie/mtg-sale-manager/internal/models"
	"golang.org/x/oauth2"
)

// eBay OAuth 2.0 endpoints.
var (
	OAuth20Endpoint = oauth2.Endpoint{
		AuthURL:  "https://auth.ebay.com/oauth2/authorize",
		TokenURL: "https://api.ebay.com/identity/v1/oauth2/token",
	}
	OAuth20SandboxEndpoint = oauth2.Endpoint{
		AuthURL:  "https://auth.sandbox.ebay.com/oauth2/authorize",
		TokenURL: "https://api.sandbox.ebay.com/identity/v1/oauth2/token",
	}
)

// BearerTokenSource forces the type of the token returned by the 'base' TokenSource to 'Bearer'.
// The eBay API will return "Application Access Token" or "User Access Token" as token_type but
// it must be set to 'Bearer' for subsequent requests.
type BearerTokenSource struct {
	base oauth2.TokenSource
}

// TokenSource returns a new BearerTokenSource.
func TokenSource(base oauth2.TokenSource) *BearerTokenSource {
	return &BearerTokenSource{base: base}
}

// Token allows BearerTokenSource to implement oauth2.TokenSource.
// func (ts *BearerTokenSource) Token() (*oauth2.Token, error) {
// 	t, err := ts.base.Token()
// 	if t != nil {
// 		t.TokenType = "Bearer"
// 	}
// 	return t, err
// }

func NewEbayCardClientsConfigProdDefault() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.EbayAppId,
		ClientSecret: config.EbayCertId,
		RedirectURL:  "Joel_Berta-JoelBert-mtglis-lhxiiqgo",
		Endpoint:     OAuth20Endpoint,

		Scopes: []string{"https://api.ebay.com/oauth/api_scope/commerce.identity.readonly",
			"https://api.ebay.com/oauth/api_scope/sell.stores.readonly",
			"https://api.ebay.com/oauth/api_scope",
			"https://api.ebay.com/oauth/api_scope/sell.marketing.readonly",
			"https://api.ebay.com/oauth/api_scope/sell.marketing",
			"https://api.ebay.com/oauth/api_scope/sell.inventory.readonly",
			"https://api.ebay.com/oauth/api_scope/sell.inventory",
			"https://api.ebay.com/oauth/api_scope/sell.account.readonly",
			"https://api.ebay.com/oauth/api_scope/sell.account",
			"https://api.ebay.com/oauth/api_scope/sell.fulfillment.readonly",
			"https://api.ebay.com/oauth/api_scope/sell.fulfillment"},
	}
}

type EbayListingsFetcherHttpClient struct {
	Client *http.Client

	TokenSource oauth2.TokenSource

	OAuth2Tokens *oauth2.Token
}

func NewEbayCardClientHttpClientProdDefault() *EbayListingsFetcherHttpClient {
	defConfig := NewEbayCardClientsConfigProdDefault()

	ctx := context.Background()

	tokenSource := defConfig.TokenSource(ctx, &oauth2.Token{
		RefreshToken: config.EbayRefreshToken,
		TokenType:    "bearer",
		Expiry:       time.Now(),
	})

	oAuth2Tokens, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("Error retrieving access token: %v", err)
	}

	tc := oauth2.NewClient(ctx, tokenSource)

	return &EbayListingsFetcherHttpClient{
		Client:       tc,
		TokenSource:  tokenSource,
		OAuth2Tokens: oAuth2Tokens,
	}
}

type EbayCardClient struct {
	EbayBrowseClient  *ebayBrowseClient.ClientWithResponses
	EbayListingClient *ebayListingClient.ClientWithResponses
}

func SetHeaders(httpClient *EbayListingsFetcherHttpClient) func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		if httpClient.OAuth2Tokens.Expiry.After(time.Now()) {
			httpClient.OAuth2Tokens, _ = httpClient.TokenSource.Token()
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Language", "en-AU")
		req.Header.Set("Accept-Language", "en-AU")
		req.Header.Set("X-EBAY-C-MARKETPLACE-ID", "EBAY_AU")
		req.Header.Set("Authorization", "Bearer "+httpClient.OAuth2Tokens.AccessToken)

		return nil
	}
}

func NewEbayCardClient(httpClient *EbayListingsFetcherHttpClient) *EbayCardClient {
	browseBaseClient := ebayBrowseClient.WithRequestEditorFn(SetHeaders(httpClient))
	listingBaseClient := ebayListingClient.WithRequestEditorFn(SetHeaders(httpClient))
	browseClient, err := ebayBrowseClient.NewClientWithResponses("https://api.ebay.com/buy/browse/v1", browseBaseClient)
	listingClient, err2 := ebayListingClient.NewClientWithResponses("https://api.ebay.com/sell/inventory/v1", listingBaseClient)
	if err != nil {
		log.Fatal(err)
	}

	if err2 != nil {
		log.Fatal(err2)
	}

	return &EbayCardClient{
		EbayBrowseClient:  browseClient,
		EbayListingClient: listingClient,
	}
}

func (e *EbayCardClient) GetListings(query string) (*ebayBrowseClient.SearchResponse, error) {
	filter := "itemLocationCountry:AU"
	return e.EbayBrowseClient.SearchWithResponse(context.Background(), &ebayBrowseClient.SearchParams{
		Q:      &query,
		Filter: &filter,
	})
}

func (e *EbayCardClient) GetActiveListings() (*ebayListingClient.GetInventoryItemsResponse, error) {
	return e.EbayListingClient.GetInventoryItemsWithResponse(context.Background(), &ebayListingClient.GetInventoryItemsParams{})
}

func (e *EbayCardClient) GetLiveOffers(sku *string) (*[]ebayListingClient.EbayOfferDetailsWithAll, error) {
	resp, err := e.EbayListingClient.GetOffersWithResponse(context.Background(), &ebayListingClient.GetOffersParams{
		Sku: sku,
	})
	if err != nil {
		return nil, err
	}
	var filtered []ebayListingClient.EbayOfferDetailsWithAll
	for _, offer := range *resp.JSON200.Offers {
		if offer.Listing.ListingStatus != nil && *offer.Listing.ListingStatus == "ACTIVE" {
			filtered = append(filtered, offer)
		}
	}
	return &filtered, nil
}
func (e *EbayCardClient) GetPublicListingsForCard(ebayListing *models.EbayListing) *ebayBrowseClient.SearchResponse {
	filter := "itemLocationCountry:AU"
	query := removePunctuation(ebayListing.EbayTitle)
	res, err := e.EbayBrowseClient.SearchWithResponse(context.Background(), &ebayBrowseClient.SearchParams{
		Q:      &query,
		Filter: &filter,
	})
	if err != nil {
		log.Fatalf("Error fetching eBay listings: %v", err)
	}

	return res
}

func (e *EbayCardClient) GetOwnedOfferForCard(ebayListing *models.EbayListing) *ebayListingClient.GetOfferResponse {
	res, err := e.EbayListingClient.GetOfferWithResponse(context.Background(), ebayListing.EbayOfferId)
	if err != nil {
		log.Fatalf("Error fetching eBay listings: %v", err)
	}
	return res
}

func removePunctuation(s string) string {
	// Define a regular expression to match common punctuation characters.
	// You can customize this pattern based on which specific punctuation
	// you want to remove.
	// `[^a-zA-Z0-9\s]` matches any character that is not an alphanumeric character or whitespace.
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	return strings.Join(strings.Fields(re.ReplaceAllString(s, "")), " ")
}
