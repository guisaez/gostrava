package oauth

import (
	strava "github.com/guisaez/go-strava"
)

// Token exchange result
type StravaOAuthResponse struct {
	AccessToken  string                `json:"access_token"`
	RefreshToken string                `json:"refresh_token"` // The refresh token for this user, to be used to get the next access token for this user. Please expect that this value can change anytime you retrieve a new access token. Once a new refresh token code has been returned, the older code will no longer work.
	TokenType    string                `json:"token_type"`    // Bearer
	ExpiresAt    uint64                `json:"expires_at"`    // The number of seconds since the epoch when the provided access token will expire
	ExpiresIn    int                   `json:"expires_in"`    // Seconds until the short-lived access token will expire
	Athlete      strava.SummaryAthlete `json:"athlete"`       // A summary of the athlete information
	Scopes 		 []string
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
	"activity_read_all",
	"activity:write",
}

var endpoints = struct {
	Auth string
	Token string
}{
	"https://www.strava.com/oauth/authorize",
	"https://www.strava.com/oauth/token",
}