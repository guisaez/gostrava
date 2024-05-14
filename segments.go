package go_strava

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type Bounds struct {
	SWLat float64
	SWLng float64
	NELat float64
	NELng float64
}

// TO-DO Be More specific here!
type ExploreSegmentsOpts struct {
	ActivityType string // Desired activity type. May take one of the following values: running, riding.
	MinCat       int    // The minimum climbing category
	MaxCat       int    // The maximum climbing category
}

// Returns the specified segment, read_all scope required in order to retrieve athlete specific segment information,
// or to retrieve private segments. 
func (sc *StravaClient) GetSegmentById(ctx context.Context, id int64) (*DetailedSegment, error) {
	path := fmt.Sprintf("/segments/%d", id)

	var res DetailedSegment
	if err := sc.get(ctx, path, nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Returns the top 10 segments matching a specified query.
func (sc *StravaClient) ExploreSegments(ctx context.Context, bounds Bounds, opt *ExploreSegmentsOpts) ([]ExplorerResponse, error) {

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
	if err := sc.get(ctx, "/segments/explore", params, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// List of the authenticated athlete's starred segments. Private segments are filtered out unless requested by a token with read_all scope.
func (sc *StravaClient) ListStarredSegments(ctx context.Context, opt *GeneralParams) ([]SummarySegment, error) {

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
	if err := sc.get(ctx, "/segments/starred", params, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// // TO-DO
// func (sc *StravaClient) StarSegment (ctx context.Context, id uint64, starred bool) (*DetailedSegment, error) {
	
// 	path := fmt.Sprintf("/segments/%d/starred")

// 	var res DetailedSegment
// 	if err := sc.do(ctx, http.MethodGet, path, params, nil, &res); err != nil {
// 		return nil, err
// 	}

// 	return &res, nil
// }


func (b *Bounds) toString() string {
	return fmt.Sprintf("%2f,%2f,%2f,%2f", b.SWLat, b.SWLng, b.NELat, b.NELng)
}