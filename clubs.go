package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ClubService service

const clubs string = "clubs"

type ClubMeta struct {
	ID            *int    `json:"id,omitempty"`             // The club's unique identifier.
	Name          *string `json:"name,omitempty"`           // The club's name.
	ResourceState *uint8  `json:"resource_state,omitempty"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
}

type ClubSummary struct {
	ClubMeta
	Admin              *bool          `json:"admin,omitempty"`                // Whether the currently logged-in athlete is an administrator of this club.
	ActivityTypes      []ActivityType `json:"activity_types,omitempty"`       // The activity types that count for a club. This takes precedence over sport_type.
	ActivityTypesIcon  *string        `json:"activity_types_icon,omitempty"`  //
	City               *string        `json:"city,omitempty"`                 // The club's city.
	Country            *string        `json:"country,omitempty"`              // The club's country.
	CoverPhoto         *string        `json:"cover_photo,omitempty"`          // URL to a ~1185x580 pixel cover photo.
	CoverPhotoSmall    *string        `json:"cover_photo_small,omitempty"`    // URL to a ~360x176 pixel cover photo.
	Dimensions         *[]string      `json:"dimensions,omitempty"`           //
	Featured           *bool          `json:"featured,omitempty"`             // Whether the club is featured or not.
	LocalizedSportType *string        `json:"localized_sport_type,omitempty"` //
	Membership         *string        `json:"membership,omitempty"`           // The membership status of the logged-in athlete. May take one of the following values: member, pending
	MemberCount        *int           `json:"member_count,omitempty"`         // The club's member count.
	Private            *bool          `json:"private,omitempty"`              // Whether the club is private.
	Profile            *string        `json:"profile,omitempty"`              //
	ProfileMedium      *string        `json:"profile_medium,omitempty"`       // URL to a 60x60 pixel profile picture.
	SportType          *ClubSportType `json:"sport_type,omitempty"`           // Deprecated. Prefer to use activity_types. May take one of the following values: ClubSportTypes.Cycling, ClubSportTypes.Running, ClubSportTypes.Triathlon,  ClubSportTypes.Other
	State              *string        `json:"state,omitempty"`                // The club's state or geographical region.
	URL                *string        `json:"url,omitempty"`                  // The club's vanity URL.
	Verified           *bool          `json:"verified,omitempty"`             // Whether the club is verified or not.
}

type ClubDetailed struct {
	ClubSummary
	ClubType       *string `json:"club_type,omitempty"`
	Description    *string `json:"description,omitempty"`     // The club's description
	FollowingCount *int    `json:"following_count,omitempty"` // The number of athletes in the club that the logged-in athlete follows.
	Owner          *bool   `json:"owner,omitempty"`           // Whether the currently logged-in athlete is the owner of this club.
	Website        *string `json:"website,omitempty"`
}

type ClubActivity struct {
	Athlete            *AthleteSummary `json:"athlete,omitempty"`              // An instance of MetaAthlete.
	Distance           *float32        `json:"distance,omitempty"`             // The activity's distance, in meters
	ElapsedTime        *int            `json:"elapsed_time,omitempty"`         // The activity's elapsed time, in seconds
	MovingTime         *int            `json:"moving_time,omitempty"`          // The activity's moving time, in seconds
	Name               *string         `json:"name,omitempty"`                 // The name of the activity
	ResourceState      *int8           `json:"resource_state,omitempty"`       // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
	SportType          *SportType      `json:"sport_type,omitempty"`           // An instance of SportType.
	Type               *ActivityType   `json:"activity_type,omitempty"`        // Deprecated. Prefer to use sport_type
	TotalElevationGain *float32        `json:"total_elevation_gain,omitempty"` // The activity's total elevation gain.
}

type ClubMember struct {
	Admin         *bool   `json:"admin,omitempty"`          // Whether the athlete is a club admin.
	FirstName     *string `json:"firstname,omitempty"`      // The athlete's first name.
	LastName      *string `json:"lastname,omitempty"`       // The athlete's last initial.
	Membership    *string `json:"membership,omitempty"`     // The membership status of the logged-in athlete. May take one of the following values: member, pending
	Owner         *bool   `json:"owner,omitempty"`          // Whether the athlete is club owner.
	ResourceState *uint8  `json:"resource_state,omitempty"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
}

// Returns a given club using its identifier
func (s *ClubService) GetById(accessToken string, id int) (*ClubDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d", clubs, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(ClubDetailed)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a list of the administrators of a given club.
// The resource state of each admin will be Summary Athlete but Strava tends to send only:
//
//	[{
//		"firstname": "John",
//		"lastname": "Doe",
//		"resource_state": 2
//	}]
func (s *ClubService) ListClubAdministrators(accessToken string, id int, p *RequestParams) ([]AthleteSummary, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/admins", clubs, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []AthleteSummary{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Retrieve recent activities from members of a specific club. The authenticated athlete must belong to the request club in order to hit this endpoint, Pagination is supported. Athlete profile
// visibility is respected for all activities.
// The reference to the Athlete for each activity will only include the following
//
//	[{
//		"athlete" : {
//			"firstname": "John",
//			"lastname": "Doe",
//			"resource_state": 2
//		},
//		...
//	}]
func (s *ClubService) ListClubActivities(accessToken string, id int, p *RequestParams) ([]ClubActivity, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/activities", clubs, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []ClubActivity{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns of list of the athletes who are members of a given club.
func (s *ClubService) ListClubMembers(accessToken string, id int, p *RequestParams) ([]ClubMember, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/members", clubs, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []ClubMember{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Return a list of the clubs whose membership includes the authenticated athlete.
func (s *ClubService) ListAthleteClubs(accessToken string, p *RequestParams) ([]ClubSummary, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%s", athlete, clubs),
		AccessToken: accessToken,
		Method:      http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	resp := []ClubSummary{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
