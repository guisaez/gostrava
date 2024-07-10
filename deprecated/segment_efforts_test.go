package gostrava

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func registerSegmentEffortResponders() {
	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/segment_efforts/1",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/segment_efforts/get_segment_effort_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/segment_efforts",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/segment_efforts/list_segment_efforts_response.json"))
		})
}

func TestGetSegmentEffort(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerSegmentEffortResponders()

	strava := NewClient(nil)

	resp, err := strava.SegmentEfforts.GetSegmentEffort("12345", 1)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestListSegmentEfforts(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerSegmentEffortResponders()

	strava := NewClient(nil)

	resp, err := strava.SegmentEfforts.ListSegmentEfforts("12345", &ListSegmentEffortOptions{})
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}
