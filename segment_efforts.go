package gostrava

import (
	"context"
	"fmt"
	"net/http"
)

// ***************Types ********************

type SegmentEffortSummary struct {
	ID             int       `json:"id"`               // The unique identifier of this effort
	ActivityID     int       `json:"activity_id"`      // The unique identifier of the activity related to this effort
	ElapsedTime    int       `json:"elapsed_time"`     // The effort's elapsed time
	StartDate      TimeStamp `json:"start_date"`       // The time at which the effort was started.
	StartDateLocal TimeStamp `json:"start_date_local"` // The time at which the effort was started in the local timezone.
	Distance       float32   `json:"distance"`         //  The effort's distance in meters
	IsKom          bool      `json:"is_kom"`           // Whether this effort is the current best on the leaderboard
}

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

// *************** Methods ********************

type SegmentEffortService service

const segmentEfforts string = "/api/v3/segment_efforts"

// Returns a segment effort from an activity that is owned by the authenticated athlete. 
//
// GET: https://www.strava.com/api/v3/segment_efforts/{id}
func (s *SegmentEffortService) GetById(ctx context.Context, accessToken string, id int) (*SegmentEffortDetailed, *http.Response, error) {
	
	urlStr := fmt.Sprintf("%s/%d", segmentEfforts, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil{
		return nil, nil, err
	}

	segmentEffort := new(SegmentEffortDetailed)
	resp, err := s.client.DoAndParse(ctx, req, segmentEffort)
	if err != nil {
		return nil, resp, err
	}

	return segmentEffort, resp, err
}

