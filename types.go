package gostrava

import (
	"encoding/json"
	"errors"
)

type RequestParams struct {
	Page    int // Page number. Defaults to 1
	PerPage int // Number of items per page. Defaults to 30
}

// A roll-up of metrics pertaining to a set of activities. Values are in seconds and meters.
type ActivityTotal struct {
	Count            *int     `json:"count,omitempty"`             // The number of activities considered in this total.
	Distance         *float32 `json:"distance,omitempty"`          // The total distance covered by the considered activities.
	MovingTime       *int     `json:"moving_time,omitempty"`       // The total moving time of the considered activities.
	ElapsedTime      *int     `json:"elapsed_time,omitempty"`      // The total elapsed time of the considered activities.
	ElevationGain    *float32 `json:"elevation_gain,omitempty"`    // The total elevation gain of the considered activities.
	AchievementCount *int     `json:"achievement_count,omitempty"` // The total number of achievements of the considered activities.
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

type ClubSportType string

const (
	Cycling   ClubSportType = "cycling"
	Running   ClubSportType = "running"
	Triathlon ClubSportType = "triathlon"
	Other     ClubSportType = "other"
)

type RouteType string

const (
	RideRoute RouteType = "ride"
	RunRoute  RouteType = "run"
)

func (rt *RouteType) UnmarshalJSON(data []byte) error {
	var routeType int
	err := json.Unmarshal(data, &routeType)
	if err != nil {
		return err
	}

	switch routeType {
	case 1:
		*rt = RideRoute
	case 2:
		*rt = RunRoute
	default:
		return errors.New("invalid route type")
	}

	return nil
}

type SegmentActivityType string

const (
	RideSegment SegmentActivityType = "Ride"
	RunSegment  SegmentActivityType = "Run"
)

type SubRouteType string

const (
	Road         SubRouteType = "road"
	MountainBike SubRouteType = "mountain_bike"
	Cross        SubRouteType = "cross"
	Trail        SubRouteType = "train"
	Mixed        SubRouteType = "mixed"
)

func (rt *SubRouteType) UnmarshalJSON(data []byte) error {
	var subRouteType int
	err := json.Unmarshal(data, &subRouteType)
	if err != nil {
		return err
	}

	switch subRouteType {
	case 1:
		*rt = Road
	case 2:
		*rt = MountainBike
	case 3:
		*rt = Cross
	case 4:
		*rt = Trail
	case 5:
		*rt = Mixed
	default:
		return errors.New("invalid sub-route type")
	}

	return nil
}

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

type HeartRateZoneRanges struct {
	CustomZones *bool       `json:"custom_zone,omitempty"` // Whether the athlete has set their own custom heart rate zones
	Zones       *ZoneRanges `json:"zones,omitempty"`       // An instance of ZoneRanges.
}

type PowerZoneRanges struct {
	Zones *ZoneRanges `json:"zones,omitempty"` // An instance of ZoneRanges.
}

type ZoneRanges []ZoneRange

type ZoneRange struct {
	Max *int `json:"max,omitempty"` // The maximum value in the range.
	Min *int `json:"min,omitempty"` // The minimum value in the range.
}

type LatLng []float64 // A collection of float objects. A pair of latitude/longitude coordinates, represented as an array of 2 floating point numbers.

type PolylineMap struct {
	ID              string  `json:"id"`                       // The identifier of the map
	Polyline        *string `json:"polyline,omitempty"`       // The polyline of the map, only returned on detailed representation of an object
	ResourceState   uint8   `json:"resource_state,omitempty"` //
	SummaryPolyline string  `json:"summary_polyline"`         // The summary polyline of the map
}

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

// A collection of #/TimedZoneRange objects. Stores the exclusive ranges representing zones and the time spent in each.
type TimedZoneDistribution []TimedZoneRange

// A union type representing the time spent in a given zone.
type TimedZoneRange struct {
	Min  int `json:"min"`  // The minimum value in the range.
	Max  int `json:"max"`  // The maximum value in the range.
	Time int `json:"time"` // The number of seconds spent in this zone
}

type Lap struct {
	ID                 int          `json:"id"`                        // The unique identifier of this lap
	Activity           ActivityMeta `json:"activity"`                  // An instance of ActivityMeta.
	Athlete            AthleteMeta  `json:"athlete"`                   // AN instance of AthleteMeta.
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
	StartDate          TimeStamp    `json:"start_date"`                // The time at which the lap was started.
	StartDateLocal     TimeStamp    `json:"start_date_local"`          // The time at which the lap was started in the local timezone.
	TotalElevationGain float32      `json:"total_elevation_gain"`      // The elevation gain of this lap, in meters
	PaceZone           *int         `json:"pace_zone,omitempty"`       // The athlete's pace zone during this lap
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

type SummaryPRSegmentEffort struct {
	PRActivityID  int       `json:"pr_activity_id"`  // The unique identifier of the activity related to the PR effort.
	PRElapsedTime int       `json:"pr_elapsed_time"` // The elapsed time ot the PR effort.
	PRDate        TimeStamp `json:"pr_date"`         //  The time at which the PR effort was started.
	EffortCount   int       `json:"effort_count"`    // Number of efforts by the authenticated athlete on this segment.
}

type SummarySegmentEffort struct {
	ID             int       `json:"id"`               // The unique identifier of this effort
	ActivityID     int       `json:"activity_id"`      // The unique identifier of the activity related to this effort
	ElapsedTime    int       `json:"elapsed_time"`     // The effort's elapsed time
	StartDate      TimeStamp `json:"start_date"`       // The time at which the effort was started.
	StartDateLocal TimeStamp `json:"start_date_local"` // The time at which the effort was started in the local timezone.
	Distance       float32   `json:"distance"`         //  The effort's distance in meters
	IsKom          bool      `json:"is_kom"`           // Whether this effort is the current best on the leaderboard
}

type Waypoint struct {
	LatLng            LatLng   `json:"latlng"`              // The location along the route that the waypoint is closest to
	TargetLatLng      LatLng   `json:"target_latlng"`       // A location off of the route that the waypoint is (optional)
	Categories        []string `json:"categories"`          // Categories that the waypoint belongs to
	Title             string   `json:"string"`              // A title for the waypoint
	Description       string   `json:"description"`         // A description of the waypoint (optional)
	DistanceIntoRoute int      `json:"distance_into_route"` // The number meters along the route that the waypoint is located
}

type SegmentSummary struct {
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
