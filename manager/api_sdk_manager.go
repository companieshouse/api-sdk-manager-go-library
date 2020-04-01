package manager

import (
	"encoding/json"
	"net/http"

	"github.com/companieshouse/api-sdk-go/apikey"
	sdk "github.com/companieshouse/api-sdk-go/companieshouseapi"
	choauth2 "github.com/companieshouse/api-sdk-go/oauth2"
	"github.com/companieshouse/api-sdk-manager-go-library/config"
	"github.com/pkg/errors"
	goauth2 "golang.org/x/oauth2"
)

var basePathOverridden = false

// GetSDK will return an instance of the Go SDK using an oauth2 authenticated
// HTTP client if possible, else an API-key authenticated HTTP client will be used
func GetSDK(req *http.Request, usePassthrough bool) (*sdk.Service, error) {

	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	// Override BasePath here to route API requests via ERIC
	if !basePathOverridden && len(cfg.APIURL) > 0 {
		sdk.BasePath = cfg.APIURL
		basePathOverridden = true
	}

	httpClient, err := getHTTPClient(req, usePassthrough)
	if err != nil {
		return nil, err
	}

	return sdk.New(httpClient)
}

// getHTTPClient returns an Http Client. It will be either Oauth2 or API-key
// authenticated depending on whether an Oauth token can be procured from the
// passthrough token
func getHTTPClient(req *http.Request, usePassthrough bool) (*http.Client, error) {
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
		httpClient, err = getAPIKeyHTTPClient(req)
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
func getAPIKeyHTTPClient(req *http.Request) (*http.Client, error) {

	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	// Initialise an apikey cfg struct
	apiKeyConfig := &apikey.Config{Key: cfg.APIKey}

	// Create an http client
	return apiKeyConfig.Client(req.Context(), cfg.APIKey), nil
}

// getOauth2HttpClient returns an Oauth2-authenticated HTTP client
func getOauth2HTTPClient(req *http.Request, tok *goauth2.Token) (*http.Client, error) {

	// Fetch oauth config
	oauth2Config, err := config.GetOauthConfig()
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
