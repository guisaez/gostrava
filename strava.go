package gostrava

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
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




var timeStampType = reflect.TypeOf(TimeStamp{})

// Stringify attempts to create a reasonable string representation of types in
// the GitHub library. It does things like resolve pointers to their values
// and omits struct fields with nil values.
func Stringify(message interface{}) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(message)
	stringifyValue(&buf, v)
	return buf.String()
}

// stringifyValue was heavily inspired by the goprotobuf library.

func stringifyValue(w *bytes.Buffer, val reflect.Value) {
	if val.Kind() == reflect.Ptr && val.IsNil() {
		w.Write([]byte("<nil>"))
		return
	}

	v := reflect.Indirect(val)

	switch v.Kind() {
	case reflect.String:
		fmt.Fprintf(w, `"%s"`, v)
	case reflect.Slice:
		w.Write([]byte{'['})
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				w.Write([]byte{' '})
			}

			stringifyValue(w, v.Index(i))
		}

		w.Write([]byte{']'})
		return
	case reflect.Struct:
		if v.Type().Name() != "" {
			w.Write([]byte(v.Type().String()))
		}

		// special handling of Timestamp values
		if v.Type() ==  timeStampType{
			fmt.Fprintf(w, "{%s}", v.Interface())
			return
		}

		w.Write([]byte{'{'})

		var sep bool
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Slice && fv.IsNil() {
				continue
			}
			if fv.Kind() == reflect.Map && fv.IsNil() {
				continue
			}

			if sep {
				w.Write([]byte(", "))
			} else {
				sep = true
			}

			w.Write([]byte(v.Type().Field(i).Name))
			w.Write([]byte{':'})
			stringifyValue(w, fv)
		}

		w.Write([]byte{'}'})
	default:
		if v.CanInterface() {
			fmt.Fprint(w, v.Interface())
		}
	}
}
