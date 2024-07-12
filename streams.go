package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type StreamsService service

const streams = "streams"

type StreamSet struct {
	AltitudeStream       *AltitudeStream       `json:"altitude,omitempty"`        // An instance of AltitudeStream.
	CadenceStream        *CadenceStream        `json:"cadence,omitempty"`         // An instance of CadenceStream.
	DistanceStream       *DistanceStream       `json:"distance,omitempty"`        // An instance of DistanceStream.
	HeartRateStream      *HeartrateStream      `json:"heartrate,omitempty"`       // An instance of HeartrateStream.
	LatLngStream         *LatLngStream         `json:"latlng,omitempty"`          // An instance of LatLngStream.
	MovingStream         *MovingStream         `json:"moving,omitempty"`          // An instance of MovingStream.
	SmoothGradeStream    *SmoothGradeStream    `json:"grade_smooth,omitempty"`    // An instance of SmoothGradeStream.
	SmoothVelocityStream *SmoothVelocityStream `json:"velocity_smooth,omitempty"` // An instance of SmoothVelocityStream.
	TempStream           *TemperatureStream    `json:"temp,omitempty"`            // An instance of TemperatureStream.
	TimeStream           *TimeStream           `json:"time,omitempty"`            // An instance of TimeStream.
	WattsStream          *PowerStream          `json:"watts,omitempty"`           // An instance of PowerStream.
}

func (ss *StreamSet) String() string {
	return Stringify(ss)
}

type baseStream struct {
	OriginalSize *int    `json:"original_size,omitempty"` // The number of data points in this stream
	Resolution   *string `json:"resolution,omitempty"`    // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   *string `json:"series_type,omitempty"`   // The base series used in the case the stream was downsampled May take one of the following values: distance, time
}
type AltitudeStream struct {
	Data []float32 `json:"data"` // The sequence of altitude values for this stream, in meters
	baseStream
}

type CadenceStream struct {
	baseStream
	Data []int `json:"data"` //  The sequence of cadence values for this stream, in rotations per minute
}

type DistanceStream struct {
	Data []float32 `json:"data"` // The sequence of distance values for this stream, in meters
	baseStream
}

type HeartrateStream struct {
	Data []int `json:"data"` // The sequence of heart rate values for this stream, in beats per minute
	baseStream
}

type LatLngStream struct {
	Data []LatLng `json:"data"` // The sequence of lat/long values for this stream
	baseStream
}

type MovingStream struct {
	Data []bool // The sequence of moving values for this stream, as boolean values
	baseStream
}

type PowerStream struct {
	Data []int `json:"data"` // The sequence of power values for this stream, in watts
	baseStream
}

type SmoothGradeStream struct {
	Data []float32 `json:"data"` // The sequence of grade values for this stream, as percents of a grade
	baseStream
}

type SmoothVelocityStream struct {
	Data []float32 `json:"data"` // The sequence of velocity values for this stream, in meters per second
	baseStream
}

type TemperatureStream struct {
	Data []int `json:"data"` // The sequence of temperature values for this stream, in celsius degrees
	baseStream
}

type TimeStream struct {
	Data []int `json:"data"` // The sequence of time values for this stream, in seconds
	baseStream
}

// Returns the given activity's streams. Requires activity:read scope. Requires activity:read_all scope for Only Me activities.
func (s *StreamsService) GetActivityStreams(accessToken string, activityID int, keys []string) (*StreamSet, error) {
	params := url.Values{}
	params.Add("keys", strings.Join(keys, ","))
	params.Add("key_by_type", "true")

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/%s", activities, activityID, streams),
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := new(StreamSet)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the given route's streams. Requires read_all scope for private routes.
func (s *StreamsService) GetRouteStreams(accessToken string, routeID int) (*StreamSet, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/%s", routes, routeID, streams),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(StreamSet)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a set of streams for a segment effort completed by the authenticated athlete. Requires read_all scope.
func (s *StreamsService) GetSegmentEffortStreams(accessToken string, segmentEffortID int, keys []string) (*StreamSet, error) {
	params := url.Values{}
	params.Add("keys", strings.Join(keys, ","))
	params.Add("key_by_type", "true")

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/%s", segment_efforts, segmentEffortID, streams),
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := new(StreamSet)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a set of streams for a segment completed by the authenticated athlete. Requires read_all scope.
func (s *StreamsService) GetSegmentStreams(accessToken string, segmentID int, keys []string) (*StreamSet, error) {
	params := url.Values{}
	params.Add("keys", strings.Join(keys, ","))
	params.Add("key_by_type", "true")

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/%s", segments, segmentID, streams),
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := new(StreamSet)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
