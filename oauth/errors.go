package oauth

type StravaOAuthError struct {
	Message string
}

func (e *StravaOAuthError) Error() string {
	return e.Message
}

var (
	InvalidCodeError = &StravaOAuthError{"invalid code"}
	AccessDenied     = &StravaOAuthError{"access_denied"}
)
