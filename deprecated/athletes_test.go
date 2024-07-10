package gostrava

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func registerAthleteResponders() {
	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/athlete",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/athlete/get_authenticated_athlete_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/athlete/zones",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/athlete/get_zones_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/athletes/1/stats",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/athlete/get_athlete_stats_response.json"))
		})

	httpmock.RegisterResponder("PUT", "https://www.strava.com/api/v3/athlete",
	func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("Authorization") == "" {
			return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
		}

		return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/athlete/update_athlete_response.json"))
	})
}

func TestCurrentAthlete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerAthleteResponders()

	strava := NewClient(nil)

	resp, err := strava.Athletes.GetAuthenticatedAthlete("1234")
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestGetZones(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerAthleteResponders()

	strava := NewClient(nil)

	resp, err := strava.Athletes.GetZones("1234")
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestGetAthleteStats(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerAthleteResponders()

	strava := NewClient(nil)

	resp, err := strava.Athletes.GetAthleteStats("12345", 1)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}



func TestAthleteUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerAthleteResponders()

	strava := NewClient(nil)

	resp, err := strava.Athletes.UpdateAthlete("1234", UpdateAthletePayload{Weight: 10.2})
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}