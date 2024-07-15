package gostrava

import (
	"encoding/json"
	"errors"
	"strconv"
)

type RouteDetailed struct {
	RouteSummary
	Segments []*SegmentSummary `json:"segments"` // The segments traversed by this route
}

type RouteSummary struct {
	Athlete             *AthleteSummary  `json:"athlete"`               // An instance of AthleteSummary.
	Description         string           `json:"description"`           // The description of the route
	Distance            float32          `json:"distance"`              // The route's distance, in meters
	ElevationGain       float32          `json:"elevation_gain"`        // The route's elevation gain.
	ID                  int              `json:"id"`                    // The unique identifier of this route
	IdStr               string           `json:"id_str"`                // The unique identifier of the route in string format
	Map                 PolylineSummmary `json:"map"`                   // An instance of PolylineMap.
	MapUrls             Urls             `json:"map_urls"`              //
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

// *****************************************************

type RouteService service

// Returns a route using its identifier. Requires read_all scope for private routes.
func (s *RouteService) GetById(accessToken string, id int) (*RouteDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "routes/" + strconv.Itoa(id),
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(RouteDetailed)
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
