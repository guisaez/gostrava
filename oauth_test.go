package gostrava

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGenerateAuthorizationURL(t *testing.T) {
	client := NewClient(nil)
	client.SetCredentials("test-client-id", "test-client-secret")

	service := &OAuthService{
		service:      service{client: client},
	}

	redirectURI := "https://example.com/callback"
	options := OAuthFlowOptions{
		State: "test-state",
		Force: true,
	}
	scopes := []Scope{Read, ActivityWrite}

	expectedURL := BuildAuthorizationURL(service.client.BaseURL.String(), service.client.clientID, redirectURI, scopes, options)
	authURL := service.AuthorizationURL(redirectURI, options, scopes...)

	if authURL != expectedURL {
		t.Errorf("AuthorizationURL() = %v, want %v", authURL, expectedURL)
	}
}

func TestBuildTokenRevocationURL(t *testing.T) {
	baseURL := "https://www.strava.com"
	clientID := "test-client-id"
	clientSecret := "test-client-secret"
	accessToken := "test-access-token"
	expectedURL := "https://www.strava.com/oauth/deauthorize?access_token=test-access-token"

	resultURL := BuildTokenRevocationURL(baseURL, clientID, clientSecret, accessToken)

	if resultURL != expectedURL {
		t.Errorf("BuildTokenRevocationURL() = %v, want %v", resultURL, expectedURL)
	}
}

func TestTokenRevocationURL(t *testing.T) {
	client := &Client{
		BaseURL: mustParseURL("https://www.strava.com"),
	}

	client.SetCredentials("test-client-id", "test-client-secret")

	service := &OAuthService{
		service:      service{client: client},
	}

	accessToken := "test-access-token"
	expectedURL := "https://www.strava.com/oauth/deauthorize?access_token=test-access-token"
	resultURL := service.TokenRevocationURL(accessToken)

	if resultURL != expectedURL {
		t.Errorf("TokenRevocationURL() = %v, want %v", resultURL, expectedURL)
	}
}

func TestExchangeAuthorizationCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("code") != "test-auth-code" {
			http.Error(w, "Invalid code", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "new-access-token",
			"refresh_token": "new-refresh-token",
			"expires_at": 1234567890,
			"expires_in": 3600,
			"token_type": "bearer",
			"athlete": {},
			"scopes": ["read", "activity:write"]
		}`))
	}))
	defer server.Close()

	client := NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)
	client.SetCredentials("test-client-id", "test-client-secret")

	service := &OAuthService{
		service:      service{client: client},
	}

	code := "test-auth-code"
	resp, httpResp, err := service.ExchangeAuthorizationCode(context.Background(), code)
	if err != nil {
		t.Fatalf("ExchangeAuthorizationCode() error = %v", err)
	}

	if resp.AccessToken != "new-access-token" {
		t.Errorf("ExchangeAuthorizationCode() AccessToken = %v, want new-access-token", resp.AccessToken)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Errorf("ExchangeAuthorizationCode() StatusCode = %v, want %v", httpResp.StatusCode, http.StatusOK)
	}
}

func TestRefreshToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"access_token": "refreshed-access-token",
			"refresh_token": "new-refresh-token",
			"expires_at": 1234567890,
			"expires_in": 3600
		}`))
	}))
	defer server.Close()

	client := NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)
	client.SetCredentials("test-client-id", "test-client-secret")

	service := &OAuthService{
		service:      service{client: client},
	}

	refreshToken := "test-refresh-token"
	resp, httpResp, err := service.RefreshToken(context.Background(), refreshToken)
	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}

	if resp.AccessToken != "refreshed-access-token" {
		t.Errorf("RefreshToken() AccessToken = %v, want refreshed-access-token", resp.AccessToken)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Errorf("RefreshToken() StatusCode = %v, want %v", httpResp.StatusCode, http.StatusOK)
	}
}

func TestDeauthorizeToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %v", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/oauth/deauthorize") {
			t.Errorf("Expected URL path to start with /oauth/deauthorize, got %v", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)
	client.SetCredentials("test-client-id", "test-client-secret")

	service := &OAuthService{
		service: service{client: client},
	}

	accessToken := "test-access-token"
	resp, err := service.RevokeToken(context.Background(), accessToken)
	if err != nil {
		t.Fatalf("RevokeToken() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("RevokeToken() StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func mustParseURL(rawURL string) *url.URL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return parsedURL
}
