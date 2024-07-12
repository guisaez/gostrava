package gostrava

import (
	"fmt"
	"net/http"
)

type GearsService service

const gear string = "gear"

// Returns an equipment using its identifier.
func (s *GearsService) GetEquipment(accessToken string, id string) (*GearDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%s", gear, id),
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
