package gostrava

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const stravaBaseURL string = "https://www.strava.com/api/v3"

type Client struct {
	// Base URL user for API request
	BaseURL *url.URL

	// Debug tool that can be used to print each http response
	Logger func([]byte)

	// HTTP Client used to communicate with the server
	httpClient *http.Client

	OAuth   *OAuthService
	Athlete *AthleteService
}

type service struct {
	client *Client
}

// Initializes a new Strava Client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseUrl, _ := url.Parse(stravaBaseURL)

	c := &Client{
		BaseURL:    baseUrl,
		httpClient: httpClient,
	}

	c.OAuth = &OAuthService{client: c}
	c.Athlete = &AthleteService{client: c}

	return c
}

type requestOpts struct {
	// HTTP Method
	Method string

	// Access Token if request needs authorization
	AccessToken string

	// URL the request
	URL *url.URL

	// Path
	Path string

	// Request Body
	Body interface{}
}

func (c *Client) newRequest(opts requestOpts) (*http.Request, error) {
	if opts.Method == "" {
		opts.Method = http.MethodGet
	}

	if opts.URL == nil {
		opts.URL = c.BaseURL
	}

	if opts.Path != "" {
		opts.URL = opts.URL.JoinPath(opts.Path)
	}

	var req *http.Request
	var err error

	if opts.Method == http.MethodPost || opts.Method == http.MethodPut {
		// If request body is url.Values then the payload is formData
		if b, ok := opts.Body.(url.Values); ok {
			req, err = http.NewRequest(opts.Method, opts.URL.String(), strings.NewReader(b.Encode()))
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else if b, ok := opts.Body.(io.ReadCloser); ok {
			req, err = http.NewRequest(opts.Method, opts.URL.String(), b)
			if err != nil {
				return nil, err
			}
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, err = http.NewRequest(opts.Method, opts.URL.String(), nil)
			if err != nil {
				return nil, err
			}
		}
	} else { // GET requests
		if opts.Body != nil {
			if p, ok := opts.Body.(url.Values); ok {
				opts.URL.RawQuery = p.Encode()
			}
		}

		req, err = http.NewRequest(opts.Method, opts.URL.String(), nil)
		if err != nil {
			return nil, err
		}
	}

	if opts.AccessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", opts.AccessToken))
	}

	req.Header.Add("Accept", "application/json")

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		resp.Body.Close()
	}()

	var buf bytes.Buffer
	r := io.TeeReader(resp.Body, &buf)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errResponse Error

		if err := json.NewDecoder(r).Decode(&errResponse); err != nil {
			return err
		}

		if c.Logger != nil {
			c.Logger(buf.Bytes())
		}

		return &errResponse
	}

	if v != nil {
		err := json.NewDecoder(r).Decode(&v)

		if c.Logger != nil {
			c.Logger(buf.Bytes())
		}

		if err != nil {
			return err
		}
	}

	return nil
}
