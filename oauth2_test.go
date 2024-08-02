package gostrava

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// Test for AuthorizationCodeURL
func TestAuthorizationCodeURL(t *testing.T) {
	client := New()
	options := OAuthFlowOptions{
		State: "test_state",
		Force: true,
	}
	expectedURL := "https://www.strava.com/oauth/authorize?approval_prompt=force&client_id=&redirect_uri=&response_type=code&scope=&state=test_state"

	got := client.AuthorizationCodeURL("", options)
	if got != expectedURL {
		t.Errorf("AuthorizationCodeURL() = %v; want %v", got, expectedURL)
	}
}

// Test for MakeAuthorizationCodeURL
func TestMakeAuthorizationCodeURL(t *testing.T) {
	clientID := "test_client_id"
	clientSecret := "test_client_secret"
	redirectURI := "http://localhost"
	scopes := []Scope{Read, ActivityWrite}
	options := OAuthFlowOptions{
		State: "test_state",
		Force: false,
	}
	expectedURL := "https://www.strava.com/oauth/authorize?client_id=test_client_id&redirect_uri=http%3A%2F%2Flocalhost&response_type=code&scope=read%2Cactivity%3Awrite&state=test_state"

	got := MakeAuthorizationCodeURL(clientID, clientSecret, redirectURI, scopes, options)
	if got != expectedURL {
		t.Errorf("MakeAuthorizationCodeURL() = %v; want %v", got, expectedURL)
	}
}

// Test for AuthorizationTokenRequest
func TestAuthorizationTokenRequest(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := AccessTokenResponse{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
			ExpiresAt:    1715689200,
			ExpiresIn:    3600,
			TokenType:    "Bearer",
			Athlete:      nil,
			Scopes:       []Scope{Read, ActivityWrite},
		}
		json.NewEncoder(w).Encode(response)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &Strava{
		client: &stravaHTTPClient{httpClient: server.Client()},
	}
	code := "auth_code"
	expected := AccessTokenResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
		ExpiresAt:    1715689200,
		ExpiresIn:    3600,
		TokenType:    "Bearer",
		Athlete:      nil,
		Scopes:       []Scope{Read, ActivityWrite},
	}

	got, err := client.AuthorizationTokenRequest(code)
	if err != nil {
		t.Fatalf("AuthorizationTokenRequest() error = %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("AuthorizationTokenRequest() = %v; want %v", got, expected)
	}
}

// Test for RefreshTokenRequest
func TestRefreshTokenRequest(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := RefreshTokenResponse{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
			ExpiresAt:    1715689200,
			ExpiresIn:    3600,
		}
		json.NewEncoder(w).Encode(response)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &Strava{
		client: &stravaHTTPClient{httpClient: server.Client()},
	}
	refreshToken := "refresh_token"
	expected := RefreshTokenResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
		ExpiresAt:    1715689200,
		ExpiresIn:    3600,
	}

	got, err := client.RefreshTokenRequest(refreshToken)
	if err != nil {
		t.Fatalf("RefreshTokenRequest() error = %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("RefreshTokenRequest() = %v; want %v", got, expected)
	}
}

// Test for joinScopes function
func TestJoinScopes(t *testing.T) {
	scopes := []Scope{Read, ActivityWrite}
	expected := "read,activity:write"

	got := joinScopes(scopes)
	if got != expected {
		t.Errorf("joinScopes() = %v; want %v", got, expected)
	}
}
