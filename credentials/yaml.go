package credentials

import (
	"os"

	"github.com/calumj/ebay-oauth-go-client/environment"
	"github.com/go-yaml/yaml"
)

//YAML - yaml credential implemetation
type YAML struct {
	envCredentialsMap map[*environment.Environment]*Credential
}

//NewYAML - Creates struct with credentials YAML file loaded in to memory
func NewYAML(file *os.File) (Credentials, error) {
	decoder := yaml.NewDecoder(file)

	wrapper := &struct {
		Sandbox    *Credential `yaml:"api.sandbox.ebay.com"`
		Production *Credential `yaml:"api.ebay.com"`
	}{}

	err := decoder.Decode(&wrapper)
	if err != nil {
		return nil, err
	}

	impl := &YAML{
		envCredentialsMap: make(map[*environment.Environment]*Credential, 2),
	}

	impl.envCredentialsMap[environment.PRODUCTION] = wrapper.Production
	impl.envCredentialsMap[environment.SANDBOX] = wrapper.Sandbox

	return impl, nil
}

//GetCredentials - get oAuth information for specified environment
func (y *YAML) GetCredentials(env *environment.Environment) *Credential {
	return y.envCredentialsMap[env]
}
