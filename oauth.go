package gostrava

import (
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
var StravaScopes = struct {
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

type OAuthOpts struct {
	ClientID     string // The application’s ID, obtained during registration.
	ClientSecret string // The application’s secret, obtained during registration.
	HttpClient   *http.Client
	CallbackURL  string
	Scopes       []string
}

type OAuth struct {
	OAuthOpts

	client *Client
}

func NewStravaOAuth(opts OAuthOpts) *OAuth {
	return &OAuth{
		OAuthOpts: opts,
		client:    NewClient(opts.HttpClient),
	}
}

// Generates the authentication url that will be used to redirect the user to
// Strava so it can initiate the OAuthFlow
// Args:
//   - force: Use force to always show the authorization prompt even if the user has already authorized the current application, default is auto.
//   - state: Returned in the redirect URI. Useful if the authentication is done from various points in an app.
//
// Returns the auth code url the user should be redirected to.
func (oauth *OAuth) MakeAuthCodeURL(force bool, state string) string {
	authCodeUrl, err := url.Parse(endpoints.Auth)
	if err != nil {
		return ""
	}

	query := authCodeUrl.Query()
	query.Set("response_type", "code")
	query.Set("client_id", oauth.ClientID)
	query.Set("client_secret", oauth.ClientSecret)
	query.Set("redirect_uri", oauth.CallbackURL)
	query.Set("scope", strings.Join(oauth.Scopes, ","))

	if force {
		query.Set("approval_prompt", "force")
	}

	if state != "" {
		query.Set("state", state)
	}

	authCodeUrl.RawQuery = query.Encode()

	return authCodeUrl.String()
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

func (oauth *OAuth) Exchange(code string, scopes []string) (*StravaOAuthResponse, error) {
	formData := url.Values{
		"client_id":     {oauth.ClientID},
		"client_secret": {oauth.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	}
	url, err := url.Parse(endpoints.Token)
	if err != nil {
		return nil, err
	}

	req, err := oauth.client.newRequest(newStravaRequestOpts{
		url:    url,
		body:   formData,
		method: http.MethodPost,
	})
	if err != nil {
		return nil, err
	}

	resp := &StravaOAuthResponse{}
	err = oauth.client.do(req, resp)
	if err != nil {
		return nil, err
	}

	if len(scopes) != 0 {
		resp.Scopes = scopes
	}

	return resp, nil
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`  // The short-lived access token
	RefreshToken string `json:"refresh_token"` // The number of seconds since the epoch when the provided access token will expire
	ExpiresAt    uint64 `json:"expires_at"`    // Seconds until the short-lived access token will expire
	ExpiresIn    int    `json:"expires_in"`    // The refresh token for this user, to be used to get the next access token for this user. Please expect that this value can change anytime you retrieve a new access token. Once a new refresh token code has been returned, the older code will no longer work.
}

func (oauth *OAuth) Refresh(refreshToken string) (*RefreshTokenResponse, error) {
	formData := url.Values{
		"client_id":     {oauth.ClientID},
		"client_secret": {oauth.ClientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	}

	url, err := url.Parse(endpoints.Token)
	if err != nil {
		return nil, err
	}
	req, err := oauth.client.newRequest(newStravaRequestOpts{
		url:    url,
		method: http.MethodPost,
		body:   formData,
	})
	if err != nil {
		return nil, err
	}

	resp := &RefreshTokenResponse{}
	if err := oauth.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (oauth *OAuth) RevokeAccess(access_token string) error {
	formData := url.Values{
		"access_token": {access_token},
	}

	url, err := url.Parse(endpoints.Unauthorize)
	if err != nil {
		return err
	}
	req, err := oauth.client.newRequest(newStravaRequestOpts{
		url:    url,
		method: http.MethodPost,
		body:   formData,
	})
	if err != nil {
		return err
	}

	return oauth.client.do(req, nil)
}

func (oauth *OAuth) HandlerFunc(
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
