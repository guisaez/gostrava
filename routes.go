package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type RouteService service

const routes string = "routes"

type Route struct {
	Athlete             AthleteSummary   `json:"athlete"`               // An instance of AthleteSummary.
	Description         string           `json:"description"`           // The description of the route
	Distance            float32          `json:"distance"`              // The route's distance, in meters
	ElevationGain       float32          `json:"elevation_gain"`        // The route's elevation gain.
	ID                  int              `json:"id"`                    // The unique identifier of this route
	IdStr               string           `json:"id_str"`                // The unique identifier of the route in string format
	Map                 PolylineMap      `json:"map"`                   // An instance of PolylineMap.
	Name                string           `json:"name"`                  // The name of this route
	Private             bool             `json:"private"`               // Whether this route is private
	Starred             bool             `json:"starred"`               // Whether this route is starred by the logged-in athlete
	Timestamp           int              `json:"timestamp"`             // An epoch timestamp of when the route was created
	Type                RouteType        `json:"type"`                  // This route's type RouteTypes.Ride, RouteTypes.Run
	SubType             SubRouteType     `json:"sub_type"`              // This route's sub-type (SubRouteTypes.Road, SubRouteTypes.MountainBike, SubRouteTypes.Cross, SubRouteTypes.Trail, SubRouteTypes.Mixed)
	CreatedAt           TimeStamp        `json:"created_at"`            // The time at which the route was created
	UpdatedAt           TimeStamp        `json:"updated_at"`            // The time at which the route was last updated
	EstimatedMovingTime int              `json:"estimated_moving_time"` // Estimated time in seconds for the authenticated athlete to complete route
	Segments            []SummarySegment `json:"segments"`              // The segments traversed by this route
	Waypoints           []Waypoint       `json:"waypoints"`             // The custom waypoints along this route
}

// Returns a route using its identifier. Requires read_all scope for private routes.
func (s *RouteService) GetById(accessToken string, id int64) (*Route, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d",routes, id),
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
