package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SegmentsService service

const segments string = "segments"

type SegmentDetailed struct {
	SummarySegment
	CreatedAt          TimeStamp   `json:"created_at"`           // The time at which the segment was created.
	UpdatedAt          TimeStamp   `json:"updated_at"`           // The time at which the segment was last updated.
	TotalElevationGain float32     `json:"total_elevation_gain"` // The segment's total elevation gain.
	Map                PolylineMap `json:"map"`                  // An instance of PolylineMap.
	EffortCount        int         `json:"effort_count"`         // The total number of efforts for this segment
	AthleteCount       int         `json:"athlete_count"`        // The number of unique athletes who have an effort for this segment
	Hazardous          bool        `json:"hazardous"`            // Whether this segment is considered hazardous
	StarCount          int         `json:"star_count"`           // The number of stars for this segment
}

type ExplorerResponse struct {
	Segments []ExplorerSegment `json:"segments"` // The set of segments matching an explorer request
}

type ExplorerSegment struct {
	ID                int     `json:"id"`                  // The unique identifier of this segment
	Name              string  `json:"name"`                // The name of this segment
	ClimbCategory     uint8   `json:"climb_category"`      // The category of the climb [0, 5]. Higher is harder ie. 5 is Hors catégorie, 0 is uncategorized in climb_category. If climb_category = 5, climb_category_desc = HC. If climb_category = 2, climb_category_desc = 3.
	ClimbCategoryDesc string  `json:"climb_category_desc"` // The description for the category of the climb May take one of the following values: NC, 4, 3, 2, 1, HC
	AvgGrade          float32 `json:"avg_grade"`           // The segment's average grade, in percents
	StartLatLng       LatLng  `json:"start_latlng"`        // An instance of LatLng.
	EndLatLng         LatLng  `json:"end_latlng"`          // An instance of LatLng.
	ElevationDiff     float32 `json:"elev_difference"`     // The segments's elevation difference, in meters
	Distance          float32 `json:"distance"`            // The segment's distance, in meters
	Points            string  `json:"points"`              // The polyline of the segment
}

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
