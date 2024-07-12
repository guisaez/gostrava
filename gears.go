package gostrava

import (
	"fmt"
	"net/http"
)

type GearsService service

const gear string = "gear"

type GearSummary struct {
	ID           string  `json:"id"`             // The gear's unique identifier.
	ResourceRate uint8   `json:"resource_state"` // Resource state, indicates level of detail. Possible values: 1 (Meta), 2 (Summary), 3 (Detailed)
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

// Returns an equipment using its identifier.
func (s *GearsService) GetEquipment(accessToken string, id int64) (*GearDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d", gear, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(GearDetailed)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
