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

type ActivityMeta struct {
	ID            int     `json:"id"`                   // The unique identifier of the activity
	ResourceState uint8   `json:"resource_state"`       //
	Visibility    *string `json:"visibility,omitempty"` //
}

type ActivitySummary struct {
	ActivityMeta
	AchievementCount           int          `json:"achievement_count"`           // The number of achievements gained during this activity
	Athlete                    AthleteMeta  `json:"athlete"`                     // An instance of AthleteMeta.
	AthleteCount               int          `json:"athlete_count"`               // The number of athletes for taking part in a group activity
	AvgHeartRate               float32      `json:"average_heartrate,omitempty"` // The activity's average heart rate, in beats per minute
	AvgSpeed                   float32      `json:"average_speed,omitempty"`     // The activity's average speed, in meters per second
	CommentCount               int          `json:"comment_count"`               // The number of comments for this activity
	Commute                    bool         `json:"commute"`                     // Whether this activity is a commute
	DisplayHideHeartRateOption bool         `json:"display_heartrate_option"`    //
	Distance                   float32      `json:"distance"`                    // The activity's distance, in meters
	ElapsedTime                int          `json:"elapsed_time"`                // The activity's elapsed time, in seconds
	ElevationHigh              float32      `json:"elev_high"`                   // The activity's highest elevation, in meters
	ElevationLow               float32      `json:"elev_low"`                    // The activity's lowest elevation, in meters
	EndLatLng                  LatLng       `json:"end_latlng"`                  // An instance of LatLng.
	ExternalID                 string       `json:"external_id"`                 // The identifier provided at upload time
	Flagged                    bool         `json:"flagged"`                     // Whether this activity is flagged
	FromAcceptedTag            bool         `json:"from_accepted_tag"`           //
	GearID                     *string      `json:"gear_id"`                     // The id of the gear for the activity
	HasHeartRate               bool         `json:"has_heartrate"`               // Indicates weather the activity has a heartrate recorder
	HasKudoed                  bool         `json:"has_kudoed"`                  // Whether the logged-in athlete has kudoed this activity
	HeartRateOptOut            bool         `json:"heartrate_opt_out"`           //
	HideFromHome               bool         `json:"hide_from_home"`              // Whether the activity is muted
	KudosCount                 int          `json:"kudos_count"`                 // The number of kudos given for this activity
	LocationCity               *string      `json:"location_city"`               //
	LocationCountry            *string      `json:"location_country"`            //
	LocationState              *string      `json:"location_state"`              //
	Manual                     bool         `json:"manual"`                      // Indicates whether this activity was manually created by the user
	Map                        *PolylineMap `json:"map"`                         // An instance of PolylineMap.
	MaxHeartRate               float32      `json:"max_heartrate"`               // The activity's max heartrate in beats per minute
	MaxSpeed                   float32      `json:"max_speed"`                   // The activity's max speed, in meters per second
	MovingTime                 int          `json:"moving_time"`                 // The activity's moving time, in seconds
	Name                       string       `json:"name"`                        // The name of the activity
	PhotoCount                 int          `json:"photo_count"`                 // The number of Instagram photos for this activity
	PRCount                    int          `json:"pr_count"`                    //
	Private                    bool         `json:"private"`                     // Whether this activity is private
	SportType                  SportType    `json:"sport_type"`                  // An instance of SportType.
	StartDate                  TimeStamp    `json:"start_date"`                  // The time at which the activity was started.
	StartDateLocal             TimeStamp    `json:"start_date_local"`            // The time at which the activity was started in the local timezone.
	StartLatLng                LatLng       `json:"start_latlng"`                // An instance of LatLng.
	SufferScore                float32      `json:"suffer_score"`                //
	Timezone                   string       `json:"timezone"`                    // The timezone of the activity
	TotalElevationGain         float32      `json:"total_elevation_gain"`        // The activity's total elevation gain.
	TotalPhotoCount            int          `json:"total_photo_count"`           // The number of Instagram and Strava photos for this activity
	Trainer                    bool         `json:"trainer"`                     // Whether this activity was recorded on a training machine
	Type                       ActivityType `json:"type"`                        // Deprecated. Prefer to use sport_type
	UploadID                   int          `json:"upload_id"`                   // The identifier of the upload that resulted in this activity
	UploadIdStr                string       `json:"upload_id_str"`               // The unique identifier of the upload in string format
	UTCOffset                  float32      `json:"utc_offset"`                  //
	WorkoutType                int          `json:"workout_type,omitempty"`      //  The activity's workout type
}

type ActivityDetailed struct {
	ActivitySummary
	AvailableZones []string                `json:"available_zones"` // Activity available zones, it can include the following values: heartrate, power
	BestEfforts    []SegmentEffortDetailed `json:"best_efforts"`    // A collection of SegmentEffortDetailed objects.
	Calories       float32                 `json:"calories"`        // The number of kilocalories consumed during this activity
	DeviceName     string                  `json:"device_name"`     // The name of the device used to record the activity
	Description    *string                 `json:"description"`     // The description of the activity
	EmbedToken     string                  `json:"embed_token"`     // The token used to embed a Strava activity
	Gear           *GearSummary            `json:"gear"`            // An instance of SummaryGear.
	Laps           []Lap                   `json:"laps"`            // A collection of Lap objects.
	Photos         *PhotosSummary          `json:"photos"`          // An instance of PhotosSummary.
	SegmentEfforts []SegmentEffortDetailed `json:"segment_efforts"` // A collection of SegmentEffortDetailed objects.
	SplitsMetric   []Split                 `json:"splits_metric"`   // The splits of this activity in metric units (for runs)
	SplitsStandard []Split                 `json:"splits_standard"` // The splits of this activity in imperial units (for runs)
}

type ActivityZone struct {
	Score               float32               `json:"score"`
	DistributionBuckets TimedZoneDistribution `json:"distribution_buckets"`
	Type                string                `json:"type"` // May take one of the following values: heartrate, power
	SensorBased         bool                  `json:"sensor_based"`
	Points              int                   `json:"points"`
	CustomZones         bool                  `json:"custom_zones"`
	Max                 int                   `json:"max"`
}

type Comment struct {
	ID         int            `json:"id"`          // The unique identifier of this comment
	ActivityID int            `json:"activity_id"` // The identifier of the activity this comment is related to
	Text       string         `json:"text"`        // The content of the comment
	Athlete    AthleteSummary `json:"athlete"`     // An instance of AthleteSummary.
	CreatedAt  TimeStamp      `json:"created_at"`  // The time at which this comment was created.
}

type NewActivity struct {
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
		Path:        athlete,
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
