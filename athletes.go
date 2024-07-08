package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
)

type AthleteAPIService apiService

type MetaAthlete struct {
	ID            int           `json:"id,omitempty"`   // The unique identifier of the athlete
	ResourceState uint8 `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
}

type SummaryAthlete struct {
	MetaAthlete
	BadgeTypeId   *uint8    `json:"badge_type_id,omitempty"`
	Bio           *string   `json:"bio,omitempty"`            // The athlete's bio.
	City          *string   `json:"city,omitempty"`           // The athlete's city.
	Country       *string   `json:"country,omitempty"`        // The athlete's country.
	CreatedAt     *DateTime `json:"created_at,omitempty"`     // The time at which the athlete was created.
	FirstName     *string   `json:"firstname,omitempty"`      // The athlete's first name.
	LastName      *string   `json:"lastname,omitempty"`       // The athlete's last name.
	Premium       *bool     `json:"premium,omitempty"`        // Deprecated. Use summit field instead. Whether the athlete has any Summit subscription.
	Profile       *string   `json:"profile,omitempty"`        // URL to a 124x124 pixel profile picture.
	ProfileMedium *string   `json:"profile_medium,omitempty"` // URL to a 62x62 pixel profile picture.
	Sex           *string   `json:"sex,omitempty"`            // The athlete's sex. May take one of the following values: M, F
	State         *string   `json:"state,omitempty"`          // The athlete's state or geographical region.
	Summit        *bool     `json:"summit,omitempty"`         // Whether the athlete has any Summit subscription.
	UpdatedAt     *DateTime `json:"updated_at,omitempty"`     // The time at which the athlete was last updated.
	Weight        *float64  `json:"weight,omitempty"`         // The athlete's weight.
}

type DetailedAthlete struct {
	SummaryAthlete
	AthleteType           int8          `json:"athlete_type"`
	Blocked               bool          `json:"blocked"`
	CanFollow             bool          `json:"can_follow"`
	DatePreference        string        `json:"date_preference"`
	FollowerCount         int           `json:"follower_count"` // The athlete's follower count.
	FriendCount           int           `json:"friend_count"`   // The athlete's friend count.
	IsWinBackViaUpload    bool          `json:"is_winback_via_upload"`
	IsWinBackViaView      bool          `json:"is_winback_via_view"`
	MeasurementPreference string        `json:"measurement_preference"` // The athlete's preferred unit system. May take one of the following values: feet, meters
	MutualFriendCount     int           `json:"mutual_friend_count"`
	PostableClubsCount    int           `json:"postable_clubs_count"`
	FTP                   *int          `json:"ftp"`   // The athlete's FTP (Functional Threshold Power).
	Clubs                 []SummaryClub `json:"clubs"` // The athlete's clubs.
	Bikes                 []SummaryGear `json:"bikes"` // The athlete's bikes.
	Shoes                 []SummaryGear `json:"shoes"` // The athlete's shoes.
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

// Returns the currently authenticated athlete. Tokens with profile:read_all scope will receive
// a detailed athlete representation; all others will receive a SummaryAthlete representation
func (s *AthleteAPIService) GetAuthenticatedAthlete(access_token string) (*DetailedAthlete, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath)

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	s.client.TestingFileName = "athlete_get_authenticated_athlete_server_response.json"

	resp := &DetailedAthlete{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the authenticated athlete's heart rate and power zones. Requires profile:read_all.
func (s *AthleteAPIService) GetZones(access_token string) (*Zones, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath, "zones")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := &Zones{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the activity stats of an athlete. Only includes data from activities set to Everyone's visibility.
func (s *AthleteAPIService) GetAthleteStats(access_token string, id int) (*ActivityStats, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletesPath, fmt.Sprint(id), "stats")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := &ActivityStats{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type UpdateAthletePayload struct {
	Weight float64 // The weigh of the athlete in kilograms.
}

// Update the currently authenticated athlete. Requires profile:write scope.
func (s *AthleteAPIService) UpdateAthlete(access_token string, p UpdateAthletePayload) (*DetailedAthlete, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath)

	params := url.Values{}

	if p.Weight > 0 {
		params.Set("weight", fmt.Sprintf("%.2f", p.Weight))
	}

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodPut,
		access_token: access_token,
		body:         params,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp := &DetailedAthlete{}
	if err := s.client.do(req, resp); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}
