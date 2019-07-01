# eBay OAuth Client Library (Go)
This is a loose implementation of the [ebay-oauth-java-client](https://github.com/eBay/ebay-oauth-java-client) for Go

# How To Install
Current Version : 0.2

Download code using `go get github.com/calumj/ebay-oauth-go-client`

# How To Use 
Initialise the default YAML credentials using `creds, err := credentials.NewYAML(file)`

Once the credentials have been initalised, all further interactions can be done using the API struct, using `api := oauth.New(creds)`


## Example Client Credentials
```
file, err := os.Open("./ebay-config.yaml")
if err != nil {
	log.Fatal(err)
}

creds, err := credentials.NewYAML(file)
if err != nil {
	log.Fatal(err)
}

api := oauth.New(creds)

token, client, err := api.GetApplicationTokenAndClient(context.Background(), environment.SANDBOX, "https://api.ebay.com/oauth/api_scope")
if err != nil {
	log.Fatal(err)
}
```
