package gostrava

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestNewStravaOAuth(t *testing.T) {
	var (
		clientId     = "1234567890"
		clientSecret = "12345678910ABCD"
		callbackURL  = "https://testing.com"
		scopes       = []string{StravaOAuthScopes.Read}
	)

	opts := StravaOAuthOpts{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		CallbackURL:  callbackURL,
		Scopes:       scopes,
	}

	NewStravaOAuth(opts)
}

func TestGenerateAuthCodeURL(t *testing.T) {
	oauth := makeOAuthClient()

	var (
		clientId     = "1234567890"
		clientSecret = "12345678910ABCD"
		callbackURL  = "https://testing.com"
	)
	expected1 := fmt.Sprintf("%s?response_type=code&client_id=%s&client_secret=%s&redirect_uri=%s&scope=%s", endpoints.Auth, clientId, clientSecret, callbackURL, StravaOAuthScopes.Read)
	actual1 := oauth.GenerateAuthCodeURL(false, "")

	if expected1 != actual1 {
		t.Errorf("expected read scope to be automatically added if not provided, got %s", actual1)
	}

	expected2 := fmt.Sprintf("%s&approval_prompt=force", expected1)
	actual2 := oauth.GenerateAuthCodeURL(true, "")

	if expected2 != actual2 {
		t.Errorf("expected %s, got %s", expected2, actual2)
	}
}

func TestOAuthExchange(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", endpoints.Token,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/oauth_exchange_success.json"))
			return resp, err
		},
	)

	oauth := makeOAuthClient()

	code := "test_code_12345"
	scopes := []string{StravaOAuthScopes.ActivityRead}

	tokenResponse, err := oauth.Exchange(code, scopes)
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}

	expectedResponse := &StravaOAuthResponse{
		AccessToken:  "a4b945687g...",
		RefreshToken: "e5n567567...",
		ExpiresIn:    21600,
		ExpiresAt:    1568775134,
		TokenType:    "Bearer",
		Athlete: SummaryAthlete{
			MetaAthlete: MetaAthlete{
				ID: 1234553636,
			},
			ResourceState: 2,
			FirstName:     "John",
			LastName:      "Doe",
			ProfileMedium: "http://example.com/profile_medium.jpg",
			Profile:       "http://example.com/profile.jpg",
			City:          "CityName",
			State:         "StateName",
			Country:       "CountryName",
			Sex:           "M",
			Premium:       false,
			Summit:        true,
			CreatedAt:     "2018-02-16T14:56:25Z",
			UpdatedAt:     "2018-02-16T14:56:25Z",
		},
	}

	if !reflect.DeepEqual(expectedResponse, tokenResponse) {
		t.Errorf("expected %v, go %v", expectedResponse, tokenResponse)
	}
}

func makeOAuthClient() *StravaOAuth {
	return NewStravaOAuth(StravaOAuthOpts{
		ClientID:     "1234567890",
		ClientSecret: "12345678910ABCD",
		CallbackURL:  "https://testing.com",
	})
}
