package gostrava

import (
	"context"
	"net/http"
)

// *************** Types ********************

type GearSummary struct {
	ID           string  `json:"id"`             // The gear's unique identifier.
	ResourceRate int8    `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
	Primary      bool    `json:"primary"`        // Whether this gear's is the owner's default one.
	Name         string  `json:"name"`           // The gear's name.
	Distance     float32 `json:"distance"`       // The distance logged with this gear.
}

type GearDetailed struct {
	GearSummary
	BrandName   string `json:"brand_name"`  // The gear's brand name.
	ModelName   string `json:"model_name"`  // The gear's model name.
	FrameType   int    `json:"frame_type"`  // The gear's frame type (bike only).
	Description string `json:"description"` // The gear's description.
}

// *************** Methods ********************

type GearService service

const gears = "/api/v3/gears"

// Returns an equipment using its identifier
//
// GET: https://www.strava.com/api/v3/gears/{id}
func (s *GearService) GetEquipment(ctx context.Context, accessToken string, id string) (*GearDetailed, *http.Response, error) {
	urlStr := gears + "/" + id

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	gear := new(GearDetailed)
	resp, err := s.client.DoAndParse(ctx, req, gear)
	if err != nil {
		return nil, resp, err
	}

	return gear, resp, nil
}