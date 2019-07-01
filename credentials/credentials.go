package credentials

import "github.com/calumj/ebay-oauth-go-client/environment"

//Credential - container for oAuth credentials
type Credential struct {
	AppID       string `yaml:"appid"`
	CertID      string `yaml:"certid"`
	DevID       string `yaml:"devid"`
	RedirectURI string `yaml:"redirecturi"`
}

//Credentials -  Interface for getting credentials
type Credentials interface {
	GetCredentials(env *environment.Environment) *Credential
}
