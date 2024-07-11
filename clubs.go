package gostrava

type ClubService service

type MetaClub struct {
	ID            int    `json:"id,omitempty"`             // The club's unique identifier.
	Name          string `json:"name,omitempty"`           // The club's name.
	ResourceState uint8  `json:"resource_state,omitempty"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
}

type SummaryClub struct {
	MetaClub
	Admin              bool           `json:"admin,omitempty"`                // Whether the currently logged-in athlete is an administrator of this club.
	ActivityTypes      []ActivityType `json:"activity_types,omitempty"`       // The activity types that count for a club. This takes precedence over sport_type.
	ActivityTypesIcon  string         `json:"activity_types_icon,omitempty"`  //
	City               string         `json:"city,omitempty"`                 // The club's city.
	Country            string         `json:"country,omitempty"`              // The club's country.
	CoverPhoto         string         `json:"cover_photo,omitempty"`          // URL to a ~1185x580 pixel cover photo.
	CoverPhotoSmall    string         `json:"cover_photo_small,omitempty"`    // URL to a ~360x176 pixel cover photo.
	Dimensions         []string       `json:"dimensions,omitempty"`           //
	Featured           bool           `json:"featured,omitempty"`             // Whether the club is featured or not.
	LocalizedSportType string         `json:"localized_sport_type,omitempty"` //
	Membership         string         `json:"membership,omitempty"`           // The membership status of the logged-in athlete. May take one of the following values: member, pending
	MemberCount        int            `json:"member_count,omitempty"`         // The club's member count.
	Private            bool           `json:"private,omitempty"`              // Whether the club is private.
	Profile            string         `json:"profile,omitempty"`              //
	ProfileMedium      string         `json:"profile_medium,omitempty"`       // URL to a 60x60 pixel profile picture.
	SportType          ClubSportType  `json:"sport_type,omitempty"`           // Deprecated. Prefer to use activity_types. May take one of the following values: ClubSportTypes.Cycling, ClubSportTypes.Running, ClubSportTypes.Triathlon,  ClubSportTypes.Other
	State              string         `json:"state,omitempty"`                // The club's state or geographical region.
	URL                string         `json:"url,omitempty"`                  // The club's vanity URL.
	Verified           bool           `json:"verified,omitempty"`             // Whether the club is verified or not.
}

type DetailedClub struct {
	SummaryClub
	ClubType       string `json:"club_type,omitempty"`
	Description    string `json:"description,omitempty"`     // The club's description
	FollowingCount int    `json:"following_count,omitempty"` // The number of athletes in the club that the logged-in athlete follows.
	Owner          bool   `json:"owner,omitempty"`           // Whether the currently logged-in athlete is the owner of this club.
	Website        string `json:"website,omitempty"`
}

type ClubActivity struct {
	Athlete            SummaryAthlete `json:"athlete,omitempty"`              // An instance of MetaAthlete.
	Distance           float32        `json:"distance,omitempty"`             // The activity's distance, in meters
	ElapsedTime        int            `json:"elapsed_time,omitempty"`         // The activity's elapsed time, in seconds
	MovingTime         int            `json:"moving_time,omitempty"`          // The activity's moving time, in seconds
	Name               string         `json:"name,omitempty"`                 // The name of the activity
	ResourceState      uint8          `json:"resource_state,omitempty"`       // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
	SportType          SportType      `json:"sport_type,omitempty"`           // An instance of SportType.
	Type               ActivityType   `json:"activity_type,omitempty"`        // Deprecated. Prefer to use sport_type
	TotalElevationGain float32        `json:"total_elevation_gain,omitempty"` // The activity's total elevation gain.
}

type ClubAthlete struct {
	Admin         *bool   `json:"admin,omitempty"`          // Whether the athlete is a club admin.
	FirstName     *string `json:"firstname,omitempty"`      // The athlete's first name.
	LastName      *string `json:"lastname,omitempty"`       // The athlete's last initial.
	Membership    *string `json:"membership,omitempty"`     // The membership status of the logged-in athlete. May take one of the following values: member, pending
	Owner         *bool   `json:"owner,omitempty"`          // Whether the athlete is club owner.
	ResourceState *uint8  `json:"resource_state,omitempty"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
}
