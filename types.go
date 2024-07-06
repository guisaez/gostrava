package gostrava

import (
	"time"
)

type DateTime string

func (dt *DateTime) Parse() (time.Time, error) {
	return time.Parse(time.RFC3339, string(*dt))
}

// A set of rolled-up statistics and totals for an athlete
type ActivityStats struct {
	AllRideTotals        ActivityTotal `json:"all_ride_totals"`    // The all time ride stats for the athlete.
	AllRunTotals         ActivityTotal `json:"all_run_totals"`     // The all time run stats for the athlete.
	AllSwimTotals        ActivityTotal `json:"all_swim_totals"`    // The all time swim stats for the athlete.
	RecentRideTotals     ActivityTotal `json:"recent_ride_totals"` // The recent (last 4 weeks) ride stats for the athlete.
	RecentRunTotals      ActivityTotal `json:"recent_run_totals"`  // The recent (last 4 weeks) run stats for the athlete.
	RecentSwimTotals     ActivityTotal `json:"recent_swim_totals"` // The recent (last 4 weeks) swim stats for the athlete.
	YearToDateRideTotals ActivityTotal `json:"ytd_ride_totals"`    // The year to date ride stats for the athlete.
	YearToDateRunTotals  ActivityTotal `json:"ytd_run_totals"`     // The year to date run stats for the athlete.
	YearToDateSwimTotals ActivityTotal `json:"ytd_swim_totals"`    // The year to date swim stats for the athlete.
}

// A roll-up of metrics pertaining to a set of activities. Values are in seconds and meters.
type ActivityTotal struct {
	Count            int     `json:"count"`             // The number of activities considered in this total.
	Distance         float32 `json:"distance"`          // The total distance covered by the considered activities.
	MovingTime       int     `json:"moving_time"`       // The total moving time of the considered activities.
	ElapsedTime      int     `json:"elapsed_time"`      // The total elapsed time of the considered activities.
	ElevationGain    float32 `json:"elevation_gain"`    // The total elevation gain of the considered activities.
	AchievementCount int     `json:"achievement_count"` // The total number of achievements of the considered activities.
}

// An enumeration of the types an activity may have. Note that this enumeration does not include new sport types
type ActivityType string

var ActivityTypes = struct {
	AlpineSki       ActivityType
	BackgroundSki   ActivityType
	Canoening       ActivityType
	Crossfit        ActivityType
	EBikeRide       ActivityType
	Elliptical      ActivityType
	Golf            ActivityType
	Handcycle       ActivityType
	Hike            ActivityType
	IceSkate        ActivityType
	InlineSkate     ActivityType
	Kayaking        ActivityType
	Kitesurf        ActivityType
	Nordicski       ActivityType
	Ride            ActivityType
	RockClimbing    ActivityType
	RollerSki       ActivityType
	Rowing          ActivityType
	Run             ActivityType
	Sail            ActivityType
	SkateBoard      ActivityType
	Snowboard       ActivityType
	Snowshoe        ActivityType
	Soccer          ActivityType
	StairStepper    ActivityType
	StandupPaddling ActivityType
	Surfing         ActivityType
	Swim            ActivityType
	Velomobile      ActivityType
	VirtualRide     ActivityType
	VirtualRun      ActivityType
	Walk            ActivityType
	WeightTraining  ActivityType
	Wheelchair      ActivityType
	Windsurf        ActivityType
	Workout         ActivityType
	Yoga            ActivityType
}{
	"AlpineSki", "BackgroundSki", "Canoening", "Crossfit", "EBikeRide", "Elliptical", "Golf", "Handcycle", "Hike",
	"IceSkate", "InlineSkate", "Kayaking", "Kitesurf", "Nordicski", "Ride", "RockClimbing", "RollerSki", "Rowing", "Run",
	"Sail", "SkateBoard", "Snowboard", "Snowshoe", "Soccer", "StairStepper", "StandupPaddling", "Surfing", "Swim", "Velomobile",
	"VirtualRide", "VirtualRun", "Walk", "WeightTraining", "Wheelchair", "Windsurf", "Workout", "Yoga",
}

type ActivityZone struct {
	Score               int                   `json:"score"`
	DistributionBuckets TimedZoneDistribution `json:"distribution_buckets"`
	Type                ActivityZoneType      `json:"type"`
	SensorBased         bool                  `json:"sensor_based"`
	Points              int                   `json:"points"`
	CustomZones         bool                  `json:"custom_zones"`
	Max                 int                   `json:"max"`
}

// May take one of the following values: ActivityZoneTypes.HeartRate, ActivityZoneTypes.Power
type ActivityZoneType string

var ActivityZoneTypes = struct {
	HeartRate ActivityZoneType
	Power     ActivityZoneType
}{
	"heartrate", "power",
}

type AltitudeStream struct {
	BaseStream
	Data []float32 `json:"data"` // The sequence of altitude values for this stream, in meters
}

type BaseStream struct {
	OriginalSize int                  `json:"original_size"`        // The number of data points in this stream
	Resolution   BaseStreamResolution `json:"BaseStreamResolution"` // The level of detail (sampling) in which this stream was returned May take one of the following values: BaseStreamResolutions.Low, BaseStreamResolution.Medium, BaseStreamResolution.High
	SeriesType   BaseStreamSeriesType `json:"series_type"`          // The base series used in the case the stream was downsampled May take one of the following values: BaseStreamSeriesTypes.Distance, BaseStreamSeriesTypes.Time
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

type CadenceStream struct {
	BaseStream
	Data []int `json:"data"` //  The sequence of cadence values for this stream, in rotations per minute
}

type ClubMembership string

var ClubMembershipStatus = struct {
	Member  ClubMembership
	Pending ClubMembership
}{
	"member", "pending",
}

type ClubSportType string

var ClubSportTypes = struct {
	Cycling   ClubSportType
	Running   ClubSportType
	Triathlon ClubSportType
	Other     ClubSportType
}{
	"cycling", "running", "triathlon", "other",
}

type Comment struct {
	ID         int64          `json:"id"`          // The unique identifier of this comment
	ActivityID int64          `json:"activity_id"` // The identifier of the activity this comment is related to
	Text       string         `json:"text"`        // The content of the comment
	Athlete    SummaryAthlete `json:"athlete"`     // An instance of SummaryAthlete.
	CreatedAt  DateTime       `json:"created_at"`  // The time at which this comment was created.
}

type DetailedGear struct {
	SummaryGear
	BrandName   string `json:"brand_name"`  // The gear's brand name.
	ModelName   string `json:"model_name"`  // The gear's model name.
	FrameType   int    `json:"frame_type"`  // The gear's frame type (bike only).
	Description string `json:"description"` // The gear's description.
}

type DetailedSegment struct {
	SummarySegment
	CreatedAt          DateTime    `json:"created_at"`           // The time at which the segment was created.
	UpdatedAt          DateTime    `json:"updated_at"`           // The time at which the segment was last updated.
	TotalElevationGain float32     `json:"total_elevation_gain"` // The segment's total elevation gain.
	Map                PolylineMap `json:"map"`                  // An instance of PolylineMap.
	EffortCount        int         `json:"effort_count"`         // The total number of efforts for this segment
	AthleteCount       int         `json:"athlete_count"`        // The number of unique athletes who have an effort for this segment
	Hazardous          bool        `json:"hazardous"`            // Whether this segment is considered hazardous
	StarCount          int         `json:"star_count"`           // The number of stars for this segment
}

type DetailedSegmentEffort struct {
	Name         string         `json:"name"`              // The name of the segment on which this effort was performed
	Activity     MetaActivity   `json:"activity"`          // An instance of MetaActivity.
	Athlete      MetaAthlete    `json:"athlete"`           // An instance of MetaAthlete.
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

type DistanceStream struct {
	BaseStream
	Data []float32 `json:"data"` // The sequence of distance values for this stream, in meters
}

type ExplorerResponse struct {
	Segments []ExplorerSegment `json:"segments"` // The set of segments matching an explorer request
}

type ExplorerSegment struct {
	ID                int64   `json:"id"`                  // The unique identifier of this segment
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

type HeartRateZoneRanges struct {
	CustomZones bool       `json:"custom_zones"` // Whether the athlete has set their own custom heart rate zones
	Zones       ZoneRanges `json:"zones"`        // An instance of ZoneRanges.
}

type HeartrateStream struct {
	BaseStream
	Data []int `json:"data"` // The sequence of heart rate values for this stream, in beats per minute
}

type Lap struct {
	ID                 int64        `json:"id"`                   // The unique identifier of this lap
	Activity           MetaActivity `json:"activity"`             // An instance of MetaActivity.
	Athlete            MetaAthlete  `json:"athlete"`              // AN instance of MetaAthlete.
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
	Split              int          `json:"split"`                // An instance of integer.
	StartIndex         int          `json:"start_index"`          // The start index of this effort in its activity's stream
	StartDate          DateTime     `json:"start_date"`           // The time at which the lap was started.
	StartDateLocal     DateTime     `json:"start_date_local"`     // The time at which the lap was started in the local timezone.
	TotalElevationGain float32      `json:"total_elevation_gain"` // The elevation gain of this lap, in meters
	PaceZone           int          `json:"pace_zone"`            // The athlete's pace zone during this lap
}

type LatLng [2]float64 // A collection of float objects. A pair of latitude/longitude coordinates, represented as an array of 2 floating point numbers.

type LatLngStream struct {
	BaseStream
	Data []LatLng `json:"data"` // The sequence of lat/long values for this stream
}

type Measurement string

var Measurements = struct {
	Feet   Measurement
	Meters Measurement
}{
	"feet", "meters",
}

type MovingStream struct {
	BaseStream
	Data []bool // The sequence of moving values for this stream, as boolean values
}

type PhotosSummary struct {
	Count   int                  `json:"count"`   // The number of photos
	Primary PhotosSummaryPrimary `json:"primary"` // An instance of PhotosSummaryPrimary.
}

type PhotosSummaryPrimary struct {
	ID       int64  `json:"id"`
	Source   int    `json:"source"`
	UniqueID string `json:"unique_id"`
	Urls     string `json:"string"`
}

type PolylineMap struct {
	ID              string `json:"id"`               // The identifier of the map
	Polyline        string `json:"polyline"`         // The polyline of the map, only returned on detailed representation of an object
	SummaryPolyline string `json:"summary_polyline"` // The summary polyline of the map
}

type PowerStream struct {
	BaseStream
	Data []int `json:"data"` // The sequence of power values for this stream, in watts
}

type PowerZoneRanges struct {
	Zones ZoneRanges `json:"zones"` // An instance of ZoneRanges.
}

type ResourceState uint8

var ResourceStates = struct {
	Meta    ResourceState
	Summary ResourceState
	Detail  ResourceState
}{
	1, 2, 3,
}

type Route struct {
	Athlete             SummaryAthlete   `json:"athlete"`               // An instance of SummaryAthlete.
	Description         string           `json:"description"`           // The description of the route
	Distance            float32          `json:"distance"`              // The route's distance, in meters
	ElevationGain       float32          `json:"elevation_gain"`        // The route's elevation gain.
	ID                  int64            `json:"id"`                    // The unique identifier of this route
	IdStr               string           `json:"id_str"`                // The unique identifier of the route in string format
	Map                 PolylineMap      `json:"map"`                   // An instance of PolylineMap.
	Name                string           `json:"name"`                  // The name of this route
	Private             bool             `json:"private"`               // Whether this route is private
	Starred             bool             `json:"starred"`               // Whether this route is starred by the logged-in athlete
	Timestamp           int64            `json:"timestamp"`             // An epoch timestamp of when the route was created
	Type                RouteType        `json:"type"`                  // This route's type RouteTypes.Ride, RouteTypes.Run
	SubType             SubRouteType     `json:"sub_type"`              // This route's sub-type (SubRouteTypes.Road, SubRouteTypes.MountainBike, SubRouteTypes.Cross, SubRouteTypes.Trail, SubRouteTypes.Mixed)
	CreatedAt           DateTime         `json:"created_at"`            // The time at which the route was created
	UpdatedAt           DateTime         `json:"updated_at"`            // The time at which the route was last updated
	EstimatedMovingTime int              `json:"estimated_moving_time"` // Estimated time in seconds for the authenticated athlete to complete route
	Segments            []SummarySegment `json:"segments"`              // The segments traversed by this route
	Waypoints           []Waypoint       `json:"waypoints"`             // The custom waypoints along this route
}

type RouteType int8

var RouteTypes = struct {
	Ride RouteType
	Run  RouteType
}{
	1, 2,
}

type SegmentActivityType string

var SegmentActivityTypes = struct {
	Ride SegmentActivityType
	Run  SegmentActivityType
}{
	"Ride", "Rune",
}

type SmoothGradeStream struct {
	BaseStream
	Data []float32 `json:"data"` // The sequence of grade values for this stream, as percents of a grade
}

type SmoothVelocityStream struct {
	BaseStream
	Data []float32 `json:"data"` // The sequence of velocity values for this stream, in meters per second
}

type Split struct {
	AvgSpeed      float32 `json:"average_speed"`        // The average speed of this split, in meters per second
	Distance      float32 `json:"distance"`             //  The distance of this split, in meters
	ElapsedTime   int     `json:"elapsed_time"`         // The elapsed time of this split, in seconds
	ElevationDiff float32 `json:"elevation_difference"` // The elevation difference of this split, in meters
	PaceZone      int     `json:"pace_zone"`            // The pacing zone of this split
	MovingTime    int     `json:"moving_time"`          // The moving time of this split, in seconds
	Split         int     `json:"split"`
}

type SportType string // An enumeration of the sport types an activity may have. Distinct from ActivityType in that it has new types (e.g. MountainBikeRide)

var SportTypes = struct {
	AlpineSki                     SportType
	BackcountrySki                SportType
	Badminton                     SportType
	Canoeing                      SportType
	Crossfit                      SportType
	EBikeRide                     SportType
	Elliptical                    SportType
	EMountainBikeRide             SportType
	Golf                          SportType
	GravelRide                    SportType
	Handcycle                     SportType
	HighIntensityIntervalTraining SportType
	Hike                          SportType
	IceSkate                      SportType
	InlineSkate                   SportType
	Kayaking                      SportType
	Kitesurf                      SportType
	MountainBikeRide              SportType
	NordicSki                     SportType
	Pickleball                    SportType
	Pilates                       SportType
	Racquetball                   SportType
	Ride                          SportType
	RockClimbing                  SportType
	RollerSki                     SportType
	Rowing                        SportType
	Run                           SportType
	Sail                          SportType
	Skateboard                    SportType
	Snowboard                     SportType
	Snowshoe                      SportType
	Soccer                        SportType
	Squash                        SportType
	StairStepper                  SportType
	StandUpPaddling               SportType
	Surfing                       SportType
	Swim                          SportType
	TableTennis                   SportType
	Tennis                        SportType
	TrailRun                      SportType
	Velomobile                    SportType
	VirtualRide                   SportType
	VirtualRow                    SportType
	VirtualRun                    SportType
	Walk                          SportType
	WeightTraining                SportType
	Wheelchair                    SportType
	Windsurf                      SportType
	Workout                       SportType
	Yoga                          SportType
}{
	"AlpineSki", "BackcountrySki", "Badminton", "Canoeing", "Crossfit", "EBikeRide", "Elliptical", "EMountainBikeRide", "Golf",
	"GravelRide", "Handcycle", "HighIntensityIntervalTraining", "Hike", "IceSkate", "InlineSkate", "Kayaking", "Kitesurf",
	"MountainBikeRide", "NordicSki", "Pickleball", "Pilates", "Racquetball", "Ride", "RockClimbing", "RollerSki", "Rowing",
	"Run", "Sail", "Skateboard", "Snowboard", "Snowshoe", "Soccer", "Squash", "StairStepper", "StandUpPaddling", "Surfing",
	"Swim", "TableTennis", "Tennis", "TrailRun", "Velomobile", "VirtualRide", "VirtualRow", "VirtualRun", "Walk",
	"WeightTraining", "Wheelchair", "Windsurf", "Workout", "Yoga",
}

type StreamSet struct {
	TimeStream           TimeStream           `json:"time"`            // An instance of TimeStream.
	DistanceStream       DistanceStream       `json:"distance"`        // An instance of DistanceStream.
	LatLngStream         LatLngStream         `json:"latlng"`          // An instance of LatLngStream.
	AltitudeStream       AltitudeStream       `json:"altitude"`        // An instance of AltitudeStream.
	SmoothVelocityStream SmoothVelocityStream `json:"velocity_smooth"` // An instance of SmoothVelocityStream.
	HeartRateStream      HeartrateStream      `json:"heartrate"`       // An instance of HeartrateStream.
	CadenceStream        CadenceStream        `json:"cadence"`         // An instance of CadenceStream.
	WattsStream          PowerStream          `json:"watts"`           // An instance of PowerStream.
	TempStream           TemperatureStream    `json:"temp"`            // An instance of TemperatureStream.
	MovingStream         MovingStream         `json:"moving"`          // An instance of MovingStream.
	SmoothGradeStream    SmoothGradeStream    `json:"grade_smooth"`    // An instance of SmoothGradeStream.
}

type SubRouteType int8

var SubRouteTypes = struct {
	Road         SubRouteType
	MountainBike SubRouteType
	Cross        SubRouteType
	Trail        SubRouteType
	Mixed        SubRouteType
}{
	1, 2, 3, 4, 5,
}

type SummaryGear struct {
	ID           string        `json:"id"`             // The gear's unique identifier.
	ResourceRate ResourceState `json:"resource_state"` // Resource state, indicates level of detail. Possible values: ResourceStates.Meta, ResourceStates.Summary, ResourceStates.Detail
	Primary      bool          `json:"primary"`        // Whether this gear's is the owner's default one.
	Name         string        `json:"name"`           // The gear's name.
	Distance     float32       `json:"distance"`       // The distance logged with this gear.
}

type SummaryPRSegmentEffort struct {
	PRActivityID  int64    `json:"pr_activity_id"`  // The unique identifier of the activity related to the PR effort.
	PRElapsedTime int      `json:"pr_elapsed_time"` // The elapsed time ot the PR effort.
	PRDate        DateTime `json:"pr_date"`         //  The time at which the PR effort was started.
	EffortCount   int      `json:"effort_count"`    // Number of efforts by the authenticated athlete on this segment.
}

type SummarySegment struct {
	ID                  int64                  `json:"id"`                     // The unique identifier of this segment
	Name                string                 `json:"name"`                   // The name of this segment
	ActivityType        SegmentActivityType    `json:"activity_type"`          // May take one of the following values: SegmentActivityTypes.Ride, SegmentActivityTypes.Run
	Distance            float32                `json:"distance"`               // The segment's distance, in meters
	AvgGrade            float32                `json:"average_grade"`          // The segment's average grade, in percents
	MaximumGrade        float32                `json:"maximum_grade"`          // The segments's maximum grade, in percents
	ElevationHigh       float32                `json:"elevation_high"`         // The segments's highest elevation, in meters
	ElevationLow        float32                `json:"elevation_low"`          // The segments's lowest elevation, in meters
	StartLatLng         LatLng                 `json:"start_latlng"`           // An instance of LatLng.
	EndLatLng           LatLng                 `json:"end_latlng"`             // An instance of LatLng.
	ClimbCategory       int8                   `json:"climb_category"`         // The category of the climb [0, 5]. Higher is harder ie. 5 is Hors catégorie, 0 is uncategorized in climb_category.
	City                string                 `json:"city"`                   // The segments's city.
	State               string                 `json:"state"`                  // The segments's state or geographical region.
	Country             string                 `json:"country"`                // The segment's country.
	Private             bool                   `json:"private"`                // Whether this segment is private.
	AthletePREffort     SummaryPRSegmentEffort `json:"athlete_pr_effort"`      // An instance of SummaryPRSegmentEffort.
	AthleteSegmentStats SummarySegmentEffort   `json:"athlete_segmentZ_stats"` // An instance ofSummarySegmentEffort.
}

type SummarySegmentEffort struct {
	ID             int64    `json:"id"`               // The unique identifier of this effort
	ActivityID     int64    `json:"activity_id"`      // The unique identifier of the activity related to this effort
	ElapsedTime    int      `json:"elapsed_time"`     // The effort's elapsed time
	StartDate      DateTime `json:"start_date"`       // The time at which the effort was started.
	StartDateLocal DateTime `json:"start_date_local"` // The time at which the effort was started in the local timezone.
	Distance       float32  `json:"distance"`         //  The effort's distance in meters
	IsKom          bool     `json:"is_kom"`           // Whether this effort is the current best on the leaderboard
}

type TemperatureStream struct {
	BaseStream
	Data []int `json:"data"` // The sequence of temperature values for this stream, in celsius degrees
}

type TimeStream struct {
	BaseStream
	Data []int `json:"data"` // The sequence of time values for this stream, in seconds
}

// A collection of #/TimedZoneRange objects. Stores the exclusive ranges representing zones and the time spent in each.
type TimedZoneDistribution []TimedZoneRange

// A union type representing the time spent in a given zone.
type TimedZoneRange struct {
	Min  int `json:"min"`  // The minimum value in the range.
	Max  int `json:"max"`  // The maximum value in the range.
	Time int `json:"time"` // The number of seconds spent in this zone
}

type Upload struct {
	ID         int64  `json:"id"`          // The unique identifier of the upload
	IdSrt      string `json:"id_str"`      // The unique identifier of the upload in string format
	ExternalID string `json:"external_id"` // The external identifier of the upload
	Error      string `json:"error"`       // The error associated with this upload
	Status     string `json:"string"`      // The status of this upload
	ActivityID int64  `json:"activity_id"` // The identifier of the activity this upload resulted into
}

type Waypoint struct {
	LatLng            LatLng   `json:"latlng"`              // The location along the route that the waypoint is closest to
	TargetLatLng      LatLng   `json:"target_latlng"`       // A location off of the route that the waypoint is (optional)
	Categories        []string `json:"categories"`          // Categories that the waypoint belongs to
	Title             string   `json:"string"`              // A title for the waypoint
	Description       string   `json:"description"`         // A description of the waypoint (optional)
	DistanceIntoRoute int      `json:"distance_into_route"` // The number meters along the route that the waypoint is located
}

type ZoneRange struct {
	Max int `json:"max"` // The maximum value in the range.
	Min int `json:"min"` // The minimum value in the range.
}

type ZoneRanges []ZoneRange

type Zones struct {
	HearRate HeartRateZoneRanges `json:"heart_rate"` // An instance of HeartRateZoneRanges.
	Power    PowerZoneRanges     `json:"power"`      // An instance of PowerZoneRanges.
}
