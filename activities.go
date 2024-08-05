package gostrava

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"github.com/google/go-querystring/query"
)

// *************** Types ********************

type ActivityMeta struct {
	ID            int     `json:"id"`                   // The unique identifier of the activity
	ResourceState int8    `json:"resource_state"`       //
	Visibility    *string `json:"visibility,omitempty"` //
}

type ActivitySummary struct {
	ActivityMeta
	Athlete                    AthleteMeta      `json:"athlete"` // An instance of AthleteMeta.
	Name                       string           `json:"name"`
	Distance                   float32          `json:"distance"` // The activity's distance, in meters
	MovingTime                 int              `json:"moving_time"`
	ElapsedTime                int              `json:"elapsed_time"`         // The activity's elapsed time, in seconds
	TotalElevationGain         float32          `json:"total_elevation_gain"` // The activity's total elevation gain.
	Type                       ActivityType     `json:"type"`                 // Deprecated. Prefer to use sport_type
	SportType                  SportType        `json:"sport_type"`
	WorkoutType                int              `json:"workout_type,omitempty"` //  The activity's workout type
	StartDate                  TimeStamp        `json:"start_date"`             // The time at which the activity was started.
	StartDateLocal             TimeStamp        `json:"start_date_local"`       // The time at which the activity was started in the local timezone.
	Timezone                   string           `json:"timezone"`               // The timezone of the activity
	UTCOffset                  float32          `json:"utc_offset"`
	LocationCity               string           `json:"location_city"`    // City where the activity took place
	LocationState              string           `json:"location_state"`   // State where the activity took place
	LocationCountry            string           `json:"location_country"` // Country where the activity took place
	AchievementCount           int              `json:"achievement_count"`
	KudosCount                 int              `json:"kudos_count"`          // The number of kudos given for this activity
	CommentCount               int              `json:"comment_count"`        // The number of comments for this activity
	AthleteCount               int              `json:"athlete_count"`        // The number of athletes for taking part in a group activity
	PhotoCount                 int              `json:"photo_count"`          // The number of Instagram photos for this activity
	Map                        PolylineSummmary `json:"map"`                  // An instance of PolylineSummary.
	Trainer                    bool             `json:"trainer"`              // Whether this activity was recorded on a training machine
	Commute                    bool             `json:"commute,omitempty"`    // Whether this activity is a commute
	Private                    bool             `json:"private"`              // Whether this activity is private
	Manual                     bool             `json:"manual"`               // Indicates whether this activity was manually created by the user
	Flagged                    bool             `json:"flagged"`              // Whether this activity is flagged
	GearID                     *string          `json:"gear_id,omitempty"`    // The id of the gear for the activity
	StartLatLng                LatLng           `json:"start_latlng"`         // An instance of LatLng.
	EndLatLng                  LatLng           `json:"end_latlng,omitempty"` // An instance of LatLng.
	AvgSpeed                   float32          `json:"average_speed"`        // The activity's average speed, in meters per second
	MaxSpeed                   float32          `json:"max_speed"`            // The activity's max speed, in meters per second
	AvgWatts                   float32          `json:"average_watts"`        // Average power output in watts during this activity. Rides only
	MaxWatts                   float32          `json:"max_watts"`            // Rides with power meter data only
	Kilojoules                 float32          `json:"kilojoules"`           // The total work done in kilojoules during this activity. Rides only
	DeviceWatts                bool             `json:"device_watts"`
	HasHeartRate               bool             `json:"has_heartrate"`     // Indicates weather the activity has a heartrate recorder
	AvgHeartRate               float32          `json:"average_heartrate"` // The activity's average heart rate, in beats per minute
	MaxHeartRate               float32          `json:"max_heartrate"`     // The activity's max heartrate in beats per minute
	HeartRateOptOut            bool             `json:"heartrate_opt_out"`
	DisplayHideHeartRateOption *bool            `json:"display_heartrate_option,omitempty"`
	ElevationHigh              float32          `json:"elev_high"`     // The activity's highest elevation, in meters
	ElevationLow               float32          `json:"elev_low"`      // The activity's lowest elevation, in meters
	UploadID                   int              `json:"upload_id"`     // The identifier of the upload that resulted in this activity
	UploadIdStr                string           `json:"upload_id_str"` // The unique identifier of the upload in string format
	ExternalID                 string           `json:"external_id"`   // The identifier provided at upload time
	FromAcceptedTag            bool             `json:"from_accepted_tag"`
	PRCount                    int              `json:"pr_count"`
	TotalPhotoCount            int              `json:"total_photo_count"` // The number of Instagram and Strava photos for this activity
	HasKudoed                  bool             `json:"has_kudoed"`        // Whether the logged-in athlete has kudoed this activity
	SufferScore                *float32         `json:"suffer_score,omitempty"`
}

type ActivityDetailed struct {
	ActivitySummary
	Map            PolylineDetailed         `json:"map"`
	Description    string                   `json:"description,omitempty"` // The description of the activity
	Calories       float32                  `json:"calories"`              // The number of kilocalories consumed during this activity
	SegmentEfforts []*SegmentEffortDetailed `json:"segment_efforts"`       // A collection of SegmentEffortDetailed objects.
	SplitsMetric   []*Split                 `json:"splits_metric"`         // The splits of this activity in metric units (for runs)
	SplitsStandard []*Split                 `json:"splits_standard"`       // The splits of this activity in imperial units (for runs)
	Laps           []*Lap                   `json:"laps"`                  // A collection of Lap objects.
	Photos         PhotosSummary            `json:"photos"`                // An instance of PhotosSummary.
	HideFromHome   bool                     `json:"hide_from_home"`        // Whether the activity is muted
	DeviceName     string                   `json:"device_name,"`          // The name of the device used to record the activity
	EmbedToken     string                   `json:"embed_token"`           // The token used to embed a Strava activity
	Gear           *GearSummary             `json:"gear,omitempty"`        // An instance of SummaryGear.
	BestEfforts    []*SegmentEffortDetailed `json:"best_efforts"`          // A collection of SegmentEffortDetailed objects.
	AvailableZones []string                 `json:"available_zones"`       // Activity available zones
}

type Split struct {
	Distance              float32 `json:"distance"`             //  The distance of this split, in meters
	ElapsedTime           int     `json:"elapsed_time"`         // The elapsed time of this split, in seconds
	ElevationDifference   float32 `json:"elevation_difference"` // The elevation difference of this split, in meters
	MovingTime            int     `json:"moving_time"`          // The moving time of this split, in seconds
	Split                 int     `json:"split"`                // Split number
	AvgSpeed              float32 `json:"average_speed"`        // The average speed of this split, in meters per second
	AvgGradeAdjustedSpeed float32 `json:"average_grade_adjusted_speed"`
	AvgHeartRate          float32 `json:"average_heartrate"` // The average heartrate of this split, in beats per minute
	PaceZone              int     `json:"pace_zone"`         // The pacing zone of this split
}

type Lap struct {
	ID                 int          `json:"id"`                   // The unique identifier of this lap
	Activity           ActivityMeta `json:"activity"`             // An instance of ActivityMeta.
	Athlete            AthleteMeta  `json:"athlete"`              // An instance of AthleteMeta.
	AvgCadence         float32      `json:"average_cadence"`      // The lap's average cadence
	AvgHeartRate       float32      `json:"average_heartrate"`    // The lap's average heartrate
	AvgSpeed           float32      `json:"average_speed"`        // The lap's average speed
	DeviceWatts        bool         `json:"device_watts"`         // Whether the watts are from a power meter, false if estimated
	Distance           float32      `json:"distance"`             // The lap's distance, in meters
	ElapsedTime        int          `json:"elapsed_time"`         // The lap's elapsed time, in seconds
	EndIndex           int          `json:"end_index"`            // The end index of this effort in its activity's stream
	LapIndex           int          `json:"lap_index"`            // The index of this lap in the activity it belongs to
	MaxHeartRate       float32      `json:"max_heartrate"`        // The maximum heartrate of this lap, in beats per minute
	MaxSpeed           float32      `json:"max_speed"`            // The maximum speed of this lat, in meters per second
	MovingTime         int          `json:"moving_time"`          // The lap's moving time, in seconds
	Name               string       `json:"name"`                 // The name of the lap
	PaceZone           int          `json:"pace_zone"`            // The athlete's pace zone during this lap
	ResourceState      int8         `json:"resource_state"`       // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
	Split              int          `json:"split"`                // An instance of integer.
	StartIndex         int          `json:"start_index"`          // The start index of this effort in its activity's stream
	StartDate          TimeStamp    `json:"start_date"`           // The time at which the lap was started.
	StartDateLocal     TimeStamp    `json:"start_date_local"`     // The time at which the lap was started in the local timezone.
	TotalElevationGain float32      `json:"total_elevation_gain"` // The elevation gain of this lap, in meters
}

type PhotosSummary struct {
	Count   int                   `json:"count"`             // The number of photos
	Primary *PhotosSummaryPrimary `json:"primary,omitempty"` // An instance of PhotosSummaryPrimary.
}

type PhotosSummaryPrimary struct {
	ID        string         `json:"unique_id"`
	Source    int            `json:"source"`
	MediaType int            `json:"media_type"`
	Urls      map[int]string `json:"urls"`
}

type ActivityZone struct {
	Score               *float32         `json:"score,omitempty"`
	DistributionBuckets []TimedZoneRange `json:"distribution_buckets,omitempty"`
	Type                *string          `json:"type,omitempty"` // May take one of the following values: heartrate, power
	SensorBased         *bool            `json:"sensor_based,omitempty"`
	Points              *int             `json:"points,omitempty"`
	CustomZones         *bool            `json:"custom_zones,omitempty"`
	Max                 *int             `json:"max,omitempty"`
}

type Zones struct {
	HearRate HeartRateZoneRanges `json:"heart_rate"` // An instance of HeartRateZoneRanges.
	Power    PowerZoneRanges     `json:"power"`      // An instance of PowerZoneRanges.
}

type HeartRateZoneRanges struct {
	CustomZones bool        `json:"custom_zone"` // Whether the athlete has set their own custom heart rate zones
	Zones       []ZoneRange `json:"zones"`       // An instance of ZoneRanges.
}

type PowerZoneRanges struct {
	Zones []ZoneRange `json:"zones"` // An instance of ZoneRanges.
}

type ZoneRange struct {
	Max int `json:"max"` // The maximum value in the range.
	Min int `json:"min"` // The minimum value in the range.
}

// A union type representing the time spent in a given zone.
type TimedZoneRange struct {
	Min  int `json:"min"`  // The minimum value in the range.
	Max  int `json:"max"`  // The maximum value in the range.
	Time int `json:"time"` // The number of seconds spent in this zone
}

type Comment struct {
	ID         int            `json:"id"`          // The unique identifier of this comment
	ActivityID int            `json:"activity_id"` // The identifier of the activity this comment is related to
	Text       string         `json:"text"`        // The content of the comment
	Athlete    AthleteSummary `json:"athlete"`     // An instance of AthleteSummary.
	CreatedAt  TimeStamp      `json:"created_at="` // The time at which this comment was created.
}

// *************** Methods ********************

type ActivityService service

type NewActivity struct {
	Name           string       `url:"name"`             // The name of the activity.
	Type           ActivityType `url:"activity_type"`    // Type of activity. For example - Run, Ride etc.
	SportType      SportType    `url:"sport_type"`       // Sport type of activity. For example - Run, MountainBikeRide, Ride, etc.
	StartDateLocal time.Time    `url:"start_date_local"` // ISO 8601 formatted date time.
	ElapsedTime    int          `url:"elapsed_time"`     // In seconds.
	Description    string       `url:"description"`      // Description of the activity.
	Distance       float64      `url:"distance"`         // In meters.
	Trainer        bool         `url:"trainer"`          // Mark as a trainer activity, 0 otherwise.
	Commute        bool         `url:"commute"`          // Mark as commuter, 0 otherwise.
}

// New creates a new manual activity for an athlete.
// Requires activity:write scope.
// POST: https://www.strava.com/api/v3/activities
func (s *ActivityService) New(ctx context.Context, accessToken string, activity NewActivity) (*ActivityDetailed, *http.Response, error) {

	values, err := query.Values(activity)

	if activity.Trainer {
		formData.Set("trainer", "1")
	}
	if activity.Commute {
		formData.Set("commute", "1")
	}

	// Create the request
	req, err := s.client.NewRequest(http.MethodPost, "/activities", values.Encode(), SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	// Execute the request and parse the response
	newActivity := new(ActivityDetailed)
	resp, err := s.client.DoAndParse(ctx, req, newActivity)
	if err != nil {
		return nil, resp, err
	}

	return newActivity, resp, nil
}

// GetByID retrieves an activity by its ID.
// GET: https://www.strava.com/api/v3/activities/{id}
// Parameters:
// - ctx: Context for the request.
// - accessToken: OAuth access token.
// - ID: ID of the activity to retrieve.
// - includeEfforts: Whether to include efforts in the response.
func (s *ActivityService) GetByID(ctx context.Context, accessToken string, id int, includeEfforts bool) (*ActivityDetailed, *http.Response, error) {
	q := url.Values{}
	q.Add("include_all_efforts", strconv.FormatBool(includeEfforts))

	// Construct the request URL
	urlStr := fmt.Sprintf("/activities/%d", id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, q, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	activity := new(ActivityDetailed)
	resp, err := s.client.DoAndParse(ctx, req, activity)
	if err != nil {
		return nil, resp, err
	}

	return activity, resp, nil
}

// ListActivityLaps retrieves the laps of a specific activity.
// GET: https://www.strava.com/api/v3/activities/{id}/laps
// Parameters:
// - ctx: Context for the request.
// - accessToken: OAuth access token.
// - id: ID of the activity whose laps are to be retrieved.
func (s *ActivityService) ListActivityLaps(ctx context.Context, accessToken string, id int) ([]Lap, *http.Response, error) {
	// Construct the URL to retrieve laps for the specified activity
	urlStr := fmt.Sprintf("/activities/%d/laps", id)

	// Create the request
	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	// Prepare a container for laps
	var laps []Lap

	// Send the request and parse the response
	resp, err := s.client.DoAndParse(ctx, req, &laps)
	if err != nil {
		return nil, resp, err
	}

	return laps, resp, nil
}

// func (s *ActivityService) GetActivityZones(ctx context.Context, accessToken string, id int)
