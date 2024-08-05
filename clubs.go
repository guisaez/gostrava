package gostrava

type ClubMeta struct {
	ID            int    `json:"id"`             // The club's unique identifier.
	Name          string `json:"name"`           // The club's name.
	ResourceState int8   `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
}

type ClubSummary struct {
	ClubMeta
	ProfileMedium      string         `json:"profile_medium"`      // URL to a 60x60 pixel profile picture.
	Profile            string         `json:"profile"`             // URL to a 124x124 pixel profile picture.
	CoverPhoto         string         `json:"cover_photo"`         // URL to a ~1185x580 pixel cover photo.
	CoverPhotoSmall    string         `json:"cover_photo_small"`   // URL to a ~360x176 pixel cover photo.
	ActivityTypes      []ActivityType `json:"activity_types"`      // The activity types that count for a club. This takes precedence over sport_type.
	ActivityTypesIcon  string         `json:"activity_types_icon"` //
	Dimensions         []string       `json:"dimensions"`
	SportType          string         `json:"sport_type"`           // Deprecated. Prefer to use activity_types. May take one of the following values: "casual_club", "racing_team", "shop", "other"
	LocalizedSportType string         `json:"localized_sport_type"` //
	City               string         `json:"city"`                 // The club's city.
	State              string         `json:"state"`                // The club's state or geographical region.
	Country            string         `json:"country"`              // The club's country.
	Private            bool           `json:"private"`              // Whether the club is private.
	MemberCount        int            `json:"member_count"`         // The club's member count.
	Featured           bool           `json:"featured,omitempty"`   // Whether the club is featured or not.
	Verified           bool           `json:"verified"`             // Whether the club is verified or not.
	URL                string         `json:"url"`                  // The club's vanity URL.
}

type ClubDetailed struct {
	ClubSummary
	Membership     string `json:"membership"` // The membership status of the logged-in athlete. May take one of the following values: member, pending
	Admin          bool   `json:"admin"`      // Whether the currently logged-in athlete is an administrator of this club.
	Owner          bool   `json:"owner"`      // Whether the currently logged-in athlete is the owner of this club.
	Description    string `json:"description"`
	Type           string `json:"club_type"`
	FollowingCount int    `json:"following_count"` // The number of athletes in the club that the logged-in athlete follows.
	Website        string `json:"website"`
}
