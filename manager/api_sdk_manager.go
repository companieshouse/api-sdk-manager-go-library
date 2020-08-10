package manager

import (
	"encoding/json"
	"net/http"

	sdk "github.com/companieshouse/api-sdk-go/companieshouseapi"

	"github.com/companieshouse/api-sdk-go/apikey"
	choauth2 "github.com/companieshouse/api-sdk-go/oauth2"
	privatesdk "github.com/companieshouse/private-api-sdk-go/companieshouseapi"
	"github.com/pkg/errors"
	goauth2 "golang.org/x/oauth2"
)

var sdkBasePathOverridden = false
var privateSdkBasePathOverridden = false

// APISDKManager struct holds the required values to provide API to API communication
type APISDKManager struct {
	APIKey string
	APIURL string
}

var sdkManager *APISDKManager
var oauthConfig *choauth2.Config

// GetSDK will return an instance of the Go SDK using an oauth2 authenticated
// HTTP client if requested, else an API-key authenticated HTTP client will be used
func (manager APISDKManager) GetSDK(req *http.Request, usePassthrough bool) (*sdk.Service, error) {

	// Override sdkBasePath here to route API requests via ERIC
	if !sdkBasePathOverridden && len(manager.APIURL) > 0 {
		sdk.BasePath = manager.APIURL
		sdkBasePathOverridden = true
	}

	httpClient, err := manager.getHTTPClient(req, usePassthrough)
	if err != nil {
		return nil, err
	}

	return sdk.New(httpClient)
}

// GetPrivateSDK will return an instance of the Private Go SDK using an oauth2 authenticated
// HTTP client if requested, else an API-key authenticated HTTP client will be used
func (manager APISDKManager) GetPrivateSDK(req *http.Request, usePassthrough bool) (*privatesdk.Service, error) {

	// Override privateSdkBasePath here to route API requests via ERIC
	if !privateSdkBasePathOverridden && len(manager.APIURL) > 0 {
		privatesdk.BasePath = manager.APIURL
		privateSdkBasePathOverridden = true
	}

	httpClient, err := manager.getHTTPClient(req, usePassthrough)
	if err != nil {
		return nil, err
	}

	return privatesdk.New(httpClient)
}

// getHTTPClient returns an Http Client. It will be either Oauth2 or API-key
// authenticated depending on whether the calling service has requested to use the passthrough token
func (manager APISDKManager) getHTTPClient(req *http.Request, usePassthrough bool) (*http.Client, error) {
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
		httpClient, err = manager.getAPIKeyHTTPClient(req)
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
func (manager APISDKManager) getAPIKeyHTTPClient(req *http.Request) (*http.Client, error) {
	//cfg := Get()
	// Initialise an apikey cfg struct
	apiKeyConfig := &apikey.Config{Key: manager.APIKey}

	// Create an http client
	return apiKeyConfig.Client(req.Context(), manager.APIKey), nil
}

// getOauth2HttpClient returns an Oauth2-authenticated HTTP client
func getOauth2HTTPClient(req *http.Request, tok *goauth2.Token) (*http.Client, error) {

	// Fetch oauth config
	oauth2Config := &choauth2.Config{}

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
