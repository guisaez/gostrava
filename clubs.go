package gostrava

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// *************** Types ********************

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

// *************** Methods ********************

type ClubService service

const clubs = "/v3/api/clubs"

// GetById retrieves a club by its id.
//
// GET https://www.strava/com/api/v3/clubs/{id}
func (s *ClubService) GetById(ctx context.Context, accessToken string, id int) (*ClubDetailed, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d", clubs, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	club := new(ClubDetailed)
	resp, err := s.client.DoAndParse(ctx, req, club)
	if err != nil {
		return nil, resp, err
	}

	return club, resp, nil
}

// ListClubActivities retrieves the list of activities for a given club.
//
// Even though an administrator is represented as a AthleteSummary as it is mentioned
// in the docs, the API only returns the first_name and last_name of the athlete.
//
// GET https://www.strava.com/api/v3/clubs/{id}/activities
func (s *ClubService) ListClubActivities(ctx context.Context, accessToken string, id int, options *ListOptions) ([]ActivitySummary, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/activities", clubs, id)

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

	return activities, resp, nil
}

// ListClubAdministrators retrieves the list of administrators for a given club.
//
// Even though an administrator is represented as a AthleteSummary as it is mentioned
// in the docs, the API only returns the first_name and last_name of the athlete.
//
// GET https://www.strava.com/api/v3/clubs/{id}/admins
func (s *ClubService) ListClubAdministrators(ctx context.Context, accessToken string, id int, options *ListOptions) ([]AthleteSummary, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/admins", clubs, id)

	q, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodGet, urlStr, q, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var admins []AthleteSummary
	resp, err := s.client.DoAndParse(ctx, req, &admins)
	if err != nil {
		return nil, resp, err
	}

	return admins, resp, nil
}

// ListClubMembers retrives the list of administrators for a given club.
//
// Even though a member is represented as a AthleteSummary as it is mentioned
// in the docs, the API only returns the first_name and last_name of the athlete.
//
// GET https://www.strava.com/api/v3/clubs/{id}/members
func (s *ClubService) ListClubMembers(ctx context.Context, accessToken string, id int, options *ListOptions) ([]AthleteSummary, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/admins", clubs, id)

	q, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodGet, urlStr, q, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var members []AthleteSummary
	resp, err := s.client.DoAndParse(ctx, req, &members)
	if err != nil {
		return nil, resp, err
	}

	return members, resp, nil
}
