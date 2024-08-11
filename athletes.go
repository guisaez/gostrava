package gostrava

import (
	"context"
	"fmt"
	"net/http"
)

// *************** Types ********************

type AthleteMeta struct {
	ID            *int  `json:"id,omitempty"`
	ResourceState *int8 `json:"resource_state,omitempty"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
}

type AthleteSummary struct {
	AthleteMeta
	Username      *string    `json:"username,omitempty"`
	FirstName     *string    `json:"firstname,omitempty"`  // The athlete's first name.
	LastName      *string    `json:"lastname,omitempty"`   // The athlete's last name.
	Bio           *string    `json:"bio,omitempty"`        // The athlete's bio.
	City          *string    `json:"city,omitempty"`       // The athlete's city.
	State         *string    `json:"state,omitempty"`      // The athlete's state or geographical region.
	Country       *string    `json:"country,omitempty"`    // The athlete's country.
	Sex           *string    `json:"sex,omitempty"`        // The athlete's sex. May take one of the following values: M, F, or empty
	Premium       *bool      `json:"premium,omitempty"`    // Deprecated. Use summit field instead. Whether the athlete has any Summit subscription.
	Summit        *bool      `json:"summit,omitempty"`     // Whether the athlete has any Summit subscription.
	CreatedAt     *TimeStamp `json:"created_at,omitempty"` // The time at which the athlete was created.
	UpdatedAt     *TimeStamp `json:"updated_at,omitempty"` // The time at which the athlete was last updated.
	BadgeTypeId   *int8      `json:"badge_type_id,omitempty"`
	ProfileMedium *string    `json:"profile_medium,omitempty"` // URL to a 62x62 pixel profile picture.
	Weight        *float64   `json:"weight,omitempty"`         // The athlete's weight in kilograms
	Profile       *string    `json:"profile,omitempty"`        // URL to a 124x124 pixel profile picture.
	Friend        *string    `json:"friend,omitempty"`         // ‘pending’, ‘accepted’, ‘blocked’ or ‘’, the authenticated athlete’s following status of this athlete
	Follower      *string    `json:"follower,omitempty"`       // this athlete’s following status of the authenticated athlete
}

type AthleteDetailed struct {
	AthleteSummary
	Blocked               *bool         `json:"blocked,omitempty"`
	CanFollow             *bool         `json:"can_follow,omitempty"`
	FollowerCount         *int          `json:"follower_count,omitempty"`         // The athlete's follower count.
	FriendCount           *int          `json:"friend_count,omitempty"`           // The athlete's friend count.
	MutualFriendCount     *int          `json:"mutual_friend_count,omitempty"`    // Number of mutual friends between the authenticated athlete and this athlete
	AthleteType           *int8         `json:"athlete_type,omitempty"`           //
	DatePreference        *string       `json:"date_preference,omitempty"`        // Athlete's date preference
	MeasurementPreference *string       `json:"measurement_preference,omitempty"` // The athlete's preferred unit system. May take one of the following values: feet, meters
	Clubs                 []ClubSummary `json:"clubs,omitempty"`                  // The athlete's clubs.
	PostableClubsCount    *int          `json:"postable_clubs_count,omitempty"`   //
	FTP                   *int          `json:"ftp,omitempty"`                    // The athlete's FTP (Functional Threshold Power).
	Bikes                 []GearSummary `json:"bikes,omitempty"`                  // The athlete's bikes.
	Shoes                 []GearSummary `json:"shoes,omitempty"`                  // The athlete's shoes.
	IsWinBackViaUpload    *bool         `json:"is_winback_via_upload,omitempty"`
	IsWinBackViaView      *bool         `json:"is_winback_via_view,omitempty"`
}

// A set of rolled-up statistics and totals for an athlete
type AthleteStats struct {
	RecentRideTotals     ActivityTotal `json:"recent_ride_totals"` // The recent (last 4 weeks) ride stats for the athlete.
	AllRideTotals        ActivityTotal `json:"all_ride_totals"`    // The all time ride stats for the athlete.
	RecentRunTotals      ActivityTotal `json:"recent_run_totals"`  // The recent (last 4 weeks) run stats for the athlete.
	AllRunTotals         ActivityTotal `json:"all_run_totals"`     // The all time run stats for the athlete.
	RecentSwimTotals     ActivityTotal `json:"recent_swim_totals"` // The recent (last 4 weeks) swim stats for the athlete.
	AllSwimTotals        ActivityTotal `json:"all_swim_totals"`    // The all time swim stats for the athlete.
	YearToDateRideTotals ActivityTotal `json:"ytd_ride_totals"`    // The year to date ride stats for the athlete.
	YearToDateSwimTotals ActivityTotal `json:"ytd_swim_totals"`    // The year to date swim stats for the athlete.
	YearToDateRunTotals  ActivityTotal `json:"ytd_run_totals"`     // The year to date run stats for the athlete.
}

// A roll-up of metrics pertaining to a set of activities. Values are in seconds and meters.
type ActivityTotal struct {
	Count         int     `json:"count"`          // The number of activities considered in this total.
	Distance      float32 `json:"distance"`       // The total distance covered by the considered activities.
	MovingTime    int     `json:"moving_time"`    // The total moving time of the considered activities.
	ElapsedTime   int     `json:"elapsed_time"`   // The total elapsed time of the considered activities.
	ElevationGain float32 `json:"elevation_gain"` // The total elevation gain of the considered activities.

	// only present for recent totals, not ytd, or all
	AchievementCount int `json:"achievement_count"` // The total number of achievements of the considered activities.
}

// *************** Methods ********************

type AthletesService service

const (
	athlete  = "/api/v3/athlete"
	athletes = "/api/v3/athletes"
)

// GetAthlete retrieves the athlete that corresponds to the provided accessToken.
// Tokes with profile_read:all scope will receive a detailed athlete representation
// all others will receive a SummaryAthlete representation
//
// GET https://www.strava.com/api/v3/athlete
func (s *AthletesService) GetAthlete(ctx context.Context, accessToken string) (*AthleteDetailed, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, athlete, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	athlete := new(AthleteDetailed)

	resp, err := s.client.DoAndParse(ctx, req, athlete)
	if err != nil {
		return nil, resp, err
	}

	return athlete, resp, err
}

// GetAthleteStats retrives the activities stats of an athlete. Only includes data
// from activities set to Everyone's visibility
//
// GET: https://www.strava.com/api/v3/athletes/{id}/stats
func (s *AthletesService) GetAthleteStats(ctx context.Context, accessToken string, id int) (*AthleteStats, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/stats", athletes, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	stats := new(AthleteStats)
	resp, err := s.client.DoAndParse(ctx, req, stats)
	if err != nil {
		return nil, resp, err
	}

	return stats, resp, nil
}

// GetAthleteRoutes retrieves the routes created by a specified athlete.
//
// GET: https://www.strava.com/api/v3/athletes/{id}/routes
func (s *AthletesService) GetAthleteRoutes(ctx context.Context, accessToken string, id int) ([]RouteSummary, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/routes", athletes, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var routes []RouteSummary
	resp, err := s.client.DoAndParse(ctx, req, &routes)
	if err != nil {
		return nil, resp, nil
	}

	return routes, resp, nil
}

// GetAthleteZones retrieves the athlete's heartrate and power zones.
// Required profile_read:all
//
// GET https://www.strava.com/api/v3/athlete/zones
func (s *AthletesService) GetAthleteZones(ctx context.Context, accessToken string) (*Zones, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/zones", athlete)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	zones := new(Zones)

	resp, err := s.client.DoAndParse(ctx, req, zones)
	if err != nil {
		return nil, resp, err
	}

	return zones, resp, nil
}
