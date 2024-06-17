package gostrava

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func registerRoutesResponders() {
	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/routes/1",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/routes/get_route_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/athletes/1/routes",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/routes/list_athlete_routes_response.json"))
		})
}

func TestGetRoute(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerRoutesResponders()

	strava := NewClient(nil)

	resp, err := strava.Routes.GetById("12357", 1)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestListRoutes(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerRoutesResponders()

	strava := NewClient(nil)

	resp, err := strava.Routes.ListAthleteRoutes("12357", 1, nil)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}
