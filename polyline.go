package gostrava

type PolylineSummary struct {
	ID              string `json:"id"`
	SummaryPolyline string `json:"summary_polyline"`
	ResourceState   int8   `json:"resource_state"`
}

type PolylineDetailed struct {
	PolylineSummary
	Polyline string `json:"polyline"`
}
