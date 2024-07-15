package gostrava

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

type RequestParams struct {
	Page    int // Page number. Defaults to 1
	PerPage int // Number of items per page. Defaults to 30
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

type Urls struct {
	Url       string `json:"url"`
	RetinaUrl string `json:"retirna_url"`
	DarkUrl   string `json:"dark_url"`
	LightUrl  string `json:"light_url"`
}


