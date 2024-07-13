package gostrava

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ActivityService service

type NewActivity struct {
	Name           string       `json:"name"`             // The name of the activity.
	Type           ActivityType `json:"activity_type"`    // Type of activity. For example - Run, Ride etc.
	SportType      SportType    `json:"sport_type"`       // Sport type of activity. For example - Run, MountainBikeRide, Ride, etc.
	StartDateLocal time.Time    `json:"start_date_local"` // ISO 8601 formatted date time.
	ElapsedTime    int          `json:"elapsed_time"`     // In seconds.
	Description    string       `json:"description"`      // Description of the activity.
	Distance       int          `json:"distance"`         // In meters.
	Trainer        bool         `json:"trainer"`          // Mark as a trainer activity, 0 otherwise.
	Commute        bool         `json:"commute"`          // Mark as commuter, 0 otherwise.
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
	// some of go drawbacks/potentials -> explicit conversions
	if body.Trainer {
		formData.Set("trainer", "1")
	}
	if body.Commute {
		formData.Set("commute", "1")
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "activities",
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
	params.Add("include_all_efforts", "true")

	req, err := s.client.newRequest(requestOpts{
		Path:        "activities/" + strconv.Itoa(id),
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

type CommentsReqParams struct {
	Page        int    // Page number. Defaults to 1
	PerPage     int    // Number of items per page. Defaults to 30
	PageSize    int    // Number of items per page. Defaults to 30
	AfterCursor string // Cursor of the las item in the previous page of results, used to request the subsequent page of results. When omitted, the first page of results is fetched.
}

// Returns the comments on the given activity. Requires activity:read for Everyone and Followers activities. Requires activity:read_all for Only Me activities.
func (s *ActivityService) ListActivityComments(accessToken string, id int, opts CommentsReqParams) ([]Comment, error) {
	params := url.Values{}
	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}
	if opts.PageSize > 0 {
		params.Set("page_size", strconv.Itoa(opts.PageSize))
	}
	if opts.AfterCursor != "" {
		params.Set("after_cursor", opts.AfterCursor)
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "activities/" + strconv.Itoa(id) + "/comments",
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
func (s *ActivityService) ListActivityKudoers(accessToken string, id int, opts RequestParams) ([]AthleteSummary, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "activities/" + strconv.Itoa(id) + "/kudos",
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
		Path:        "activities/" + strconv.Itoa(id) + "/laps",
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

type GetActivityOpts struct {
	Page    int // Page number. Defaults to 1
	PerPage int // Number of items per page. Defaults to 30
	Before  int // An epoch timestamp to use for filtering activities that have taken place before that certain time.
	After   int // An epoch timestamp to use for filtering activities that have taken place after a certain time.
}

// Returns the activities of an athlete for a specific identifier. Requires activity:read, OnlyMe activities will be filtered out unless
// requested by a token with activity_read:all.
func (s *ActivityService) ListAthleteActivities(accessToken string, opts GetActivityOpts) ([]ActivitySummary, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page_size", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("page_size", strconv.Itoa(opts.PerPage))
	}
	if opts.Before > 0 {
		params.Set("before", strconv.Itoa(opts.Before))
	}
	if opts.After > 0 {
		params.Set("after", strconv.Itoa(opts.After))
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "athlete/activities",
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

// Summit Feature. Returns the zones of a given activity.
// Requires activity:read for Everyone and Followers activities.
// Requires activity:read_all for Only Me activities.
func (s *ActivityService) GetActivityZones(accessToken string, id int) ([]ActivityZone, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "activities/" + strconv.Itoa(id) + "/zones",
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
		Path:        "activities/" + strconv.Itoa(id),
		AccessToken: accessToken,
		Method:      http.MethodPut,
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
