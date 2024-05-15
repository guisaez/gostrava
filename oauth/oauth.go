package oauth

import (
	"fmt"
	"net/http"
	"strings"
)

type StravaOAtuh struct {
	ClientID     string  // The application’s ID, obtained during registration.
	ClientSecret string  // The application’s secret, obtained during registration.
	CallbackURL  string  // URL to which the user will be redirected after authentication. Must be within the callback domain specified by the application. localhost and 127.0.0.1 are white-listed.
	Scopes       []string // Requested scopes, as a comma delimited string, e.g. "activity:read_all,activity:write". Applications should request only the scopes required for the application to function normally. The scope activity:read is required for activity webhooks.

	// Generates an HTTP Client for making requests during the OAuth token exchange process.
	RequestClient func(r *http.Request) *http.Client
}

// Generates the authentication url that will be used to redirect the user to
// Strava so it can initiate the OAuthFlow
// State will be returned by the redirect URI. Useful if the authentication is done from
// various points in the app.
// Force will always show the authorization prompt even if the user has already authorized the current
// application
func (oauth *StravaOAtuh) AuthCodeURL(state string, force bool) string {

	url := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&client_secret=%s&redirect_uri=%s&scope=%s&state=%s",
		endpoints.Authurl,
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


