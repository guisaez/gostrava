package gostrava

import (
	"net/http"
	"testing"
)

type RoundTripFunc func(r *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r), nil
}

func NewTestingClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestNewStravaClient(t *testing.T) {

	opts := StravaClientOpts{}

	client := NewStravaClient(opts)

	if client.HttpClient != http.DefaultClient {
		t.Errorf("expected HTTPClient to be automatically define if its not provided")
	}
}





