package gostrava

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
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
	if err != nil {
		return nil, nil, err
	}

	segmentEffort := new(SegmentEffortDetailed)
	resp, err := s.client.DoAndParse(ctx, req, segmentEffort)
	if err != nil {
		return nil, resp, err
	}

	return segmentEffort, resp, err
}

// GetSegmentEffortStreams returns a set of streams for a segment effort completed by the authenticated athlete.
// Requires read_all scope. // By default returns the primary stream of the of the segment effort.
//
// GET https://www.strava.com/api/v3/segment_efforts/{id}/streams
func (s *SegmentEffortService) GetSegmentEffortStreams(ctx context.Context, accessToken string, id int, streamTypes []StreamType) ([]Stream, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/stream", segmentEfforts, id)

	v := url.Values{}

	typesSlice := make([]string, len(streamTypes))
	for i, v := range streamTypes {
		typesSlice[i] = string(v)
	}
	v.Add("keys", strings.Join(typesSlice, ","))
	v.Add("keys_by_type", "true")

	req, err := s.client.NewRequest(http.MethodGet, urlStr, v, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var streams []Stream
	resp, err := s.client.DoAndParse(ctx, req, &streams)
	if err != nil {
		return nil, resp, err
	}
	return streams, resp, nil
}

type ListSegmentEffortOptions struct {
	Page           int       `url:"page,omitempty"`
	PerPage        int       `url:"per_page,omitempty"`
	StartDateLocal time.Time `url:"start_date_local,omitempty"`
	EndDateLocal   time.Time `url:"end_date_local,omitempty"`
}

// ListSegmentEfforts returns a set containing the corresponding athlete's segment efforts for a given segment.
//
//	GET: https://www.strava.com/api/v3/segment_efforts
func (s *SegmentEffortService) ListSegmentEfforts(ctx context.Context, accessToken string, segmentId int, options *ListSegmentEffortOptions) ([]SegmentEffortDetailed, *http.Response, error) {
	urlStr := segmentEfforts

	q, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}
	q.Add("segment_id", strconv.Itoa(segmentId))

	req, err := s.client.NewRequest(http.MethodGet, urlStr, q, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var segmentEfforts []SegmentEffortDetailed
	resp, err := s.client.DoAndParse(ctx, req, &segmentEfforts)
	if err != nil {
		return nil, resp, err
	}

	return segmentEfforts, resp, nil
}
