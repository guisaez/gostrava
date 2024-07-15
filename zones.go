package gostrava

type Zones struct {
	HearRate HeartRateZoneRanges `json:"heart_rate"` // An instance of HeartRateZoneRanges.
	Power    PowerZoneRanges     `json:"power"`      // An instance of PowerZoneRanges.
}

type HeartRateZoneRanges struct {
	CustomZones bool        `json:"custom_zone"` // Whether the athlete has set their own custom heart rate zones
	Zones       []ZoneRange `json:"zones"`       // An instance of ZoneRanges.
}

type PowerZoneRanges struct {
	Zones []ZoneRange `json:"zones"` // An instance of ZoneRanges.
}

type ZoneRange struct {
	Max int `json:"max"` // The maximum value in the range.
	Min int `json:"min"` // The minimum value in the range.
}

// A union type representing the time spent in a given zone.
type TimedZoneRange struct {
	Min  int `json:"min"`  // The minimum value in the range.
	Max  int `json:"max"`  // The maximum value in the range.
	Time int `json:"time"` // The number of seconds spent in this zone
}
