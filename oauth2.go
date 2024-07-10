package gostrava

import "net/http"

const oauthBaseUrl = "https://www.strava.com/oauth/authorize"

type Scope string

const (
	Read            Scope = "read"
	ReadAll         Scope = "read_all"
	ProfileReadAll  Scope = "profile:read_all"
	ProfileWrite    Scope = "profile:write"
	ActivityRead    Scope = "activity:read"
	ActivityReadAll Scope = "activity:read_all"
	ActivityWrite   Scope = "write"
)

type OAuthService struct {
	// The application's ID, obtained during registration
	ClientID string

	// The application's secret, obtained during registration
	ClientSecret string

	// CallbackURL
	CallbackURL string

	// Scopes the application will be trying to access
	Scopes     []Scope

	httpClient *http.Client
}



