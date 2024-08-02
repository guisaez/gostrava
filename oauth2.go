package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuthBaseURL is the base URL for OAuth2 operations with Strava.
const OAuthBaseURL string = "https://www.strava.com/oauth/"

// OAuth request slugs.
const (
	oauthAuthorizationRequestSlug string = "authorize"
	oauthGenerateTokenRequestSlug string = "token"
	oauthRevokeTokenRequestSlug   string = "revoke"
)

// Scope represents the different scopes of access that can be requested from Strava.
type Scope string

// Different OAuth scopes for accessing Strava data.
const (
	Read            Scope = "read"
	ReadAll         Scope = "read_all"
	ProfileReadAll  Scope = "profile:read_all"
	ProfileWrite    Scope = "profile:write"
	ActivityRead    Scope = "activity:read"
	ActivityReadAll Scope = "activity:read_all"
	ActivityWrite   Scope = "activity:write"
)

// OAuth contains the client credentials and scopes for OAuth2 authentication.
type OAuth struct {
	scopes       []Scope
	clientID     string
	clientSecret string
}

// RefreshTokenResponse represents the response returned when a refresh token is exchanged for a new access token.
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int64  `json:"expires_in"`
}

// AccessTokenResponse represents the response returned when an authorization code is exchanged for an access token.
type AccessTokenResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    int64       `json:"expires_at"`
	ExpiresIn    int64       `json:"expires_in"`
	TokenType    string      `json:"token_type"`
	Athlete      interface{} `json:"athlete"`
	Scopes       []Scope     `json:"scopes"`
}

// IsExpired checks if the access token is expired.
func (a *AccessTokenResponse) IsExpired() bool {
	return time.Unix(a.ExpiresAt, 0).Before(time.Now())
}

// OAuthFlowOptions holds configuration options for the OAuth2 flow.
type OAuthFlowOptions struct {
	State string // A state parameter used to maintain state between the request and callback.
	Force bool   // If true, forces the user to re-authenticate.
}

// AuthorizationCodeURL generates the authentication URL to initiate the OAuth2 flow.
func (c *Strava) AuthorizationCodeURL(redirectURI string, options OAuthFlowOptions) string {
	return MakeAuthorizationCodeURL(c.oauth.clientID, c.oauth.clientSecret, redirectURI, c.oauth.scopes, options)
}

// MakeAuthorizationCodeURL constructs the URL to initiate the OAuth2 flow.
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

// AuthorizationTokenURL generates the URL for requesting an access token using an authorization code.
func (c *Strava) AuthorizationTokenURL(code string) string {
	return MakeAuthorizationTokenURL(c.oauth.clientID, c.oauth.clientSecret, code)
}

// MakeAuthorizationTokenURL constructs the URL for requesting an access token from the OAuth server.
func MakeAuthorizationTokenURL(clientID, clientSecret, code string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("code", code)
	q.Set("grant_type", "authorization_code")

	return fmt.Sprintf("%s%s?%s", OAuthBaseURL, oauthGenerateTokenRequestSlug, q.Encode())
}

// RefreshTokenURL generates a URL for requesting a new access token using a refresh token.
func (c *Strava) RefreshTokenURL(refreshToken string) string {
	return MakeRefreshTokenURL(c.ClientID(), c.ClientSecret(), refreshToken)
}

// MakeRefreshTokenURL constructs the URL for refreshing an access token using a refresh token.
func MakeRefreshTokenURL(clientID, clientSecret, refreshToken string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("refresh_token", refreshToken)
	q.Set("grant_type", "refresh_token")

	return fmt.Sprintf("%s%s?%s", OAuthBaseURL, oauthRevokeTokenRequestSlug, q.Encode())
}

// AuthorizationTokenRequest exchanges an authorization code for an access token.
func (s *Strava) AuthorizationTokenRequest(code string) (AccessTokenResponse, error) {
	authTokenURL := s.AuthorizationTokenURL(code)

	authResp := new(AccessTokenResponse)

	err := s.client.NewRequest(RequestOptions{
		URL:    authTokenURL,
		Method: http.MethodPost,
	}, authResp)
	if err != nil {
		return AccessTokenResponse{}, err
	}

	return *authResp, nil
}

// RefreshTokenRequest exchanges a refresh token for a new access token.
func (s *Strava) RefreshTokenRequest(refreshToken string) (RefreshTokenResponse, error) {
	tokenURL := s.RefreshTokenURL(refreshToken)

	refreshTokenResp := new(RefreshTokenResponse)

	err := s.client.NewRequest(RequestOptions{
		URL:    tokenURL,
		Method: http.MethodPost,
	}, refreshTokenResp)
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	return *refreshTokenResp, nil
}

// --------- Helper ---------

// joinScopes joins multiple scopes into a single comma-separated string.
func joinScopes(scopes []Scope) string {
	stringScopes := make([]string, len(scopes))
	for i, scope := range scopes {
		stringScopes[i] = string(scope)
	}

	return strings.Join(stringScopes, ",")
}
