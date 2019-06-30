package oauth

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/calumj/ebay-oauth-go-client/credentials"
	"github.com/calumj/ebay-oauth-go-client/environment"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

//API - oAuth api helper class
type API struct {
	appAccessMap sync.Map
}

//New - returns API
func New() *API {
	return &API{
		appAccessMap: sync.Map{},
	}
}

//GetApplicationToken - returns a client credentials token
func (a *API) GetApplicationToken(ctx context.Context, environment *environment.Environment, scopes ...string) (*oauth2.Token, error) {

	if token, ok := a.appAccessMap.Load(environment); ok {
		log.Debug("application access token returned from cache")
		return token.(oauth2.TokenSource).Token()
	}

	credentials := credentials.GetCredentials(environment)
	if credentials == nil {
		return nil, fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	config := generateClientCredentialsConfig(credentials, environment, scopes)

	ts := config.TokenSource(ctx)

	a.appAccessMap.Store(environment, ts)

	token, err := ts.Token()
	if err != nil {
		return nil, err
	}

	return token, nil
}

//GetApplicationTokenClient - returns a http client with a refreshing credentials token in the transport
func (a *API) GetApplicationTokenClient(ctx context.Context, environment *environment.Environment, scopes ...string) (*http.Client, error) {
	credentials := credentials.GetCredentials(environment)
	if credentials == nil {
		return nil, fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	config := generateClientCredentialsConfig(credentials, environment, scopes)

	client := config.Client(ctx)

	return client, nil
}

func generateClientCredentialsConfig(credentials *credentials.Credentials, environment *environment.Environment, scopes []string) *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     credentials.AppID,
		ClientSecret: credentials.CertID,
		TokenURL:     environment.GetAPIEndpoint(),
		Scopes:       scopes,
	}
}
