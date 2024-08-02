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

type myResponse struct {
	Key string `json:"key"`
}

func TestNewRequest_Success(t *testing.T) {
	// Create a mock server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := myResponse{Key: "value"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &stravaHTTPClient{
		httpClient: server.Client(),
	}

	// Initialize ResponseData with a pointer to an empty struct
	var responseData myResponse
	options := RequestOptions{
		Method: http.MethodGet,
		URL:    server.URL,
	}

	err := client.NewRequest(options, &responseData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := myResponse{Key: "value"}

	// Verify that the response data matches the expected value
	if !reflect.DeepEqual(responseData, expected) {
		t.Errorf("expected %v, got %v", expected, responseData)
	}
}

func TestNewRequest_ErrorResponse(t *testing.T) {
	// Create a mock server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fault := Fault{
			Message: "Bad Request",
			Errors: []Error{
				{Code: "invalid_request", Field: "field", Resource: "resource"},
			},
		}
		json.NewEncoder(w).Encode(fault)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &stravaHTTPClient{
		httpClient: server.Client(),
	}

	var responseData struct{} // An empty struct to test error handling
	options := RequestOptions{
		Method: http.MethodGet,
		URL:    server.URL,
	}

	err := client.NewRequest(options, &responseData)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedError := `status_code: 400, {"errors":[{"code":"invalid_request","field":"field","resource":"resource"}],"message":"Bad Request"}`
	if err.Error() != expectedError {
		t.Errorf("expected %v, got %v", expectedError, err)
	}
}

func TestNewRequest_BodyHandling(t *testing.T) {
	// Create a mock server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodPost && r.Header.Get("Content-Type") == contentTypeApplicationJSON {
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("error decoding request body: %v", err)
			}
			if body["key"] != "value" {
				t.Errorf("expected key=value, got %v", body)
			}
		} else {
			t.Errorf("unexpected method or content type")
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &stravaHTTPClient{
		httpClient: server.Client(),
	}

	body := map[string]string{"key": "value"}
	bodyReader := io.NopCloser(bytes.NewReader(mustMarshal(t, body)))
	options := RequestOptions{
		Method:      http.MethodPost,
		URL:         server.URL,
		Body:        bodyReader,
		ContentType: contentTypeApplicationJSON,
	}

	// Use an empty struct to satisfy the `responseData` requirement
	var responseData []byte
	err := client.NewRequest(options, &responseData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Helper function to marshal data to JSON and handle errors
func mustMarshal(t *testing.T, v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("error marshaling JSON: %v", err)
	}
	return data
}
