package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ClubsAPIService apiService

type MetaClub struct {
	ID            int64         `json:"id"`             // The club's unique identifier.
	ResourceState ResourceState `json:"resource_state"` // Resource state, indicates level of detail. Possible values: ResourceStates.Meta, ResourceStates.Summary, ResourceStates.Detail
	Name          string        `json:"name"`           // The club's name.
}

type SummaryClub struct {
	MetaClub
	Admin              bool           `json:"admin"`          // Whether the currently logged-in athlete is an administrator of this club.
	ActivityTypes      []ActivityType `json:"activity_types"` // The activity types that count for a club. This takes precedence over sport_type.
	ActivityTypesIcon  string         `json:"activity_types_icon"`
	City               string         `json:"city"`              // The club's city.
	Country            string         `json:"country"`           // The club's country.
	CoverPhoto         string         `json:"cover_photo"`       // URL to a ~1185x580 pixel cover photo.
	CoverPhotoSmall    string         `json:"cover_photo_small"` // URL to a ~360x176 pixel cover photo.
	Dimensions         []string       `json:"dimensions"`
	Featured           bool           `json:"featured"` // Whether the club is featured or not.
	LocalizedSportType string         `json:"localized_sport_type"`
	Membership         ClubMembership `json:"membership"`   // The membership status of the logged-in athlete. May take one of the following values: member, pending
	MemberCount        int            `json:"member_count"` // The club's member count.
	Private            bool           `json:"private"`      // Whether the club is private.
	Profile            string         `json:"profile"`
	ProfileMedium      string         `json:"profile_medium"` // URL to a 60x60 pixel profile picture.
	SportType          ClubSportType  `json:"sport_type"`     // Deprecated. Prefer to use activity_types. May take one of the following values: ClubSportTypes.Cycling, ClubSportTypes.Running, ClubSportTypes.Triathlon,  ClubSportTypes.Other
	State              string         `json:"state"`          // The club's state or geographical region.
	URL                string         `json:"url"`            // The club's vanity URL.
	Verified           bool           `json:"verified"`       // Whether the club is verified or not.
}

type DetailedClub struct {
	SummaryClub
	ClubType       string `json:"club_type"`
	Description    string `json:"description"`     // The club's description
	FollowingCount int    `json:"following_count"` // The number of athletes in the club that the logged-in athlete follows.
	Owner          bool   `json:"owner"`           // Whether the currently logged-in athlete is the owner of this club.
	Website        string `json:"website"`
}

type ClubActivity struct {
	Athlete            SummaryAthlete `json:"athlete"`              // An instance of MetaAthlete.
	Distance           float32        `json:"distance"`             // The activity's distance, in meters
	ElapsedTime        int            `json:"elapsed_time"`         // The activity's elapsed time, in seconds
	MovingTime         int            `json:"moving_time"`          // The activity's moving time, in seconds
	Name               string         `json:"name"`                 // The name of the activity
	ResourceState      ResourceState  `json:"resource_state"`       // Resource state, indicates level of detail. Possible values: ResourceStates.Meta, ResourceStates.Summary, ResourceStates.Detail
	SportType          SportType      `json:"sport_type"`           // An instance of SportType.
	Type               ActivityType   `json:"activity_type"`        // Deprecated. Prefer to use sport_type
	TotalElevationGain float32        `json:"total_elevation_gain"` // The activity's total elevation gain.
}

type ClubAthlete struct {
	Admin         bool           `json:"admin"`          // Whether the athlete is a club admin.
	FirstName     string         `json:"firstname"`      // The athlete's first name.
	LastName      string         `json:"lastname"`       // The athlete's last initial.
	Membership    ClubMembership `json:"membership"`     // The membership status of the logged-in athlete. May take one of the following values: member, pending
	Owner         bool           `json:"owner"`          // Whether the athlete is club owner.
	ResourceState ResourceState  `json:"resource_state"` // Resource state, indicates level of detail. Possible values: ResourceStates.Meta, ResourceStates.Summary, ResourceStates.Detail
}

// Returns a given club using its identifier
func (s *ClubsAPIService) GetById(access_token string, id int64) (*DetailedClub, error) {
	requestUrl := s.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id))

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := &DetailedClub{}
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
func (s *ClubsAPIService) ListClubAdministrators(access_token string, id int64, p *RequestParams) ([]SummaryAthlete, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	requestUrl := s.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id), "/admins")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		body:         params,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := []SummaryAthlete{}
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
func (s *ClubsAPIService) ListClubActivities(access_token string, id int64, p *RequestParams) ([]ClubActivity, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	requestUrl := s.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id), "/activities")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		body:         params,
		access_token: access_token,
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
func (s *ClubsAPIService) ListClubMembers(access_token string, id int64, p *RequestParams) ([]ClubAthlete, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	requestUrl := s.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id), "/members")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		body:         params,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := []ClubAthlete{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Return a list of the clubs whose membership includes the authenticated athlete.
func (s *ClubsAPIService) ListAthleteClubs(access_token string, p *RequestParams) ([]SummaryClub, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath, clubsPath)

	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		body:         params,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := []SummaryClub{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
