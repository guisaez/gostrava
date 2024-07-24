package gostrava

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type ScopeString string

const OAuthBaseURL string = "https://www.strava.com/oauth/"

const (
	oauthAuthorizationRequestSlug string = "authoriza"
	oauthGenerateTokenRequestSlug string = "token"
	oauthRevokeTokenRequestSlug   string = "revoke"
)

type OAuth struct {
	clientID     string
	clientSecret string
}

func (c *Client) SetClientID(clientID string) {
	c.oauth.clientID = clientID
}

func (c *Client) SetClientSecret(clientSecret string) {
	c.oauth.clientSecret = clientSecret
}

const (
	Read            ScopeString = "read"
	ReadAll         ScopeString = "read_all"
	ProfileReadAll  ScopeString = "profile:read_all"
	ProfileWrite    ScopeString = "profile:write"
	ActivityRead    ScopeString = "activity:read"
	ActivityReadAll ScopeString = "activity:read_all"
	ActivityWrite   ScopeString = "activity:write"
)

// AuthorizationCodeURL generates the authentication URL that the user will be redirected to
// in order to initiate the OAuthFlow.
//
// Args:
//   - redirectURI: The URI to which the authorization server will redirect the user-agent
//     after the user grants or denies permission. This URI must be registered with the
//     authorization server as part of the client registration.
//   - state: An opaque value used by the client to maintain state between the request and callback.
//     Typically used to prevent CSRF attacks and to maintain user state.
//   - force: If true, it will force the authorization server to prompt the user for consent
//     even if they have already done so for the current application.
//   - scopes: A list of scopes (permissions) that the application requests access to.
//
// Returns:
//   - string: The fully formed URL that the user-agent should be redirected to initiate the OAuthFlow.
func (c *Client) AuthorizationCodeURL(redirectURI, state string, force bool, scopes []ScopeString) string {
	return MakeAuthorizationCodeURL(c.oauth.clientID, c.oauth.clientSecret, redirectURI, state, force, scopes)
}

type serverMessage struct {
	Type    string
	Content string
	Scopes  []string
}

func (c *Client) AuthorizationCodeRequest(redirectURI string, state string, force bool, scopes []ScopeString) error {

	authURL := c.AuthorizationCodeURL(redirectURI, state, force, scopes)

	srvChan := make(chan int8)
	srvResponse := make(chan serverMessage, 1)
	errChan := make(chan serverMessage, 1)

	var srv *http.Server

	localRedirect := strings.Contains(redirectURI, "localhost") || strings.Contains(redirectURI, "127.0.0.1")
	if localRedirect {
		// Starts a localhost server that will handle the redirect url
		u, err := url.Parse(redirectURI)
		if err != nil {
			return fmt.Errorf("failed to parse redirect URI: %s", err)
		}
		_, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return fmt.Errorf("failed to split the redirect uri into host and port segments: %s", err)
		}

		srv := &http.Server{Addr: ":" + port}

		redirectPath := u.EscapedPath()
		http.HandleFunc(redirectPath, func(w http.ResponseWriter, r *http.Request) {

			q := r.URL.Query()

			errs := q.Get("error")
			if errs == "" {
				w.Write([]byte("Got an error from the server: " + errs))
				errChan <- serverMessage{Type: "error", Content: errs}
				return
			}

			code, scopes := q.Get("code"), q.Get("scope")
			
			// TO-DO 

			w.Write([]byte("if code is present, it has been retrieved! you can close this window"))
		})
	}
}

// MakeAuthorizationCodeURL generates the authentication URL that the user will be redirected to
// in order to initiate the OAuthFlow.
//
// Args:
//   - clientID: The ID assigned to the client application by the authorization server.
//   - clientSecret: The client secret used for authentication with the authorization server.
//   - redirectURI: The URI to which the authorization server will redirect the user-agent
//     after the user grants or denies permission. This URI must be registered with the
//     authorization server as part of the client registration.
//   - state: An opaque value used by the client to maintain state between the request and callback.
//     Typically used to prevent CSRF attacks and to maintain user state.
//   - force: If true, it will force the authorization server to prompt the user for consent
//     even if they have already done so for the current application.
//   - scopes: A list of scopes (permissions) that the application requests access to.
//
// Returns:
//   - string: The fully formed URL that the user-agent should be redirected to initiate the OAuthFlow.
func MakeAuthorizationCodeURL(
	clientID, clientSecret, redirectURI, state string,
	force bool, scopes []ScopeString,
) string {
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", joinScopes(scopes))
	q.Set("state", state)

	if force {
		q.Set("approval_prompt", "force")
	}

	return fmt.Sprintf("%s%s?%s", OAuthBaseURL, oauthAuthorizationRequestSlug, q.Encode())
}

func joinScopes(scopes []ScopeString) string {
	stringScopes := make([]string, len(scopes))
	for i, scope := range scopes {
		stringScopes[i] = string(scope)
	}

	return strings.Join(stringScopes, ",")
}
