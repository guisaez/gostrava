package go_strava

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"net/url"
	"strings"
)

type StravaOAuth struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
	Scopes       []string

	// Generates an HTTP Client for making requests during the OAuth token exchange process.
	RequestClient func(r *http.Request) *http.Client
}

var endpoint = struct {
	AuthURL string
	Token   string
}{
	"https://www.strava.com/oauth/authorize",
	"https://www.strava.com/oauth/token",
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
	"activity_read_all",
	"activity:write",
}

// Token exchange result
type StravaOAuthResp struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	TokenType    string         `json:"token_type"`
	ExpiresAt    uint64         `json:"expires_at"`
	ExpiresIn    int            `json:"expires_in"`
	Athlete      SummaryAthlete `json:"athlete"`
}

type TokenRefreshResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    uint64 `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
}

// Generates the authentication url that will be used to redirect the user to
// Strava so it can initiate the OAuthFlow
// State will be returned by the redirect URI. Useful if the authentication is done from
// various points in the app.
// Force will always show the authorization prompt even if the user has already authorized the current
// application
func (auth *StravaOAuth) AuthCodeURL(state string, force bool) string {

	url := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&client_secret=%s&redirect_uri=%s&scope=%s&state=%s",
		endpoint.AuthURL,
		auth.ClientID,
		auth.ClientSecret,
		auth.CallbackURL,
		strings.Join(auth.Scopes, ","),
		state,
	)

	if force {
		return url + "&approval_prompt=force"
	}

	return url
}

// Handles the token exchange portion of the OAuth2 flow.
func (auth *StravaOAuth) Exchange(code string, client *http.Client) (*StravaOAuthResp, error) {

	if code == "" {
		return nil, OAuthInvalidCodeError
	}

	if client == nil {
		client = http.DefaultClient
	}

	
	resp, err := client.PostForm(endpoint.Token, url.Values{
		"client_id":     {auth.ClientID},
		"client_secret": {auth.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var response StravaOAuthResp
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

func (auth *StravaOAuth) Refresh(refreshToken string, client *http.Client) (*TokenRefreshResp, error) {

	if client == nil {
		client = http.DefaultClient
	}

	jsonData, err := json.Marshal(RefreshTokenReq{
		ClientID:     auth.ClientID,
		ClientSecret: auth.ClientSecret,
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(endpoint.Token, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.New("internal server error", )
	}

	var response TokenRefreshResp
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (auth *StravaOAuth) HandlerFunc(
	handleSuccess func(tokens *StravaOAuthResp, w http.ResponseWriter, r *http.Request),
	handleError func(err error, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if errParam := r.URL.Query().Get("error"); errParam == "access_denied" {
			handleError(OAuthAccessDeniedErr, w, r)
			return
		}

		client := http.DefaultClient
		if auth.RequestClient != nil {
			client = auth.RequestClient(r)
		}

		tokens, err := auth.Exchange(r.URL.Query().Get("code"), client)
		if err != nil {
			handleError(err, w, r)
			return
		}

		handleSuccess(tokens, w, r)
	}
}




