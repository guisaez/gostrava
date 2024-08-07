package gostrava

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// *************** Types ********************

type RouteSummary struct {
	Athlete             *AthleteSummary  `json:"athlete"`               // An instance of AthleteSummary.
	Description         string           `json:"description"`           // The description of the route
	Distance            float32          `json:"distance"`              // The route's distance, in meters
	ElevationGain       float32          `json:"elevation_gain"`        // The route's elevation gain.
	ID                  int              `json:"id"`                    // The unique identifier of this route
	IdStr               string           `json:"id_str"`                // The unique identifier of the route in string format
	Map                 PolylineSummary `json:"map"`                   // An instance of PolylineMap.
	MapUrls             URL              `json:"map_urls"`              //
	Name                string           `json:"name"`                  // The name of this route
	Private             bool             `json:"private"`               // Whether this route is private
	ResourceState       int8             `json:"resource_state"`        //
	Starred             bool             `json:"starred"`               // Whether this route is starred by the logged-in athlete
	SubType             SubRouteType     `json:"sub_type"`              // This route's sub-type: "road" / "mountain_bike" / "cross" / "trail" / "mixed"
	CreatedAt           TimeStamp        `json:"created_at"`            // The time at which the route was created
	UpdatedAt           TimeStamp        `json:"updated_at"`            // The time at which the route was last updated
	Timestamp           int              `json:"timestamp"`             // An epoch timestamp of when the route was created
	Type                RouteType        `json:"type"`                  // This route's type "ride" / "run"
	EstimatedMovingTime int              `json:"estimated_moving_time"` // Estimated time in seconds for the authenticated athlete to complete route
	Waypoints           []Waypoint       `json:"waypoints"`             // The custom waypoints along this route
}

type RouteDetailed struct {
	RouteSummary
	Segments []*SegmentSummary `json:"segments"` // The segments traversed by this route
}

type Waypoint struct {
	LatLng            LatLng   `json:"latlng"`              // The location along the route that the waypoint is closest to
	TargetLatLng      LatLng   `json:"target_latlng"`       // A location off of the route that the waypoint is (optional)
	Categories        []string `json:"categories"`          // Categories that the waypoint belongs to
	Title             string   `json:"string"`              // A title for the waypoint
	Description       string   `json:"description"`         // A description of the waypoint (optional)
	DistanceIntoRoute int      `json:"distance_into_route"` // The number meters along the route that the waypoint is located
}

type SubRouteType string

const (
	Road         SubRouteType = "road"
	MountainBike SubRouteType = "mountain_bike"
	Cross        SubRouteType = "cross"
	Trail        SubRouteType = "train"
	Mixed        SubRouteType = "mixed"
)

func (rt *SubRouteType) UnmarshalJSON(data []byte) error {
	var subRouteType int
	err := json.Unmarshal(data, &subRouteType)
	if err != nil {
		return err
	}

	switch subRouteType {
	case 1:
		*rt = Road
	case 2:
		*rt = MountainBike
	case 3:
		*rt = Cross
	case 4:
		*rt = Trail
	case 5:
		*rt = Mixed
	default:
		return errors.New("invalid sub-route type")
	}

	return nil
}

type RouteType string

const (
	RideRoute RouteType = "ride"
	RunRoute  RouteType = "run"
)

func (rt *RouteType) UnmarshalJSON(data []byte) error {
	var routeType int
	err := json.Unmarshal(data, &routeType)
	if err != nil {
		return err
	}

	switch routeType {
	case 1:
		*rt = RideRoute
	case 2:
		*rt = RunRoute
	default:
		return errors.New("invalid route type")
	}

	return nil
}

// *************** Methods********************

type RoutesService service

const routes string = "/api/v3/routes"

// GetById retrieves a route by its id. Requires read_all scope for private routes.
//
// GET: https://www.strava.com/api/v3/routes{id}
func (s *RoutesService) GetById(ctx context.Context, accessToken string, id int) (*RouteDetailed, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d", routes, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	route := new(RouteDetailed)
	resp, err := s.client.DoAndParse(ctx, req, route)
	if err != nil {
		return nil, resp, err
	}

	return route, resp, nil
}

// GetRouteStreams returns the corresponding route streams.
// Requires read_all scope for private routes.
//
// GET: https://www.strava.com/api/v3/routes/{id}/streams
func (s *RoutesService) GetRouteStreams(ctx context.Context, accessToken string, id int) ([]Stream, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/streams", routes, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var streams []Stream
	resp, err := s.client.DoAndParse(ctx, req, &streams)
	if err != nil {
		return nil, resp, err
	}

	return streams, resp, nil
}

// ExportRouteGPX returns a GPX binary of the route. Required read_all scope for private routes.
//
// GET: https://www.strava.com/api/v3/routes/{id}/export_gpx
func (s *RoutesService) ExportRouteGPX(ctx context.Context, accessToken string, id int) ([]byte, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/export_gpx", routes, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	gpx := []byte{}
	resp, err := s.client.DoAndParse(ctx, req, &gpx)
	if err != nil {
		return nil, resp, err
	}

	return gpx, resp, nil
}

// ExportRouteTCX returns a TCX file of the route. Required read_all scope for private routes.
//
// GET: https://www.strava.com/api/v3/routes/{id}/export_tcx
func (s *RoutesService) ExportRouteTCX(ctx context.Context, accessToken string, id int) ([]byte, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/export_tcx", routes, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	tcx := []byte{}
	resp, err := s.client.DoAndParse(ctx, req, &tcx)
	if err != nil {
		return nil, resp, err
	}

	return tcx, resp, nil
}
