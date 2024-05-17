package gostrava

import "encoding/json"

type Error struct {
	Errors  []ErrorContent `json:"errors"`
	Message string         `json:"message"`
}

type ErrorContent struct {
	Code     string `json:"code"`
	Field    string `json:"field"`
	Resource string `json:"resource"`
}

func (e *Error) Error() string {
	err, _ := json.Marshal(e)
	return string(err)
}

type StravaOAuthError struct {
	Message string
}

func (e *StravaOAuthError) Error() string {
	return e.Message
}

var (
	InvalidCodeError = &StravaOAuthError{"invalid code"}
	AccessDeniedError     = &StravaOAuthError{"access_denied"}
	InternalServerError = &Error{Message: "internal server error"}
)