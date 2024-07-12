package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ClubService service

const clubs string = "clubs"

// Returns a given club using its identifier
func (s *ClubService) GetById(accessToken string, id int) (*ClubDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d", clubs, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := new(ClubDetailed)
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a list of the administrators of a given club.
// It uses a predefined ClubAthlete struct and not AthleteSummary because it currently sends only FirstName and LastName
func (s *ClubService) ListAdministrators(accessToken string, id int, p *RequestParams) ([]ClubAthlete, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/admins", clubs, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []ClubAthlete{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Retrieve recent activities from members of a specific club. The authenticated athlete must belong to the request club in order to hit this endpoint, Pagination is supported. Athlete profile
// visibility is respected for all activities.
// It uses a predefined ClubAthlete struct and not AthleteSummary because it currently sends only FirstName and LastName
func (s *ClubService) ListActivities(accessToken string, id int, p *RequestParams) ([]ClubActivity, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/activities", clubs, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []ClubActivity{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns of list of the athletes who are members of a given club.
func (s *ClubService) ListMembers(accessToken string, id int, p *RequestParams) ([]Member, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%d/members", clubs, id),
		Method:      http.MethodGet,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	resp := []Member{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Return a list of the clubs whose membership includes the authenticated athlete.
func (s *ClubService) ListAthleteClubs(accessToken string, p *RequestParams) ([]ClubSummary, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        fmt.Sprintf("%s/%s", athlete, clubs),
		AccessToken: accessToken,
		Method:      http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	resp := []ClubSummary{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
