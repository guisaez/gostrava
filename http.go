package gostrava

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	contentTypeApplicationJSON = "application/json"
	contentTypeFormURLEncoded  = "application/x-www-form-urlencoded"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// // stravaHTTPClient handles HTTP requests for the Strava API.
// type stravaHTTPClient struct {
// 	httpClient HTTPClient
// }

// // RequestOptions holds the configuration for a single HTTP request.
// type RequestOptions struct {
// 	Method      string            // HTTP method to be used for the request.
// 	URL         string            // The full URL for the request.
// 	URLParams   *url.Values       // Optional URL parameters to be included in the request.
// 	Body        interface{}       // Optional request body, which can be of type io.Reader or other.
// 	Headers     map[string]string // Optional headers to be included in the request.
// 	ContentType string            // Optional content type for the request body.
// 	Context     context.Context   // Optional context to control cancellation and timeouts.
// }

// // NewRequest creates and sends an HTTP request based on the provided options and response data.
// //
// // Parameters:
// //   - options (RequestOptions): Configuration for the HTTP request, including method, URL, URL parameters, request body, headers, content type, and context.
// //   - responseData (interface{}): A pointer to a variable where the response data will be stored. This can be a pointer to any type that matches the expected response format.
// //     If the response is JSON, the value pointed to should be of a type that matches the JSON structure.
// //     If the response is not JSON, the pointer should be of type *[]byte to capture the raw response body.
// //
// // Returns:
// //   - error: Returns an error if the request fails or if there are issues with decoding the response. The error includes details about the HTTP request failure,
// //     or if there are issues with decoding JSON or reading the response body.
// func (c *stravaHTTPClient) NewRequest(options RequestOptions, responseData interface{}) error {
// 	if responseData == nil || reflect.TypeOf(responseData).Kind() != reflect.Ptr {
// 		return fmt.Errorf("responseData must be a non-nil pointer")
// 	}

// 	// Set default context if not provided
// 	if options.Context == nil {
// 		options.Context = context.Background()
// 	}

// 	// Set default content type if not provided
// 	if options.ContentType == "" {
// 		options.ContentType = contentTypeFormURLEncoded
// 	}

// 	var req *http.Request
// 	var err error

// 	// Create request based on HTTP method and content type
// 	if options.Method == http.MethodPost || options.Method == http.MethodPut {
// 		if options.ContentType == contentTypeApplicationJSON {
// 			if body, ok := options.Body.(io.Reader); ok {
// 				req, err = http.NewRequestWithContext(options.Context, options.Method, options.URL, body)
// 				if err != nil {
// 					return err
// 				}
// 				req.Header.Set("Content-Type", contentTypeApplicationJSON)
// 			} else {
// 				return fmt.Errorf("body must be of type io.Reader for JSON content type")
// 			}
// 		} else {
// 			if options.Body != nil {
// 				var reqBody io.Reader
// 				if params, ok := options.Body.(url.Values); ok {
// 					reqBody = strings.NewReader(params.Encode())
// 					req, err = http.NewRequestWithContext(options.Context, options.Method, options.URL, reqBody)
// 					if err != nil {
// 						return err
// 					}
// 					req.Header.Set("Content-Type", contentTypeFormURLEncoded)
// 				} else {
// 					return fmt.Errorf("body must be url.Values for form-encoded content type")
// 				}
// 			} else {
// 				req, err = http.NewRequestWithContext(options.Context, options.Method, options.URL, nil)
// 				if err != nil {
// 					return err
// 				}
// 			}
// 		}
// 	} else {
// 		req, err = http.NewRequestWithContext(options.Context, options.Method, options.URL, nil)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// Set default headers
// 	req.Header.Set("Accept", contentTypeApplicationJSON)

// 	// Add or overwrite headers
// 	for key, value := range options.Headers {
// 		req.Header.Set(key, value)
// 	}

// 	// Perform the HTTP request
// 	resp, err := c.httpClient.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("error sending HTTP request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	// Handle non-successful HTTP response
// 	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
// 		var fault Fault
// 		if err := json.NewDecoder(resp.Body).Decode(&fault); err != nil {
// 			return fmt.Errorf("error decoding error response body: %w", err)
// 		}
// 		return fmt.Errorf("status_code: %d, %s", resp.StatusCode, &fault)
// 	}

// 	// Handle successful response and decode it
// 	if strings.Contains(resp.Header.Get("Content-Type"), contentTypeApplicationJSON) {
// 		if err := json.NewDecoder(resp.Body).Decode(responseData); err != nil {
// 			return fmt.Errorf("error decoding response body: %w", err)
// 		}
// 	} else {
// 		respBytes, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			return fmt.Errorf("error reading response body: %w", err)
// 		}
// 		if byteSlice, ok := responseData.(*[]byte); ok {
// 			*byteSlice = respBytes
// 		} else {
// 			return fmt.Errorf("responseData must be of type *[]byte for non-JSON content")
// 		}
// 	}

// 	return nil
// }
