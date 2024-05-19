package gostrava

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var baseURL string = "https://www.strava.com/api/v3"

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

type StravaClient struct {
	ClientID     string // The application’s ID, obtained during registration.
	ClientSecret string // The application’s secret, obtained during registration.

	client HTTPClient

	Activities     *StravaActivities
	Athletes       *StravaAthletes
	Clubs          *StravaClubs
	Gears          *StravaGears
	Routes         *StravaRoutes
	SegmentEfforts *StravaSegmentEfforts
	Segments       *StravaSegments
	Streams        *StravaStreams
}

type baseModule struct {
	client *StravaClient
}

func NewStravaClient(clientID, clientSecret string, customClient HTTPClient) *StravaClient {
	if customClient == nil {
		customClient = http.DefaultClient
	}

	client := &StravaClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		client:       customClient,
	}

	// Initialize modules˝
	client.Activities = &StravaActivities{client: client}
	client.Athletes = &StravaAthletes{client: client}
	client.Clubs = &StravaClubs{client: client}
	client.Gears = &StravaGears{client: client}
	client.Routes = &StravaRoutes{client: client}
	client.SegmentEfforts = &StravaSegmentEfforts{client: client}
	client.Segments = &StravaSegments{client: client}
	client.Streams = &StravaStreams{client: client}

	return client
}

type RequestOption struct {
	Params url.Values
	Body   io.Reader
}

func (sc *StravaClient) do_request(r *http.Request, v interface{}) error {

	r.Header.Add("Accept", "application/json")

	resp, err := sc.client.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = handleBadResponse(resp)
	if err != nil {
		return err
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return &Error{
				Message: "could not parse response body, invalid JSON",
			}
		}
	}

	return nil
}

func (sc *StravaClient) get(ctx context.Context, access_token, path string, params url.Values, v interface{}) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s?%v", baseURL, path, params.Encode()), nil)
	if err != nil {
		return err
	}

	if access_token != "" {
		req.Header.Add("Authorization", "Bearer "+access_token)
	}

	return sc.do_request(req, v)
}

func (sc *StravaClient) postForm(ctx context.Context, access_token, path string, params url.Values, v interface{}) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+path, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-url/encoded")

	if access_token != "" {
		req.Header.Add("Authorization", "Bearer "+access_token)
	}

	if access_token != "" {
		req.Header.Add("Authorization", "Bearer "+access_token)
	}

	return sc.do_request(req, v)
}

func (sc *StravaClient) put(ctx context.Context, access_token, path, contentType string, body interface{}, v interface{}) error {

	var reqBody io.Reader
	switch contentType {
	case "application/json":
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(jsonData)
	case "application/x-www-form-urlencoded":
		formData, ok := body.(url.Values)
		if !ok {
			return errors.New("body is not of type url.Values")
		}
		reqBody = strings.NewReader(formData.Encode())
	default:
		return errors.New("unsupported content type")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, baseURL+path, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)

	if access_token != "" {
		req.Header.Add("Authorization", "Bearer "+access_token)
	}

	return sc.do_request(req, v)
}
