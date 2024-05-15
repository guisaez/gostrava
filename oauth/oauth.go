package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	strava "github.com/guisaez/go-strava"
)

type StravaOAuth struct {
	ClientID     string  // The application’s ID, obtained during registration.
	ClientSecret string  // The application’s secret, obtained during registration.
	CallbackURL  string  // URL to which the user will be redirected after authentication. Must be within the callback domain specified by the application. localhost and 127.0.0.1 are white-listed.
	Scopes       []string // Requested scopes, as a comma delimited string, e.g. "activity:read_all,activity:write". Applications should request only the scopes required for the application to function normally. The scope activity:read is required for activity webhooks.

	// Custom httpClient
	RequestClient *http.Client
}

// Generates the authentication url that will be used to redirect the user to
// Strava so it can initiate the OAuthFlow
// State will be returned by the redirect URI. Useful if the authentication is done from
// various points in the app.
// Force will always show the authorization prompt even if the user has already authorized the current
// application
func (oauth *StravaOAuth) AuthCodeURL(state string, force bool) string {

	url := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&client_secret=%s&redirect_uri=%s&scope=%s&state=%s",
		endpoints.Auth,
		oauth.ClientID,
		oauth.ClientSecret,
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

	client := http.DefaultClient

	if oauth.RequestClient != nil {
		client = oauth.RequestClient
	}

	resp, err := client.PostForm(endpoints.Token, url.Values{
		"client_id": { oauth.ClientID },
		"client_secret": {oauth.ClientSecret },
		"code": { code },
		"grant_type": { "authorization_code "},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := strava.HandleBadResponse(resp); err != nil {
		return nil, err
	}
	
	var response StravaOAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	
	return &response, nil
}

// Handles generating a new set of access and refresh tokens based on a previous refresh token.
func (oauth *StravaOAuth) Refresh(refreshToken string) (*RefreshTokenResponse, error) {

	client := http.DefaultClient

	if oauth.RequestClient != nil {
		client = oauth.RequestClient
	}

	payload, err := json.Marshal(RefreshTokenPayload{
		ClientID: oauth.ClientID,
		ClientSecret: oauth.ClientSecret,
		RefreshToken: refreshToken,
		GrantType: "refresh_token",
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(endpoints.Token, "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := strava.HandleBadResponse(resp); err != nil {
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

		if errParam := r.URL.Query().Get("error"); errParam == "access_denied" {
			handleError(AccessDenied, w, r)
			return
		}

		tokens, err := oauth.Exchange(r.URL.Query().Get("code"))
		if err != nil {
			handleError(err, w, r)
			return
		}

		tokens.withScopes(r.URL.Query().Get("scope"))

		handleSuccess(tokens, w, r)
	}
}

func (oauthResponse *StravaOAuthResponse) withScopes(scopes string)  {
	oauthResponse.Scopes = strings.Split(scopes, ",")
}