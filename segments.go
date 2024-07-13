package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SegmentsService service

// Returns the specified segment, read_all scope required in order to retrieve athlete specific segment information,
// or to retrieve private segments.
func (s *SegmentsService) GetById(accessToken string, id int) (*SegmentDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "segments/" + strconv.Itoa(id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(SegmentDetailed)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type Bounds struct {
	SWLat float32
	SWLng float32
	NELat float32
	NELng float32
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
func (s *SegmentsService) ExploreSegments(accessToken string, bounds Bounds, opts ExploreSegmentsOpts) (*ExplorerResponse, error) {
	params := url.Values{}
	params.Set("bounds", bounds.String())

	if opts.ActivityType != "" {
		params.Set("activity_type", opts.ActivityType)
	}
	if opts.MinCat > 0 {
		params.Set("min_cat", fmt.Sprintf("%d", opts.MinCat))
	}
	if opts.MaxCat > 0 {
		params.Set("max_cat", fmt.Sprintf("%d", opts.MaxCat))
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "segments/explore",
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := new(ExplorerResponse)
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// List of the authenticated athlete's starred segments. Private segments are filtered out unless requested by a token with read_all scope.
func (s *SegmentsService) ListStarredSegments(accessToken string, opts RequestParams) ([]SegmentSummary, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.Page))
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "segments/starred",
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := []SegmentSummary{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Stars/Unstars the given segment for the authenticated athlete. Requires profile:write scope.
func (s *SegmentsService) StarSegment(accessToken string, id int, starred bool) (*SegmentDetailed, error) {
	formData := url.Values{}
	formData.Add("starred", fmt.Sprint(starred))

	req, err := s.client.newRequest(requestOpts{
		Path:        "segments/" + strconv.Itoa(id) + "/starred",
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
