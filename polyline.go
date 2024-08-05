package gostrava

type PolylineSummmary struct {
	ID              string `json:"id"`
	SummaryPolyline string `json:"summary_polyline"`
	ResourceState   int8   `json:"resource_state"`
}

type PolylineDetailed struct {
	PolylineSummmary
	Polyline string `json:"polyline"`
}
