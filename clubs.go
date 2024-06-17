package gostrava

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ClubsAPIService apiService

// Returns a given club using its identifier
func (s *ClubsAPIService) GetById(access_token string, id int64) (*DetailedClub, error) {
	requestUrl := s.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id))

	req, err := s.client.newRequest(clientRequestOpts{
		url: requestUrl,
		method: http.MethodGet,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := &DetailedClub{}
	if err := s.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a list of the administrators of a given club.
func (s *ClubsAPIService) ListClubAdministrators(access_token string, id int64, p *RequestParams) ([]SummaryAthlete, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	requestUrl := s.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id), "/admins")

	req, err := s.client.newRequest(clientRequestOpts{
		url: requestUrl,
		method: http.MethodGet,
		body: params,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := []SummaryAthlete{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Retrieve recent activities from members of a specific club. The authenticated athlete must belong to the request club in order to hit this endpoint, Pagination is supported. Athlete profile
// visibility is respected for all activities.
func (s *ClubsAPIService) ListClubActivities(access_token string, id int64, p *RequestParams) ([]ClubActivity, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	requestUrl := s.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id), "/activities")

	req, err := s.client.newRequest(clientRequestOpts{
		url: requestUrl,
		method: http.MethodGet,
		body: params,
		access_token: access_token,
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
func (s *ClubsAPIService) ListClubMembers(access_token string, id int64, p *RequestParams) ([]ClubAthlete, error) {
	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	requestUrl := s.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id), "/members")

	req, err := s.client.newRequest(clientRequestOpts{
		url: requestUrl,
		method: http.MethodGet,
		body: params,
		access_token: access_token,
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

// Return a list of the clubs whose membership includes the authenticated athlete.
func (s *ClubsAPIService) ListAthleteClubs(access_token string, p *RequestParams) ([]SummaryClub, error) {
	requestUrl := s.client.BaseURL.JoinPath(athletePath, clubsPath)

	params := url.Values{}
	if p != nil {
		if p.Page > 0 {
			params.Set("page", strconv.Itoa(p.Page))
		}
		if p.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(p.PerPage))
		}
	}

	req, err := s.client.newRequest(clientRequestOpts{
		url:          requestUrl,
		method:       http.MethodGet,
		body:         params,
		access_token: access_token,
	})
	if err != nil {
		return nil, err
	}

	resp := []SummaryClub{}
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
