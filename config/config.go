package config

import (
	choauth2 "github.com/companieshouse/api-sdk-go/oauth2"
	"github.com/ian-kent/gofigure"
	goauth2 "golang.org/x/oauth2"
)

// Config holds the configuration values required for the `api-sdk-manager-go-library`
type Config struct {
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

var oauthConfig *choauth2.Config
var config *Config

// Get returns an instance of the config struct, with an error where appropriate
func Get() (*Config, error) {

	if config != nil {
		return config, nil
	}

	config = &Config{}
	err := gofigure.Gofigure(config)
	return config, err
}

// GetOauthConfig returns an instance of a Companies House oauth config struct
// and an error wheere appropriate
func GetOauthConfig() (*choauth2.Config, error) {

	if oauthConfig != nil {
		return oauthConfig, nil
	}

	config, err := Get()
	if err != nil {
		return nil, err
	}

	oauthConfig = &choauth2.Config{}
	oauthConfig.ClientID = config.ClientID
	oauthConfig.ClientSecret = config.ClientSecret
	oauthConfig.RedirectURL = config.RedirectURL
	oauthConfig.Scopes = config.Scopes
	oauthConfig.Endpoint = goauth2.Endpoint{
		AuthURL:  config.AuthURL,
		TokenURL: config.TokenURL,
	}

	return oauthConfig, nil
}
