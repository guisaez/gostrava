package gostrava

import (
	"net/url"
	"strconv"
	"time"
)

// *****************************************************

type SegmentEffortsService service

// Returns a segment effort from an activity that is owned by the authenticated athlete. Requires subscription.
func (s *SegmentEffortsService) GetSegmentEffort(accessToken string, id int) (*SegmentEffortDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "segment_efforts/" + strconv.Itoa(id),
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
	Page           int // Page number. Defaults to 1
	PerPage        int // Number of items per page. Defaults to 30
	StartDateLocal time.Time
	EndDateLocal   time.Time
}

// Returns a set of the authenticated athlete's segment efforts for a given segment. Requires subscription
func (s *SegmentEffortsService) ListSegmentEfforts(accessToken string, segmentID int, opts ListSegmentEffortOptions) ([]SegmentEffortDetailed, error) {
	params := url.Values{}

	params.Set("segment_id", strconv.Itoa(segmentID))

	if opts.StartDateLocal.IsZero() {
		params.Set("start_date_local", opts.StartDateLocal.Format(time.RFC3339))
	}
	if opts.EndDateLocal.IsZero() {
		params.Set("end_date_local", opts.EndDateLocal.Format(time.RFC3339))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "segment_efforts",
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
