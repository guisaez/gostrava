package oauth2

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/guisaez/gostrava"
)

const oauthBaseUrl string = "https://www.strava.com/oauth"

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

type OAuth struct {
	// The application's ID, obtained during registration
	ClientID string

	// The application's secret, obtained during registration
	ClientSecret string

	// Scopes the application will be trying to access
	Scopes []Scope

	client *gostrava.Client
}

func Register(clientId, clientSecret string, scopes ...Scope) *OAuth {
	return &OAuth{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       scopes,

		client: gostrava.NewClient(nil),
	}
}

type Authorization struct {
	AccessToken  string                   `json:"access_token"`
	ExpiresAt    int64                    `json:"expires_at"`           // The number of seconds since the epoch when the provided access token will expire
	ExpiresIn    int                      `json:"expires_in"`           // Seconds until the short-lived access token will expire
	RefreshToken string                   `json:"refresh_token"`        // The refresh token for this user, to be used to get the next access token for this user. Please expect that this value can change anytime you retrieve a new access token. Once a new refresh token code has been returned, the older code will no longer work
	TokenType    *string                  `json:"token_type,omitempty"` // Bearer
	Athlete      *gostrava.AthleteSummary `json:"athlete,omitempty"`    // A summary of the athlete information
	Scopes       []Scope                  `json:"scopes,omitempty"`     // Scopes the user accepted
}

// This function handles the exchange step of an authorization code for an acces token in the
// OAuth 2.0 authorization code grant flow.
//
// POST: "https://www.strava.com/oauth/token"
func (oauth *OAuth) Exchange(code string, scopes string) (*Authorization, error) {
	formData := url.Values{
		"client_id":     {oauth.ClientID},
		"client_secret": {oauth.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	}

	url, _ := url.Parse(oauthBaseUrl)

	req, err := oauth.client.NewRequest(gostrava.RequestOpts{
		URL:    url,
		Path:   "token",
		Body:   formData,
		Method: http.MethodPost,
	})
	if err != nil {
		return nil, err
	}

	auth := new(Authorization)
	if err := oauth.client.Do(req, auth); err != nil {
		return nil, err
	}

	auth.Scopes = splitScopes(scopes)

	return auth, nil
}

// Generates the authentication url that the user will be redirected to
// in order to initiate the OAuthFlow
// Args:
//   - force: If true, it will always show the authorization prompt even if the user has already
//     authorized the current application.
//   - state: Returned in the redirect URI. Useful if the authentication is done from various points in an app
func (oauth *OAuth) MakeAuthCodeURL(callbackUrl string, force bool, state ...string) *url.URL {
	authorizeUrl, _ := url.Parse(oauthBaseUrl)
	authorizeUrl = authorizeUrl.JoinPath("authorize")

	queryParams := authorizeUrl.Query()
	queryParams.Set("response_type", "code")
	queryParams.Set("client_id", oauth.ClientID)
	queryParams.Set("client_secret", oauth.ClientSecret)
	queryParams.Set("redirect_uri", callbackUrl)
	queryParams.Set("scope", joinScopes(oauth.Scopes))

	if force {
		queryParams.Set("approval_prompt", "force")
	} else {
		queryParams.Set("approval_prompt", "auto")
	}

	if len(state) != 0 {
		queryParams.Set("state", state[0])
	}

	authorizeUrl.RawQuery = queryParams.Encode()

	return authorizeUrl
}

// This function handles the process of using a refresh token to obtain a new access token in the
// OAuth 2.0 authorization flow.
//
// POST "https://www.strava.com/oauth/refresh"
func (oauth *OAuth) Refresh(refreshToken string) (*Authorization, error) {
	formData := url.Values{
		"client_id":     {oauth.ClientID},
		"client_secret": {oauth.ClientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	}

	url, _ := url.Parse(oauthBaseUrl)

	req, err := oauth.client.NewRequest(gostrava.RequestOpts{
		URL:    url,
		Path:   "token",
		Body:   formData,
		Method: http.MethodPost,
	})
	if err != nil {
		return nil, err
	}

	refresh := new(Authorization)
	if err := oauth.client.Do(req, refresh); err != nil {
		return nil, err
	}

	return refresh, nil
}

// This function will invalidate all refresh_tokens and access_tokens that the application has for the athlete.
//
// POST "https://www.strava.com/oauth/deathorize"
func (oauth *OAuth) RevokeAccess(accessToken string) error {
	formData := url.Values{
		"access_token": {accessToken},
	}

	url, _ := url.Parse(oauthBaseUrl)

	req, err := oauth.client.NewRequest(gostrava.RequestOpts{
		URL:    url,
		Path:   "deauthorize",
		Method: http.MethodPost,
		Body:   formData,
	})
	if err != nil {
		return err
	}

	return oauth.client.Do(req, nil)
}

type OAuthError struct {
	Message string
}

func (e *OAuthError) Error() string {
	return e.Message
}

// OAuthHandler returns an HTTP handler function for handling OAuth authorization responses.
// This handler processes the incoming HTTP request to extract OAuth parameters and invokes
// appropriate callback functions based on success or error scenarios.
//
// Parameters:
//   - onSuccess: A callback function that will be invoked when the OAuth authorization is successful.
//     It accepts an authorization token of type *Authorization, an HTTP response writer (http.ResponseWriter),
//     and an HTTP request (http.Request).
//   - onError: A callback function that will be invoked when there is an error during the OAuth authorization process.
//     It accepts an error (error), an HTTP response writer (http.ResponseWriter), and an HTTP request (http.Request).
//
// Returns:
// - An HTTP handler function (http.HandlerFunc) that processes the OAuth authorization response.
func (oauth *OAuth) OAuthHandler(
	onSuccess func(token *Authorization, w http.ResponseWriter, r *http.Request),
	onError func(err error, w http.ResponseWriter, r *http.Request),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		errParam := query.Get("error")

		if errParam == "" {
			onError(&OAuthError{Message: errParam}, w, r)
			return
		}

		code := query.Get("code")
		scopes := query.Get("scope")

		auth, err := oauth.Exchange(code, scopes)
		if err != nil {
			onError(err, w, r)
		}

		onSuccess(auth, w, r)
	}
}

// ------- Utils --------

func splitScopes(scopes string) []Scope {
	splittedScopes := strings.Split(scopes, ",")
	parsedScopes := make([]Scope, len(splittedScopes))
	for i, scope := range splittedScopes {
		parsedScopes[i] = Scope(scope)
	}
	return parsedScopes
}

func joinScopes(scopes []Scope) string {
	stringScopes := make([]string, len(scopes))
	for i, scope := range scopes {
		stringScopes[i] = string(scope)
	}

	return strings.Join(stringScopes, ",")
}
