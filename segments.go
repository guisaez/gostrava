package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SegmentsAPIService apiService

// Returns the specified segment, read_all scope required in order to retrieve athlete specific segment information,
// or to retrieve private segments.
func (s *SegmentsAPIService) GetSegment(access_token string, id int64) (*DetailedSegment, error) {
	requestUrl := s.client.BaseURL.JoinPath(segmentsPath, fmt.Sprint(id))

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := &DetailedSegment{}
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
func (s *SegmentsAPIService) ExploreSegments(access_token string, bounds Bounds, opt *ExploreSegmentsOpts) ([]ExplorerResponse, error) {
	requestUrl := s.client.BaseURL.JoinPath(segmentsPath, "explore")
	
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

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
		body: params,
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
func (s *SegmentsAPIService) ListStarredSegments(access_token string, opt *RequestParams) ([]SummarySegment, error) {
	requestUrl := s.client.BaseURL.JoinPath(segmentsPath, "starred")
	
	params := url.Values{}
	if opt != nil {
		if opt.Page > 0 {
			params.Set("page", strconv.Itoa(opt.Page))
		}
		if opt.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(opt.Page))
		}
	}

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
		body: params,
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
func (s *SegmentsAPIService) StarSegment(access_token string, id int64, starred bool) (*DetailedSegment, error){
	requestUrl := s.client.BaseURL.JoinPath(segmentsPath, fmt.Sprint(id), "starred")
	
	formData := url.Values{}
	formData.Add("starred", fmt.Sprint(starred))

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodPut,
		access_token: access_token,
		body: formData,
	})
	if err != nil {
		return nil, err
	}

	resp := &DetailedSegment{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
