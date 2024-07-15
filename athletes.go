package gostrava

import (
	"net/url"
	"strconv"
)

type AthleteDetailed struct {
	AthleteSummary
	AthleteType           int8           `json:"athlete_type"`
	Blocked               bool           `json:"blocked"`
	CanFollow             bool           `json:"can_follow"`
	DatePreference        string         `json:"date_preference"` // Athlete's date preference
	FollowerCount         int            `json:"follower_count"`  // The athlete's follower count.
	FriendCount           int            `json:"friend_count"`    // The athlete's friend count.
	Friend                string         `json:"friend"`          // ‘pending’, ‘accepted’, ‘blocked’ or ‘null’, the authenticated athlete’s following status of this athlete
	Follower              string         `json:"follower"`        // this athlete’s following status of the authenticated athlete
	IsWinBackViaUpload    bool           `json:"is_winback_via_upload"`
	IsWinBackViaView      bool           `json:"is_winback_via_view"`
	MeasurementPreference string         `json:"measurement_preference"` // The athlete's preferred unit system. May take one of the following values: feet, meters
	MutualFriendCount     int            `json:"mutual_friend_count"`    // Number of mutual friends between the authenticated athlete and this athlete
	PostableClubsCount    int            `json:"postable_clubs_count"`
	FTP                   int            `json:"ftp"`   // The athlete's FTP (Functional Threshold Power).
	Clubs                 []*ClubSummary `json:"clubs"` // The athlete's clubs.
	Bikes                 []*GearSummary `json:"bikes"` // The athlete's bikes.
	Shoes                 []*GearSummary `json:"shoes"` // The athlete's shoes.
}

type AthleteSummary struct {
	AthleteMeta
	BadgeTypeId   int8      `json:"badge_type_id"`
	Bio           string    `json:"bio"`            // The athlete's bio.
	City          string    `json:"city"`           // The athlete's city.
	Country       string    `json:"country"`        // The athlete's country.
	CreatedAt     TimeStamp `json:"created_at"`     // The time at which the athlete was created.
	FirstName     string    `json:"firstname"`      // The athlete's first name.
	LastName      string    `json:"lastname"`       // The athlete's last name.
	Premium       bool      `json:"premium"`        // Deprecated. Use summit field instead. Whether the athlete has any Summit subscription.
	Profile       string    `json:"profile"`        // URL to a 124x124 pixel profile picture.
	ProfileMedium string    `json:"profile_medium"` // URL to a 62x62 pixel profile picture.
	Sex           string    `json:"sex"`            // The athlete's sex. May take one of the following values: M, F, or empty
	ResourceState int8      `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
	State         string    `json:"state"`          // The athlete's state or geographical region.
	Summit        bool      `json:"summit"`         // Whether the athlete has any Summit subscription.
	UpdatedAt     TimeStamp `json:"updated_at"`     // The time at which the athlete was last updated.
	Weight        float64   `json:"weight"`         // The athlete's weight in kilograms
}

type AthleteMeta struct {
	ID int `json:"id"`
}

// A set of rolled-up statistics and totals for an athlete
type AthleteStats struct {
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
	Count         int     `json:"count"`          // The number of activities considered in this total.
	Distance      float32 `json:"distance"`       // The total distance covered by the considered activities.
	MovingTime    int     `json:"moving_time"`    // The total moving time of the considered activities.
	ElapsedTime   int     `json:"elapsed_time"`   // The total elapsed time of the considered activities.
	ElevationGain float32 `json:"elevation_gain"` // The total elevation gain of the considered activities.

	// only correct for recent totals, not ytd or all
	AchievementCount int `json:"achievement_count"` // The total number of achievements of the considered activities.
}

type AthleteService service

// Returns the activity stats of an athlete. Only includes data from activities set to Everyone's visibility.
func (s *AthleteService) GetAthleteStats(accessToken string, id int) (*AthleteStats, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "athletes/" + strconv.Itoa(id) + "/stats",
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(AthleteStats)
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a list of the routes created by the authenticated athlete. Private routes are filtered out
// unless request by a token with read_all scope.
func (s *RouteService) ListRoutes(accessToken string, id int, opts RequestParams) ([]Route, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "athletes/" + strconv.Itoa(id) + "/routes",
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := []Route{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
