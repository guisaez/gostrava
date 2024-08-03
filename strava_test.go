package gostrava

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestDoAndParseSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	client := NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)

	req, _ := client.NewRequest(http.MethodGet, "", nil)

	var result map[string]string
	_, err := client.DoAndParse(context.Background(), req, &result)
	if err != nil {
		t.Fatalf("DoAndParse() error = %v", err)
	}
	if result["message"] != "success" {
		t.Errorf("DoAndParse() result = %v, want %v", result["message"], "success")
	}
}

func TestDoAndParseErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"errors":null,"message":"internal server error"}`))
	}))
	defer server.Close()

	client := NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)

	req, _ := client.NewRequest(http.MethodGet, "/", nil)
	var result map[string]string
	_, err := client.DoAndParse(context.Background(), req, &result)
	if err == nil {
		t.Fatal("DoAndParse() expected an error but got none")
	}

	fault, ok := err.(*Fault)
	if !ok {
		t.Fatalf("DoAndParse() error = %v, want type *Fault", err)
	}
	if fault.Error() != `{"errors":null,"message":"internal server error"}` {
		t.Errorf("DoAndParse() fault error = %v, want %v", fault.Error(), `{"errors":null,"message":"internal server error"}`)
	}
}

func TestDoAndParseEmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)

	req, _ := client.NewRequest("GET", "/", nil)
	var result map[string]interface{}
	_, err := client.DoAndParse(context.Background(), req, &result)
	if err != nil {
		t.Fatalf("DoAndParse() error = %v", err)
	}
	if len(result) != 0 {
		t.Errorf("DoAndParse() result = %v, want empty map", result)
	}
}

func TestDoAndParseInvalidJSON(t *testing.T) {
	// Create a test server that returns malformed JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"`)) // Malformed JSON
	}))
	defer server.Close()

	// Create a new client and set the base URL to the test server's URL
	client := NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)

	// Create a new request
	req, err := client.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}

	// Define a variable to hold the decoded response
	var result map[string]string

	// Execute the request and decode the response
	_, err = client.DoAndParse(context.Background(), req, &result)
	if err == nil {
		t.Fatal("DoAndParse() expected an error but got none")
	}
}

func TestDoAndParseContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a long response
		select {
		case <-r.Context().Done():
			return
		case <-time.After(time.Second):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"success"}`))
		}
	}))
	defer server.Close()

	client := NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)

	req, _ := client.NewRequest(http.MethodGet, "/", nil)

	// Create a context that will be cancelled immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.DoAndParse(ctx, req, nil)
	if err == nil {
		t.Fatal("DoAndParse() expected an error due to cancelled context but got none")
	}
	if err != context.Canceled {
		t.Errorf("DoAndParse() error = %v, want %v", err, context.Canceled)
	}
}
