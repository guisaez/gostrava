package gostrava

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var defaultBaseURL string = "https://www.strava.com/api/v3"

const (
	activitiesPath     = "activities"
	athletePath        = "athlete"
	athletesPath       = "athletes"
	clubsPath          = "clubs"
	gearPath           = "gear"
	routesPath         = "routes"
	segmentEffortsPath = "segment_efforts"
	segmentsPath       = "segments"
	streamPath         = "streams"
	uploadsPath        = "uploads"
)

type Client struct {
	// HTTP client used to communicate with API
	client *http.Client

	// Base URL for API requests
	BaseURL *url.URL

	// Clubs API service
	Activities     ActivitiesAPIService
	Athletes       AthleteAPIService
	Clubs          ClubsAPIService
	Gears          GearsAPIService
	Routes         RoutesAPIService
	SegmentEfforts SegmentEffortsAPIService
	Segments       SegmentsAPIService
}

type apiService struct {
	client *Client
}

type RequestParams struct {
	Page    int // Page number. Defaults to 1
	PerPage int // Number of items per page. Defaults to 30
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:  httpClient,
		BaseURL: baseURL,
	}
	c.Activities = ActivitiesAPIService{c}
	c.Athletes = AthleteAPIService{c}
	c.Clubs = ClubsAPIService{c}
	c.Gears = GearsAPIService{c}
	c.Routes = RoutesAPIService{c}
	c.SegmentEfforts = SegmentEffortsAPIService{c}
	c.Segments = SegmentsAPIService{c}

	return c
}

func (c *Client) do(req *http.Request, v interface{}) error {
	req.Header.Add("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errResp Error
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return err
		}
		return &errResp
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(&v)
	}

	return nil
}

type clientRequestOpts struct {
	// HTTP Method
	method string

	// Access Token if request needs authorization
	access_token string

	// URL the request will be made to
	url *url.URL

	// FormData to be sent in the request
	body interface{}

	// Context
	ctx context.Context
}

func (c *Client) newRequest(opts clientRequestOpts) (*http.Request, error) {
	if opts.method == "" {
		opts.method = http.MethodGet
	}

	if opts.ctx == nil {
		opts.ctx = context.Background()
	}

	if opts.url == nil {
		opts.url = c.BaseURL
	}

	var req *http.Request
	var err error

	if opts.method == http.MethodPost || opts.method == http.MethodPut {
		if b, ok := opts.body.(url.Values); ok {
			req, err = http.NewRequestWithContext(opts.ctx, opts.method, opts.url.String(), strings.NewReader(b.Encode()))
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else if b, ok := opts.body.(io.ReadCloser); ok {
			req, err = http.NewRequestWithContext(opts.ctx, opts.method, opts.url.String(), b)
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, err = http.NewRequestWithContext(opts.ctx, opts.method, opts.url.String(), nil)
			if err != nil {
				return nil, err
			}
		}
	} else { // GET Request, url values are appended to the query
		if opts.body != nil {
			if p, ok := opts.body.(url.Values); ok {
				opts.url.RawQuery = p.Encode()
			}
		}

		req, err = http.NewRequestWithContext(opts.ctx, opts.method, opts.url.String(), nil)
		if err != nil {
			return nil, err
		}
	}

	if opts.access_token != "" {
		req.Header.Set("Authorization", "Bearer "+opts.access_token)
	}

	return req, nil
}
