package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SegmentEffortsAPIService apiService

// Returns a segment effort from an activity that is owned by the authenticated athlete. Requires subscription.
func (s *SegmentEffortsAPIService) GetSegmentEffort(access_token string, id int64) (*DetailedSegmentEffort, error) {
	requestUrl := s.client.BaseURL.JoinPath(segmentEffortsPath, fmt.Sprint(id))

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := &DetailedSegmentEffort{}
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
func (s *SegmentEffortsAPIService) ListSegmentEfforts(access_token string, opt *ListSegmentEffortOptions) ([]DetailedSegmentEffort, error) {
	requestUrl := s.client.BaseURL.JoinPath(segmentEffortsPath)

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

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := []DetailedSegmentEffort{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
