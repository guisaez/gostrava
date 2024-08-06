package gostrava

import (
	"context"
	"fmt"
	"net/http"
)

// *************** Types ********************

// A set of rolled-up statistics and totals for an athlete
type AthleteStats struct {
	RecentRideTotals     ActivityTotal `json:"recent_ride_totals"` // The recent (last 4 weeks) ride stats for the athlete.
	AllRideTotals        ActivityTotal `json:"all_ride_totals"`    // The all time ride stats for the athlete.
	RecentRunTotals      ActivityTotal `json:"recent_run_totals"`  // The recent (last 4 weeks) run stats for the athlete.
	AllRunTotals         ActivityTotal `json:"all_run_totals"`     // The all time run stats for the athlete.
	RecentSwimTotals     ActivityTotal `json:"recent_swim_totals"` // The recent (last 4 weeks) swim stats for the athlete.
	AllSwimTotals        ActivityTotal `json:"all_swim_totals"`    // The all time swim stats for the athlete.
	YearToDateRideTotals ActivityTotal `json:"ytd_ride_totals"`    // The year to date ride stats for the athlete.
	YearToDateSwimTotals ActivityTotal `json:"ytd_swim_totals"`    // The year to date swim stats for the athlete.
	YearToDateRunTotals  ActivityTotal `json:"ytd_run_totals"`     // The year to date run stats for the athlete.
}

// A roll-up of metrics pertaining to a set of activities. Values are in seconds and meters.
type ActivityTotal struct {
	Count         int     `json:"count"`          // The number of activities considered in this total.
	Distance      float32 `json:"distance"`       // The total distance covered by the considered activities.
	MovingTime    int     `json:"moving_time"`    // The total moving time of the considered activities.
	ElapsedTime   int     `json:"elapsed_time"`   // The total elapsed time of the considered activities.
	ElevationGain float32 `json:"elevation_gain"` // The total elevation gain of the considered activities.

	// only present for recent totals, not ytd, or all
	AchievementCount int `json:"achievement_count"` // The total number of achievements of the considered activities.
}

// *************** Methods ********************

type AthletesService service

const athletes = "/api/v3/athletes"

// GetAthleteStats retrives the activities stats of an athlete. Only includes data
// from activities set to Everyone's visibility
//
// GET: https://www.strava.com/api/v3/athletes/{id}/stats
func (s *AthletesService) GetAthleteStats(ctx context.Context, accessToken string, id int) (*AthleteStats, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/stats", athletes, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	stats := new(AthleteStats)
	resp, err := s.client.DoAndParse(ctx, req, stats)
	if err != nil {
		return nil, resp, err
	}

	return stats, resp, nil
}

// GetAthleteRoutes retrieves the routes created by a specified athlete.
//
// GET: https://www.strava.com/api/v3/athletes/{id}/routes
func (s *AthletesService) GetAthleteRoutes(ctx context.Context, accessToken string, id int) ([]RouteSummary, *http.Response, error) {
	urlStr := fmt.Sprintf("%s/%d/routes", athletes, id)

	req, err := s.client.NewRequest(http.MethodGet, urlStr, nil, SetAuthorizationHeader(accessToken))
	if err != nil {
		return nil, nil, err
	}

	var routes []RouteSummary
	resp, err := s.client.DoAndParse(ctx, req, &routes)
	if err != nil {
		return nil, resp, nil
	}

	return routes, resp, nil
}
