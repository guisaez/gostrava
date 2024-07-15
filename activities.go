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

type ActivityDetailed struct {
	ActivitySummary
	Calories        float32                  `json:"calories"`               // The number of kilocalories consumed during this activity 
	Description     string                   `json:"description,omitempty"`  // The description of the activity
	AvailableZones  []string                 `json:"available_zones"`        // Activity available zones
	BestEfforts     []SegmentEffortDetailed  `json:"best_efforts,omitempty"` // A collection of SegmentEffortDetailed objects.
	SportType       SportType                `json:"sport_type"`
	StartDate       TimeStamp                `json:"start_date"`       // The time at which the activity was started.
	StartDateLocal  TimeStamp                `json:"start_date_local"` // The time at which the activity was started in the local timezone.
	Timezone        string                   `json:"timezone"`         // The timezone of the activity
	UTCOffset       float32                  `json:"utc_offset"`
	DeviceName      string                   `json:"device_name,"`          // The name of the device used to record the activity
	LocationCity    string                   `json:"location_city"`         // City where the activity took place
	LocationState   string                   `json:"location_state"`        // State where the activity took place
	LocationCountry string                   `json:"location_country"`      // Country where the activity took place



	EmbedToken      string                   `json:"embed_token,omitempty"` // The token used to embed a Strava activity
	Gear            *GearSummary             `json:"gear,omitempty"`        // An instance of SummaryGear.
	Laps            []Lap                    `json:"laps,omitempty"`        // A collection of Lap objects.
	Photos          PhotosSummary            `json:"photos"`                // An instance of PhotosSummary.
	SegmentEfforts  []*SegmentEffortDetailed `json:"segment_efforts"`       // A collection of SegmentEffortDetailed objects.
	SplitsMetric    []*Split                 `json:"splits_metric"`         // The splits of this activity in metric units (for runs)
	SplitsStandard  []*Split                 `json:"splits_standard"`       // The splits of this activity in imperial units (for runs)
}

type ActivitySummary struct {
	ActivityMeta
	Name               string       `json:"name"`
	Distance           float32      `json:"distance"`             // The activity's distance, in meters
	MovingTime         int          `json:"moving_time"`          // The activity's moving time, in seconds
	ElapsedTime        int          `json:"elapsed_time"`         // The activity's elapsed time, in seconds
	TotalElevationGain float32      `json:"total_elevation_gain"` // The activity's total elevation gain.
	Type               ActivityType `json:"type"`                 // Deprecated. Prefer to use sport_type

	
	Athlete                    AthleteMeta `json:"athlete"`                            // An instance of AthleteMeta.
	AthleteCount               int         `json:"athlete_count"`                      // The number of athletes for taking part in a group activity
	AvgHeartRate               float32     `json:"average_heartrate"`                  // The activity's average heart rate, in beats per minute
	AvgSpeed                   float32     `json:"average_speed"`                      // The activity's average speed, in meters per second
	CommentCount               int         `json:"comment_count"`                      // The number of comments for this activity
	Commute                    bool        `json:"commute,omitempty"`                  // Whether this activity is a commute
	DisplayHideHeartRateOption *bool       `json:"display_heartrate_option,omitempty"` //

	ElevationHigh   float32 `json:"elev_high"`                // The activity's highest elevation, in meters
	ElevationLow    float32 `json:"elev_low"`                 // The activity's lowest elevation, in meters
	EndLatLng       LatLng  `json:"end_latlng,omitempty"`     // An instance of LatLng.
	ExternalID      string  `json:"external_id"`              // The identifier provided at upload time
	Flagged         bool    `json:"flagged"`                  // Whether this activity is flagged
	FromAcceptedTag bool    `json:"from_accepted_tag"`        //
	GearID          *string `json:"gear_id,omitempty"`        // The id of the gear for the activity
	HasHeartRate    bool    `json:"has_heartrate"`            // Indicates weather the activity has a heartrate recorder
	HasKudoed       bool    `json:"has_kudoed"`               // Whether the logged-in athlete has kudoed this activity
	HeartRateOptOut bool    `json:"heartrate_opt_out"`        //
	HideFromHome    *bool   `json:"hide_from_home,omitempty"` // Whether the activity is muted
	KudosCount      int     `json:"kudos_count"`              // The number of kudos given for this activity

	Manual       bool        `json:"manual"`        // Indicates whether this activity was manually created by the user
	Map          PolylineMap `json:"map"`           // An instance of PolylineMap.
	MaxHeartRate float32     `json:"max_heartrate"` // The activity's max heartrate in beats per minute
	MaxSpeed     float32     `json:"max_speed"`     // The activity's max speed, in meters per second

	PhotoCount int  `json:"photo_count"` // The number of Instagram photos for this activity
	PRCount    int  `json:"pr_count"`    //
	Private    bool `json:"private"`     // Whether this activity is private

	StartLatLng LatLng   `json:"start_latlng"`           // An instance of LatLng.
	SufferScore *float32 `json:"suffer_score,omitempty"` //

	TotalPhotoCount int  `json:"total_photo_count"` // The number of Instagram and Strava photos for this activity
	Trainer         bool `json:"trainer"`           // Whether this activity was recorded on a training machine

	UploadID    int    `json:"upload_id"`     // The identifier of the upload that resulted in this activity
	UploadIdStr string `json:"upload_id_str"` // The unique identifier of the upload in string format

	WorkoutType *int `json:"workout_type,omitempty"` //  The activity's workout type
}

type ActivityMeta struct {
	ID            int     `json:"id"`                   // The unique identifier of the activity
	ResourceState int8    `json:"resource_state"`       //
	Visibility    *string `json:"visibility,omitempty"` //
}

// *****************************************************

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
