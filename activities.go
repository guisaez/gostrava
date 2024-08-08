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

// *************** Types ********************

type ActivityMeta struct {
	ID            *int    `json:"id,omitempty"`             // The unique identifier of the activity
	ResourceState *int8   `json:"resource_state,omitempty"` //
	Visibility    *string `json:"visibility,omitempty"`     //
}

type ActivitySummary struct {
	ActivityMeta
	Athlete                    *AthleteMeta      `json:"athlete"`              // An instance of AthleteMeta.
	Name                       string            `json:"name"`                 // Activity name
	Distance                   float32           `json:"distance"`             // The activity's distance, in meters
	MovingTime                 int               `json:"moving_time"`          // Total moving time
	ElapsedTime                int               `json:"elapsed_time"`         // The activity's elapsed time, in seconds
	TotalElevationGain         float32           `json:"total_elevation_gain"` // The activity's total elevation gain.
	Type                       ActivityType      `json:"type"`                 // Deprecated. Prefer to use sport_type
	SportType                  SportType         `json:"sport_type"`
	WorkoutType                int               `json:"workout_type,omitempty"` //  The activity's workout type
	StartDate                  TimeStamp         `json:"start_date"`             // The time at which the activity was started.
	StartDateLocal             TimeStamp         `json:"start_date_local"`       // The time at which the activity was started in the local timezone.
	Timezone                   string            `json:"timezone"`               // The timezone of the activity
	UTCOffset                  float32           `json:"utc_offset"`
	LocationCity               string            `json:"location_city"`     // City where the activity took place
	LocationState              string            `json:"location_state"`    // State where the activity took place
	LocationCountry            string            `json:"location_country"`  // Country where the activity took place
	AchievementCount           int               `json:"achievement_count"` // Activity Achievement count
	KudosCount                 int               `json:"kudos_count"`       // The number of kudos given for this activity
	CommentCount               int               `json:"comment_count"`     // The number of comments for this activity
	AthleteCount               int               `json:"athlete_count"`     // The number of athletes for taking part in a group activity
	PhotoCount                 int               `json:"photo_count"`       // The number of Instagram photos for this activity
	Map                        *PolylineSummary `json:"map"`               // An instance of PolylineSummary.
	Trainer                    bool              `json:"trainer"`           // Whether this activity was recorded on a training machine
	Commute                    bool              `json:"commute"`           // Whether this activity is a commute
	Private                    bool              `json:"private"`           // Whether this activity is private
	Manual                     bool              `json:"manual"`            // Indicates whether this activity was manually created by the user
	Flagged                    bool              `json:"flagged"`           // Whether this activity is flagged
	GearID                     *string           `json:"gear_id"`           // The id of the gear for the activity
	StartLatLng                LatLng            `json:"start_latlng"`      // An instance of LatLng.
	EndLatLng                  LatLng            `json:"end_latlng"`        // An instance of LatLng.
	AvgSpeed                   float32           `json:"average_speed"`     // The activity's average speed, in meters per second
	MaxSpeed                   float32           `json:"max_speed"`         // The activity's max speed, in meters per second
	AvgWatts                   float32           `json:"average_watts"`     // Average power output in watts during this activity. Rides only
	MaxWatts                   float32           `json:"max_watts"`         // Rides with power meter data only
	Kilojoules                 float32           `json:"kilojoules"`        // The total work done in kilojoules during this activity. Rides only
	DeviceWatts                bool              `json:"device_watts"`
	HasHeartRate               bool              `json:"has_heartrate"`     // Indicates weather the activity has a heartrate recorder
	AvgHeartRate               float32           `json:"average_heartrate"` // The activity's average heart rate, in beats per minute
	MaxHeartRate               float32           `json:"max_heartrate"`     // The activity's max heartrate in beats per minute
	HeartRateOptOut            bool              `json:"heartrate_opt_out"`
	DisplayHideHeartRateOption bool              `json:"display_hide_heartrate_option,omitempty"`
	ElevationHigh              float32           `json:"elev_high"`     // The activity's highest elevation, in meters
	ElevationLow               float32           `json:"elev_low"`      // The activity's lowest elevation, in meters
	UploadID                   int               `json:"upload_id"`     // The identifier of the upload that resulted in this activity
	UploadIdStr                string            `json:"upload_id_str"` // The unique identifier of the upload in string format
	ExternalID                 string            `json:"external_id"`   // The identifier provided at upload time
	FromAcceptedTag            bool              `json:"from_accepted_tag"`
	PRCount                    int               `json:"pr_count"`
	TotalPhotoCount            int               `json:"total_photo_count"` // The number of Instagram and Strava photos for this activity
	HasKudoed                  bool              `json:"has_kudoed"`        // Whether the logged-in athlete has kudoed this activity
	SufferScore                *float32          `json:"suffer_score,omitempty"`
}

type ActivityDetailed struct {
	ActivitySummary
	Map             *PolylineDetailed        `json:"map"`
	Description     string                   `json:"description,omitempty"` // The description of the activity
	Calories        float32                  `json:"calories"`              // The number of kilocalories consumed during this activity
	SegmentEfforts  []SegmentEffortDetailed  `json:"segment_efforts"`       // A collection of SegmentEffortDetailed objects.
	SplitsMetric    []Split                  `json:"splits_metric"`         // The splits of this activity in metric units (for runs)
	SplitsStandard  []Split                  `json:"splits_standard"`       // The splits of this activity in imperial units (for runs)
	Laps            []Lap                    `json:"laps"`                  // A collection of Lap objects.
	Photos          PhotosSummary            `json:"photos"`                // An instance of PhotosSummary.
	HideFromHome    bool                     `json:"hide_from_home"`        // Whether the activity is muted
	DeviceName      string                   `json:"device_name,"`          // The name of the device used to record the activity
	EmbedToken      string                   `json:"embed_token"`           // The token used to embed a Strava activity
	Gear            *GearSummary             `json:"gear,omitempty"`        // An instance of SummaryGear.
	BestEfforts     []*SegmentEffortDetailed `json:"best_efforts"`          // A collection of SegmentEffortDetailed objects.
	AvailableZones  []string                 `json:"available_zones"`       // Activity available zones
	StatsVisibility []StatVisibility         `json:"stats_visibility"`
}

type StatVisibility struct {
	Type       string `json:"type"`
	Visibility string `json:"visibility"`
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
	ID                 int           `json:"id"`                   // The unique identifier of this lap
	Activity           *ActivityMeta `json:"activity"`             // An instance of ActivityMeta.
	Athlete            *AthleteMeta  `json:"athlete"`              // An instance of AthleteMeta.
	AvgCadence         float32       `json:"average_cadence"`      // The lap's average cadence
	AvgHeartRate       float32       `json:"average_heartrate"`    // The lap's average heartrate
	AvgSpeed           float32       `json:"average_speed"`        // The lap's average speed
	DeviceWatts        bool          `json:"device_watts"`         // Whether the watts are from a power meter, false if estimated
	Distance           float32       `json:"distance"`             // The lap's distance, in meters
	ElapsedTime        int           `json:"elapsed_time"`         // The lap's elapsed time, in seconds
	EndIndex           int           `json:"end_index"`            // The end index of this effort in its activity's stream
	LapIndex           int           `json:"lap_index"`            // The index of this lap in the activity it belongs to
	MaxHeartRate       float32       `json:"max_heartrate"`        // The maximum heartrate of this lap, in beats per minute
	MaxSpeed           float32       `json:"max_speed"`            // The maximum speed of this lat, in meters per second
	MovingTime         int           `json:"moving_time"`          // The lap's moving time, in seconds
	Name               string        `json:"name"`                 // The name of the lap
	PaceZone           int           `json:"pace_zone"`            // The athlete's pace zone during this lap
	ResourceState      int8          `json:"resource_state"`       // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
	Split              int           `json:"split"`                // An instance of integer.
	StartIndex         int           `json:"start_index"`          // The start index of this effort in its activity's stream
	StartDate          TimeStamp     `json:"start_date"`           // The time at which the lap was started.
	StartDateLocal     TimeStamp     `json:"start_date_local"`     // The time at which the lap was started in the local timezone.
	TotalElevationGain float32       `json:"total_elevation_gain"` // The elevation gain of this lap, in meters
}

type PhotosSummary struct {
	Count   int                   `json:"count"`             // The number of photos
	Primary *PhotosSummaryPrimary `json:"primary,omitempty"` // An instance of PhotosSummaryPrimary.
}

type PhotosSummaryPrimary struct {
	ID        string `json:"unique_id"`
	Source    int    `json:"source"`
	MediaType int    `json:"media_type"`
	Urls      URL    `json:"urls"`
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

type ActivityType string

const (
	AlpineSki       ActivityType = "AlpineSki"
	BackgroundSki   ActivityType = "BackgroundSki"
	Canoeing        ActivityType = "Canoeing"
	Crossfit        ActivityType = "Crossfit"
	EBikeRide       ActivityType = "EBikeRide"
	Elliptical      ActivityType = "Elliptical"
	Golf            ActivityType = "Golf"
	Handcycle       ActivityType = "Handcycle"
	Hike            ActivityType = "Hike"
	IceSkate        ActivityType = "IceSkate"
	InlineSkate     ActivityType = "InlineSkate"
	Kayaking        ActivityType = "Kayaking"
	Kitesurf        ActivityType = "Kitesurf"
	NordicSki       ActivityType = "NordicSki"
	Ride            ActivityType = "Ride"
	RockClimbing    ActivityType = "RockClimbing"
	RollerSki       ActivityType = "RollerSki"
	Rowing          ActivityType = "Rowing"
	Run             ActivityType = "Run"
	Sail            ActivityType = "Sail"
	SkateBoard      ActivityType = "SkateBoard"
	Snowboard       ActivityType = "Snowboard"
	Snowshoe        ActivityType = "Snowshoe"
	Soccer          ActivityType = "Soccer"
	StairStepper    ActivityType = "StairStepper"
	StandupPaddling ActivityType = "StandupPaddling"
	Surfing         ActivityType = "Surfing"
	Swim            ActivityType = "Swim"
	Velomobile      ActivityType = "Velomobile"
	VirtualRide     ActivityType = "VirtualRide"
	VirtualRun      ActivityType = "VirtualRun"
	Walk            ActivityType = "Walk"
	WeightTraining  ActivityType = "WeightTraining"
	Wheelchair      ActivityType = "Wheelchair"
	Windsurf        ActivityType = "Windsurf"
	Workout         ActivityType = "Workout"
	Yoga            ActivityType = "Yoga"
)

type SportType string

const (
	AlpineSkiSport                     SportType = "AlpineSki"
	BackcountrySkiSport                SportType = "BackcountrySki"
	BadmintonSport                     SportType = "Badminton"
	CanoeingSport                      SportType = "Canoeing"
	CrossfitSport                      SportType = "Crossfit"
	EBikeRideSport                     SportType = "EBikeRide"
	EllipticalSport                    SportType = "Elliptical"
	EMountainBikeRideSport             SportType = "EMountainBikeRide"
	GolfSport                          SportType = "Golf"
	GravelRideSport                    SportType = "GravelRide"
	HandcycleSport                     SportType = "Handcycle"
	HighIntensityIntervalTrainingSport SportType = "HighIntensityIntervalTraining"
	HikeSport                          SportType = "Hike"
	IceSkateSport                      SportType = "IceSkate"
	InlineSkateSport                   SportType = "InlineSkate"
	KayakingSport                      SportType = "Kayaking"
	KitesurfSport                      SportType = "Kitesurf"
	MountainBikeRideSport              SportType = "MountainBikeRide"
	NordicSkiSport                     SportType = "NordicSki"
	PickleballSport                    SportType = "Pickleball"
	PilatesSport                       SportType = "Pilates"
	RacquetballSport                   SportType = "Racquetball"
	RideSport                          SportType = "Ride"
	RockClimbingSport                  SportType = "RockClimbing"
	RollerSkiSport                     SportType = "RollerSki"
	RowingSport                        SportType = "Rowing"
	RunSport                           SportType = "Run"
	SailSport                          SportType = "Sail"
	Skateboard                         SportType = "Skateboard"
	SnowboardSport                     SportType = "Snowboard"
	SnowshoeSport                      SportType = "Snowshoe"
	SoccerSport                        SportType = "Soccer"
	SquashSport                        SportType = "Squash"
	StairStepperSport                  SportType = "StairStepper"
	StandUpPaddlingSport               SportType = "StandUpPaddling"
	SurfingSport                       SportType = "Surfing"
	SwimSport                          SportType = "Swim"
	TableTennisSport                   SportType = "TableTennis"
	TennisSport                        SportType = "Tennis"
	TrailRunSport                      SportType = "TrailRun"
	VelomobileSport                    SportType = "Velomobile"
	VirtualRideSport                   SportType = "VirtualRide"
	VirtualRowSport                    SportType = "VirtualRow"
	VirtualRunSport                    SportType = "VirtualRun"
	WalkSportType                      SportType = "Walk"
	WeightTrainingSport                SportType = "WeightTraining"
	WheelchairSport                    SportType = "Wheelchair"
	WindsurfSport                      SportType = "Windsurf"
	WorkoutSport                       SportType = "Workout"
	YogaSport                          SportType = "Yoga"
)

// *************** Methods ********************

type ActivityService service

const activitiesPath string = "/api/v3/activities"

type ActivityPayload struct {
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
func (s *ActivityService) New(ctx context.Context, accessToken string, activity ActivityPayload) (*ActivityDetailed, *http.Response, error) {

	values, err := query.Values(activity)
	if err != nil {
		return nil, nil, err
	}

	if activity.Trainer {
		values.Set("trainer", "1")
	} else {
		values.Set("trainer", "0")
	}

	if activity.Commute {
		values.Set("commute", "1")
	} else {
		values.Set("commute", "0")
	}

	// Create the request
	req, err := s.client.NewRequest(http.MethodPost, activitiesPath, values, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	fmt.Println(req)

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
func (s *ActivityService) GetByID(ctx context.Context, accessToken string, id int, includeEfforts bool) (*ActivityDetailed, *http.Response, error) {
	q := url.Values{}
	q.Add("include_all_efforts", strconv.FormatBool(includeEfforts))

	urlStr := fmt.Sprintf("%s/%d", activitiesPath, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, q, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	fmt.Println(req)

	activity := new(ActivityDetailed)
	resp, err := s.client.DoAndParse(ctx, req, activity)
	if err != nil {
		return nil, resp, err
	}

	return activity, resp, nil
}

// ListActivityLaps retrieves the laps of a specific activity.
// GET: https://www.strava.com/api/v3/activities/{id}/laps
func (s *ActivityService) ListActivityLaps(ctx context.Context, accessToken string, id int) ([]Lap, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/laps", activitiesPath, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var laps []Lap

	// Send the request and parse the response
	resp, err := s.client.DoAndParse(ctx, req, &laps)
	if err != nil {
		return nil, resp, err
	}

	return laps, resp, nil
}

// GetActivityZones retrieves the zones of a given activity
// Requires actvity:read scope for Everyone and Followers activities
// Required activity:read_all for OnlyMe activities
//
// GET: https://www.strava.com/api/v3/activities/{id}/zones
func (s *ActivityService) GetActivityZones(ctx context.Context, accessToken string, id int) ([]ActivityZone, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/zones", activitiesPath, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var zones []ActivityZone

	resp, err := s.client.DoAndParse(ctx, req, &zones)
	if err != nil {
		return nil, resp, err
	}

	return zones, resp, nil
}

type ListCommentOptions struct {
	Page        int    `url:"page,omitempty"`         // Page number. Defaults to 1.
	PerPage     int    `url:"per_page,omitmepty"`     // Number of items per page. Defaults to 30
	PageSize    int    `url:"page_size,omitempty"`    // Number of items per page. Defaults to 30
	AfterCursor string `url:"after_cursor,omitempty"` // Cursor of the last item in the previous page of results, used to request the subsequent page of results. When omitted, the first page of results is fetched.
}

// ListActivityComments retrives the comments of a given activity.
// Requires activity:read for Everyone and Followers activities.
// Requires actvity:read_all for Only Me activities.
//
// GET: https://www.strava.com/api/v3/activities/{id}/comments
func (s *ActivityService) ListActivityComments(ctx context.Context, accessToken string, id int, opts *ListCommentOptions) ([]Comment, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/comments", activitiesPath, id)

	v, err := query.Values(opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, urlStr, v, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var comments []Comment

	resp, err := s.client.DoAndParse(ctx, req, &comments)
	if err != nil {
		return nil, resp, err
	}

	return comments, resp, nil
}

// ListActivityKudoers retrieves the athletes who kudoed a given activity.
// Requires activity:read scope for Everyone and Followers activities
// Required activity:read_all scope for Only Me activities
//
// GET: https://www.strava.com/api/v3/activities/{id}/kudos
func (s *ActivityService) ListActivityKudoers(ctx context.Context, accessToken string, id int, opts *ListOptions) ([]AthleteSummary, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/kudos", activitiesPath, id)

	v, err := query.Values(opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, urlStr, v, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var kudoers []AthleteSummary
	resp, err := s.client.DoAndParse(ctx, req, &kudoers)
	if err != nil {
		return nil, resp, err
	}

	return kudoers, resp, nil
}

// GetActivityStreamTypes retrieves the activity's streams.
// Requires activity:read scope. Required activity:read_all scope for Only Me activities.
// By default return the primary stream of the activity.
func (s *ActivityService) GetActivityStreams(ctx context.Context, accessToken string, id int, streamTypes []StreamType) ([]Stream, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/streams", activitiesPath, id)

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
