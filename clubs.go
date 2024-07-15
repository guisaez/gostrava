package gostrava

import (
	"net/url"
	"strconv"
)

type ClubService service

// Returns a given club using its identifier
func (s *ClubService) GetById(accessToken string, id int) (*ClubDetailed, error) {
	req, err := s.client.newRequest(requestOpts{
		Path:        "clubs/" + strconv.Itoa(id),
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
func (s *ClubService) ListAdministrators(accessToken string, id int, opts RequestParams) ([]ClubAthlete, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "clubs/" + strconv.Itoa(id) + "/admins",
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
func (s *ClubService) ListActivities(accessToken string, id int, opts RequestParams) ([]ClubActivity, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "clubs/" + strconv.Itoa(id) + "/activities",
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
func (s *ClubService) ListMembers(accessToken string, id int, opts RequestParams) ([]Member, error) {
	params := url.Values{}

	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params.Set("per_page", strconv.Itoa(opts.PerPage))
	}

	req, err := s.client.newRequest(requestOpts{
		Path:        "clubs/" + strconv.Itoa(id) + "/members",
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


