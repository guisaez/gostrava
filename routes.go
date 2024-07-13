package gostrava

import (
	"net/url"
	"strconv"
)

type RouteService service

// Returns a route using its identifier. Requires read_all scope for private routes.
func (s *RouteService) GetById(accessToken string, id int) (*Route, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "routes/" + strconv.Itoa(id),
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
func (s *RouteService) ExportRouteGPX(accessToken string, id int) ([]byte, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "routes/" + strconv.Itoa(id) + "/export_gpx",
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
func (s *RouteService) ExportRouteTCX(accessToken string, id int) ([]byte, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "routes/" + strconv.Itoa(id) + "/export_tcx",
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
func (s *RouteService) ListAthleteRoutes(accessToken string, athleteID int, opts RequestParams) ([]Route, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "athletes/" + strconv.Itoa(athleteID) + "/routes",
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
