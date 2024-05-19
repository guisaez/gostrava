package gostrava

import (
	"context"
	"fmt"
)

type StravaRoutes baseModule

// Returns a route using its identifier. Requires read_all scope for private routes.
func (sc *StravaRoutes) GetById(ctx context.Context, access_token string, id int64) (*Route, error) {

	path := fmt.Sprintf("/routes/%d", id)

	var resp Route
	if err := sc.client.get(ctx, access_token, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Returns a GPX file of the route. Required read_all scope for private routes.
// ExportRouteGPX returns a GPX file of the route.
func (sc *StravaRoutes) ExportRouteGPX(ctx context.Context, access_token string, routeID int64) ([]byte, error) {
    path := fmt.Sprintf("/routes/%d/export_gpx", routeID)

    var gpxData []byte
    err := sc.client.get(ctx, access_token, path, nil, &gpxData)
    if err != nil {
        return nil, err
    }

    return gpxData, nil
}

// Returns a TCX file of the route.. Requires read_all scope for private routes.
func (sc *StravaRoutes) ExportRouteTCX(ctx context.Context, access_token string, routeID int64) ([]byte, error) {
    path := fmt.Sprintf("/routes/%d/export_tcx", routeID)

    // Create an empty params URL.Values since there are no query parameters

    var tcxData []byte
    err := sc.client.get(ctx, access_token, path, nil, &tcxData)
    if err != nil {
        return nil, err
    }

    return tcxData, nil
}