package gostrava

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var endpoints = struct {
	Auth        string
	Unauthorize string
	Token       string
}{
	"https://www.strava.com/oauth/authorize",
	"https://www.strava.com/oauth/deauthorize",
	"https://www.strava.com/oauth/token",
}

// Defines the level of access or permission that the application requests from
// and user's Strava Account when the user authorizes the application. Determines
// what actions the application can perform on behalf of the user and what data it can access.
var StravaOAuthScopes = struct {
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

type StravaOAuthOpts struct {
	ClientID     string // The application’s ID, obtained during registration.
	ClientSecret string // The application’s secret, obtained during registration.
	HTTPClient   HTTPClient
	CallbackURL  string
	Scopes       []string
}

type StravaOAuth struct {
	StravaOAuthOpts

	client *StravaClient
}

func NewStravaOAuth(opts StravaOAuthOpts) *StravaOAuth {
	clientOpts := StravaClientOpts{
		HttpClient: opts.HTTPClient,
	}

	if opts.Scopes == nil || len(opts.Scopes) == 0 {
		opts.Scopes = []string{StravaOAuthScopes.Read}
	}

	return &StravaOAuth{
		StravaOAuthOpts: opts,
		client:          NewStravaClient(clientOpts),
	}
}

// Generates the authentication url that will be used to redirect the user to
// Strava so it can initiate the OAuthFlow
// Args:
//   - force: Use force to always show the authorization prompt even if the user has already authorized the current application, default is auto.
//   - state: Returned in the redirect URI. Useful if the authentication is done from various points in an app.
//
// Returns the auth code url the user should be redirected to.
func (oauth *StravaOAuth) GenerateAuthCodeURL(force bool, state string) string {
	url := fmt.Sprintf("%s?response_type=code&client_id=%s&client_secret=%s&redirect_uri=%s&scope=%s", endpoints.Auth, oauth.ClientID, oauth.ClientSecret, oauth.CallbackURL, strings.Join(oauth.Scopes, ","))

	if force {
		url = fmt.Sprintf("%s&approval_prompt=force", url)
	}

	if len(state) > 0 {
		url = fmt.Sprintf("%s&state=%s", url, state)
	}

	return url
}

type StravaOAuthResponse struct {
	AccessToken  string         `json:"access_token"`
	ExpiresAt    uint64         `json:"expires_at"`    // The number of seconds since the epoch when the provided access token will expire
	ExpiresIn    int            `json:"expires_in"`    // Seconds until the short-lived access token will expire
	RefreshToken string         `json:"refresh_token"` // The refresh token for this user, to be used to get the next access token for this user. Please expect that this value can change anytime you retrieve a new access token. Once a new refresh token code has been returned, the older code will no longer work.
	TokenType    string         `json:"token_type"`    // Bearer
	Athlete      SummaryAthlete `json:"athlete"`       // A summary of the athlete information
	Scopes       []string       // Current Scopes the users agreed on. Not part of the JSON payload sent by Strava
}

func (oauth *StravaOAuth) Exchange(code string, scopes []string) (*StravaOAuthResponse, error) {
	formData := url.Values{
		"client_id":     {oauth.ClientID},
		"client_secret": {oauth.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?%s", endpoints.Token, formData.Encode()), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-url-encoded")

	var resp StravaOAuthResponse
	err = oauth.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	if len(scopes) == 0 {
		resp.Scopes = scopes
	}

	
	return &resp, nil
}

type refreshTokenPayload struct {
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`  // The short-lived access token
	RefreshToken string `json:"refresh_token"` // The number of seconds since the epoch when the provided access token will expire
	ExpiresAt    uint64 `json:"expires_at"`    // Seconds until the short-lived access token will expire
	ExpiresIn    int    `json:"expires_in"`    // The refresh token for this user, to be used to get the next access token for this user. Please expect that this value can change anytime you retrieve a new access token. Once a new refresh token code has been returned, the older code will no longer work.
}

func (oauth *StravaOAuth) Refresh(refreshToken string) (*RefreshTokenResponse, error) {
	buf, err := json.Marshal(refreshTokenPayload{
		ClientID:     oauth.ClientID,
		ClientSecret: oauth.ClientSecret,
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, endpoints.Token, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	var resp RefreshTokenResponse
	if err := oauth.client.do(req, resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (oauth *StravaOAuth) RevokeAccess(access_token string) error {
	formData := url.Values{
		"access_token": {access_token},
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?%s", endpoints.Unauthorize, formData.Encode()), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-url-encoded")

	return oauth.client.do(req, nil)
}

func (oauth *StravaOAuth) HandlerFunc(
	handleSuccess func(tokens *StravaOAuthResponse, w http.ResponseWriter, r *http.Request),
	handleError func(err error, w http.ResponseWriter, r *http.Request),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		errParam := params.Get("error")

		if errParam == "access_denied" {
			handleError(&StravaOAuthError{}, w, r)
			return
		}

		tokens, err := oauth.Exchange(params.Get("code"), strings.Split(params.Get("scope"), ","))
		if err != nil {
			handleError(err, w, r)
		}

		handleSuccess(tokens, w, r)
	}
}
