package gostrava

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
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

// *************** Methods ********************

type CurrentAthleteService service

const athlete = "/api/v3/athlete"

// GetAthlete retrieves the athlete that corresponds to the provided accessToken.
// Tokes with profile_read:all scope will receive a detailed athlete representation
// all others will recevie a SummaryAthlete representation
//
// GET https://www.strava.com/api/v3/athlete
func (s *CurrentAthleteService) GetAthlete(ctx context.Context, accessToken string) (*AthleteDetailed, *http.Response, error) {

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

// GetAthleteZones retreives the athlete's heartrate and power zones.
// Required profile_read:all
//
// GET https://www.strava.com/api/v3/athlete/zones
func (s *CurrentAthleteService) GetAthleteZones(ctx context.Context, accessToken string) (*Zones, *http.Response, error) {
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

// ListAthleteClubs retrieves a list of all the clubs the current athlete
// is a member of.
//
// GET https://www.strava.com/api/v3/athlete/clubs
func (s *CurrentAthleteService) ListAthleteClubs(ctx context.Context, accessToken string, options *ListOptions) ([]ClubSummary, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/clubs", athlete)

	q, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, urlStr, q, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var clubs []ClubSummary
	resp, err := s.client.DoAndParse(ctx, req, &clubs)
	if err != nil {
		return nil, resp, err
	}

	return clubs, resp, err
}

type ListActivityOptions struct {
	Page    int   `url:"page,omitempty"`     // Defaults to 1
	PerPage int   `url:"per_page,omitempty"` // Defaults to 30
	Before  int64 `url:"before,omitempty"`   // An epoch timestamp to use for filtering activities that have taken place before that certain time.
	After   int64 `url:"after,omitempty"`    // An epoch timestamp to use for filtering activities that have taken place after a certain time.
}

// ListAthleteActivities retrieves a list of the activities recorded by the authenticated athlete
//
// GET https://www.strava.com/api/v3/athlete/activities
func (s *CurrentAthleteService) ListAthleteActivities(ctx context.Context, accessToken string, options *ListActivityOptions) ([]ActivitySummary, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/activities", athlete)

	q, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, urlStr, q, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var activities []ActivitySummary
	resp, err := s.client.DoAndParse(ctx, req, &activities)
	if err != nil {
		return nil, resp, err
	}

	return activities, resp, err
}
