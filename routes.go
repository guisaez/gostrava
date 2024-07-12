package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type RouteService service

const routes string = "routes"



// Returns a route using its identifier. Requires read_all scope for private routes.
func (s *RouteService) GetById(accessToken string, id int64) (*Route, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d", routes, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(Route)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a GPX file of the route. Required read_all scope for private routes.
// ExportRouteGPX returns a GPX file of the route.
func (s *RouteService) ExportRouteGPX(accessToken string, id int64) ([]byte, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/export_gpx", routes, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
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
func (s *RouteService) ExportRouteTCX(accessToken string, id int64) ([]byte, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/export_tcx", routes, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
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
func (s *RouteService) ListAthleteRoutes(accessToken string, athleteID int64, p *RequestParams) ([]Route, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/%s", athletes, athleteID, routes),
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
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
