package gostrava

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type StravaSegmentEfforts baseModule

// Returns a segment effort from an activity that is owned by the authenticated athlete. Requires subscription.
func (sc *StravaSegmentEfforts) GetById(ctx context.Context,access_token string, id int64) (*DetailedSegmentEffort, error) {

	path := fmt.Sprintf("/segment_efforts/%d", id)

	var resp DetailedSegmentEffort
	if err := sc.client.get(ctx, access_token, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type ListSegmentEffortOptions struct {
	StartDateLocal time.Time
	EndDateLocal   time.Time
	PerPage        int // Number of items per page. Defaults to 30.
}

// Returns a set of the authenticated athlete's segment efforts for a given segment. Requires subscription
func (sc *StravaSegmentEfforts) List(ctx context.Context, access_token string, opt *ListSegmentEffortOptions) (*DetailedSegmentEffort, error) {

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
	if err := sc.client.get(ctx, access_token, "/segment_efforts", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
