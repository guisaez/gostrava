package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
)

type AthleteService service

const (
	athlete  string = "athlete"
	athletes string = "athletes"
)

type AthleteMeta struct {
	ID            *int  `json:"id"`             // The unique identifier of the athlete
	ResourceState *int8 `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
}

func (a *AthleteMeta) String() string {
	return Stringify(a)
}

type AthleteSummary struct {
	AthleteMeta
	BadgeTypeId   uint8      `json:"badge_type_id,omitempty"`
	Bio           *string    `json:"bio,omitempty"`            // The athlete's bio.
	City          *string    `json:"city,omitempty"`           // The athlete's city.
	Country       *string    `json:"country,omitempty"`        // The athlete's country.
	CreatedAt     *TimeStamp `json:"created_at,omitempty"`     // The time at which the athlete was created.
	FirstName     *string    `json:"firstname,omitempty"`      // The athlete's first name.
	LastName      *string    `json:"lastname,omitempty"`       // The athlete's last name.
	Premium       *bool      `json:"premium,omitempty"`        // Deprecated. Use summit field instead. Whether the athlete has any Summit subscription.
	Profile       *string    `json:"profile,omitempty"`        // URL to a 124x124 pixel profile picture.
	ProfileMedium *string    `json:"profile_medium,omitempty"` // URL to a 62x62 pixel profile picture.
	Sex           *string    `json:"sex,omitempty"`            // The athlete's sex. May take one of the following values: M, F
	State         *string    `json:"state,omitempty"`          // The athlete's state or geographical region.
	Summit        *bool      `json:"summit,omitempty"`         // Whether the athlete has any Summit subscription.
	UpdatedAt     *TimeStamp `json:"updated_at,omitempty"`     // The time at which the athlete was last updated.
	Weight        *float64   `json:"weight,omitempty"`         // The athlete's weight.
}

func (a *AthleteSummary) String() string {
	return Stringify(a)
}

type AthleteDetailed struct {
	AthleteSummary
	AthleteType           *int8         `json:"athlete_type,omitempty"`
	Blocked               *bool         `json:"blocked,omitempty"`
	CanFollow             *bool         `json:"can_follow,omitempty"`
	DatePreference        *string       `json:"date_preference,omitempty"`
	FollowerCount         *int          `json:"follower_count,omitempty"` // The athlete's follower count.
	FriendCount           *int          `json:"friend_count,omitempty"`   // The athlete's friend count.
	IsWinBackViaUpload    *bool         `json:"is_winback_via_upload,omitempty"`
	IsWinBackViaView      *bool         `json:"is_winback_via_view,omitempty"`
	MeasurementPreference *string       `json:"measurement_preference,omitempty"` // The athlete's preferred unit system. May take one of the following values: feet, meters
	MutualFriendCount     *int          `json:"mutual_friend_count,omitempty"`
	PostableClubsCount    *int          `json:"postable_clubs_count,omitempty"`
	FTP                   *int          `json:"ftp,omitempty"`   // The athlete's FTP (Functional Threshold Power).
	Clubs                 []ClubSummary `json:"clubs,omitempty"` // The athlete's clubs.
	// Bikes                 []SummaryGear `json:"bikes,omitempty"` // The athlete's bikes.
	// Shoes                 []SummaryGear `json:"shoes,omitempty"` // The athlete's shoes.
}

func (a *AthleteDetailed) String() string {
	return Stringify(a)
}

// Returns the currently authenticated athlete. Tokens with profile:read_all scope will receive
// a detailed athlete representation; all others will receive a SummaryAthlete representation
func (s *AthleteService) GetAuthenticatedAthlete(accessToken string) (*AthleteDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        athlete,
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(AthleteDetailed)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, err
}

type Zones struct {
	HearRate *HeartRateZoneRanges `json:"heart_rate,omitempty"` // An instance of HeartRateZoneRanges.
	Power    *PowerZoneRanges     `json:"power,omitempty"`      // An instance of PowerZoneRanges.
}

// Returns the current athlete's heart rate and power zones. Requires profile:read_all.
func (s *AthleteService) GetZones(accessToken string) (*Zones, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/zones", athlete),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(Zones)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// A set of rolled-up statistics and totals for an athlete
type ActivityStats struct {
	AllRideTotals        *ActivityTotal `json:"all_ride_totals"`    // The all time ride stats for the athlete.
	AllRunTotals         *ActivityTotal `json:"all_run_totals"`     // The all time run stats for the athlete.
	AllSwimTotals        *ActivityTotal `json:"all_swim_totals"`    // The all time swim stats for the athlete.
	RecentRideTotals     *ActivityTotal `json:"recent_ride_totals"` // The recent (last 4 weeks) ride stats for the athlete.
	RecentRunTotals      *ActivityTotal `json:"recent_run_totals"`  // The recent (last 4 weeks) run stats for the athlete.
	RecentSwimTotals     *ActivityTotal `json:"recent_swim_totals"` // The recent (last 4 weeks) swim stats for the athlete.
	YearToDateRideTotals *ActivityTotal `json:"ytd_ride_totals"`    // The year to date ride stats for the athlete.
	YearToDateRunTotals  *ActivityTotal `json:"ytd_run_totals"`     // The year to date run stats for the athlete.
	YearToDateSwimTotals *ActivityTotal `json:"ytd_swim_totals"`    // The year to date swim stats for the athlete.
}

func (as *ActivityStats) String() string {
	return Stringify(as)
}

// Returns the activity stats of an athlete. Only includes data from activities set to Everyone's visibility.
func (s *AthleteService) GetAthleteStats(accessToken string, id int) (*ActivityStats, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/stats", athletes, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(ActivityStats)
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type UpdatedAthlete struct {
	Weight float32 // The weigh of the athlete in kilograms.
}

// Updates the authenticated user. Requires profile:write scope
func (s *AthleteService) Update(accessToken string, updatedAthlete UpdatedAthlete) (*AthleteDetailed, error) {
	params := url.Values{}

	params.Set("weight", fmt.Sprintf("%.2f", updatedAthlete.Weight))

	req, err := s.client.newRequest(requestOpts{
		Path:        athlete,
		Method:      http.MethodPut,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := new(AthleteDetailed)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, err
}
