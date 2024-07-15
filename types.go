package gostrava

import (
	"encoding/json"
	"errors"
)










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

type ActivityZone struct {
	Score               *float32         `json:"score,omitempty"`
	DistributionBuckets []TimedZoneRange `json:"distribution_buckets,omitempty"`
	Type                *string          `json:"type,omitempty"` // May take one of the following values: heartrate, power
	SensorBased         *bool            `json:"sensor_based,omitempty"`
	Points              *int             `json:"points,omitempty"`
	CustomZones         *bool            `json:"custom_zones,omitempty"`
	Max                 *int             `json:"max,omitempty"`
}






type ClubMeta struct {
	ID            int    `json:"id"`             // The club's unique identifier.
	Name          string `json:"name"`           // The club's name.
	ResourceState int8   `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
}

type ClubSummary struct {
	ClubMeta
	Admin              bool           `json:"admin"`                // Whether the currently logged-in athlete is an administrator of this club.
	ActivityTypes      []ActivityType `json:"activity_types"`       // The activity types that count for a club. This takes precedence over sport_type.
	ActivityTypesIcon  string         `json:"activity_types_icon"`  //
	City               string         `json:"city"`                 // The club's city.
	Country            string         `json:"country"`              // The club's country.
	CoverPhoto         string         `json:"cover_photo"`          // URL to a ~1185x580 pixel cover photo.
	CoverPhotoSmall    string         `json:"cover_photo_small"`    // URL to a ~360x176 pixel cover photo.
	Dimensions         []string       `json:"dimensions"`           //
	Featured           bool           `json:"featured,omitempty"`   // Whether the club is featured or not.
	LocalizedSportType string         `json:"localized_sport_type"` //
	Membership         string         `json:"membership"`           // The membership status of the logged-in athlete. May take one of the following values: member, pending
	MemberCount        int            `json:"member_count"`         // The club's member count.
	Private            bool           `json:"private"`              // Whether the club is private.
	Profile            string         `json:"profile"`              //
	ProfileMedium      string         `json:"profile_medium"`       // URL to a 60x60 pixel profile picture.
	SportType          ClubSportType  `json:"sport_type"`           // Deprecated. Prefer to use activity_types. May take one of the following values: ClubSportTypes.Cycling, ClubSportTypes.Running, ClubSportTypes.Triathlon,  ClubSportTypes.Other
	State              string         `json:"state"`                // The club's state or geographical region.
	URL                string         `json:"url"`                  // The club's vanity URL.
	Verified           bool           `json:"verified"`             // Whether the club is verified or not.
}

type ClubDetailed struct {
	ClubSummary
	ClubType       string `json:"club_type"`
	Description    string `json:"description"`     // The club's description
	FollowingCount int    `json:"following_count"` // The number of athletes in the club that the logged-in athlete follows.
	Owner          bool   `json:"owner"`           // Whether the currently logged-in athlete is the owner of this club.
	Website        string `json:"website"`         // Club Website
}

type ClubActivity struct {
	Athlete            ClubAthlete  `json:"athlete"`              // An instance of MetaAthlete.
	Distance           float32      `json:"distance"`             // The activity's distance, in meters
	ElapsedTime        int          `json:"elapsed_time"`         // The activity's elapsed time, in seconds
	MovingTime         int          `json:"moving_time"`          // The activity's moving time, in seconds
	Name               string       `json:"name"`                 // The name of the activity
	ResourceState      int8         `json:"resource_state"`       // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
	SportType          SportType    `json:"sport_type"`           // An instance of SportType.
	Type               ActivityType `json:"activity_type"`        // Deprecated. Prefer to use sport_type
	TotalElevationGain float32      `json:"total_elevation_gain"` // The activity's total elevation gain.
}

type ClubAthlete struct {
	FirstName     string `json:"firstname"`      // Athlete First Name
	LastName      string `json:"lastname"`       // Athlete Last Name
	ResourceState int8   `json:"resource_state"` //
}

type ClubSportType string

const (
	Cycling   ClubSportType = "cycling"
	Running   ClubSportType = "running"
	Triathlon ClubSportType = "triathlon"
	Other     ClubSportType = "other"
)

type Comment struct {
	ID         int            `json:"id"`          // The unique identifier of this comment
	ActivityID int            `json:"activity_id"` // The identifier of the activity this comment is related to
	Text       string         `json:"text"`        // The content of the comment
	Athlete    AthleteSummary `json:"athlete"`     // An instance of AthleteSummary.
	CreatedAt  TimeStamp      `json:"created_at="` // The time at which this comment was created.
}

type ExplorerResponse struct {
	Segments []ExplorerSegment `json:"segments"` // The set of segments matching an explorer request
}

type ExplorerSegment struct {
	ID                 int     `json:"id"`                   // The unique identifier of this segment
	Name               string  `json:"name"`                 // The name of this segment
	ClimbCategory      int8    `json:"climb_category"`       // The category of the climb [0, 5]. Higher is harder ie. 5 is Hors catégorie, 0 is uncategorized in climb_category. If climb_category = 5, climb_category_desc = HC. If climb_category = 2, climb_category_desc = 3.
	ClimbCategoryDesc  string  `json:"climb_category_desc"`  // The description for the category of the climb May take one of the following values: NC, 4, 3, 2, 1, HC
	AvgGrade           float32 `json:"avg_grade"`            // The segment's average grade, in percents
	StartLatLng        LatLng  `json:"start_latlng"`         // An instance of LatLng.
	EndLatLng          LatLng  `json:"end_latlng"`           // An instance of LatLng.
	ElevationDiff      float32 `json:"elev_difference"`      // The segments's elevation difference, in meters
	Distance           float32 `json:"distance"`             // The segment's distance, in meters
	Polyline           string  `json:"points"`               // The polyline of the segment
	LocalLegendEnabled bool    `json:"local_legend_enabled"` //
}

type GearSummary struct {
	ID           string  `json:"id"`             // The gear's unique identifier.
	ResourceRate int8    `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
	Primary      bool    `json:"primary"`        // Whether this gear's is the owner's default one.
	Name         string  `json:"name"`           // The gear's name.
	Distance     float32 `json:"distance"`       // The distance logged with this gear.
}

type GearDetailed struct {
	GearSummary
	BrandName   string `json:"brand_name"`  // The gear's brand name.
	ModelName   string `json:"model_name"`  // The gear's model name.
	FrameType   int    `json:"frame_type"`  // The gear's frame type (bike only).
	Description string `json:"description"` // The gear's description.
}

type HeartRateZoneRanges struct {
	CustomZones bool        `json:"custom_zone"` // Whether the athlete has set their own custom heart rate zones
	Zones       []ZoneRange `json:"zones"`       // An instance of ZoneRanges.
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

type LatLng [2]float32 // A collection of float objects. A pair of latitude/longitude coordinates, represented as an array of 2 floating point numbers.

type Member struct {
	Admin         bool   `json:"admin"`          // Whether the athlete is a club admin.
	FirstName     string `json:"firstname"`      // The athlete's first name.
	LastName      string `json:"lastname"`       // The athlete's last initial.
	Membership    string `json:"membership"`     // The membership status of the logged-in athlete. May take one of the following values: member, pending
	Owner         bool   `json:"owner"`          // Whether the athlete is club owner.
	ResourceState int8   `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
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

type PolylineMap struct {
	ID              string `json:"id"` // The identifier of the map
	Polyline        string `json:"polyline"`
	ResourceState   int8   `json:"resource_state"`   //
	SummaryPolyline string `json:"summary_polyline"` // The summary polyline of the map
}

type PowerZoneRanges struct {
	Zones []ZoneRange `json:"zones"` // An instance of ZoneRanges.
}

type RequestParams struct {
	Page    int // Page number. Defaults to 1
	PerPage int // Number of items per page. Defaults to 30
}

type Route struct {
	Athlete             *AthleteSummary  `json:"athlete,omitempty"`               // An instance of AthleteSummary.
	Description         *string          `json:"description,omitempty"`           // The description of the route
	Distance            *float32         `json:"distance,omitempty"`              // The route's distance, in meters
	ElevationGain       *float32         `json:"elevation_gain,omitempty"`        // The route's elevation gain.
	ID                  *int             `json:"id,omitempty"`                    // The unique identifier of this route
	IdStr               *string          `json:"id_str,omitempty"`                // The unique identifier of the route in string format
	Map                 *PolylineMap     `json:"map,omitempty"`                   // An instance of PolylineMap.
	Name                *string          `json:"name,omitempty"`                  // The name of this route
	Private             *bool            `json:"private,omitempty"`               // Whether this route is private
	Starred             *bool            `json:"starred,omitempty"`               // Whether this route is starred by the logged-in athlete
	Timestamp           *int             `json:"timestamp,omitempty"`             // An epoch timestamp of when the route was created
	Type                *RouteType       `json:"type,omitempty"`                  // This route's type RouteTypes.Ride, RouteTypes.Run
	SubType             *SubRouteType    `json:"sub_type,omitempty"`              // This route's sub-type (SubRouteTypes.Road, SubRouteTypes.MountainBike, SubRouteTypes.Cross, SubRouteTypes.Trail, SubRouteTypes.Mixed)
	CreatedAt           *TimeStamp       `json:"created_at,omitempty"`            // The time at which the route was created
	UpdatedAt           *TimeStamp       `json:"updated_at,omitempty"`            // The time at which the route was last updated
	EstimatedMovingTime *int             `json:"estimated_moving_time,omitempty"` // Estimated time in seconds for the authenticated athlete to complete route
	Segments            []SegmentSummary `json:"segments,omitempty"`              // The segments traversed by this route
	Waypoints           []Waypoint       `json:"waypoints,omitempty"`             // The custom waypoints along this route
}

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

type Split struct {
	AvgGradeAdjustedSpeed float32 `json:"average_grade_adjusted_speed"`
	AvgHeartRate          float32 `json:"average_heartrate"`    // The average heartrate of this split, in beats per minute
	AvgSpeed              float32 `json:"average_speed"`        // The average speed of this split, in meters per second
	Distance              float32 `json:"distance"`             //  The distance of this split, in meters
	ElapsedTime           int     `json:"elapsed_time"`         // The elapsed time of this split, in seconds
	ElevationDiff         float32 `json:"elevation_difference"` // The elevation difference of this split, in meters
	MovingTime            int     `json:"moving_time"`          // The moving time of this split, in seconds
	PaceZone              int     `json:"pace_zone"`            // The pacing zone of this split
	Split                 int     `json:"split"`                // Split number
}

type SummaryPRSegmentEffort struct {
	PRActivityID  int       `json:"pr_activity_id"`  // The unique identifier of the activity related to the PR effort.
	PRElapsedTime int       `json:"pr_elapsed_time"` // The elapsed time ot the PR effort.
	PRDate        TimeStamp `json:"pr_date"`         //  The time at which the PR effort was started.
	EffortCount   int       `json:"effort_count"`    // Number of efforts by the authenticated athlete on this segment.
}

type SegmentSummaryEffort struct {
	ID             int       `json:"id"`               // The unique identifier of this effort
	ActivityID     int       `json:"activity_id"`      // The unique identifier of the activity related to this effort
	ElapsedTime    int       `json:"elapsed_time"`     // The effort's elapsed time
	StartDate      TimeStamp `json:"start_date"`       // The time at which the effort was started.
	StartDateLocal TimeStamp `json:"start_date_local"` // The time at which the effort was started in the local timezone.
	Distance       float32   `json:"distance"`         //  The effort's distance in meters
	IsKom          bool      `json:"is_kom"`           // Whether this effort is the current best on the leaderboard
}

type SegmentActivityType string

const (
	RideSegment SegmentActivityType = "Ride"
	RunSegment  SegmentActivityType = "Run"
)

type SegmentSummary struct {
	ID              int                     `json:"id"`                          // The unique identifier of this segment
	Name            string                  `json:"name"`                        // The name of this segment
	ActivityType    SegmentActivityType     `json:"activity_type"`               // May take one of the following values: SegmentActivityTypes.Ride, SegmentActivityTypes.Run
	Distance        float32                 `json:"distance"`                    // The segment's distance, in meters
	AvgGrade        float32                 `json:"average_grade"`               // The segment's average grade, in percents
	MaximumGrade    float32                 `json:"maximum_grade"`               // The segments's maximum grade, in percents
	ElevationHigh   float32                 `json:"elevation_high"`              // The segments's highest elevation, in meters
	ElevationLow    float32                 `json:"elevation_low"`               // The segments's lowest elevation, in meters
	StartLatLng     LatLng                  `json:"start_latlng"`                // An instance of LatLng.
	EndLatLng       LatLng                  `json:"end_latlng"`                  // An instance of LatLng.
	ClimbCategory   int8                    `json:"climb_category"`              // The category of the climb [0, 5]. Higher is harder ie. 5 is Hors catégorie, 0 is uncategorized in climb_category.
	City            string                  `json:"city"`                        // The segments's city.
	State           string                  `json:"state"`                       // The segments's state or geographical region.
	Country         string                  `json:"country"`                     // The segment's country.
	Private         bool                    `json:"private"`                     // Whether this segment is private.
	AthletePREffort *SummaryPRSegmentEffort `json:"athlete_pr_effort,omitempty"` // An instance of SummaryPRSegmentEffort.
}

type SegmentDetailed struct {
	SegmentSummary
	CreatedAt           TimeStamp             `json:"created_at"`                      // The time at which the segment was created.
	UpdatedAt           TimeStamp             `json:"updated_at"`                      // The time at which the segment was last updated.
	TotalElevationGain  float32               `json:"total_elevation_gain"`            // The segment's total elevation gain.
	ElevationProfiles   ElevationProfiles     `json:"elevation_profiles"`              //
	Map                 *PolylineMap          `json:"map,omitempty"`                   // An instance of PolylineMap.
	EffortCount         int                   `json:"effort_count"`                    // The total number of efforts for this segment
	EffortDescription   string                `json:"effort_description"`              //
	AthleteCount        int                   `json:"athlete_count"`                   // The number of unique athletes who have an effort for this segment
	Hazardous           bool                  `json:"hazardous"`                       // Whether this segment is considered hazardous
	StarCount           int                   `json:"star_count"`                      // The number of stars for this segment
	AthleteSegmentStats *SegmentSummaryEffort `json:"athlete_segment_stats,omitempty"` // An instance ofSegmentSummaryEffort.
	Xoms                *Xoms                 `json:"xoms,omitempty"`                  //
	LocalLegend         *LocalLegend          `json:"local_legend,omitempty"`          //
}

type LocalLegend struct {
	AthleteID         int    `json:"athlete_id"`
	Title             string `json:"title"`
	Profile           string `json:"profile"`
	EffortDescription string `json:"effort_description"`
	EffortCount       string `json:"effort_count"`
	Destination       string `json:"destination"`
}

type Xoms struct {
	Kom         string `json:"kom"`
	Qom         string `json:"qom"`
	Overall     string `json:"overall"`
	Destination struct {
		Href string `json:"href"`
		Type string `json:"type"`
		Name string `json:"name"`
	}
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

type ElevationProfiles struct {
	DarkUrl  string `json:"dark_url"`
	LightUrl string `json:"light_url"`
}

// A union type representing the time spent in a given zone.
type TimedZoneRange struct {
	Min  int `json:"min"`  // The minimum value in the range.
	Max  int `json:"max"`  // The maximum value in the range.
	Time int `json:"time"` // The number of seconds spent in this zone
}

type Zones struct {
	HearRate HeartRateZoneRanges `json:"heart_rate"` // An instance of HeartRateZoneRanges.
	Power    PowerZoneRanges     `json:"power"`      // An instance of PowerZoneRanges.
}

type ZoneRange struct {
	Max int `json:"max"` // The maximum value in the range.
	Min int `json:"min"` // The minimum value in the range.
}

type Waypoint struct {
	LatLng            LatLng   `json:"latlng"`              // The location along the route that the waypoint is closest to
	TargetLatLng      LatLng   `json:"target_latlng"`       // A location off of the route that the waypoint is (optional)
	Categories        []string `json:"categories"`          // Categories that the waypoint belongs to
	Title             string   `json:"string"`              // A title for the waypoint
	Description       string   `json:"description"`         // A description of the waypoint (optional)
	DistanceIntoRoute int      `json:"distance_into_route"` // The number meters along the route that the waypoint is located
}
