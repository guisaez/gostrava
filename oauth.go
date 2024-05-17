package gostrava

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type StravaOAuth struct {
	CallbackURL string   // URL to which the user will be redirected after authentication. Must be within the callback domain specified by the application. localhost and 127.0.0.1 are white-listed.
	Scopes      []string // Requested scopes, as a comma delimited string, e.g. "activity:read_all,activity:write". Applications should request only the scopes required for the application to function normally. The scope activity:read is required for activity webhooks.

	*StravaClient
}

// Token exchange result
type StravaOAuthResponse struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"` // The refresh token for this user, to be used to get the next access token for this user. Please expect that this value can change anytime you retrieve a new access token. Once a new refresh token code has been returned, the older code will no longer work.
	TokenType    string         `json:"token_type"`    // Bearer
	ExpiresAt    uint64         `json:"expires_at"`    // The number of seconds since the epoch when the provided access token will expire
	ExpiresIn    int            `json:"expires_in"`    // Seconds until the short-lived access token will expire
	Athlete      SummaryAthlete `json:"athlete"`       // A summary of the athlete information
	Scopes       []string
}

func (oauthResponse *StravaOAuthResponse) withScopes(scopes string)  {
	oauthResponse.Scopes = strings.Split(scopes, ",")
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`  // The short-lived access token
	RefreshToken string `json:"refresh_token"` // The number of seconds since the epoch when the provided access token will expire
	ExpiresAt    uint64 `json:"expires_at"`    // Seconds until the short-lived access token will expire
	ExpiresIn    int    `json:"expires_in"`    // The refresh token for this user, to be used to get the next access token for this user. Please expect that this value can change anytime you retrieve a new access token. Once a new refresh token code has been returned, the older code will no longer work.
}

type RefreshTokenPayload struct {
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

// Defines the level of access or permission that the application requests from
// and user's Strava Account when the user authorizes the application. Determines
// what actions the application can perform on behalf of the user and what data it can access.
var Scopes = struct {
	Read            string
	ReadAll         string
	ProfileReadAll  string
	ProfileWrite    string
	ActivityRead    string
	ActivityReadAll string
	ActivityWrite   string
}{
	"read",
	"read_all",
	"profile:read_all",
	"profile:write",
	"activity:read",
	"activity:read_all",
	"activity:write",
}

var endpoints = struct {
	Auth        string
	Unauthorize string
	Token       string
}{
	"https://www.strava.com/oauth/authorize",
	"https://www.strava.com/oauth/deauthorize",
	"https://www.strava.com/oauth/token",
}

// Generates the authentication url that will be used to redirect the user to
// Strava so it can initiate the OAuthFlow
// State will be returned by the redirect URI. Useful if the authentication is done from
// various points in the app.
// Force will always show the authorization prompt even if the user has already authorized the current
// application
func (oauth *StravaOAuth) AuthCodeURL(force bool, state string) string {

	url := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&client_secret=%s&redirect_uri=%s&scope=%s&state=%s",
		endpoints.Auth,
		oauth.clientID,
		oauth.clientSecret,
		oauth.CallbackURL,
		strings.Join(oauth.Scopes, ","),
		state,
	)

	if force {
		return url + "&approval_prompt=force"
	}

	return url
}

// Handles the token exchange portion of the OAuth2 flow.
func (oauth *StravaOAuth) Exchange(code string) (*StravaOAuthResponse, error) {

	if code == "" {
		return nil, InvalidCodeError
	}

	var response StravaOAuthResponse
	err := oauth.postForm(context.Background(), "", endpoints.Token, url.Values{
		"client_id": { oauth.clientID },
		"client_secret": { oauth.clientSecret },
		"code": { code },
		"grant_type": { "authorization_code" },
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (oauth *StravaOAuth) RevokeAccess(access_token string) error {

	err := oauth.postForm(context.Background(), "", endpoints.Unauthorize, url.Values{
		"access_token": { access_token }}, nil)
	if err != nil {
		return err
	}
	return nil
}

// Handles generating a new set of access and refresh tokens based on a previous refresh token.
func (oauth *StravaOAuth) Refresh(refreshToken string) (*RefreshTokenResponse, error) {

	payload, err := json.Marshal(RefreshTokenPayload{
		ClientID: oauth.clientID,
		ClientSecret: oauth.clientSecret,
		RefreshToken: refreshToken,
		GrantType: "refresh_token",
	})
	if err != nil {
		return nil, err
	}

	resp, err := oauth.client.Post(endpoints.Token, "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleBadResponse(resp); err != nil {
		return nil, err
	}

	var response RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (oauth *StravaOAuth) HandlerFunc(
	handleSuccess func(tokens *StravaOAuthResponse, w http.ResponseWriter, r *http.Request),
	handleError func(err error, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := r.URL.Query()
		if errParam := params.Get("error"); errParam == "access_denied" {
			handleError(AccessDeniedError, w, r)
			return
		}

		tokens, err := oauth.Exchange(r.URL.Query().Get("code"))
		if err != nil {
			handleError(err, w, r)
			return
		}

		tokens.withScopes(params.Get("scope"))

		handleSuccess(tokens, w, r)
	}
}