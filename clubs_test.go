package gostrava

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func registerClubResponders() {
	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/clubs/1",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/clubs/get_club_response.json"))
		},
	)

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/clubs/1/admins",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/clubs/list_club_administrators_response.json"))
		},
	)

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/clubs/1/activities",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/clubs/list_club_activities_response.json"))
		},
	)

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/clubs/1/members",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/clubs/list_club_members_response.json"))
		},
	)

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/athlete/clubs",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/clubs/list_athlete_clubs_response.json"))
		})
}

func TestGetClubById(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerClubResponders()

	strava := NewClient(nil)

	_, err := strava.Clubs.GetById("", 1)
	if err == nil {
		t.Error("expected and error")
	}

	_, err = strava.Clubs.GetById("12345", 1)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}
}

func TestGetAdministrators(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerClubResponders()

	strava := NewClient(nil)

	_, err := strava.Clubs.GetById("", 1)
	if err == nil {
		t.Error("expected and error")
	}

	_, err = strava.Clubs.ListClubAdministrators("12345", 1, &RequestParams{})
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}
}

func TestGetClubActivities(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerClubResponders()

	strava := NewClient(nil)

	_, err := strava.Clubs.GetById("", 1)
	if err == nil {
		t.Error("expected and error")
	}

	_, err = strava.Clubs.ListClubActivities("12345", 1, &RequestParams{})
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}
}

func TestGetMembers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerClubResponders()

	strava := NewClient(nil)

	_, err := strava.Clubs.GetById("", 1)
	if err == nil {
		t.Error("expected and error")
	}

	_, err = strava.Clubs.ListClubMembers("12345", 1, &RequestParams{Page: 1, PerPage: 2})
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}
}

func TestListAthleteClubs(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerClubResponders()

	strava := NewClient(nil)

	resp, err := strava.Clubs.ListAthleteClubs("123456", nil)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}
