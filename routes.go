package go_strava

import (
	"context"
	"fmt"
)

// Returns a route using its identifier. Requires read_all scope for private routes.
func (sc *StravaClient) GetRoute(ctx context.Context, id int64) (*Route, error) {

	path := fmt.Sprintf("/routes/%d", id)

	var resp Route
	if err := sc.get(ctx,path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Returns a GPX file of the route. Required read_all scope for private routes.
// ExportRouteGPX returns a GPX file of the route.
func (sc *StravaClient) ExportRouteGPX(ctx context.Context, routeID int64) ([]byte, error) {
    path := fmt.Sprintf("/routes/%d/export_gpx", routeID)

    var gpxData []byte
    err := sc.get(ctx, path, nil, &gpxData)
    if err != nil {
        return nil, err
    }

    return gpxData, nil
}

// Returns a TCX file of the route.. Requires read_all scope for private routes.
func (sc *StravaClient) ExportRouteTCX(ctx context.Context, routeID int64) ([]byte, error) {
    path := fmt.Sprintf("/routes/%d/export_tcx", routeID)

    // Create an empty params URL.Values since there are no query parameters

    var tcxData []byte
    err := sc.get(ctx, path, nil, &tcxData)
    if err != nil {
        return nil, err
    }

    return tcxData, nil
}