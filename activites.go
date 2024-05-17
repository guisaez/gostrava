package gostrava

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type StravaActivities struct {
	AccessToken string
	*StravaClient
}

type GeneralParams struct {
	Page    int // Page number. Defaults to 1
	PerPage int // Number of items per page. Defaults to 30
}

type NewActivityRequest struct {
	Name           string       `json:"name"`             // The name of the activity.
	Type           ActivityType `json:"type"`             // Type of activity. For example - Run, Ride etc.
	SportType      SportType    `json:"sport_type"`       // Sport type of activity. For example - Run, MountainBikeRide, Ride, etc.
	StartDateLocal time.Time    `json:"start_date_local"` // ISO 8601 formatted date time.
	ElapsedTime    int          `json:"elapsed_time"`     // In seconds.
	Description    string       `json:"description"`      // Description of the activity.
	Distance       int          `json:"distance"`         // In meters.
	Trainer        int8         `json:"trainer"`          // Set to 1 to mark as a trainer activity.
	Commute        int8         `json:"commute"`          // Set to 1 to mark as commute.
}

// Creates a manual activity for an athlete, requires activity:write scope.
func (sc *StravaActivities) New(ctx context.Context, payload NewActivityRequest) (*DetailedActivity, error) {

	params := url.Values{}
	params.Set("name", payload.Name)
	params.Set("type", string(payload.Type))
	params.Set("sport_type", string(payload.SportType))
	params.Set("start_date_local", payload.StartDateLocal.Format(time.RFC3339)) // Assuming RFC3339 format
	params.Set("elapsed_time", strconv.Itoa(payload.ElapsedTime))
	params.Set("description", payload.Description)
	params.Set("distance", strconv.Itoa(payload.Distance))
	params.Set("trainer", strconv.Itoa(int(payload.Trainer)))
	params.Set("commute", strconv.Itoa(int(payload.Commute)))

	path := "/activities"

	var detailedActivity DetailedActivity
	err := sc.postForm(ctx, sc.AccessToken, path, params, &detailedActivity)
	if err != nil {
		return nil, err
	}

	return &detailedActivity, nil
}

// Returns the given activity that is owned by the authenticated athlete.
// Requires activity:read for Everyone and Followers activities.
// Requires activity:read_all for Only Me activities.
func (sc *StravaActivities) ByID(ctx context.Context, activityID int64, includeEfforts bool) (*DetailedActivity, error) {

	path := fmt.Sprintf("/activities/%d", activityID)

	params := url.Values{}
	params.Add("include_all_efforts", fmt.Sprintf("%v", includeEfforts))

	var resp DetailedActivity
	err := sc.get(ctx, sc.AccessToken, path, params, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Summit Feature. Returns the zones of a given activity.
// Requires activity:read for Everyone and Followers activities.
// Requires activity:read_all for Only Me activities.
func (sc *StravaActivities) GetZones(ctx context.Context, activityID int64) ([]ActivityZone, error) {

	path := fmt.Sprintf("/activities/%d", activityID)

	var resp []ActivityZone
	err := sc.get(ctx, sc.AccessToken, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}


type GetCommentsOptions struct {
	PageSize    int    // Number of items per page. Defaults to 30
	AfterCursor string // Cursor of the las item in the previous page of results, used to request the subsequent page of results. When omitted, the first page of results is fetched.
}

// Returns the comments on the given activity. Requires activity:read for Everyone and Followers activities. Requires activity:read_all for Only Me activities.
func (sc *StravaActivities) GetComments(ctx context.Context, activityID int64, opt *GetCommentsOptions) ([]Comment, error) {

	path := fmt.Sprintf("/activities/%d/comments", activityID)

	params := url.Values{}
	if opt != nil {
		if opt.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opt.PageSize))
		}
		if opt.AfterCursor != "" {
			params.Set("after_cursor", opt.AfterCursor)
		}
	}

	var resp []Comment
	err := sc.get(ctx, sc.AccessToken, path, params, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the athletes who kudoed an activity identified by an identifier. Requires activity:read for Everyone and Followers activities.
// Requires activity:read_all for OnlyMe Activities
func (sc *StravaActivities) GetKudoers(ctx context.Context, activityID int64, opt *GeneralParams) ([]SummaryAthlete, error) {

	path := fmt.Sprintf("/activities/%d/kudos", activityID)

	params := url.Values{}
	if opt != nil {
		if opt.Page > 0 {
			params.Set("page", strconv.Itoa(opt.Page))
		}
		if opt.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(opt.Page))
		}
	}

	var resp []SummaryAthlete
	err := sc.get(ctx, sc.AccessToken, path, params, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the laps of an activity identified by an identifier. Requires activity:read for Everyone and
// Follower activities. Required activity:read_all for OnlyMeActivities.
func (sc *StravaActivities) GetLaps(ctx context.Context, activityID int64) ([]Lap, error) {

	path := fmt.Sprintf("/activities/%d/laps", activityID)

	var resp []Lap
	err := sc.get(ctx, sc.AccessToken, path, nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Updates the given activity that is owned by the authenticated athlete. Requires activity:write. Also requires activity:read_all in order
// to update only me activities.
func (sc *StravaActivities) UpdateActivity(ctx context.Context, activityID int64, ua UpdatableActivity) (*DetailedActivity, error) {

	path := fmt.Sprintf("/activities/%d", activityID)

	var resp DetailedActivity
	err := sc.put(ctx, sc.AccessToken, path, "application/json", ua, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}