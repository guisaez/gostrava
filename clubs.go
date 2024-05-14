package go_strava

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// Returns a given club using its identifier
func (sc *StravaClient) GetClubById(ctx context.Context, id int64) (*DetailedClub, error) {

	path := fmt.Sprintf("/clubs/%d", id)

	var resp DetailedClub
	if err := sc.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Returns a list of the administrators of a given club.
func (sc *StravaClient) ListClubAdministrators(ctx context.Context, id int64, opt *GeneralParams) ([]SummaryAthlete, error) {

	path := fmt.Sprintf("/clubs/%d/admins", id)

	params := url.Values{}
	if opt != nil {
		if opt.Page > 0 {
			params.Set("page", strconv.Itoa(opt.Page))
		}
		if opt.Page > 0 {
			params.Set("per_page", strconv.Itoa(opt.PerPage))
		}
	}

	var resp []SummaryAthlete
	if err := sc.get(ctx, path, params, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Retrieve recent activities from members of a specific club. The authenticated athlete must belong to the request club in order to hit this endpoint, Pagination is supported. Athlete profile
// visibility is respected for all activities.
func (sc *StravaClient) ListClubActivities(ctx context.Context, id int64, opt *GeneralParams) ([]ClubActivity, error) {

	path := fmt.Sprintf("/clubs/%d/activities", id)

	params := url.Values{}
	if opt != nil {
		if opt.Page > 0 {
			params.Set("page", strconv.Itoa(opt.Page))
		}
		if opt.Page > 0 {
			params.Set("per_page", strconv.Itoa(opt.PerPage))
		}
	}

	var resp []ClubActivity
	if err := sc.get(ctx, path, params, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns of list of the athletes who are members of a given club.
func (sc *StravaClient) ListClubMembers(ctx context.Context, id int64, opt *GeneralParams) ([]ClubAthlete, error) {

	path := fmt.Sprintf("/clubs/%d/members", id)

	params := url.Values{}
	if opt != nil {
		if opt.Page > 0 {
			params.Set("page", strconv.Itoa(opt.Page))
		}
		if opt.Page > 0 {
			params.Set("per_page", strconv.Itoa(opt.PerPage))
		}
	}

	var resp []ClubAthlete
	if err := sc.get(ctx,path, params, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}