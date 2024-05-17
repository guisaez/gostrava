package gostrava

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type RoundTripFunc func(r *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestDoRequest(t *testing.T) {
	
	errorPayload := &Error{
		Message: "Authorization Error",
		Errors: []ErrorContent{
			{
				Code: "invalid",
				Field: "access_token",
				Resource: "Athlete",
			},
		},
	}

	jsonM, err := json.Marshal(errorPayload)
	if err != nil {
		t.Errorf("error marshalling payload - not related to test case")
	}
	
	fn := func(req *http.Request) *http.Response {

		if v := req.Header.Get("Authorization"); v == "" {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Status: "401 Unauthorized",
				Body: io.NopCloser(bytes.NewBuffer(jsonM)),
			}
		}
		
		return &http.Response{
			StatusCode: http.StatusOK,
			Status: "200 Ok",
			Body: io.NopCloser(bytes.NewBufferString("{\"message\":\"OK\"}")),
		}
	}

	// No Access Token
	strava := &StravaClient{
		client: NewTestClient(fn),
	}

	req, _ := http.NewRequest(http.MethodGet, "http://gostravatest.com", nil)

	var response interface{}
	err = strava.do_request(req, response)
	if err == nil {
		t.Errorf("expected an error %s, got nil", errorPayload.Error())
	}

	// With Access Token
	accessToken := "123"
	req.Header.Add("Authorization", "Bearer " + accessToken)
	var response2 interface{}
	err = strava.do_request(req, response2)
	if err != nil {
		t.Errorf("did not expect an error: %s", err.Error())
	}
}