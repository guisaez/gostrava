package gostrava

// *****************************************************

type GearsService service

// Returns an equipment using its identifier.
func (s *GearsService) GetEquipment(accessToken string, id string) (*GearDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "gear/" + id,
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
