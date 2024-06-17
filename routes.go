package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type RoutesAPIService apiService

// Returns a route using its identifier. Requires read_all scope for private routes.
func (s *RoutesAPIService) GetById(access_token string, id int64) (*Route, error) {
	requestUrl := s.client.BaseURL.JoinPath(routesPath, fmt.Sprint(id))
	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := &Route{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a GPX file of the route. Required read_all scope for private routes.
// ExportRouteGPX returns a GPX file of the route.
func (s *RoutesAPIService) ExportRouteGPX(access_token string, id int64) ([]byte, error) {
	requestUrl := s.client.BaseURL.JoinPath(routesPath, fmt.Sprint(id), "export_gpx")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := []byte{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a TCX file of the route.. Requires read_all scope for private routes.
func (s *RoutesAPIService) ExportRouteTCX(access_token string, id int64) ([]byte, error) {
	requestUrl := s.client.BaseURL.JoinPath(routesPath, fmt.Sprint(id), "export_tcx")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := []byte{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a list of the routes created by the authenticated athlete. Private routes are filtered out
// unless request by a token with read_all scope.
func (s *RoutesAPIService) ListAthleteRoutes(access_token string, id int64, p *RequestParams) ([]Route, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletesPath, fmt.Sprint(id), routesPath)

	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
		body:         params,
	})
	if err != nil {
		return nil, err
	}

	resp := []Route{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
