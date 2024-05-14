package go_strava

import (
	"context"
	"fmt"
	"net/url"
)

// Returns the currently authenticated athlete. Tokens with profile:read_all scope will receive
// a detailed athlete representation; all others will receive a SummaryAthlete representation
func (sc *StravaClient) GetAuthenticatedAthlete(ctx context.Context) (*DetailedAthlete, error) {

	var resp DetailedAthlete
	if err := sc.get(ctx,"/athlete", nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Returns the authenticated athlete's heart rate and power zones. Requires profile:read_all.
func (sc *StravaClient) GetZones(ctx context.Context) (*Zones, error) {
	
	var resp Zones
	if err := sc.get(ctx, "/athlete", nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Returns the activity stats of an athlete. Only includes data from activities set to Everyone's visibility.
func (sc *StravaClient) GetAthleteStats(ctx context.Context, id int64) (*ActivityStats, error) {
	
	path := fmt.Sprintf("/athlete/%d/stats", id)

	var resp ActivityStats
	if err := sc.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Return a list of the clubs whose membership includes the authenticated athlete.
func (sc *StravaClient) ListAthleteClubs(ctx context.Context) ([]SummaryClub, error) {

	var resp []SummaryClub
	if err := sc.get(ctx, "/athlete/clubs", nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns a list of the routes created by the authenticated athlete. Private routes are filtered out
// unless request by a token with read_all scope.
func (sc *StravaClient) ListAthleteRoutes(ctx context.Context) ([]Route, error) {

	var resp []Route
	if err := sc.get(ctx, "/athlete/routes", nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}


type UpdateAthleteReqParams struct {
	Weight float64 // The weigh of the athlete in kilograms.
}

// Update the currently authenticated athlete. Requires profile:write scope.
func (sc *StravaClient) UpdateAthlete(ctx context.Context, p *UpdateAthleteReqParams) (*DetailedAthlete, error) {

	var params = url.Values{}

	if p != nil && p.Weight > 0 {
		params.Set("weight", fmt.Sprintf("%.2f", p.Weight))
	}

	var resp DetailedAthlete
	if err := sc.put(ctx, "/athlete", "application/x-www-form-urlencoded", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}


