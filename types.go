package gostrava

import (
	"time"
)

type DateTime string

func (dt *DateTime) Parse() (time.Time, error) {
	return time.Parse(time.RFC3339, string(*dt))
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

type ClubSportType string

var ClubSportTypes = struct {
	Cycling   ClubSportType
	Running   ClubSportType
	Triathlon ClubSportType
	Other     ClubSportType
}{
	"cycling", "running", "triathlon", "other",
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

type ExplorerResponse struct {
	Segments []ExplorerSegment `json:"segments"` // The set of segments matching an explorer request
}

type ExplorerSegment struct {
	ID                int     `json:"id"`                  // The unique identifier of this segment
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

type Lap struct {
	ID                 int          `json:"id"`                        // The unique identifier of this lap
	Activity           MetaActivity `json:"activity"`                  // An instance of MetaActivity.
	Athlete            MetaAthlete  `json:"athlete"`                   // AN instance of MetaAthlete.
	AvgCadence         float32      `json:"average_cadence,omitempty"` // The lap's average cadence
	AvgHeartRate       float32      `json:"average_heartrate"`         // The lap's average heartrate
	AvgSpeed           float32      `json:"average_speed"`             // The lap's average speed
	DeviceWatts        bool         `json:"device_watts"`              // Whether the watts are from a power meter, false if estimated
	Distance           float32      `json:"distance"`                  // The lap's distance, in meters
	ElapsedTime        int          `json:"elapsed_time"`              // The lap's elapsed time, in seconds
	EndIndex           int          `json:"end_index"`                 // The end index of this effort in its activity's stream
	LapIndex           int          `json:"lap_index"`                 // The index of this lap in the activity it belongs to
	MaxHeartRate       float32      `json:"max_heartrate"`             // The maximum heartrate of this lap, in beats per minute
	MaxSpeed           float32      `json:"max_speed"`                 // The maximum speed of this lat, in meters per second
	MovingTime         int          `json:"moving_time"`               // The lap's moving time, in seconds
	Name               string       `json:"name"`                      // The name of the lap
	ResourceState      uint8        `json:"resource_state"`            // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
	Split              int          `json:"split"`                     // An instance of integer.
	StartIndex         int          `json:"start_index"`               // The start index of this effort in its activity's stream
	StartDate          DateTime     `json:"start_date"`                // The time at which the lap was started.
	StartDateLocal     DateTime     `json:"start_date_local"`          // The time at which the lap was started in the local timezone.
	TotalElevationGain float32      `json:"total_elevation_gain"`      // The elevation gain of this lap, in meters
	PaceZone           *int         `json:"pace_zone,omitempty"`       // The athlete's pace zone during this lap
}

type LatLng []float64 // A collection of float objects. A pair of latitude/longitude coordinates, represented as an array of 2 floating point numbers.

type PhotosSummary struct {
	Count   int                   `json:"count"`   // The number of photos
	Primary *PhotosSummaryPrimary `json:"primary"` // An instance of PhotosSummaryPrimary.
}

type PhotosSummaryPrimary struct {
	ID       int    `json:"id"`
	Source   int    `json:"source"`
	UniqueID string `json:"unique_id"`
	Urls     string `json:"string"`
}

type PolylineMap struct {
	ID              string  `json:"id"`                       // The identifier of the map
	Polyline        *string `json:"polyline,omitempty"`       // The polyline of the map, only returned on detailed representation of an object
	ResourceState   uint8   `json:"resource_state,omitempty"` //
	SummaryPolyline string  `json:"summary_polyline"`         // The summary polyline of the map
}

type PowerZoneRanges struct {
	Zones ZoneRanges `json:"zones"` // An instance of ZoneRanges.
}

type Route struct {
	Athlete             SummaryAthlete   `json:"athlete"`               // An instance of SummaryAthlete.
	Description         string           `json:"description"`           // The description of the route
	Distance            float32          `json:"distance"`              // The route's distance, in meters
	ElevationGain       float32          `json:"elevation_gain"`        // The route's elevation gain.
	ID                  int              `json:"id"`                    // The unique identifier of this route
	IdStr               string           `json:"id_str"`                // The unique identifier of the route in string format
	Map                 PolylineMap      `json:"map"`                   // An instance of PolylineMap.
	Name                string           `json:"name"`                  // The name of this route
	Private             bool             `json:"private"`               // Whether this route is private
	Starred             bool             `json:"starred"`               // Whether this route is starred by the logged-in athlete
	Timestamp           int              `json:"timestamp"`             // An epoch timestamp of when the route was created
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
	"Ride", "Run",
}

type Split struct {
	AvgGradeAdjustedSpeed *float32 `json:"average_grade_adjusted_speed"`
	AvgHeartRate          float32  `json:"average_heartrate,omitempty"` // The average heartrate of this split, in beats per minute
	AvgSpeed              float32  `json:"average_speed"`               // The average speed of this split, in meters per second
	Distance              float32  `json:"distance"`                    //  The distance of this split, in meters
	ElapsedTime           int      `json:"elapsed_time"`                // The elapsed time of this split, in seconds
	ElevationDiff         float32  `json:"elevation_difference"`        // The elevation difference of this split, in meters
	MovingTime            int      `json:"moving_time"`                 // The moving time of this split, in seconds
	PaceZone              int      `json:"pace_zone"`                   // The pacing zone of this split
	Split                 int      `json:"split"`
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
	ID           string  `json:"id"`             // The gear's unique identifier.
	ResourceRate uint8   `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
	Primary      bool    `json:"primary"`        // Whether this gear's is the owner's default one.
	Name         string  `json:"name"`           // The gear's name.
	Distance     float32 `json:"distance"`       // The distance logged with this gear.
}

type SummaryPRSegmentEffort struct {
	PRActivityID  int      `json:"pr_activity_id"`  // The unique identifier of the activity related to the PR effort.
	PRElapsedTime int      `json:"pr_elapsed_time"` // The elapsed time ot the PR effort.
	PRDate        DateTime `json:"pr_date"`         //  The time at which the PR effort was started.
	EffortCount   int      `json:"effort_count"`    // Number of efforts by the authenticated athlete on this segment.
}

type SummarySegment struct {
	ID                  int                    `json:"id"`                     // The unique identifier of this segment
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
	ID             int      `json:"id"`               // The unique identifier of this effort
	ActivityID     int      `json:"activity_id"`      // The unique identifier of the activity related to this effort
	ElapsedTime    int      `json:"elapsed_time"`     // The effort's elapsed time
	StartDate      DateTime `json:"start_date"`       // The time at which the effort was started.
	StartDateLocal DateTime `json:"start_date_local"` // The time at which the effort was started in the local timezone.
	Distance       float32  `json:"distance"`         //  The effort's distance in meters
	IsKom          bool     `json:"is_kom"`           // Whether this effort is the current best on the leaderboard
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
	ID         int    `json:"id"`          // The unique identifier of the upload
	IdSrt      string `json:"id_str"`      // The unique identifier of the upload in string format
	ExternalID string `json:"external_id"` // The external identifier of the upload
	Error      string `json:"error"`       // The error associated with this upload
	Status     string `json:"string"`      // The status of this upload
	ActivityID int    `json:"activity_id"` // The identifier of the activity this upload resulted into
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
