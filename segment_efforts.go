package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SegmentEffortsService service

const segment_efforts string = "segment_efforts"

type SegmentEffortDetailed struct {
	Name         string         `json:"name"`              // The name of the segment on which this effort was performed
	Activity     ActivityMeta   `json:"activity"`          // An instance of MetaActivity.
	Athlete      AthleteMeta    `json:"athlete"`           // An instance of MetaAthlete.
	MovingTime   int            `json:"moving_time"`       // The effort's moving time
	StartIndex   int            `json:"start_index"`       // The start index of this effort in its activity's stream
	EndIndex     int            `json:"end_index"`         // The end index of this effort in its activity's stream
	AvgCadence   float32        `json:"average_cadence"`   // The effort's average cadence
	AverageWatts float32        `json:"average_watts"`     // The average wattage of this effort
	DeviceWatts  bool           `json:"device_watts"`      // For riding efforts, whether the wattage was reported by a dedicated recording device
	AvgHeartRate bool           `json:"average_heartrate"` // The heart heart rate of the athlete during this effort
	MaxHeartRate float32        `json:"max_heartrate"`     // The maximum heart rate of the athlete during this effort
	Segment      SummarySegment `json:"segment"`           // An instance of SummarySegment.
	KomRank      int            `json:"kom_rank"`          // The rank of the effort on the global leaderboard if it belongs in the top 10 at the time of upload
	PRRank       int            `json:"pr_rank"`           // The rank of the effort on the athlete's leaderboard if it belongs in the top 3 at the time of upload
	Hidden       bool           `json:"hidden"`            // Whether this effort should be hidden when viewed within an activity
}

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
func (s *SegmentEffortsService) ListSegmentEfforts(accessToken string, opt *ListSegmentEffortOptions) ([]SegmentEffortDetailed, error) {
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
