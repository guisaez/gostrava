package gostrava

import (
	"net/url"
	"strconv"
)

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

type ClubMeta struct {
	ID            int    `json:"id"`             // The club's unique identifier.
	Name          string `json:"name"`           // The club's name.
	ResourceState int8   `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
}

// ClubAthlete is not an official Strava type since it originally corresponds to AthleteSummary, but while checking actual requests,
// it only returns this fields along with resource_state = 2.
type ClubAthlete struct {
	FirstName string `json:"firstname"` // Athlete First Name
	LastName  string `json:"lastname"`  // Athlete Last Name
}

type ClubActivity struct {
	Athlete            ClubAthlete  `json:"athlete"`              // An instance of MetaAthlete.
	Name               string       `json:"name"`                 // The name of the activity
	Distance           float32      `json:"distance"`             // The activity's distance, in meters
	MovingTime         int          `json:"moving_time"`          // The activity's moving time, in seconds
	ElapsedTime        int          `json:"elapsed_time"`         // The activity's elapsed time, in seconds
	TotalElevationGain float32      `json:"total_elevation_gain"` // The activity's total elevation gain.
	Type               ActivityType `json:"activity_type"`        // Deprecated. Prefer to use sport_type
	SportType          SportType    `json:"sport_type"`           // An instance of SportType.
}

type Member struct {
	Admin      bool   `json:"admin"`      // Whether the athlete is a club admin.
	FirstName  string `json:"firstname"`  // The athlete's first name.
	LastName   string `json:"lastname"`   // The athlete's last initial.
	Membership string `json:"membership"` // The membership status of the logged-in athlete. May take one of the following values: member, pending
	Owner      bool   `json:"owner"`      // Whether the athlete is club owner.
}

// *****************************************************

type ClubService service

// Returns a given club using its identifier
func (s *ClubService) GetById(accessToken string, id int) (*ClubDetailed, error) {
	req, err := s.client.NewRequest(RequestOpts{
		Path:        "clubs/" + strconv.Itoa(id),
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(ClubDetailed)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a list of the administrators of a given club.
// It uses a predefined ClubAthlete struct and not AthleteSummary because it currently sends only FirstName and LastName
func (s *ClubService) ListAdministrators(accessToken string, id int, opts RequestParams) ([]ClubAthlete, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	req, err := s.client.NewRequest(RequestOpts{
		Path:        "clubs/" + strconv.Itoa(id) + "/admins",
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []ClubAthlete{}
	if err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Retrieve recent activities from members of a specific club. The authenticated athlete must belong to the request club in order to hit this endpoint, Pagination is supported. Athlete profile
// visibility is respected for all activities.
// It uses a predefined ClubAthlete struct and not AthleteSummary because it currently sends only FirstName and LastName
func (s *ClubService) ListActivities(accessToken string, id int, opts RequestParams) ([]ClubActivity, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	req, err := s.client.NewRequest(RequestOpts{
		Path:        "clubs/" + strconv.Itoa(id) + "/activities",
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []ClubActivity{}
	if err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns of list of the athletes who are members of a given club.
func (s *ClubService) ListMembers(accessToken string, id int, opts RequestParams) ([]Member, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	req, err := s.client.NewRequest(RequestOpts{
		Path:        "clubs/" + strconv.Itoa(id) + "/members",
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []Member{}
	if err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
