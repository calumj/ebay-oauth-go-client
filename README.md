# eBay OAuth Client Library (Go)
This is an implementation of the [ebay-oauth-java-client](https://github.com/eBay/ebay-oauth-java-client) for Go

# How To Install
Current Version : 0.1
Download code using `go get github.com/calumj/ebay-oauth-go-client`

# How To Use 
Initialise the credentials using `credentials.Load(<io.File>)`
Once the credentials have been initalised, all further interactions can be done using the API struct,  using `api := oauth.New()`

## Example Client Credentials
```
file, err := os.Open("./ebay-config.yaml")
if err != nil {
	log.Fatal(err)
}

err = credentials.Load(file)
if err != nil {
	log.Fatal(err)
}

api := oauth.New()

token, err := api.GetApplicationToken(context.Background(), environment.SANDBOX, "https://api.ebay.com/oauth/api_scope")
if err != nil {
	log.Fatal(err)
}
```
