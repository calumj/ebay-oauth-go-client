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
	credentials  credentials.Credentials
}

//New - returns API
func New(credentials credentials.Credentials) *API {
	return &API{
		appAccessMap: sync.Map{},
		credentials:  credentials,
	}
}

type applicationTokenSource struct {
	source oauth2.TokenSource
}

func (a *applicationTokenSource) Token() (*oauth2.Token, error) {

	token, err := a.source.Token()

	if err != nil {
		return nil, err
	}

	token.TokenType = ""

	return token, nil

}

//GetApplicationTokenAndClient - returns a client credentials token and authenticated http client
func (a *API) GetApplicationTokenAndClient(ctx context.Context, environment *environment.Environment, scopes ...string) (*oauth2.Token, *http.Client, error) {

	if ts, ok := a.appAccessMap.Load(environment); ok {
		log.Debug("application access token returned from cache")
		tokenSource := ts.(oauth2.TokenSource)
		token, err := tokenSource.Token()
		if err != nil {
			return nil, nil, fmt.Errorf("unable to generate token for %s environment: %s", environment.GetConfigIdentifier(), err.Error())
		}

		return token, oauth2.NewClient(ctx, tokenSource), nil
	}

	credentials := a.credentials.GetCredentials(environment)
	if credentials == nil {
		return nil, nil, fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	config := generateClientCredentialsConfig(credentials, environment, scopes)

	tokenSource := &applicationTokenSource{source: config.TokenSource(ctx)}

	a.appAccessMap.Store(environment, tokenSource)

	token, err := tokenSource.Token()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to generate token for %s environment: %s", environment.GetConfigIdentifier(), err.Error())
	}

	return token, oauth2.NewClient(ctx, tokenSource), nil
}

//GenerateUserAuthorizationURL - returns authorization URL to send user to
func (a *API) GenerateUserAuthorizationURL(ctx context.Context, environment *environment.Environment, state string, scopes ...string) (string, error) {
	credentials := a.credentials.GetCredentials(environment)
	if credentials == nil {
		return "", fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	config := generateAuthConfig(credentials, environment, scopes)

	return config.AuthCodeURL(state), nil

}

//ExchangeCodeForAccessTokenAndClient - exchange access token for oAuth token and authenticated http client
func (a *API) ExchangeCodeForAccessTokenAndClient(ctx context.Context, environment *environment.Environment, code string) (*oauth2.Token, *http.Client, error) {
	credentials := a.credentials.GetCredentials(environment)
	if credentials == nil {
		return nil, nil, fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	config := generateAuthConfig(credentials, environment, nil)
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to generate token for %s environment: %s", environment.GetConfigIdentifier(), err.Error())
	}

	tokenSource := config.TokenSource(ctx, token)

	return token, oauth2.NewClient(ctx, tokenSource), nil

}

//GetAccessTokenAndClient - get token and authenticated http client with access token from refresh token
func (a *API) GetAccessTokenAndClient(ctx context.Context, environment *environment.Environment, refreshToken string, scopes ...string) (*oauth2.Token, *http.Client, error) {
	credentials := a.credentials.GetCredentials(environment)
	if credentials == nil {
		return nil, nil, fmt.Errorf("unable to retrieve credentials for %s environment", environment.GetConfigIdentifier())
	}

	config := generateAuthConfig(credentials, environment, scopes)

	token := &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now(),
	}

	tokenSource := config.TokenSource(ctx, token)

	token, err := tokenSource.Token()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to generate token for %s environment: %s", environment.GetConfigIdentifier(), err.Error())
	}

	return token, oauth2.NewClient(ctx, tokenSource), nil

}

func generateClientCredentialsConfig(credentials *credentials.Credential, environment *environment.Environment, scopes []string) *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     credentials.AppID,
		ClientSecret: credentials.CertID,
		TokenURL:     environment.GetAPIEndpoint(),
		Scopes:       scopes,
	}
}

func generateAuthConfig(credentials *credentials.Credential, environment *environment.Environment, scopes []string) *oauth2.Config {
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
