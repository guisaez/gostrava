package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SegmentsService service

const segments string = "segments"

// Returns the specified segment, read_all scope required in order to retrieve athlete specific segment information,
// or to retrieve private segments.
func (s *SegmentsService) GetById(accessToken string, id int64) (*SegmentDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d", segments, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := &SegmentDetailed{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type Bounds struct {
	SWLat float64
	SWLng float64
	NELat float64
	NELng float64
}

func (b *Bounds) String() string {
	return fmt.Sprintf("%2f,%2f,%2f,%2f", b.SWLat, b.SWLng, b.NELat, b.NELng)
}

type ExploreSegmentsOpts struct {
	ActivityType string // Desired activity type. May take one of the following values: running, riding.
	MinCat       int    // The minimum climbing category
	MaxCat       int    // The maximum climbing category
}

// Returns the top 10 segments matching a specified query.
func (s *SegmentsService) ExploreSegments(accessToken string, bounds Bounds, opt *ExploreSegmentsOpts) ([]ExplorerResponse, error) {
	params := url.Values{}
	params.Set("bounds", bounds.String())

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

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/explore", segments),
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := []ExplorerResponse{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// List of the authenticated athlete's starred segments. Private segments are filtered out unless requested by a token with read_all scope.
func (s *SegmentsService) ListStarredSegments(accessToken string, opt *RequestParams) ([]SummarySegment, error) {
	params := url.Values{}
	if opt != nil {
		if opt.Page > 0 {
			params.Set("page", strconv.Itoa(opt.Page))
		}
		if opt.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(opt.Page))
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/starred", segments),
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := []SummarySegment{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Stars/Unstars the given segment for the authenticated athlete. Requires profile:write scope.
func (s *SegmentsService) StarSegment(accessToken string, id int64, starred bool) (*SegmentDetailed, error) {
	formData := url.Values{}
	formData.Add("starred", fmt.Sprint(starred))

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/starred", segments, id),
		Method:      http.MethodPut,
		AccessToken: accessToken,
		Body:        formData,
	})
	if err != nil {
		return nil, err
	}

	resp := new(SegmentDetailed)
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
