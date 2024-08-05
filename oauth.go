package gostrava

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// OAuth endpoint slugs.
const (
	authorizationEndpoint string = "oauth/authorize"
	tokenExchangeEndpoint string = "oauth/token"
	tokenRevokeEndpoint   string = "oauth/deauthorize"
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

// OAuthService provides methods for interacting with Strava's OAuth2 API.
type OAuthService struct {
	service
	scopes []Scope
}

// SetScopes configures the scopes for OAuth2. Defaults to "read" if no scopes are provided.
func (s *OAuthService) SetScopes(scopes ...Scope) {
	if len(scopes) == 0 {
		s.scopes = []Scope{Read}
	}
	s.scopes = scopes
}

// OAuthFlowOptions holds options for the OAuth2 flow.
type OAuthFlowOptions struct {
	State string // A state parameter used to maintain state between the request and callback.
	Force bool   // If true, forces the user to re-authenticate.
}

// AuthorizationURL creates the URL to redirect users to Strava's authorization page.
func (s *OAuthService) AuthorizationURL(redirectURI string, options OAuthFlowOptions, scopes []Scope) string {
	if len(scopes) == 0 {
		return BuildAuthorizationURL(s.client.BaseURL.String(), s.client.clientID, s.client.clientSecret, redirectURI, s.scopes, options)
	}
	return BuildAuthorizationURL(s.client.BaseURL.String(), s.client.clientID, s.client.clientSecret, redirectURI, scopes, options)
}

// BuildAuthorizationURL constructs the URL for initiating the OAuth2 flow.
func BuildAuthorizationURL(
	baseURL, clientID, clientSecret, redirectURI string, scopes []Scope, options OAuthFlowOptions,
) string {
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", joinScopes(scopes))
	if options.State != "" {
		q.Set("state", options.State)
	}
	if options.Force {
		q.Set("approval_prompt", "force")
	}

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return fmt.Sprintf("%s%s?%s", baseURL, authorizationEndpoint, q.Encode())
}

// TokenExchangeURL creates the URL for exchanging an authorization code for an access token.
func (s *OAuthService) TokenExchangeURL(code string) string {
	return BuildTokenExchangeURL(s.client.BaseURL.String(), s.client.clientID, s.client.clientSecret, code)
}

// BuildTokenExchangeURL constructs the URL for exchanging an authorization code for an access token.
func BuildTokenExchangeURL(baseURL, clientID, clientSecret, code string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("code", code)
	q.Set("grant_type", "authorization_code")

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return fmt.Sprintf("%s%s?%s", baseURL, tokenExchangeEndpoint, q.Encode())
}

// TokenRefreshURL creates the URL for refreshing an expired access token.
func (s *OAuthService) TokenRefreshURL(refreshToken string) string {
	return BuildTokenRefreshURL(s.client.BaseURL.String(), s.client.clientID, s.client.clientSecret, refreshToken)
}

// TokenRefreshURL constructs the URL for refreshing an expired access token.
func BuildTokenRefreshURL(baseURL, clientID, clientSecret, refreshToken string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("refresh_token", refreshToken)
	q.Set("grant_type", "refresh_token")

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return fmt.Sprintf("%s%s?%s", baseURL, tokenExchangeEndpoint, q.Encode())
}

// TokenRevocationURL creates the URL for revoking an access token.
func (s *OAuthService) TokenRevocationURL(accessToken string) string {
	return BuildTokenRevocationURL(s.client.BaseURL.String(), s.client.clientID, s.client.clientSecret, accessToken)
}

// BuildTokenRevocationURL constructs the URL for revoking an access token.
func BuildTokenRevocationURL(baseURL, clientID, clientSecret, accessToken string) string {
	q := url.Values{}
	q.Set("access_token", accessToken)

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return fmt.Sprintf("%s%s?%s", baseURL, tokenRevokeEndpoint, q.Encode())
}

// AuthorizationResponse represents the response from exchanging an authorization code for an access token.
type AuthorizationResponse struct {
	AccessToken  string      `json:"access_token"`  // The access token to be used for authenticated API requests.
	RefreshToken string      `json:"refresh_token"` // The refresh token to obtain a new access token when the current one expires.
	ExpiresAt    int64       `json:"expires_at"`    // Unix timestamp indicating when the access token expires.
	ExpiresIn    int64       `json:"expires_in"`    // Duration in seconds from the time of the response until the access token expires.
	TokenType    string      `json:"token_type"`    // The type of the access token (e.g., "bearer").
	Athlete      interface{} `json:"athlete"`       // Information about the authenticated athlete, if provided by the API.
	Scopes       []Scope     `json:"scopes"`        // Scopes granted to the access token. Not returned by the API; user-specified scopes can be passed to the request.
}

// ExchangeAuthorizationCode exchanges an authorization code for an access token.
func (s *OAuthService) ExchangeAuthorizationCode(context context.Context, code string, scopes []Scope) (*AuthorizationResponse, *http.Response, error) {
	tokenExchangeURL := s.TokenExchangeURL(code)

	req, err := s.client.NewRequest(http.MethodPost, tokenExchangeURL, nil, func(req *http.Request) error {
		req.URL, _ = url.Parse(tokenExchangeURL)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	tokenResponse := new(AuthorizationResponse)

	resp, err := s.client.DoAndParse(context, req, tokenResponse)
	if err != nil {
		return nil, resp, err
	}

	tokenResponse.Scopes = scopes

	return tokenResponse, resp, nil
}

// RefreshTokenResponse represents the response from refreshing an access token.
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`  // The access token to be used for authenticated API requests.
	RefreshToken string `json:"refresh_token"` // The refresh token to obtain a new access token when the current one expires.
	ExpiresAt    int64  `json:"expires_at"`    // Unix timestamp indicating when the access token expires.
	ExpiresIn    int64  `json:"expires_in"`    // Duration in seconds from the time of the response until the access token expires.
}

// RefreshToken refreshes an expired access token using a refresh token.
func (s *OAuthService) RefreshToken(context context.Context, refreshToken string) (*RefreshTokenResponse, *http.Response, error) {
	refreshTokenURL := s.TokenRefreshURL(refreshToken)

	req, err := s.client.NewRequest(http.MethodPost, refreshTokenURL, nil, func(req *http.Request) error {
		req.URL, _ = url.Parse(refreshTokenURL)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	refreshResponse := new(RefreshTokenResponse)

	resp, err := s.client.DoAndParse(context, req, refreshResponse)
	if err != nil {
		return nil, resp, err
	}

	return refreshResponse, resp, nil
}

// RevokeToken revokes an access token, invalidating it and any associated refresh tokens.
func (s *OAuthService) RevokeToken(context context.Context, accessToken string) (*http.Response, error) {
	revokeTokenURL := s.TokenRevocationURL(accessToken)

	req, err := s.client.NewRequest(http.MethodPost, revokeTokenURL, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(context, req)
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
