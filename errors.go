package go_strava

import "encoding/json"

type Error struct {
	Code     string `json:"code"`     // The code associated with this error.
	Field    string `json:"field"`    // The specific field or aspect of the resource associated with this error.
	Resource string `json:"resource"` // The type of resource associated with this error.
}

type Fault struct {
	Errors  []Error `json:"errors"`  // The set of specific errors associated with this fault, if any.
	Message string  `json:"message"` // The message of the fault.
}

func (f *Fault) Error() string {
	b, _ := json.Marshal(f)
	return string(b)
}

type StravaOAuthError struct {
	message string
}

// Implements error type
func (e *StravaOAuthError) Error() string {
	return e.message
}

type StravaClientError struct {
	message string
}

func (e *StravaClientError) Error() string {
	return e.message
}

var (
	OAuthAccessDeniedErr         = &StravaOAuthError{"access denied"}
	OAuthInvalidCredentialsError = &StravaOAuthError{"invalid client_id or client_secret"}
	OAuthInvalidCodeError        = &StravaOAuthError{"invalid code"}
	OAuthInternalError           = &StravaOAuthError{"internal server error"}
)
