package gostrava

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const baseURLv3 = "https://www.strava.com/api/v3"

type Client struct {
	oauth OAuth
}

type OAuth struct {
	scopes       []Scope
	clientID     string
	clientSecret string
	token        AccessTokenResponse
	baseURL      string
}

const OAuthBaseURL string = "https://www.strava.com/oauth/"

const (
	oauthAuthorizationRequestSlug string = "authoriza"
	oauthGenerateTokenRequestSlug string = "token"
	oauthRevokeTokenRequestSlug   string = "revoke"
)

type AccessTokenResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    int64       `json:"expires_at"`
	ExpiresIn    int64       `json:"expires_in"`
	TokenType    string      `json:"token_type"`
	Athlete      interface{} `json:"athlete"`
	Scopes       []Scope     `json:"scopes"`
}

type Scope string

const (
	Read            Scope = "read"
	ReadAll         Scope = "read_all"
	ProfileReadAll  Scope = "profile:read_all"
	ProfileWrite    Scope = "profile:write"
	ActivityRead    Scope = "activity:read"
	ActivityReadAll Scope = "activity:read_all"
	ActivityWrite   Scope = "activity:write"
)

// SetCredentials updates the OAuth client ID and client secret for the Client.
//
// This method sets the client ID and client secret used for authenticating
// and authorizing the Client with the OAuth2 server.
//
// Parameters:
// - clientID: The OAuth client ID provided by the authorization server.
// - clientSecret: The OAuth client secret provided by the authorization server.
//
// Usage example:
//
//	client := &Client{}
//	client.SetCredentials("my-client-id", "my-client-secret")
func (c *Client) SetCredentials(clientID string, clientSecret string) {
	c.oauth.clientID = clientID
	c.oauth.clientSecret = clientSecret
}

// ClientID retrieves the current client ID used by the OAuth2 client.
//
// This method returns the client ID that has been set for the OAuth2 client.
// It allows you to access the client ID that is currently in use.
//
// Returns:
// - A string representing the OAuth2 client's client ID.
//
// Usage example:
//
//	client := &Client{}
//	id := client.ClientID()
func (c *Client) ClientID() string {
	return c.oauth.clientID
}

// ClientSecret retrieves the client secret used by the OAuth2 client.
//
// This method returns the client secret that has been set for the OAuth2 client.
// It allows you to access the client secret that is currently in use.
//
// Returns:
// - A string representing the OAuth2 client's client secret.
//
// Usage example:
//
//	client := &Client{}
//	secret := client.ClientSecret()
func (c *Client) ClientSecret() string {
	return c.oauth.clientSecret
}

// SetScopes sets the scopes for the OAuth2 client.
//
// This method configures the scopes that define the level of access requested
// by the OAuth2 client.
//
// Parameters:
// - scopes: A slice of scopes that define the level of access requested.
//
// Usage example:
//
//	client := &Client{}
//	client.SetScopes([]Scope{"read", "write"})
func (c *Client) SetScopes(scopes []Scope) {
	c.oauth.scopes = scopes
}

// Scopes retrieves the current scopes set for the OAuth2 client.
//
// This method returns the scopes that have been set for the OAuth2 client,
// which define the level of access requested.
//
// Returns:
// - A slice of scopes that define the level of access requested.
//
// Usage example:
//
//	client := &Client{}
//	scopes := client.Scopes()
func (c *Client) Scopes() []Scope {
	return c.oauth.scopes
}

// OAuthFlowOptions holds configuration options for the OAuth2 flow.
//
// It includes options such as state and force which are used during the authorization process.
type OAuthFlowOptions struct {
	State string // A state parameter used to maintain state between the request and callback.
	Force bool   // If true, forces the user to re-authenticate.
}

// AuthorizationCodeURL generates the authentication URL that the user will be redirected to
// in order to initiate the OAuth2 flow.
//
// Parameters:
//   - redirectURI: The URI to which the authorization server will redirect the user-agent
//     after the user grants or denies permission. This URI must be registered with the
//     authorization server as part of the client registration.
//   - options: Configuration options for the OAuth2 flow, including state and force.
//
// Returns:
// - string: The fully formed URL that the user-agent should be redirected to initiate the OAuth2 flow.
//
// Usage example:
//
//	url := client.AuthorizationCodeURL("http://localhost:8080/callback", OAuthFlowOptions{State: "state123", Force: false})
func (c *Client) AuthorizationCodeURL(redirectURI string, options OAuthFlowOptions) string {
	return MakeAuthorizationCodeURL(c.oauth.clientID, c.oauth.clientSecret, redirectURI, c.oauth.scopes, options)
}

// MakeAuthorizationCodeURL generates the authentication URL that the user will be redirected to
// in order to initiate the OAuth2 flow.
//
// Parameters:
//   - clientID: The ID assigned to the client application by the authorization server.
//   - clientSecret: The client secret used for authentication with the authorization server.
//   - redirectURI: The URI to which the authorization server will redirect the user-agent
//     after the user grants or denies permission. This URI must be registered with the
//     authorization server as part of the client registration.
//   - scopes: A list of scopes (permissions) that the application requests access to.
//   - options: Configuration options for the OAuth2 flow, including state and force.
//
// Returns:
// - string: The fully formed URL that the user-agent should be redirected to initiate the OAuth2 flow.
//
// Usage example:
//
//	url := MakeAuthorizationCodeURL("my-client-id", "my-client-secret", "http://localhost:8080/callback", []Scope{"read", "write"}, OAuthFlowOptions{State: "state123", Force: false})
func MakeAuthorizationCodeURL(
	clientID, clientSecret, redirectURI string, scopes []Scope, options OAuthFlowOptions,
) string {
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", joinScopes(scopes))
	if options.State != "" {
		q.Set("state", options.State)
	}
	if options.Force {
		q.Set("approval_prompt", "force")
	}

	return fmt.Sprintf("%s%s?%s", OAuthBaseURL, oauthAuthorizationRequestSlug, q.Encode())
}

type serverMessage struct {
	Type    string
	Message string
	Scopes  []Scope
}

func (c *Client) StartOAuth2Authrorization(redirectURI string, options OAuthFlowOptions) (*AccessTokenResponse, error) {

	srvChan := make(chan int8)
	srvRespChan := make(chan serverMessage, 1)


	var srv *http.Server

	isLocalRedirect := strings.Contains(redirectURI, "localhost") || strings.Contains(redirectURI, "127.0.0.1")
	if isLocalRedirect {
		// Start a localhost server that will handle the redirect url
		u, err := url.Parse(redirectURI)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the redirectURI: %s, %s", redirectURI, err)
		}
		_, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return nil, fmt.Errorf("failed to split the redirect url into host and port: %s", err)
		}

		srv := &http.Server{Addr: ":" + port}

		redirectPath := u.EscapedPath()
		http.HandleFunc(redirectPath, func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()

			errs := q.Get("error")
			if err != "" {
				w.Write([]byte("Got and error from the server: " + errs))
				srvRespChan <- serverMessage{Type: "error", Message: errs}
				return
			}

			code, scopeStrings := q.Get("code"), strings.Split(q.Get("scope"), ",")
			w.Write([]byte("if you see this, code and scopes have been retrieved! you can close this window"))

			scopes := make([]Scope, len(scopeStrings))
			for i, s := range scopeStrings {
				scopes[i] = Scope(s)
			}

			srvRespChan <- serverMessage{
				Type:    "code",
				Message: code,
				Scopes:  scopes,
			}
		})

		go func() {
			srvChan <- 1
			err := srv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				fmt.Printf("error while serving locally: %s\n", err)
			}
		}()

		<-srvChan
	}

	authURL := c.AuthorizationCodeURL(redirectURI, options)

	fmt.Println("User must be redirected to the following url: %s", authURL)

	code := ""
	scopes := []Scope{}

	if isLocalRedirect {
		// wait for code to be returned by the server
		srvMessage := <- srvRespChan

		if srvMessage.Type == "error" {
			return nil, fmt.Errorf("strava server responded with error: %s", srvMessage.Message)
		}

		code = srvMessage.Message
		scopes = srvMessage.Scopes

		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			fmt.Printf("error while shutting down local server: %s\n", err)
		}else{
			fmt.Printf("Paste code and press enter:\n")
			_, err := fmt.Scan(&code)
			if err != nil {
				return nil, fmt.Errorf("Failed to read code from input: %s", err)
			}

			fmt.Printf("Paste scopes and press enter:\n")
			scopeStrings := ""
			_, err = fmt.Scan(&scopeStrings)
			if err != nil {
				return nil, fmt.Errorf("Failed to read scopes from input: %s", err)
			}
			for _, scope := range strings.Split(scopeStrings, ",") {
				scopes = append(scopes, Scope(scope))
			}
		}
	}

	if code == "" {
		return nil, fmt.Errorf("no coce was recevied from oAuth2Flow")
	}

	if len(scopes) == 0 {
		scopes = []Scope{Read}
	}

	// Do the exchange process!
	err = 
}

// --------- Helper ---------

func joinScopes(scopes []Scope) string {
	stringScopes := make([]string, len(scopes))
	for i, scope := range scopes {
		stringScopes[i] = string(scope)
	}

	return strings.Join(stringScopes, ",")
}
