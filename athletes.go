package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
)

type AthleteAPIService apiService

// Returns the currently authenticated athlete. Tokens with profile:read_all scope will receive
// a detailed athlete representation; all others will receive a SummaryAthlete representation
func (s *AthleteAPIService) GetAuthenticatedAthlete(access_token string) (*DetailedAthlete, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath)

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := &DetailedAthlete{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the authenticated athlete's heart rate and power zones. Requires profile:read_all.
func (s *AthleteAPIService) GetZones(access_token string) ([]ActivityZone, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath, "zones")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := []ActivityZone{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the activity stats of an athlete. Only includes data from activities set to Everyone's visibility.
func (s *AthleteAPIService) GetAthleteStats(access_token string, id int64) (*ActivityStats, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletesPath, fmt.Sprint(id), "stats")

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := &ActivityStats{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type UpdateAthletePayload struct {
	Weight float64 // The weigh of the athlete in kilograms.
}

// Update the currently authenticated athlete. Requires profile:write scope.
func (s *AthleteAPIService) UpdateAthlete(access_token string, p UpdateAthletePayload) (*DetailedAthlete, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath)

	params := url.Values{}

	if p.Weight > 0 {
		params.Set("weight", fmt.Sprintf("%.2f", p.Weight))
	}

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodPut,
		access_token: access_token,
		body:         params,
	})
	if err != nil {
		return nil, err
	}

	resp := &DetailedAthlete{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
