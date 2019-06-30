package environment

import "strings"

//Environment - struct containing environment information for authentication
type Environment struct {
	configIdentifier string
	webEndpoint      string
	apiEndpoint      string
}

//GetWebEndpoint - return the web endpoint
func (e *Environment) GetWebEndpoint() string {
	return e.webEndpoint
}

//GetAPIEndpoint - return the api endpoint
func (e *Environment) GetAPIEndpoint() string {
	return e.apiEndpoint
}

//GetConfigIdentifier - return the config identifier
func (e *Environment) GetConfigIdentifier() string {
	return e.configIdentifier
}

var (
	//PRODUCTION - eBay production environment
	PRODUCTION = &Environment{configIdentifier: "api.ebay.com", webEndpoint: "https://auth.ebay.com/oauth2/authorize", apiEndpoint: "https://api.ebay.com/identity/v1/oauth2/token"}
	//SANDBOX - eBay sandbox environment
	SANDBOX = &Environment{configIdentifier: "api.sandbox.ebay.com", webEndpoint: "https://auth.sandbox.ebay.com/oauth2/authorize", apiEndpoint: "https://api.sandbox.ebay.com/identity/v1/oauth2/token"}

	environmentsByIdentiier = make(map[string]*Environment, 2)
)

func init() {
	environmentsByIdentiier[PRODUCTION.configIdentifier] = PRODUCTION
	environmentsByIdentiier[SANDBOX.configIdentifier] = SANDBOX
}

// LookupBy - Find environment by config identifier
func LookupBy(configIdentifier string) *Environment {
	return environmentsByIdentiier[strings.ToLower(configIdentifier)]

}
