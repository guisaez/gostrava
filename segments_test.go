package gostrava

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func registerSegmentsResponders() {
	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/segments/explore",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/segments/explore_segments_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/segments/starred",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/segments/list_starred_segments_response.json"))
		})

	httpmock.RegisterResponder("GET", "https://www.strava.com/api/v3/segments/1",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/segments/get_segment_response.json"))
		})

	httpmock.RegisterResponder("PUT", "https://www.strava.com/api/v3/segments/1/starred",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "" {
				return httpmock.NewJsonResponse(http.StatusUnauthorized, httpmock.File("./mock/api_error_response.json"))
			}

			return httpmock.NewJsonResponse(http.StatusOK, httpmock.File("./mock/segments/star_segment_response.json"))
		})
}

func TestGetExploreSegments(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerSegmentsResponders()

	strava := NewClient(nil)

	resp, err := strava.Segments.ExploreSegments("123345", Bounds{SWLat: 12.2, SWLng: 12.1, NELat: 12.1, NELng: 12.1}, &ExploreSegmentsOpts{ActivityType: "running"})
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestListStarredSegments(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerSegmentsResponders()

	strava := NewClient(nil)

	resp, err := strava.Segments.ListStarredSegments("12345", &RequestParams{})
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestGetSegment(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerSegmentsResponders()

	strava := NewClient(nil)

	resp, err := strava.Segments.GetSegment("2135", 1)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}

func TestStarSegment(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerSegmentsResponders()

	strava := NewClient(nil)

	resp, err := strava.Segments.StarSegment("1234", 1, true)
	if err != nil {
		t.Errorf("error not expected, got %v", err.Error())
	}

	json, _ := json.MarshalIndent(resp, "", "\t")

	fmt.Println(string(json))
}