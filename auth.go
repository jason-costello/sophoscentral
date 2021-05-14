package sophoscentral

import (
	"context"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func GetAuthToken(ctx context.Context, clientID, clientSecret, tokenURL string) (*oauth2.Token, error) {
	config := &clientcredentials.Config{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		Scopes:         []string{"token"},
		TokenURL:       tokenURL,
		EndpointParams: url.Values{},
	}
	return config.TokenSource(ctx).Token()
}
