package go_strava

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

type StravaClient struct {
	accessToken  string
	baseURL      string
	client       *http.Client
}

func NewStravaClient(accessToken string, customClient *http.Client) *StravaClient {

	if customClient != nil {
		return &StravaClient{
			accessToken: accessToken,
			baseURL:     baseURL,
			client:      customClient,
		}
	}

	return &StravaClient{
		accessToken: accessToken,
		baseURL:     baseURL,
		client:      http.DefaultClient,
	}

}

type RequestOption struct {
	Params url.Values
	Body   io.Reader
}

func (sc *StravaClient) get(ctx context.Context, path string, params url.Values, v interface{}) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s?%v", baseURL, path, params.Encode()), nil)
	if err != nil {
		return err
	}

	return sc.do_request(req, v)
}

func (sc *StravaClient) postForm(ctx context.Context, path string, params url.Values, v interface{}) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+path, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-url/encoded")

	return sc.do_request(req, v)
}

func (sc *StravaClient) put(ctx context.Context, path, contentType string, body interface{}, v interface{}) error {

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

	return sc.do_request(req, v)
}

func (sc *StravaClient) do_request(r *http.Request, v interface{}) error {

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sc.accessToken))
	r.Header.Add("Accept", "application/json")

	resp, err := sc.client.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errResponse Fault
		if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
			return fmt.Errorf("unknown error, status code: %d", resp.StatusCode)
		}
		return &errResponse
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("unknown error, status code: %d", resp.StatusCode)
	}

	return nil
}
