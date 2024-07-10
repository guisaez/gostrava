package gostrava

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func registerGearResponders() {
	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/gear/1",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/gears/get_equipment_response.json"))
		},
	)
}

func TestGetEquipment(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerGearResponders()

	strava := NewClient(nil)

	_, err := strava.Gears.GetEquipment("12345", 1)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}
}
