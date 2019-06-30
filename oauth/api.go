package oauth

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

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

//GetApplicationClient - returns a http client with a refreshing credentials token in the transport
func (a *API) GetApplicationClient(ctx context.Context, environment *environment.Environment, scopes ...string) (*http.Client, error) {
	credentials := credentials.GetCredentials(environment)
	if credentials == nil {
		return nil, fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	config := generateClientCredentialsConfig(credentials, environment, scopes)

	client := config.Client(ctx)

	return client, nil
}

//GenerateUserAuthorizationURL - returns authorization URL to send user to
func (a *API) GenerateUserAuthorizationURL(ctx context.Context, environment *environment.Environment, state string, scopes ...string) (string, error) {
	credentials := credentials.GetCredentials(environment)
	if credentials == nil {
		return "", fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	config := generateAuthConfig(credentials, environment, scopes)

	return config.AuthCodeURL(state), nil

}

//ExchangeCodeForAccessToken - exchange access token for oAuth token
func (a *API) ExchangeCodeForAccessToken(ctx context.Context, environment *environment.Environment, code string) (*oauth2.Token, error) {
	credentials := credentials.GetCredentials(environment)
	if credentials == nil {
		return nil, fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	config := generateAuthConfig(credentials, environment, nil)
	return config.Exchange(ctx, code)

}

//ExchangeCodeForAccessClient - exchange access token for http client with oAuth token
func (a *API) ExchangeCodeForAccessClient(ctx context.Context, environment *environment.Environment, code string) (*http.Client, error) {
	credentials := credentials.GetCredentials(environment)
	if credentials == nil {
		return nil, fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	token, err := a.ExchangeCodeForAccessToken(ctx, environment, code)
	if err == nil {
		return nil, err
	}
	config := generateAuthConfig(credentials, environment, nil)

	return config.Client(ctx, token), nil

}

//GetAccessToken - get access token from refresh token
func (a *API) GetAccessToken(ctx context.Context, environment *environment.Environment, refreshToken string, scopes ...string) (*oauth2.Token, error) {
	credentials := credentials.GetCredentials(environment)
	if credentials == nil {
		return nil, fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	config := generateAuthConfig(credentials, environment, scopes)

	token := &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now(),
	}
	tokenSource := config.TokenSource(ctx, token)

	return tokenSource.Token()

}

//GetAccessTokenClient - get http client with access token from refresh token
func (a *API) GetAccessTokenClient(ctx context.Context, environment *environment.Environment, refreshToken string, scopes ...string) (*http.Client, error) {
	credentials := credentials.GetCredentials(environment)
	if credentials == nil {
		return nil, fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	config := generateAuthConfig(credentials, environment, scopes)

	token := &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now(),
	}

	return config.Client(ctx, token), nil

}

func generateClientCredentialsConfig(credentials *credentials.Credentials, environment *environment.Environment, scopes []string) *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     credentials.AppID,
		ClientSecret: credentials.CertID,
		TokenURL:     environment.GetAPIEndpoint(),
		Scopes:       scopes,
	}
}

func generateAuthConfig(credentials *credentials.Credentials, environment *environment.Environment, scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     credentials.AppID,
		ClientSecret: credentials.CertID,
		Endpoint: oauth2.Endpoint{
			AuthURL:  environment.GetWebEndpoint(),
			TokenURL: environment.GetAPIEndpoint(),
		},
		RedirectURL: credentials.RedirectURI,
		Scopes:      scopes,
	}

}
