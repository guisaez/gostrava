package gostrava

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type StravaSegments baseModule

// Returns the specified segment, read_all scope required in order to retrieve athlete specific segment information,
// or to retrieve private segments.
func (sc *StravaSegments) GetById(ctx context.Context, access_token string, id int64) (*DetailedSegment, error) {
	path := fmt.Sprintf("/segments/%d", id)

	var res DetailedSegment
	if err := sc.client.get(ctx, access_token, path, nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type Bounds struct {
	SWLat float64
	SWLng float64
	NELat float64
	NELng float64
}

type ExploreSegmentsOpts struct {
	ActivityType string // Desired activity type. May take one of the following values: running, riding.
	MinCat       int    // The minimum climbing category
	MaxCat       int    // The maximum climbing category
}

// Returns the top 10 segments matching a specified query.
func (sc *StravaSegments) Explore(ctx context.Context, access_token string, bounds Bounds, opt *ExploreSegmentsOpts) ([]ExplorerResponse, error) {

	params := url.Values{}
	params.Set("bounds", bounds.toString())

	if opt != nil {
		if opt.ActivityType != "" {
			params.Set("activity_type", opt.ActivityType)
		}
		if opt.MinCat > 0 {
			params.Set("min_cat", fmt.Sprintf("%d", opt.MinCat))
		}
		if opt.MaxCat > 0 {
			params.Set("max_cat", fmt.Sprintf("%d", opt.MaxCat))
		}
	}

	var resp []ExplorerResponse
	if err := sc.client.get(ctx, access_token, "/segments/explore", params, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// List of the authenticated athlete's starred segments. Private segments are filtered out unless requested by a token with read_all scope.
func (sc *StravaSegments) ListStarred(ctx context.Context, access_token string, opt *GeneralParams) ([]SummarySegment, error) {

	params := url.Values{}
	if opt != nil {
		if opt.Page > 0 {
			params.Set("page", strconv.Itoa(opt.Page))
		}
		if opt.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(opt.Page))
		}
	}

	var resp []SummarySegment
	if err := sc.client.get(ctx, access_token, "/segments/starred", params, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// // TO-DO
// func (sc *StravaSegments) StarSegment (ctx context.Context, id uint64, starred bool) (*DetailedSegment, error) {
	
// 	path := fmt.Sprintf("/segments/%d/starred")

// 	var res DetailedSegment
// 	if err := sc.do(ctx, http.MethodGet, sc.AccessToken, path, params, nil, &res); err != nil {
// 		return nil, err
// 	}

// 	return &res, nil
// }
