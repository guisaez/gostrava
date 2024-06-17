package gostrava

import (
	"fmt"
	"net/http"
)

type GearsAPIService apiService

// Returns an equipment using its identifier.
func (s *GearsAPIService) GetEquipment(access_token string, id int64) (*DetailedGear, error) {
	requestUrl := s.client.BaseURL.JoinPath(gearPath, fmt.Sprint(id))

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := &DetailedGear{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
