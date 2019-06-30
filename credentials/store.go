package credentials

import (
	"os"

	"github.com/calumj/ebay-oauth-go-client/environment"
	"github.com/go-yaml/yaml"
)

var envCredentialsMap = make(map[*environment.Environment]*Credentials, 2)

//Credentials - container for oAuth credentials
type Credentials struct {
	AppID       string `yaml:"appid"`
	CertID      string `yaml:"certid"`
	DevID       string `yaml:"devid"`
	RedirectURI string `yaml:"redirecturi"`
}

//Load - Loads credenti	"github.com/calumj/ebay-oauth-go-client/environment"als YAML file in to memory
func Load(file *os.File) error {

	decoder := yaml.NewDecoder(file)

	wrapper := &struct {
		Sandbox    *Credentials `yaml:"api.sandbox.ebay.com"`
		Production *Credentials `yaml:"api.ebay.com"`
	}{}

	err := decoder.Decode(&wrapper)
	if err != nil {
		return err
	}

	envCredentialsMap[environment.PRODUCTION] = wrapper.Production
	envCredentialsMap[environment.SANDBOX] = wrapper.Sandbox

	return nil

}

//GetCredentials - get oAuth information for specified environment
func GetCredentials(env *environment.Environment) *Credentials {
	return envCredentialsMap[env]
}
