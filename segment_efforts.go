package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SegmentEffortsService service

const segment_efforts string = "segment_efforts"

// Returns a segment effort from an activity that is owned by the authenticated athlete. Requires subscription.
func (s *SegmentEffortsService) GetSegmentEffort(accessToken string, id int64) (*SegmentEffortDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d", segment_efforts, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(SegmentEffortDetailed)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type ListSegmentEffortOptions struct {
	StartDateLocal time.Time
	EndDateLocal   time.Time
	RequestParams
}

// Returns a set of the authenticated athlete's segment efforts for a given segment. Requires subscription
func (s *SegmentEffortsService) ListSegmentEfforts(accessToken string, segmentID int, opt *ListSegmentEffortOptions) ([]SegmentEffortDetailed, error) {
	params := url.Values{}

	params.Set("segment_id", fmt.Sprintf("%d", segmentID))

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

	req, err := s.client.newRequest(requestOpts{
		Path:        segment_efforts,
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []SegmentEffortDetailed{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
