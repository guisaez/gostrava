package gostrava

type AthleteMeta struct {
	ID int `json:"id"`
}

type AthleteSummary struct {
	AthleteMeta
	Username      string    `json:"username"`
	ResourceState int8      `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
	FirstName     string    `json:"firstname"`      // The athlete's first name.
	LastName      string    `json:"lastname"`       // The athlete's last name.
	Bio           string    `json:"bio"`            // The athlete's bio.
	City          string    `json:"city"`           // The athlete's city.
	State         string    `json:"state"`          // The athlete's state or geographical region.
	Country       string    `json:"country"`        // The athlete's country.
	Sex           string    `json:"sex"`            // The athlete's sex. May take one of the following values: M, F, or empty
	Premium       bool      `json:"premium"`        // Deprecated. Use summit field instead. Whether the athlete has any Summit subscription.
	Summit        bool      `json:"summit"`         // Whether the athlete has any Summit subscription.
	CreatedAt     TimeStamp `json:"created_at"`     // The time at which the athlete was created.
	UpdatedAt     TimeStamp `json:"updated_at"`     // The time at which the athlete was last updated.
	BadgeTypeId   int8      `json:"badge_type_id"`
	ProfileMedium string    `json:"profile_medium"` // URL to a 62x62 pixel profile picture.
	Weight        float64   `json:"weight"`         // The athlete's weight in kilograms
	Profile       string    `json:"profile"`        // URL to a 124x124 pixel profile picture.
	Friend        string    `json:"friend"`         // ‘pending’, ‘accepted’, ‘blocked’ or ‘’, the authenticated athlete’s following status of this athlete
	Follower      string    `json:"follower"`       // this athlete’s following status of the authenticated athlete
}

type AthleteDetailed struct {
	AthleteSummary
	Blocked               bool           `json:"blocked"`
	CanFollow             bool           `json:"can_follow"`
	FollowerCount         int            `json:"follower_count"`         // The athlete's follower count.
	FriendCount           int            `json:"friend_count"`           // The athlete's friend count.
	MutualFriendCount     int            `json:"mutual_friend_count"`    // Number of mutual friends between the authenticated athlete and this athlete
	AthleteType           int8           `json:"athlete_type"`           //
	DatePreference        string         `json:"date_preference"`        // Athlete's date preference
	MeasurementPreference string         `json:"measurement_preference"` // The athlete's preferred unit system. May take one of the following values: feet, meters
	Clubs                 []*ClubSummary `json:"clubs"`                  // The athlete's clubs.
	PostableClubsCount    int            `json:"postable_clubs_count"`   //
	FTP                   int            `json:"ftp"`                    // The athlete's FTP (Functional Threshold Power).
	Bikes                 []*GearSummary `json:"bikes"`                  // The athlete's bikes.
	Shoes                 []*GearSummary `json:"shoes"`                  // The athlete's shoes.
	IsWinBackViaUpload    bool           `json:"is_winback_via_upload"`
	IsWinBackViaView      bool           `json:"is_winback_via_view"`
}