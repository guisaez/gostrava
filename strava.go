package gostrava

import "net/http"

// BaseURLv3 is the base URL for the Strava API v3.
const BaseURLv3 = "https://www.strava.com/api/v3"

// Strava represents the main client for interacting with the Strava API.
type Strava struct {
	oauth      OAuth
	client *stravaHTTPClient
}

// New creates a new instance of the Strava client with default settings.
//
// It initializes a new Strava client with the default HTTP client.
//
// Returns:
// - *Strava: A new instance of the Strava client.
func New() *Strava {
	stravaHTTPClient := &stravaHTTPClient{
		httpClient: http.DefaultClient,
	}
	return &Strava{
		client: stravaHTTPClient,
	}
}

// SetCredentials updates the OAuth client ID and client secret for the Strava client.
//
// This method sets the client ID and client secret used for authenticating
// and authorizing the client with the OAuth2 server.
//
// Parameters:
// - clientID: The OAuth client ID provided by the authorization server.
// - clientSecret: The OAuth client secret provided by the authorization server.
func (c *Strava) SetCredentials(clientID string, clientSecret string) {
	c.oauth.clientID = clientID
	c.oauth.clientSecret = clientSecret
}

// ClientID retrieves the current client ID used by the OAuth2 client.
//
// This method returns the client ID that has been set for the OAuth2 client.
// It allows you to access the client ID that is currently in use.
//
// Returns:
// - string: A string representing the OAuth2 client's client ID.
func (c *Strava) ClientID() string {
	return c.oauth.clientID
}

// ClientSecret retrieves the client secret used by the OAuth2 client.
//
// This method returns the client secret that has been set for the OAuth2 client.
// It allows you to access the client secret that is currently in use.
//
// Returns:
// - string: A string representing the OAuth2 client's client secret.
func (c *Strava) ClientSecret() string {
	return c.oauth.clientSecret
}

// SetScopes sets the scopes for the OAuth2 client.
//
// This method configures the scopes that define the level of access requested
// by the OAuth2 client.
//
// Parameters:
// - scopes: A slice of scopes that define the level of access requested.
func (c *Strava) SetScopes(scopes []Scope) {
	c.oauth.scopes = scopes
}

// Scopes retrieves the current scopes set for the OAuth2 client.
//
// This method returns the scopes that have been set for the OAuth2 client,
// which define the level of access requested.
//
// Returns:
// - []Scope: A slice of scopes that define the level of access requested.
func (c *Strava) Scopes() []Scope {
	return c.oauth.scopes
}

// UseCustomHTTPClient sets a custom HTTP client for the Strava client.
//
// This allows you to provide a custom `http.Client` with specific configurations,
// such as timeouts, transport settings, or proxies, instead of using the default
// HTTP client. This is useful when you need to control the behavior of HTTP requests
// made by the Strava client.
//
// Parameters:
// - client (*http.Client): The custom HTTP client to be used for making HTTP requests.
func (c *Strava) UseCustomHTTPClient(client HTTPClient) {
	c.client.httpClient = client
}

// CustomRequest sends an HTTP request based on the provided RequestOptions and returns the result.
//
// This method leverages the `NewRequest` method of the `stravaHTTPClient` instance associated
// with the `Strava` struct to handle the request creation and execution. It provides a flexible
// way to configure and execute HTTP requests, including specifying the HTTP method, URL, URL parameters,
// request body, headers, content type, and context for request execution.
//
// Parameters:
// - options (RequestOptions): Configuration for the HTTP request, including the HTTP method,
//   URL, URL parameters, request body, headers, content type, and context. The `RequestOptions`
//   struct should be used to specify all necessary settings for the request.
// - responseData (interface{}): A pointer to a variable where the response data will be stored.
//   This variable should be a pointer to a struct or slice that matches the expected response format.
//   For non-JSON responses, use a pointer to a byte slice to capture raw response data.
//
// Returns:
// - error: Returns an error if the request fails or if there are issues with the response. 
//   If the request is successful and there are no errors, it returns nil.
func (c *Strava) CustomRequest(options RequestOptions, responseData interface{}) error {
	return c.client.NewRequest(options, responseData)
}