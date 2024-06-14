package gostrava

import (
	"encoding/json"
	"io"
	"net/http"
)

var baseURL string = "https://www.strava.com/api/v3"

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

type StravaClientOpts struct {
	HttpClient HTTPClient
}

type StravaClient struct {
	StravaClientOpts
}

func NewStravaClient(opts StravaClientOpts) *StravaClient {
	if opts.HttpClient == nil {
		opts.HttpClient = http.DefaultClient
	}

	return &StravaClient{
		opts,
	}
}

func (c *StravaClient) do(r *http.Request, v interface{}) error {
	r.Header.Add("Accept", "application/json")

	resp, err := c.HttpClient.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

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
