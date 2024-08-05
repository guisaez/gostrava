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

const defaultBaseURL string = "https://www.strava.com/"

var errBadResponse = errors.New("got a bad response. check response object")

type service struct {
	client *Client
}

type ListOptions struct {
	Page    int `url:"page,omitempty"`
	PerPage int `url:"per_page,omitempty"`
}

type Client struct {
	clientID     string
	clientSecret string

	BaseURL *url.URL
	client  *http.Client

	common service // Use a single service struct instead of allocating one for each service on the heap.

	OAuth2     OAuthService
	Activities ActivityService
	Athletes   AthletesService
}

// NewClient creates a new Client instance with the given HTTP client. If no HTTP client is provided,
// the default HTTP client is used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	c := &Client{client: httpClient}
	c.initialize()
	return c
}

// SetCredentials configures the client ID, client secret, and optionally scopes for OAuth2.
func (s *Client) SetCredentials(clientID, clientSecret string, scopes ...Scope) *Client {
	s.clientID = clientID
	s.clientSecret = clientSecret

	s.OAuth2.scopes = scopes
	return s
}

func (c *Client) initialize() {
	if c.client == nil {
		c.client = http.DefaultClient
	}
	if c.BaseURL == nil {
		c.BaseURL, _ = url.Parse(defaultBaseURL)
	}

	c.common = service{client: c}

	c.OAuth2 = OAuthService{service: c.common}
	c.Activities = ActivityService(c.common)
	c.Athletes = AthletesService(c.common)
}

// RequestOption is a function that modifies an HTTP request.
type RequestOption func(req *http.Request) error

func SetAuthorizationHeader(accessToken string) RequestOption {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+accessToken)
		return nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
//
// The function handles different types of request bodies:
// - If the body is nil, no body is attached to the request (useful for GET requests).
// - If the body is of type url.Values, it is encoded as `application/x-www-form-urlencoded`.
// - If the body is a string, it is treated as raw JSON and the content type is set to `application/json`.
// - For other types (e.g., structs), the body is marshaled to JSON, and the content type is set to `application/json`.
//
// Parameters:
// - method: The HTTP method to use for the request (e.g., "GET", "POST").
// - urlStr: A relative URL to be resolved against the client's BaseURL.
// - body: The request body. Can be nil, url.Values, a JSON string, or any type that can be marshaled to JSON.
// - opts: Variadic arguments representing request options to further configure the request.
//
// Returns:
// - An *http.Request representing the created request.
// - An error if there is a problem creating the request or setting its properties.
func (c *Client) NewRequest(method, urlStr string, body interface{}, opts ...RequestOption) (*http.Request, error) {
	// Ensure BaseURL ends with a trailing slash
	if c.BaseURL == nil {
		return nil, fmt.Errorf("BaseURL is not set")
	}
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		c.BaseURL.Path += "/"
	}

	// Determine if urlStr is an absolute URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}

	var fullURL *url.URL
	if parsedURL.IsAbs() {
		// If the provided URL is absolute, use it directly
		fullURL = parsedURL
	} else {
		// Otherwise, resolve it relative to the BaseURL
		fullURL, err = c.BaseURL.Parse(urlStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing relative URL: %v", err)
		}
	}

	var (
		contentType string
		buf         io.Reader
	)

	if body != nil {
		switch v := body.(type) {
		case url.Values:
			if method == http.MethodGet {
				fullURL.RawQuery = v.Encode()
			} else {
				// Handle form data
				contentType = "application/x-www-form-urlencoded"
				buf = strings.NewReader(v.Encode())
			}
		case string:
			// Handle JSON body, expecting the body to be a JSON string
			contentType = "application/json"
			buf = strings.NewReader(v)
		default:
			// Handle struct and other types
			contentType = "application/json"
			data, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			buf = bytes.NewReader(data)
		}
	}

	// Create the HTTP request
	req, err := http.NewRequest(method, fullURL.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set the Content-Type header if needed
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("Accept", "application/json")

	// Apply request options
	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Do executes the given HTTP request using the client's HTTP client and returns the HTTP response.
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			// If the context was cancelled, return the context's error.
			return nil, ctx.Err()
		}
		return nil, err
	}

	return resp, nil
}

// DoAndParse executes the HTTP request and decodes the response into the given variable.
// It handles error responses and successful responses.
func (c *Client) DoAndParse(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle response based on status code
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		f := new(Fault)
		decodeErr := json.NewDecoder(resp.Body).Decode(f)
		if decodeErr == io.EOF {
			return resp, errBadResponse
		}
		if decodeErr != nil {
			return resp, fmt.Errorf("error decoding JSON response: %w", decodeErr)
		}
		return resp, f
	}

	// Handle successful response
	switch v := v.(type) {
	case nil:
		// Do nothing if v is nil
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decodeErr := json.NewDecoder(resp.Body).Decode(v)
		if decodeErr == io.EOF {
			// An empty response body is acceptable
			decodeErr = nil
		}
		if decodeErr != nil {
			return resp, fmt.Errorf("error decoding JSON response: %w", decodeErr)
		}
	}

	return resp, err
}
