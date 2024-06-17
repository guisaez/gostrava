package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestNewStravaOAuth(t *testing.T) {
	var (
		clientId     = "1234567890"
		clientSecret = "12345678910ABCD"
		callbackURL  = "https://testing.com"
		scopes       = []string{StravaScopes.Read}
	)

	opts := OAuthOpts{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		CallbackURL:  callbackURL,
		Scopes:       scopes,
	}

	NewStravaOAuth(opts)
}

func TestMakeAuthCodeURL(t *testing.T) {
	oauth := makeOAuthClient()

	oauth.Scopes = []string{StravaScopes.Read}

	authCodeUrl, err := url.Parse(oauth.MakeAuthCodeURL(false, ""))
	if err != nil {
		t.Error("error not expected")
	}

	query := authCodeUrl.Query()
	if !query.Has("client_id") || !query.Has("client_secret") || !query.Has("redirect_uri") || !query.Has("response_type") || !query.Has("scope") {
		t.Errorf("missing query parameters, got %s", query)
	}
}

func TestOAuthExchange(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockStravaOAuthServer()

	oauth := NewStravaOAuth(OAuthOpts{})

	valid_code := "valid_code_12345"
	invalid_code := "invalid_code_12345"

	scopes := []string{StravaScopes.Read}

	// Empty Client Id or Wrong Client Id
	_, err := oauth.Exchange(valid_code, scopes)
	if err == nil {
		t.Error("expected an error")
	}

	fmt.Println(err.Error())

	oauth.ClientID = "12345"
	oauth.ClientSecret = "12345"

	// Empty Callback URL
	_, err = oauth.Exchange(invalid_code, scopes)
	if err == nil {
		t.Error("expected an error")
	}

	fmt.Println(err.Error())

	// Invalid Code
	_, err = oauth.Exchange(invalid_code, scopes)
	if err == nil {
		t.Error("expected an error")
	}

	_, err = oauth.Exchange(valid_code, scopes)
	if err != nil {
		t.Errorf("error not expected, got %v", err)
	}
}

func TestOAuthRefresh(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockStravaOAuthServer()

	refreshToken := ""

	oauth := makeOAuthClient()

	// Empty // Wrong refresh Token

	_, err := oauth.Refresh(refreshToken)
	if err == nil {
		t.Error("expected and error")
	}

	refreshToken = "asd123.."
	tokens, err := oauth.Refresh(refreshToken)
	if err != nil {
		t.Errorf("error not expected - %v", err)
	}

	expectedResponse := &RefreshTokenResponse{
		AccessToken:  "a9b723...",
		RefreshToken: "b5c569...",
		ExpiresAt:    1568775134,
		ExpiresIn:    20566,
	}

	if !reflect.DeepEqual(expectedResponse, tokens) {
		t.Errorf("expected %v, got %v", expectedResponse, tokens)
	}
}

func makeOAuthClient() *OAuth {
	return NewStravaOAuth(OAuthOpts{
		ClientID:     "1234567890",
		ClientSecret: "12345678910ABCD",
		CallbackURL:  "https://testing.com",
	})
}

func mockStravaOAuthServer() {
	httpmock.RegisterResponder("POST", endpoints.Token,
		func(req *http.Request) (*http.Response, error) {
			err := req.ParseForm()
			if err != nil {
				return nil, err
			}

			query := req.Form

			if query.Get("client_id") == "" {
				return httpmock.NewJsonResponse(http.StatusBadRequest, httpmock.File("./mock/oauth/oauth_exchange_bad_request_client_id.json"))
			}

			if query.Get("client_secret") == "" {
				return httpmock.NewJsonResponse(http.StatusBadRequest, httpmock.File("./mock/oauth/oauth_exchange_bad_request_client_id.json"))
			}

			switch query.Get("grant_type") {
			case "authorization_code":
				if query.Get("code") == "" || query.Get("code") == "invalid_code_12345" {
					return httpmock.NewJsonResponse(http.StatusBadRequest, httpmock.File("./mock/oauth/oauth_exchange_bad_request_invalid_code.json"))
				}
				return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/oauth/oauth_exchange_success.json"))
			case "refresh_token":
				if query.Get("refresh_token") == "" {
					return httpmock.NewJsonResponse(http.StatusBadRequest, httpmock.File("./mock/oauth/oauth_refresh_token_bad_request_refresh_token.json"))
				}
				return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/oauth/oauth_refresh_token_success.json"))
			default:
				return httpmock.NewJsonResponse(http.StatusInternalServerError, nil)
			}
		},
	)
}
