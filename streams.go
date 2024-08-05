package gostrava

type StreamType string

var StreamTypes = struct {
	All            StreamType
	Time           StreamType
	Distance       StreamType
	LatLng         StreamType
	Altitude       StreamType
	VelocitySmooth StreamType
	HeartRate      StreamType
	Cadence        StreamType
	Watts          StreamType
	Temperature    StreamType
	Moving         StreamType
	GradeSmooth    StreamType
}{
	"time,distance,latlng,altitude,velocity_smooth,heartrate,cadence,watts,temperature,moving,grade_smooth",
	"time", "distance", "latlng", "altitude", "velocity_smooth",
	"heartrate", "cadence", "watts", "temperature", "moving", "grade_smooth",
}

type Stream struct {
	Type string        `json:"type"`
	Data []interface{} `json:"data"`

	// Not available for GetRouteStreams()
	OriginalSize int    `json:"original_size"` // The number of data points in this stream
	Resolution   string `json:"resolution"`    // The level of detail (sampling) in which this stream was returned May take one of the following values: low, medium, high
	SeriesType   string `json:"series_type"`   // The base series used in the case the stream was downsampled May take one of the following values: distance, time
}

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

type AltitudeStream struct {
	Data []float32 `json:"data"` // The sequence of altitude values for this stream, in meters
	Stream
}

type CadenceStream struct {
	Data []int `json:"data"` //  The sequence of cadence values for this stream, in rotations per minute
	Stream
}

type DistanceStream struct {
	Data []float32 `json:"data"` // The sequence of distance values for this stream, in meters
	Stream
}

type HeartrateStream struct {
	Data []int `json:"data"` // The sequence of heart rate values for this stream, in beats per minute
	Stream
}

type LatLngStream struct {
	Data []LatLng `json:"data"` // The sequence of lat/long values for this stream
	Stream
}

type MovingStream struct {
	Data []bool // The sequence of moving values for this stream, as boolean values
	Stream
}

type PowerStream struct {
	Data []int `json:"data"` // The sequence of power values for this stream, in watts
	Stream
}

type SmoothGradeStream struct {
	Data []float32 `json:"data"` // The sequence of grade values for this stream, as percents of a grade
	Stream
}

type SmoothVelocityStream struct {
	Data []float32 `json:"data"` // The sequence of velocity values for this stream, in meters per second
	Stream
}

type TemperatureStream struct {
	Data []int `json:"data"` // The sequence of temperature values for this stream, in celsius degrees
	Stream
}

type TimeStream struct {
	Data []int `json:"data"` // The sequence of time values for this stream, in seconds
	Stream
}
