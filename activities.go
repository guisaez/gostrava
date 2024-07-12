package gostrava

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ActivityService service

const activities string = "activities"

type NewActivity struct {
	Name           string       // The name of the activity.
	Type           ActivityType // Type of activity. For example - Run, Ride etc.
	SportType      SportType    // Sport type of activity. For example - Run, MountainBikeRide, Ride, etc.
	StartDateLocal time.Time    // ISO 8601 formatted date time.
	ElapsedTime    int          // In seconds.
	Description    string       // Description of the activity.
	Distance       int          // In meters.
	Trainer        int8         // Set to 1 to mark as a trainer activity.
	Commute        int8         // Set to 1 to mark as commute.
}

// Creates a manual activity for an athlete, requires activity:write scope.
func (s *ActivityService) New(accessToken string, body NewActivity) (*ActivityDetailed, error) {
	formData := url.Values{}
	formData.Set("name", body.Name)
	formData.Set("type", string(body.Type))
	formData.Set("sport_type", string(body.SportType))
	formData.Set("start_date_local", body.StartDateLocal.Format(time.RFC3339)) // Assuming RFC3339 format
	formData.Set("elapsed_time", strconv.Itoa(body.ElapsedTime))
	formData.Set("description", body.Description)
	formData.Set("distance", strconv.Itoa(body.Distance))
	formData.Set("trainer", strconv.Itoa(int(body.Trainer)))
	formData.Set("commute", strconv.Itoa(int(body.Commute)))

	req, err := s.client.newRequest(requestOpts{
		Path:        activities,
		Method:      http.MethodPost,
		AccessToken: accessToken,
		Body:        formData,
	})
	if err != nil {
		return nil, err
	}

	resp := new(ActivityDetailed)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the given activity that is owned by the authenticated athlete.
// Requires activity:read for Everyone and Followers activities.
// Requires activity:read_all for Only Me activities.
func (s *ActivityService) GetByID(accessToken string, id int, includeEfforts bool) (*ActivityDetailed, error) {
	params := url.Values{}
	params.Add("include_all_efforts", fmt.Sprintf("%v", includeEfforts))

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d", activities, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := &ActivityDetailed{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Summit Feature. Returns the zones of a given activity.
// Requires activity:read for Everyone and Followers activities.
// Requires activity:read_all for Only Me activities.
func (s *ActivityService) ListActivityZones(accessToken string, id int) ([]ActivityZone, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/zones", activities, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []ActivityZone{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type CommentsReqParams struct {
	RequestParams
	PageSize    int    // Number of items per page. Defaults to 30
	AfterCursor string // Cursor of the las item in the previous page of results, used to request the subsequent page of results. When omitted, the first page of results is fetched.
}

// Returns the comments on the given activity. Requires activity:read for Everyone and Followers activities. Requires activity:read_all for Only Me activities.
func (s *ActivityService) ListActivityComments(accessToken string, id int, p *CommentsReqParams) ([]Comment, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
		if p.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(p.PageSize))
		}
		if p.AfterCursor != "" {
			params.Set("after_cursor", p.AfterCursor)
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/comments", activities, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := []Comment{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the athletes who kudoed an activity identified by an identifier. Requires activity:read for Everyone and Followers activities.
// Requires activity:read_all for OnlyMe Activities
func (s *ActivityService) ListActivityKudoers(accessToken string, id int, p *RequestParams) ([]AthleteSummary, error) {
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
		Path:        fmt.Sprintf("%s/%d/kudos", activities, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
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

// Returns the laps of an activity identified by an identifier. Requires activity:read for Everyone and
// Follower activities. Required activity:read_all for OnlyMeActivities.
func (s *ActivityService) ListActivityLaps(accessToken string, id int) ([]Lap, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/laps", activities, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []Lap{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type GetActivityParams struct {
	RequestParams
	Before int // An epoch timestamp to use for filtering activities that have taken place before that certain time.
	After  int // An epoch timestamp to use for filtering activities that have taken place after a certain time.
}

// Returns the activities of an athlete for a specific identifier. Requires activity:read, OnlyMe activities will be filtered out unless
// requested by a token with activity_read:all.
func (s *ActivityService) ListAthleteActivities(accessToken string, p *GetActivityParams) ([]ActivitySummary, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page_size", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.Page))
		}
		if p.Before > 0 {
			params.Set("before", strconv.Itoa(p.Before))
		}
		if p.After > 0 {
			params.Set("after", strconv.Itoa(p.After))
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%s", athlete, activities),
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	var resp []ActivitySummary
	err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ActivityService) GetActivityZones(accessToken string, id int) ([]ActivityZone, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/zones", activities, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}
	resp := []ActivityZone{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type UpdatedActivity struct {
	Commute      bool         `json:"commute,omitempty"`        // Whether this activity is a commute
	Trainer      bool         `json:"trainer,omitempty"`        // Whether this activity was recorded on a training machine
	HideFromHome bool         `json:"hide_from_home,omitempty"` // Whether this activity is muted
	Description  string       `json:"description,omitempty"`    // The description of the activity
	Name         string       `json:"name,omitempty"`           // The name of the activity
	Type         ActivityType `json:"type,omitempty"`           // Deprecated. Prefer to use sport_type. In a request where both type and sport_type are present, this field will be ignored
	SportType    SportType    `json:"sport_type,omitempty"`     // An instance of SportType.
	GearID       string       `json:"gear_id,omitempty"`        // Identifier for the gear associated with the activity. ‘none’ clears gear from activity
}

// Updates the given activity that is owned by the authenticated athlete. Requires activity:write. Also requires activity:read_all in order
// to update only me activities.
func (s *ActivityService) Update(accessToken string, id int, body UpdatedActivity) (*ActivityDetailed, error) {
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d", activities, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
		Body:        io.NopCloser(bytes.NewReader(json)),
	})
	if err != nil {
		return nil, err
	}

	resp := new(ActivityDetailed)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
