package clients

import (
	"context"
	"log"
	"net/http"

	generatedClients "github.com/m8ypie/mtg-sale-manager/internal/clients/ebayBrowse"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
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
func (ts *BearerTokenSource) Token() (*oauth2.Token, error) {
	t, err := ts.base.Token()
	if t != nil {
		t.TokenType = "Bearer"
	}
	return t, err
}

func NewEbayListingsFetcherConfigProdDefault() clientcredentials.Config {
	return clientcredentials.Config{
		ClientID:     "your client id",
		ClientSecret: "your client secret",
		TokenURL:     OAuth20Endpoint.TokenURL,
		Scopes:       []string{ebay.ScopeRoot /* your scopes */},
	}
}

func NewEbayListingsFetcherHttpClientProdDefault() *http.Client {
	config := NewEbayListingsFetcherConfigProdDefault()
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ebay.TokenSource(cfg.TokenSource(ctx)))
	c := &http.Client{Transport: tc, BaseURL: "https://api.ebay.com/"}

	return c
}

type EbayListingsFetcher struct {
	ClientWithResponses *generatedClients.ClientWithResponses
}

func NewEbayListingsFetcher(httpClient *http.Client) *EbayListingsFetcher {
	c, err := generatedClients.NewClientWithResponses("https://api.ebay.com/", generatedClients.WithHTTPClient(httpClient))
	if err != nil {
		log.Fatal(err)
	}
	return &EbayListingsFetcher{
		ClientWithResponses: c,
	}
}

func (e *EbayListingsFetcher) getListings(itemID string) (*generatedClients.GetItemsResponse, error) {
	return e.ClientWithResponses.GetItemsWithResponse(context.Background(), &generatedClients.GetItemsParams{
		ItemIds: &itemID,
	})
}
