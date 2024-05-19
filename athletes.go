package gostrava

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type StravaAthletes baseModule

// Returns the currently authenticated athlete. Tokens with profile:read_all scope will receive
// a detailed athlete representation; all others will receive a SummaryAthlete representation
func (sc *StravaAthletes) CurrentAthlete(ctx context.Context, access_token string) (*DetailedAthlete, error) {

	var resp DetailedAthlete
	if err := sc.client.get(ctx, access_token, "/athlete", nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Returns the authenticated athlete's heart rate and power zones. Requires profile:read_all.
func (sc *StravaAthletes) GetZones(ctx context.Context, access_token string) (*Zones, error) {
	var resp Zones
	if err := sc.client.get(ctx, access_token, "/athlete/zones", nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Returns the activity stats of an athlete. Only includes data from activities set to Everyone's visibility.
func (sc *StravaAthletes) GetAthleteStats(ctx context.Context, access_token string, id int64) (*ActivityStats, error) {
	
	path := fmt.Sprintf("/athletes/%d/stats", id)

	var resp ActivityStats
	if err := sc.client.get(ctx, access_token, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Return a list of the clubs whose membership includes the authenticated athlete.
func (sc *StravaAthletes) ListClubs(ctx context.Context, access_token string) ([]SummaryClub, error) {

	var resp []SummaryClub
	if err := sc.client.get(ctx, access_token, "/athlete/clubs", nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a list of the routes created by the authenticated athlete. Private routes are filtered out
// unless request by a token with read_all scope.
func (sc *StravaAthletes) ListRoutes(ctx context.Context, access_token string) ([]Route, error) {

	var resp []Route
	if err := sc.client.get(ctx, access_token, "/athlete/routes", nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type UpdateAthleteReqParams struct {
	Weight float64 // The weigh of the athlete in kilograms.
}

// Update the currently authenticated athlete. Requires profile:write scope.
func (sc *StravaAthletes) Update(ctx context.Context, access_token string, p *UpdateAthleteReqParams) (*DetailedAthlete, error) {

	var params = url.Values{}

	if p != nil && p.Weight > 0 {
		params.Set("weight", fmt.Sprintf("%.2f", p.Weight))
	}

	var resp DetailedAthlete
	if err := sc.client.put(ctx, access_token, "/athlete", "application/x-www-form-urlencoded", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type ListAthleteActivitiesOptions struct {
	GeneralParams
	Before int // An epoch timestamp to use for filtering activities that have taken place before that certain time.
	After  int // An epoch timestamp to use for filtering activities that have taken place after a certain time.
}

// Returns the activities of an athlete for a specific identifier. Requires activity:read, OnlyMe activities will be filtered out unless
// requested by a token with activity_read:all.
func (sc *StravaAthletes) GetActivities(ctx context.Context, access_token string, opt *ListAthleteActivitiesOptions) ([]SummaryActivity, error) {

	params := url.Values{}
	if opt != nil {
		if opt.Page > 0 {
			params.Set("page_size", strconv.Itoa(opt.Page))
		}
		if opt.PerPage > 0 {
			params.Set("per_page", strconv.Itoa(opt.Page))
		}
		if opt.Before > 0 {
			params.Set("before", strconv.Itoa(opt.Before))
		}
		if opt.After > 0 {
			params.Set("after", strconv.Itoa(opt.After))
		}
	}

	var resp []SummaryActivity
	err := sc.client.get(ctx, access_token, "athlete/activities", params, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}