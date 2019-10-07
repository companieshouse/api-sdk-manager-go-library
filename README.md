# api-sdk-manager-go-library

Go implementation of the api-sdk-manager-java-library. Users of the SDK will interact with the SDK manager, 
which wraps up some useful functionality that determines whether to initialise an OAuth2 or API Key-authenticated
http client.

## Requirements

In order to build this library locally you will need the following:
- [Go](https://golang.org/)
- [Git](https://git-scm.com/downloads)

## Getting started

The library is built using the following commands:
```
go get ./...
go build
```

## Testing
The library can be tested by running the following in the command line (in the `api-sdk-manager-go-library` directory):
```
goconvey
```

Note: this library is not a standalone service, and can only be used within services or other libraries.

## Environment Variables
The following environment variables are required when integrating the SDK manager into any Go service.

Note: These are OAuth2 config items, and are standard when using OAuth2.

Key | Description | Scope | Mandatory
----|-------------|-------|-----------
OAUTH2_CLIENT_ID | The application ID of the client | Config | Y
OAUTH2_CLIENT_SECRET | The application secret of the client | Config | Y
OAUTH2_REDIRECT_URI | The URL that OAuth2 will redirect to after authorisation | Config | Y
OAUTH2_AUTH_URI | The authorisation endpoint  | Config | Y
OAUTH2_TOKEN_URI | The token endpoint | Config | Y
SCOPE | Optional requested permissions | Config | N
API_KEY | The application access key for the API | Config | Y
API_URL | The application endpoint for the API | Config | Y

## Example library usage

To use the Manager package, add the following to the relevant package import:
- `"github.com/companieshouse/api-sdk-manager-go-library/manager"`
