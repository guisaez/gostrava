package gostrava

import "fmt"

type AthleteAPIService apiService

const (
	athletePath = "/athlete"
)

// Returns the currently authenticated athlete. Tokens with profile:read_all scope will receive
// a detailed athlete representation; all others will receive a SummaryAthlete representation
func (s *AthleteAPIService) CurrentAthlete(access_token string) (*DetailedAthlete, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath)

	req, err := s.client.get(requestUrl, nil, access_token)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+access_token)

	resp := &DetailedAthlete{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the authenticated athlete's heart rate and power zones. Requires profile:read_all.
func (s *AthleteAPIService) GetZones(access_token string) (*Zones, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath, "/zones")

	req, err := s.client.get(requestUrl, nil, access_token)
	if err != nil {
		return nil, err
	}

	resp := &Zones{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the activity stats of an athlete. Only includes data from activities set to Everyone's visibility.
func (s *AthleteAPIService) GetAthleteStats(access_token string, id int64) (*ActivityStats, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath, fmt.Sprint(id), "stats")

	req, err := s.client.get(requestUrl, nil, access_token)
	if err != nil {
		return nil, err
	}

	resp := &ActivityStats{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Return a list of the clubs whose membership includes the authenticated athlete.
func (s *AthleteAPIService) ListClubs(access_token string) ([]SummaryClub, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath, "clubs")

	req, err := s.client.get(requestUrl, nil, access_token)
	if err != nil {
		return nil, err
	}

	resp := []SummaryClub{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a list of the routes created by the authenticated athlete. Private routes are filtered out
// unless request by a token with read_all scope.
func (s *AthleteAPIService) ListRoutes(access_token string) ([]Route, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath, "routes")

	req, err := s.client.get(requestUrl, nil, access_token)
	if err != nil {
		return nil, err
	}

	resp := []Route{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
