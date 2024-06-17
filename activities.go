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

type ActivitiesAPIService apiService

type Activity struct {
	Name           string       `json:"name"`                  // The name of the activity.
	Type           ActivityType `json:"type,omitempty"`        // Type of activity. For example - Run, Ride etc.
	SportType      SportType    `json:"sport_type"`            // Sport type of activity. For example - Run, MountainBikeRide, Ride, etc.
	StartDateLocal time.Time    `json:"start_date_local"`      // ISO 8601 formatted date time.
	ElapsedTime    int          `json:"elapsed_time"`          // In seconds.
	Description    string       `json:"description,omitempty"` // Description of the activity.
	Distance       int          `json:"distance,omitempty"`    // In meters.
	Trainer        int8         `json:"trainer,omitempty"`     // Set to 1 to mark as a trainer activity.
	Commute        int8         `json:"commute,omitempty"`     // Set to 1 to mark as commute.
}

// Creates a manual activity for an athlete, requires activity:write scope.
func (s *ActivitiesAPIService) New(access_token string, payload Activity) (*DetailedActivity, error) {
	requestUrl := s.client.BaseURL.JoinPath(activitiesPath)

	formData := url.Values{}
	formData.Set("name", payload.Name)
	formData.Set("type", string(payload.Type))
	formData.Set("sport_type", string(payload.SportType))
	formData.Set("start_date_local", payload.StartDateLocal.Format(time.RFC3339)) // Assuming RFC3339 format
	formData.Set("elapsed_time", strconv.Itoa(payload.ElapsedTime))
	formData.Set("description", payload.Description)
	formData.Set("distance", strconv.Itoa(payload.Distance))
	formData.Set("trainer", strconv.Itoa(int(payload.Trainer)))
	formData.Set("commute", strconv.Itoa(int(payload.Commute)))

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodPost,
		access_token: access_token,
		body:         formData,
	})
	if err != nil {
		return nil, err
	}

	resp := &DetailedActivity{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the given activity that is owned by the authenticated athlete.
// Requires activity:read for Everyone and Followers activities.
// Requires activity:read_all for Only Me activities.
func (s *ActivitiesAPIService) GetByID(access_token string, activityID int64, includeEfforts bool) (*DetailedActivity, error) {
	requestUrl := s.client.BaseURL.JoinPath(activitiesPath, fmt.Sprint(activityID))

	params := url.Values{}
	params.Add("include_all_efforts", fmt.Sprintf("%v", includeEfforts))

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
		body:         params,
	})
	if err != nil {
		return nil, err
	}

	resp := &DetailedActivity{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Summit Feature. Returns the zones of a given activity.
// Requires activity:read for Everyone and Followers activities.
// Requires activity:read_all for Only Me activities.
func (s *ActivitiesAPIService) ListActivityZones(access_token string, activityID int64) ([]ActivityZone, error) {
	requestUrl := s.client.BaseURL.JoinPath(activitiesPath, fmt.Sprint(activityID), "zones")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
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
func (s *ActivitiesAPIService) ListActivityComments(access_token string, activityID int64, p *CommentsReqParams) ([]Comment, error) {
	requestUrl := s.client.BaseURL.JoinPath(activitiesPath, fmt.Sprint(activityID), "comments")

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

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
		body:         params,
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
func (s *ActivitiesAPIService) ListActivityKudoers(access_token string, activityID int64, p *RequestParams) ([]SummaryAthlete, error) {
	requestUrl := s.client.BaseURL.JoinPath(activitiesPath, fmt.Sprint(activityID), "kudos")

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
		access_token: access_token,
		body:         params,
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

// Returns the laps of an activity identified by an identifier. Requires activity:read for Everyone and
// Follower activities. Required activity:read_all for OnlyMeActivities.
func (s *ActivitiesAPIService) ListActivityLaps(access_token string, activityID int64) ([]Lap, error) {
	requestUrl := s.client.BaseURL.JoinPath(activitiesPath, fmt.Sprint(activityID), "laps")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
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
func (s *AthleteAPIService) ListAthleteActivities(access_token string, p *GetActivityParams) ([]SummaryActivity, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath, activitiesPath)
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

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	var resp []SummaryActivity
	err = s.client.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ActivitiesAPIService) GetActivityZones(access_token string, activityID int64) ([]ActivityZone, error) {
	requestUrl := s.client.BaseURL.JoinPath(activitiesPath, fmt.Sprint(activityID), "zones")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
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
	Commute      bool         `json:"commute"`        // Whether this activity is a commute
	Trainer      bool         `json:"trainer"`        // Whether this activity was recorded on a training machine
	HideFromHome bool         `json:"hide_from_home"` // Whether this activity is muted
	Description  string       `json:"description"`    // The description of the activity
	Name         string       `json:"name"`           // The name of the activity
	Type         ActivityType `json:"type"`           // Deprecated. Prefer to use sport_type. In a request where both type and sport_type are present, this field will be ignored
	SportType    SportType    `json:"sport_type"`     // An instance of SportType.
	GearID       string       `json:"gear_id"`        // Identifier for the gear associated with the activity. ‘none’ clears gear from activity
}

// Updates the given activity that is owned by the authenticated athlete. Requires activity:write. Also requires activity:read_all in order
// to update only me activities.
func (s *ActivitiesAPIService) UpdateActivity(access_token string, activityID int64, payload UpdatedActivity) (*DetailedActivity, error) {
	requestUrl := s.client.BaseURL.JoinPath(activitiesPath, fmt.Sprint(activityID))

	json, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
		body:         io.NopCloser(bytes.NewReader(json)),
	})
	if err != nil {
		return nil, err
	}

	resp := &DetailedActivity{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
