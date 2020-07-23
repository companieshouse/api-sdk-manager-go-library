package manager

import (
	"encoding/json"
	"net/http"

	"github.com/companieshouse/api-sdk-go/apikey"
	sdk "github.com/companieshouse/api-sdk-go/companieshouseapi"
	choauth2 "github.com/companieshouse/api-sdk-go/oauth2"
	privatesdk "github.com/companieshouse/private-api-sdk-go/companieshouseapi"
	"github.com/pkg/errors"
	goauth2 "golang.org/x/oauth2"
)

var sdkBasePathOverridden = false
var privateSdkBasePathOverridden = false

type SDKManager struct {
	gofigure        interface{} `order:"env,flag"`
	ClientID        string      `env:"OAUTH2_CLIENT_ID"     flag:"oauth2-client-id"     flagDesc:"Client ID"`
	ClientSecret    string      `env:"OAUTH2_CLIENT_SECRET" flag:"oauth2-client-secret" flagDesc:"Client Secret"`
	RedirectURL     string      `env:"OAUTH2_REDIRECT_URI"  flag:"oauth2-redirect-uri"  flagDesc:"Oauth2 Redirect URI"`
	AuthURL         string      `env:"OAUTH2_AUTH_URI"      flag:"oauth2-auth-uri"      flagDesc:"Oauth2 Auth URI"`
	TokenURL        string      `env:"OAUTH2_TOKEN_URI"     flag:"oauth2-token-uri"     flagDesc:"Oauth2 Token URI"`
	Scopes          []string    `env:"SCOPE"                flag:"scope"                flagDesc:"Scope"`
	APIKey          string      `env:"API_KEY"              flag:"api-key"              flagDesc:"API Key"`
	APIURL          string      `env:"API_URL"              flag:"api-url"              flagDesc:"API URL"`
	PostcodeService string      `env:"POSTCODE_SERVICE"     flag:"postcode-service"     flagDesc:"Postcode Service"`
}

var sdkManager *SDKManager
var oauthConfig *choauth2.Config

func Get() (*SDKManager) {

	if sdkManager != nil {
		return sdkManager
	}

	sdkManager = &SDKManager{}
	//err := gofigure.Gofigure(config)
	return sdkManager//, err
}

// GetOauthConfig returns an instance of a Companies House oauth config struct
// and an error wheere appropriate
func GetOauthConfig() (*choauth2.Config, error) {

	if oauthConfig != nil {
		return oauthConfig, nil
	}

	cfg := Get()

	oauthConfig = &choauth2.Config{}
	oauthConfig.ClientID = cfg.ClientID
	oauthConfig.ClientSecret = cfg.ClientSecret
	oauthConfig.RedirectURL = cfg.RedirectURL
	oauthConfig.Scopes = cfg.Scopes
	oauthConfig.Endpoint = goauth2.Endpoint{
		AuthURL:  cfg.AuthURL,
		TokenURL: cfg.TokenURL,
	}

	return oauthConfig, nil
}

// GetSDK will return an instance of the Go SDK using an oauth2 authenticated
// HTTP client if requested, else an API-key authenticated HTTP client will be used
func (manager SDKManager) GetSDK(req *http.Request, usePassthrough bool, APIKey string) (*sdk.Service, error) {

	//cfg, err := config.Get()
	//if err != nil {
	//	return nil, err
	//}
	cfg := Get()
	//cfg.APIKey = APIKey

	// Override sdkBasePath here to route API requests via ERIC
	if !sdkBasePathOverridden && len(cfg.APIURL) > 0 {
		sdk.BasePath = cfg.APIURL
		sdkBasePathOverridden = true
	}

	httpClient, err := getHTTPClient(req, usePassthrough, APIKey)
	if err != nil {
		return nil, err
	}

	return sdk.New(httpClient)
}

// GetPrivateSDK will return an instance of the Private Go SDK using an oauth2 authenticated
// HTTP client if requested, else an API-key authenticated HTTP client will be used
func (manager SDKManager) GetPrivateSDK(req *http.Request, usePassthrough bool, APIKey string) (*privatesdk.Service, error) {

	//cfg, err := config.Get()
	//if err != nil {
	//	return nil, err
	//}

	cfg := Get()

	// Override privateSdkBasePath here to route API requests via ERIC
	if !privateSdkBasePathOverridden && len(cfg.APIURL) > 0 {
		privatesdk.BasePath = cfg.APIURL
		privatesdk.PostcodeBasePath = cfg.PostcodeService
		privateSdkBasePathOverridden = true
	}

	httpClient, err := getHTTPClient(req, usePassthrough, APIKey)
	if err != nil {
		return nil, err
	}

	return privatesdk.New(httpClient)
}

// getHTTPClient returns an Http Client. It will be either Oauth2 or API-key
// authenticated depending on whether the calling service has requested to use the passthrough token
func getHTTPClient(req *http.Request, usePassthrough bool, APIKey string) (*http.Client, error) {
	var httpClient *http.Client
	var err error

	// If passthrough token is preferred, get the passthrough token and get an HTTP client
	if usePassthrough {
		// Check the token exists
		decodedPassthroughToken, err := decodePassthroughHeader(req)
		if err != nil {
			return nil, err
		}
		// If it exists, we'll use it to return an authenticated HTTP client
		httpClient, err = getOauth2HTTPClient(req, decodedPassthroughToken)
	} else {
		// Otherwise, we'll use API-key authentication
		httpClient, err = getAPIKeyHTTPClient(req, APIKey)
	}

	if err != nil {
		return nil, err
	}

	return httpClient, nil
}

//Returns a decoded passthrough token or nil if no token present
func decodePassthroughHeader(req *http.Request) (*goauth2.Token, error) {

	passthroughHeader := req.Header.Get("Eric-Access-Token")

	if passthroughHeader != "" {

		decodedPassthrough := &goauth2.Token{}
		err := json.Unmarshal([]byte(passthroughHeader), decodedPassthrough)
		if err != nil {
			return nil, err
		}

		return decodedPassthrough, nil
	}

	return nil, nil

}

// getAPIKeyHttpClient returns an API-key-authenticated HTTP client
func getAPIKeyHTTPClient(req *http.Request, APIKey string) (*http.Client, error) {

	//cfg, err := config.Get()
	//if err != nil {
	//	return nil, err
	//}

	cfg := Get()
	cfg.APIKey = APIKey

	// Initialise an apikey cfg struct
	apiKeyConfig := &apikey.Config{Key: cfg.APIKey}

	// Create an http client
	return apiKeyConfig.Client(req.Context(), cfg.APIKey), nil
}

// getOauth2HttpClient returns an Oauth2-authenticated HTTP client
func getOauth2HTTPClient(req *http.Request, tok *goauth2.Token) (*http.Client, error) {

	// Fetch oauth config
	oauth2Config, err := GetOauthConfig()
	if err != nil {
		return nil, err
	}

	// Initialise the callback function to be fired on session expiry
	var fn choauth2.NotifyFunc = AccessTokenChangedCallback

	// Create an http client
	return oauth2Config.Client(req.Context(), tok, fn, ""), nil
}

// AccessTokenChangedCallback is the callback to get a new access token
// As there is no session, a new token should never be acquired
func AccessTokenChangedCallback(newToken *goauth2.Token, private interface{}) error {
	return errors.New("Token expired")
}
