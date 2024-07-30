package gostrava

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type stravaHTTP struct {
	HTTPClient *http.Client
}

const contentTypeApplicationJSON = "application/json"
const contentTypeFormURLEnconded = "application/x-www-form-urlencoded; charset-UTF-8"

type requestOpts struct {
	Method       string
	URL          string
	Body         interface{}
	AccessToken  string
	ContentType  string
	URLParams    *url.Values
	ResponseData interface{}
}

func (sh *stravaHTTP) NewRequest(opts *requestOpts) error {
	return sh.NewRequestWithContext(context.Background(), opts)
}

func (sh *stravaHTTP) NewRequestWithContext(ctx context.Context, opts *requestOpts) error {
	if reflect.TypeOf(opts.ResponseData).Kind() != reflect.Ptr {
		return fmt.Errorf("failed, you must pass a pointer in the ResponseData field of the endpoint")
	}

	endpointURL := opts.URL

	var (
		req *http.Request
		err error
	)

	if opts.Method == http.MethodPost || opts.Method == http.MethodPut && opts.ContentType == contentTypeApplicationJSON{
		if body, ok := opts.Body.(io.ReadCloser); ok {
			req, err = http.NewRequestWithContext(ctx, opts.Method, opts.URL, body)
			if err != nil {
				return err
			}
		}
	} else{
		var reqBody *strings.Reader
		if opts.URLParams != nil {
			reqBody = strings.NewReader(opts.URLParams.Encode())
		
		}
		req, err = http.NewRequestWithContext(ctx, opts.Method,opts.URL, reqBody)
		if err != nil {
			return err
		}
		
		resp, err := sh.HTTPClient.Do(req)
		if err != nil {
			return err
		}

		
		
	}

	if opts.AccessToken != "" {
		req.Header.Add("Authorization", "Bearer " + opts.AccessToken)
	}
	req.Header.Add("Content-Type", opts.ContentType)
	req.Header.Add("Accept", "application/json")

}

type HTTPMethod string
