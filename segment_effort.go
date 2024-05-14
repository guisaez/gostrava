package go_strava

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type ListSegmentEffortOptions struct {
	StartDateLocal time.Time
	EndDateLocal   time.Time
	PerPage        int // Number of items per page. Defaults to 30.
}

// Returns a segment effort from an activity that is owned by the authenticated athlete. Requires subscription.
func (sc *StravaClient) GetSegmentEffort(ctx context.Context, id int64) (*DetailedSegmentEffort, error) {

	path := fmt.Sprintf("/segment_efforts/%d", id)

	var resp DetailedSegmentEffort
	if err := sc.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Returns a set of the authenticated athlete's segment efforts for a given segment. Requires subscription
func (sc *StravaClient) ListSegmentEfforts(ctx context.Context, opt *ListSegmentEffortOptions) (*DetailedSegmentEffort, error) {

	params := url.Values{}

	if opt != nil {
		if opt.StartDateLocal.IsZero() {
			params.Set("start_date_local", opt.StartDateLocal.Format(time.RFC3339))
		}
		if opt.EndDateLocal.IsZero() {
			params.Set("end_date_local", opt.EndDateLocal.Format(time.RFC3339))
		}
		if opt.PerPage > 0 {
			params.Set("per_page", fmt.Sprintf("%d", opt.PerPage))
		}
	}

	var resp DetailedSegmentEffort
	if err := sc.get(ctx, "/segment_efforts", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
