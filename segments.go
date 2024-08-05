package gostrava

type SegmentSummary struct {
	ID              int                     `json:"id"`                          // The unique identifier of this segment
	ResourceState   int8                    `json:"resource_state"`              //
	Name            string                  `json:"name"`                        // The name of this segment
	ActivityType    string                  `json:"activity_type"`               // May take one of the following values: "Ride" / "Run"
	Distance        float32                 `json:"distance"`                    // The segment's distance, in meters
	AvgGrade        float32                 `json:"average_grade"`               // The segment's average grade, in percents
	MaxGrade        float32                 `json:"maximum_grade"`               // The segments's maximum grade, in percents
	ElevationHigh   float32                 `json:"elevation_high"`              // The segments's highest elevation, in meters
	ElevationLow    float32                 `json:"elevation_low"`               // The segments's lowest elevation, in meters
	StartLatLng     LatLng                  `json:"start_latlng"`                // An instance of LatLng.
	EndLatLng       LatLng                  `json:"end_latlng"`                  // An instance of LatLng.
	ClimbCategory   int8                    `json:"climb_category"`              // The category of the climb [0, 5]. Higher is harder ie. 5 is Hors catégorie, 0 is uncategorized in climb_category.
	City            string                  `json:"city"`                        // The segments's city.
	State           string                  `json:"state"`                       // The segments's state or geographical region.
	Country         string                  `json:"country"`                     // The segment's country.
	Private         bool                    `json:"private"`                     // Whether this segment is private.
	Starred         bool                    `json:"starred"`                     // Whether this segment has been starred by the current athlete.
	AthletePREffort *SummaryPRSegmentEffort `json:"athlete_pr_effort,omitempty"` // An instance of SummaryPRSegmentEffort.
}

type SegmentDetailed struct {
	SegmentSummary
	Hazardous          bool      `json:"hazardous"`            // Whether this segment is considered hazardous
	TotalElevationGain float32   `json:"total_elevation_gain"` // The segment's total elevation gain.
	CreatedAt          TimeStamp `json:"created_at"`           // The time at which the segment was created.
	UpdatedAt          TimeStamp `json:"updated_at"`           // The time at which the segment was last updated.

	ElevationProfile  string `json:"elevation_profile"`
	ElevationProfiles URL    `json:"elevation_profiles"`

	Map          *PolylineSummmary `json:"map,omitempty"`          // An instance of PolylineMap.
	EffortCount  int               `json:"effort_count"`           // The total number of efforts for this segment
	AthleteCount int               `json:"athlete_count"`          // The number of unique athletes who have an effort for this segment
	StarCount    int               `json:"star_count"`             // The number of stars for this segment
	Xoms         *Xoms             `json:"xoms,omitempty"`         //
	LocalLegend  *LocalLegend      `json:"local_legend,omitempty"` //

	// AthleteSegmentStats *SegmentSummaryEffort `json:"athlete_segment_stats,omitempty"` // An instance ofSegmentSummaryEffort.
}

type SummaryPRSegmentEffort struct {
	PRActivityID  int       `json:"pr_activity_id"`  // The unique identifier of the activity related to the PR effort.
	PRElapsedTime int       `json:"pr_elapsed_time"` // The elapsed time ot the PR effort.
	PRDate        TimeStamp `json:"pr_date"`         //  The time at which the PR effort was started.
	EffortCount   int       `json:"effort_count"`    // Number of efforts by the authenticated athlete on this segment.
}

type LocalLegend struct {
	AthleteID         int    `json:"athlete_id"`
	Title             string `json:"title"`
	Profile           string `json:"profile"`
	EffortDescription string `json:"effort_description"`
	EffortCount       string `json:"effort_count"`
	Destination       string `json:"destination"`
}

type Xoms struct {
	Kom         string `json:"kom"`
	Qom         string `json:"qom"`
	Overall     string `json:"overall"`
	Destination struct {
		Href string `json:"href"`
		Type string `json:"type"`
		Name string `json:"name"`
	}
}

type ExplorerResponse struct {
	Segments []*ExplorerSegment `json:"segments"` // The set of segments matching an explorer request
}

type ExplorerSegment struct {
	ID                 int     `json:"id"`                  // The unique identifier of this segment
	ResourceState      int8    `json:"resource_state"`      //
	Name               string  `json:"name"`                // The name of this segment
	ClimbCategory      int8    `json:"climb_category"`      // The category of the climb [0, 5]. Higher is harder ie. 5 is Hors catégorie, 0 is uncategorized in climb_category. If climb_category = 5, climb_category_desc = HC. If climb_category = 2, climb_category_desc = 3.
	ClimbCategoryDesc  string  `json:"climb_category_desc"` // The description for the category of the climb May take one of the following values: NC, 4, 3, 2, 1, HC
	AvgGrade           float32 `json:"avg_grade"`           // The segment's average grade, in percents
	StartLatLng        LatLng  `json:"start_latlng"`        // An instance of LatLng.
	EndLatLng          LatLng  `json:"end_latlng"`          // An instance of LatLng.
	ElevationDiff      float32 `json:"elev_difference"`     // The segments's elevation difference, in meters
	Distance           float32 `json:"distance"`            // The segment's distance, in meters
	Points             string  `json:"points"`              // The polyline of the segment
	Starred            bool    `json:"starred"`
	ElevationProfile   string  `json:"elevation_profile"`
	LocalLegendEnabled bool    `json:"local_legend_enabled"`
}
