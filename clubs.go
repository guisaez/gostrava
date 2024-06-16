package gostrava

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	clubsPath = "/clubs"
)

type ClubsAPIService apiService

type ClubRequestParams generalParams

// Returns a given club using its identifier
func (c *ClubsAPIService) GetById(access_token string, id int64) (*DetailedClub, error) {
	requestUrl := c.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id))

	req, err := c.client.get(requestUrl, nil, access_token)
	if err != nil {
		return nil, err
	}

	resp := &DetailedClub{}
	if err := c.client.do(req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a list of the administrators of a given club.
func (c *ClubsAPIService) GetAdministrators(access_token string, id int64, opts *ClubRequestParams) ([]SummaryAthlete, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			params.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(opts.PerPage))
		}
	}

	request_url := c.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id), "/admins")

	req, err := c.client.get(request_url, params, access_token)
	if err != nil {
		return nil, err
	}

	resp := []SummaryAthlete{}
	if err := c.client.do(req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Retrieve recent activities from members of a specific club. The authenticated athlete must belong to the request club in order to hit this endpoint, Pagination is supported. Athlete profile
// visibility is respected for all activities.
func (c *ClubsAPIService) GetActivities(access_token string, id int64, opts *ClubRequestParams) ([]ClubActivity, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			params.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(opts.PerPage))
		}
	}

	request_url := c.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id), "/activities")

	req, err := c.client.get(request_url, params, access_token)
	if err != nil {
		return nil, err
	}

	resp := []ClubActivity{}
	if err := c.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns of list of the athletes who are members of a given club.
func (c *ClubsAPIService) GetMembers(access_token string, id int64, opts *ClubRequestParams) ([]ClubAthlete, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			params.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(opts.PerPage))
		}
	}

	request_url := c.client.BaseURL.JoinPath(clubsPath, fmt.Sprint(id), "/members")

	req, err := c.client.get(request_url, params, access_token)
	if err != nil {
		return nil, err
	}

	resp := []ClubAthlete{}
	if err := c.client.do(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
