package manager

import (
	"fmt"
	"net/http"

	"github.com/companieshouse/api-sdk-go/apikey"
	publicAPI "github.com/companieshouse/api-sdk-go/companieshouseapi"
	choauth2 "github.com/companieshouse/api-sdk-go/oauth2"
	"github.com/companieshouse/api-sdk-manager-go-library/config"
	"github.com/companieshouse/api-sdk-manager-go-library/model"
	"github.com/pkg/errors"
	goauth2 "golang.org/x/oauth2"
)

func GetSDK(req *http.Request) *publicAPI.Service {

	token := decodePassthroughHeader(req)

	httpClient, err := getHttpClient(token, req)

	apiClient, err := publicAPI.New(httpClient)
	if err != nil {
		fmt.Println(err)
	}

	return apiClient
}

func getHttpClient(passthroughToken *model.PassthroughToken, req *http.Request) (*http.Client, error) {

	var httpClient *http.Client
	var err error

	// Check the token exists because we prefer oauth
	if passthroughToken != nil {
		// If it exists, we'll use it to return an authenticated HTTP client
		tok := &goauth2.Token{
			AccessToken:passthroughToken.AccessToken,
			TokenType: passthroughToken.TokenType,
		}

		httpClient, err = getOauth2HTTPClient(req, tok)
	} else {
		// Otherwise, we'll use API-key authetication
		httpClient, err = getAPIKeyHTTPClient(req)
	}

	if err != nil {
		return nil, err
	}

	return httpClient, nil
}

//Returns a decoded passthrough token or nil if no token present
func decodePassthroughHeader(req *http.Request) *model.PassthroughToken {
	/*decodedPassthrough := &model.PassthroughToken{}
	json.Unmarshal([]byte(passthroughHeader), decodedPassthrough)

	return decodedPassthrough
	*/
	return nil
}

// getAPIKeyHttpClient returns an API-key-authenticated HTTP client
func getAPIKeyHTTPClient(req *http.Request) (*http.Client, error) {

	config, err := config.Get()
	if err != nil {
		return nil, err
	}

	// Initialise an apikey config struct
	apiKeyConfig := &apikey.Config{Key: config.APIKey}

	// Create an http client
	return apiKeyConfig.Client(req.Context(), config.APIKey), nil
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
	return oauth2Config.Client(req.Context(), tok, fn, nil), nil
}

// AccessTokenChangedCallback will be called when attempting to make an API call
// from an expired session. This function will refresh the access token on the
// session
func AccessTokenChangedCallback(newToken *goauth2.Token, private interface{}) error {
	return errors.New("Token expired")
}
