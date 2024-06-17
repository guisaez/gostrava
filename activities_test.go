package gostrava

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func registerActivityResponders() {
	httpmock.RegisterResponder("POST", "https://www.strava.com/api/v3/activities",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/activities/create_an_activity_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/activities/1",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/activities/get_activity_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/activities/1/comments",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/activities/list_activity_comments_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/activities/1/kudos",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/activities/list_activity_kudoers_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/activities/1/laps",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/activities/list_activity_laps_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/athlete/activities",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/activities/list_athlete_activities_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/activities/1/zones",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/activities/get_activity_zones_response.json"))
		})

	httpmock.RegisterResponder("PUT", "https://www.strava.com/api/v3/activities/1",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/activities/update_activity_response.json"))
		})
}

func TestNewActivity(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerActivityResponders()

	strava := NewClient(nil)

	resp, err := strava.Activities.New("1234", Activity{
		Name:           "test",
		Type:           ActivityTypes.Canoening,
		SportType:      SportTypes.Pickleball,
		StartDateLocal: time.Now(),
		Description:    "Test",
	})
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestGetActivityById(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerActivityResponders()

	strava := NewClient(nil)

	resp, err := strava.Activities.GetByID("12345", 1, false)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestListActivityComments(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerActivityResponders()

	strava := NewClient(nil)

	resp, err := strava.Activities.ListActivityComments("12345", 1, nil)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestListActivityKudoers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerActivityResponders()

	strava := NewClient(nil)

	resp, err := strava.Activities.ListActivityKudoers("12345", 1, nil)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestListActivityLaps(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerActivityResponders()

	strava := NewClient(nil)

	resp, err := strava.Activities.ListActivityLaps("12345", 1)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestListAthleteActivities(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerActivityResponders()

	strava := NewClient(nil)

	resp, err := strava.Athletes.ListAthleteActivities("12345", nil)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestGetActivityZones(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerActivityResponders()

	strava := NewClient(nil)

	resp, err := strava.Activities.GetActivityZones("12345", 1)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}


func TestUpdateActivity(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerActivityResponders()

	strava := NewClient(nil)

	resp, err := strava.Activities.UpdateActivity("12345", 1, UpdatedActivity{
		Description: "Test",
	})
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}
