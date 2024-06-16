package gostrava

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

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

	_, err = strava.Clubs.GetAdministrators("12345", 1, &ClubRequestParams{})
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}
}

func TestGetActivities(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerClubResponders()

	strava := NewClient(nil)

	_, err := strava.Clubs.GetById("", 1)
	if err == nil {
		t.Error("expected and error")
	}

	_, err = strava.Clubs.GetActivities("12345", 1, &ClubRequestParams{})
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

	_, err = strava.Clubs.GetMembers("12345", 1, &ClubRequestParams{Page: 1, PerPage: 2})
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}
}

func registerClubResponders() {
	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/clubs/1",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/clubs/clubs_get_by_id_response.json"))
		},
	)

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/clubs/1/admins",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/clubs/clubs_get_administrators_response.json"))
		},
	)

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/clubs/1/activities",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/clubs/clubs_get_activities_response.json"))
		},
	)

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/clubs/1/members",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/clubs/clubs_get_members_response.json"))
		},
	)
}
