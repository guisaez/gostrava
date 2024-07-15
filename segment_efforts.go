package gostrava

import (
	"net/url"
	"strconv"
	"time"
)
type SegmentEffortDetailed struct {
	Name         *string         `json:"name,omitempty"`              // The name of the segment on which this effort was performed
	Activity     *ActivityMeta   `json:"activity,omitempty"`          // An instance of MetaActivity.
	Athlete      *AthleteMeta    `json:"athlete,omitempty"`           // An instance of MetaAthlete.
	MovingTime   *int            `json:"moving_time,omitempty"`       // The effort's moving time
	StartIndex   *int            `json:"start_index,omitempty"`       // The start index of this effort in its activity's stream
	EndIndex     *int            `json:"end_index,omitempty"`         // The end index of this effort in its activity's stream
	AvgCadence   *float32        `json:"average_cadence,omitempty"`   // The effort's average cadence
	AverageWatts *float32        `json:"average_watts,omitempty"`     // The average wattage of this effort
	DeviceWatts  *bool           `json:"device_watts,omitempty"`      // For riding efforts, whether the wattage was reported by a dedicated recording device
	AvgHeartRate *bool           `json:"average_heartrate,omitempty"` // The heart heart rate of the athlete during this effort
	MaxHeartRate *float32        `json:"max_heartrate,omitempty"`     // The maximum heart rate of the athlete during this effort
	Segment      *SegmentSummary `json:"segment,omitempty"`           // An instance of SegmentSummary.
	KomRank      *int            `json:"kom_rank,omitempty"`          // The rank of the effort on the global leaderboard if it belongs in the top 10 at the time of upload
	PRRank       *int            `json:"pr_rank,omitempty"`           // The rank of the effort on the athlete's leaderboard if it belongs in the top 3 at the time of upload
	Hidden       *bool           `json:"hidden,omitempty"`            // Whether this effort should be hidden when viewed within an activity
}

type SegmentEffortSummary struct {
	ID             int       `json:"id"`               // The unique identifier of this effort
	ActivityID     int       `json:"activity_id"`      // The unique identifier of the activity related to this effort
	ElapsedTime    int       `json:"elapsed_time"`     // The effort's elapsed time
	StartDate      TimeStamp `json:"start_date"`       // The time at which the effort was started.
	StartDateLocal TimeStamp `json:"start_date_local"` // The time at which the effort was started in the local timezone.
	Distance       float32   `json:"distance"`         //  The effort's distance in meters
	IsKom          bool      `json:"is_kom"`           // Whether this effort is the current best on the leaderboard
}

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
