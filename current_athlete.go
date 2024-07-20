package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type CurrentAthleteService service

// Returns the currently authenticated athlete. Tokens with profile:read_all scope will receive
// a detailed athlete representation; all others will receive a SummaryAthlete representation
func (s *CurrentAthleteService) GetAthlete(accessToken string) (*AthleteDetailed, error) {
	req, err := s.client.NewRequest(RequestOpts{
		Path:        "athlete",
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(AthleteDetailed)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, err
}

// Returns the current athlete's heart rate and power zones. Requires profile:read_all.
func (s *CurrentAthleteService) GetZones(accessToken string) (*Zones, error) {
	req, err := s.client.NewRequest(RequestOpts{
		Path:        "athlete/zones",
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(Zones)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type UpdatedAthlete struct {
	Weight float32 // The weigh of the athlete in kilograms.
}

// Updates the authenticated user. Requires profile:write scope
func (s *CurrentAthleteService) Update(accessToken string, updatedAthlete UpdatedAthlete) (*AthleteDetailed, error) {
	params := url.Values{}
	params.Set("weight", fmt.Sprintf("%.2f", updatedAthlete.Weight))

	req, err := s.client.NewRequest(RequestOpts{
		Path:        "athlete",
		Method:      http.MethodPut,
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	resp := new(AthleteDetailed)
	if err := s.client.Do(req, resp); err != nil {
		return nil, err
	}

	return resp, err
}

// Return a list of the clubs whose membership includes the authenticated athlete.
func (s *CurrentAthleteService) ListClubs(accessToken string, opts RequestParams) ([]ClubSummary, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	req, err := s.client.NewRequest(RequestOpts{
		Path:        "athlete/clubs",
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []ClubSummary{}
	if err := s.client.Do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns the activities of an athlete for a specific identifier. Requires activity:read, OnlyMe activities will be filtered out unless
// requested by a token with activity_read:all.
func (s *CurrentAthleteService) ListActivities(accessToken string, opts GetActivityOpts) ([]ActivitySummary, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page_size", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}
	if opts.Before > 0 {
		params.Set("before", strconv.Itoa(opts.Before))
	}
	if opts.After > 0 {
		params.Set("after", strconv.Itoa(opts.After))
	}

	req, err := s.client.NewRequest(RequestOpts{
		Path:        "athlete/activities",
		AccessToken: accessToken,
		Body:        params,
	})
	if err != nil {
		return nil, err
	}

	var resp []ActivitySummary
	err = s.client.Do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
