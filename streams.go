package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type StreamsAPIService apiService

type StreamSet struct {
	AltitudeStream       AltitudeStream       `json:"altitude"`        // An instance of AltitudeStream.
	CadenceStream        CadenceStream        `json:"cadence"`         // An instance of CadenceStream.
	DistanceStream       DistanceStream       `json:"distance"`        // An instance of DistanceStream.
	HeartRateStream      HeartrateStream      `json:"heartrate"`       // An instance of HeartrateStream.
	LatLngStream         LatLngStream         `json:"latlng"`          // An instance of LatLngStream.
	MovingStream         MovingStream         `json:"moving"`          // An instance of MovingStream.
	SmoothGradeStream    SmoothGradeStream    `json:"grade_smooth"`    // An instance of SmoothGradeStream.
	SmoothVelocityStream SmoothVelocityStream `json:"velocity_smooth"` // An instance of SmoothVelocityStream.
	TempStream           TemperatureStream    `json:"temp"`            // An instance of TemperatureStream.
	TimeStream           TimeStream           `json:"time"`            // An instance of TimeStream.
	WattsStream          PowerStream          `json:"watts"`           // An instance of PowerStream.
}

type BaseStream struct {
	OriginalSize int                  `json:"original_size"` // The number of data points in this stream
	Resolution   BaseStreamResolution `json:"resolution"`    // The level of detail (sampling) in which this stream was returned May take one of the following values: BaseStreamResolutions.Low, BaseStreamResolution.Medium, BaseStreamResolution.High
	SeriesType   BaseStreamSeriesType `json:"series_type"`   // The base series used in the case the stream was downsampled May take one of the following values: BaseStreamSeriesTypes.Distance, BaseStreamSeriesTypes.Time
}

type BaseStreamResolution string

type BaseStreamSeriesType string

var BaseResolutions = struct {
	Low    BaseStreamResolution
	Medium BaseStreamResolution
	High   BaseStreamResolution
}{
	"low", "medium", "high",
}

var BaseStreamSeriesTypes = struct {
	Distance BaseStreamSeriesType
	Time     BaseStreamSeriesType
}{
	"distance", "time",
}

type AltitudeStream struct {
	Data []float32 `json:"data"` // The sequence of altitude values for this stream, in meters
	BaseStream
}

type DistanceStream struct {
	Data []float32 `json:"data"` // The sequence of distance values for this stream, in meters
	BaseStream
}

type HeartrateStream struct {
	Data []int `json:"data"` // The sequence of heart rate values for this stream, in beats per minute
	BaseStream
}

type LatLngStream struct {
	Data []LatLng `json:"data"` // The sequence of lat/long values for this stream
	BaseStream
}

type MovingStream struct {
	Data []bool // The sequence of moving values for this stream, as boolean values
	BaseStream
}

type PowerStream struct {
	Data []int `json:"data"` // The sequence of power values for this stream, in watts
	BaseStream
}

type SmoothGradeStream struct {
	Data []float32 `json:"data"` // The sequence of grade values for this stream, as percents of a grade
	BaseStream
}

type SmoothVelocityStream struct {
	Data []float32 `json:"data"` // The sequence of velocity values for this stream, in meters per second
	BaseStream
}

type TemperatureStream struct {
	Data []int `json:"data"` // The sequence of temperature values for this stream, in celsius degrees
	BaseStream
}

type TimeStream struct {
	Data []int `json:"data"` // The sequence of time values for this stream, in seconds
	BaseStream
}

// Returns the given activity's streams. Requires activity:read scope. Requires activity:read_all scope for Only Me activities.
func (s *StreamsAPIService) GetActivityStreams(access_token string, activityID int, keys []string) (*StreamSet, error) {
	requestUrl := s.client.BaseURL.JoinPath(activitiesPath, fmt.Sprint(activityID), streamPath)

	params := url.Values{}
	params.Add("keys", strings.Join(keys, ","))
	params.Add("key_by_type", "true")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
		body:         params,
	})
	if err != nil {
		return nil, err
	}

	s.client.TestingFileName = "streams_get_activity_streams_server_response.json"

	resp := &StreamSet{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the given route's streams. Requires read_all scope for private routes.
func (s *StreamsAPIService) GetRouteStreams(access_token string, routeID int) (*StreamSet, error) {
	requestUrl := s.client.BaseURL.JoinPath(routesPath, fmt.Sprint(routeID), streamPath)

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	s.client.TestingFileName = "streams_get_route_streams_server_response.json"

	resp := &StreamSet{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a set of streams for a segment effort completed by the authenticated athlete. Requires read_all scope.
func (s *StreamsAPIService) GetSegmentEffortStreams(access_token string, segmentEffortID int, keys []string) (*StreamSet, error) {
	requestUrl := s.client.BaseURL.JoinPath(segmentEffortsPath, fmt.Sprint(segmentEffortID), streamPath)

	params := url.Values{}
	params.Add("keys", strings.Join(keys, ","))
	params.Add("key_by_type", "true")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
		body:         params,
	})
	if err != nil {
		return nil, err
	}

	s.client.TestingFileName = "streams_get_segment_effort_streams_server_response.json"

	resp := &StreamSet{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a set of streams for a segment effort completed by the authenticated athlete. Requires read_all scope.
func (s *StreamsAPIService) GetSegmentStreams(access_token string, segmentEffortID int, keys []string) (*StreamSet, error) {
	requestUrl := s.client.BaseURL.JoinPath(segmentsPath, fmt.Sprint(segmentEffortID), streamPath)

	params := url.Values{}
	params.Add("keys", strings.Join(keys, ","))
	params.Add("key_by_type", "true")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
		body:         params,
	})
	if err != nil {
		return nil, err
	}

	s.client.TestingFileName = "streams_get_segment_streams_server_response.json"

	resp := &StreamSet{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
