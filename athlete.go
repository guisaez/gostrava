package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type AthleteService service

// Returns the currently authenticated athlete. Tokens with profile:read_all scope will receive
// a detailed athlete representation; all others will receive a SummaryAthlete representation
func (s *AthleteService) GetAuthenticatedAthlete(accessToken string) (*AthleteDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "athlete",
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(AthleteDetailed)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, err
}

// Returns the current athlete's heart rate and power zones. Requires profile:read_all.
func (s *AthleteService) GetZones(accessToken string) (*Zones, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "athlete/zones",
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(Zones)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the activity stats of an athlete. Only includes data from activities set to Everyone's visibility.
func (s *AthleteService) GetAthleteStats(accessToken string, id int) (*ActivityStats, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "athletes/" + strconv.Itoa(id) + "/stats",
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(ActivityStats)
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type UpdatedAthlete struct {
	Weight float32 // The weigh of the athlete in kilograms.
}

// Updates the authenticated user. Requires profile:write scope
func (s *AthleteService) Update(accessToken string, updatedAthlete UpdatedAthlete) (*AthleteDetailed, error) {
	params := url.Values{}
	params.Set("weight", fmt.Sprintf("%.2f", updatedAthlete.Weight))

	req, err := s.client.newRequest(requestOpts{
		Path:        "athlete",
		Method:      http.MethodPut,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := new(AthleteDetailed)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, err
}
