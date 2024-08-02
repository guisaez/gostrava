package gostrava

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestStrava_New(t *testing.T) {
	strava := New()
	if strava == nil {
		t.Fatalf("expected non-nil client")
	}
	if strava.client.httpClient == nil || strava.client.httpClient != http.DefaultClient {
		t.Fatalf("expected default HTTP client")
	}
}

func TestStrava_SetCredentials(t *testing.T) {
	client := New()
	client.SetCredentials("clientID", "clientSecret")

	if client.ClientID() != "clientID" {
		t.Errorf("expected clientID to be 'clientID', got '%s'", client.ClientID())
	}
	if client.ClientSecret() != "clientSecret" {
		t.Errorf("expected clientSecret to be 'clientSecret', got '%s'", client.ClientSecret())
	}
}

func TestStrava_SetScopes(t *testing.T) {
	client := New()
	scopes := []Scope{"read", "write"}
	client.SetScopes(scopes)

	if !reflect.DeepEqual(client.Scopes(), scopes) {
		t.Errorf("expected scopes %v, got %v", scopes, client.Scopes())
	}
}

func TestStrava_UseCustomHTTPClient(t *testing.T) {
	client := New()
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"key":"value"}`))),
			}, nil
		},
	}
	client.UseCustomHTTPClient(mockClient)

	options := RequestOptions{
		Method: http.MethodGet,
		URL:    "https://www.example.com",
	}

	var responseData []byte
	err := client.CustomRequest(options, &responseData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStrava_CustomRequest_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]string{"key": "value"}
		json.NewEncoder(w).Encode(response)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &stravaHTTPClient{
		httpClient: server.Client(),
	}

	var responseData map[string]string
	options := RequestOptions{
		Method: http.MethodGet,
		URL:    server.URL,
	}

	err := client.NewRequest(options, &responseData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]string{"key": "value"}
	if !reflect.DeepEqual(responseData, expected) {
		t.Errorf("expected %v, got %v", expected, responseData)
	}
}

func TestStrava_CustomRequest_Error(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := Fault{
			Errors: []Error{
				Error{Code: "test_code", Field: "test_field", Resource: "test_resource"},
				
			},
		}
		json.NewEncoder(w).Encode(errorResponse)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &stravaHTTPClient{
		httpClient: server.Client(),
	}

	var responseData map[string]string
	options := RequestOptions{
		Method: http.MethodGet,
		URL:    server.URL,
	}

	err := client.NewRequest(options, &responseData)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedError := `status_code: 400, {"errors":[{"code":"test_code","field":"test_field","resource":"test_resource"}],"message":""}`
	if err.Error() != expectedError {
		t.Errorf("expected %v, got %v", expectedError, err)
	}
}

func TestStrava_CustomRequest_NonJSONResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("plain text response"))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &stravaHTTPClient{
		httpClient: server.Client(),
	}

	var responseData []byte
	options := RequestOptions{
		Method: http.MethodGet,
		URL:    server.URL,
	}

	err := client.NewRequest(options, &responseData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []byte("plain text response")
	if !bytes.Equal(responseData, expected) {
		t.Errorf("expected %s, got %s", expected, responseData)
	}
}
