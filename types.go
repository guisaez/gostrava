package gostrava

import (
	"encoding/json"
	"errors"
	"time"
)

type TimeStamp struct {
	Time time.Time
}

func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	var timeStr string
	err := json.Unmarshal(data, &timeStr)
	if err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return err
	}

	t.Time = parsedTime

	return nil
}

type ActivityType string

const (
	AlpineSki       ActivityType = "AlpineSki"
	BackgroundSki   ActivityType = "BackgroundSki"
	Canoeing        ActivityType = "Canoeing"
	Crossfit        ActivityType = "Crossfit"
	EBikeRide       ActivityType = "EBikeRide"
	Elliptical      ActivityType = "Elliptical"
	Golf            ActivityType = "Golf"
	Handcycle       ActivityType = "Handcycle"
	Hike            ActivityType = "Hike"
	IceSkate        ActivityType = "IceSkate"
	InlineSkate     ActivityType = "InlineSkate"
	Kayaking        ActivityType = "Kayaking"
	Kitesurf        ActivityType = "Kitesurf"
	NordicSki       ActivityType = "NordicSki"
	Ride            ActivityType = "Ride"
	RockClimbing    ActivityType = "RockClimbing"
	RollerSki       ActivityType = "RollerSki"
	Rowing          ActivityType = "Rowing"
	Run             ActivityType = "Run"
	Sail            ActivityType = "Sail"
	SkateBoard      ActivityType = "SkateBoard"
	Snowboard       ActivityType = "Snowboard"
	Snowshoe        ActivityType = "Snowshoe"
	Soccer          ActivityType = "Soccer"
	StairStepper    ActivityType = "StairStepper"
	StandupPaddling ActivityType = "StandupPaddling"
	Surfing         ActivityType = "Surfing"
	Swim            ActivityType = "Swim"
	Velomobile      ActivityType = "Velomobile"
	VirtualRide     ActivityType = "VirtualRide"
	VirtualRun      ActivityType = "VirtualRun"
	Walk            ActivityType = "Walk"
	WeightTraining  ActivityType = "WeightTraining"
	Wheelchair      ActivityType = "Wheelchair"
	Windsurf        ActivityType = "Windsurf"
	Workout         ActivityType = "Workout"
	Yoga            ActivityType = "Yoga"
)

type ClubSportType string

const (
	Cycling   ClubSportType = "cycling"
	Running   ClubSportType = "running"
	Triathlon ClubSportType = "triathlon"
	Other     ClubSportType = "other"
)

type RouteType string

const (
	RideRoute RouteType = "ride"
	RunRoute  RouteType = "run"
)

func (rt *RouteType) UnmarshalJSON(data []byte) error {
	var routeType int
	err := json.Unmarshal(data, &routeType)
	if err != nil {
		return err
	}

	switch routeType {
	case 1:
		*rt = RideRoute
	case 2:
		*rt = RunRoute
	default:
		return errors.New("invalid route type")
	}

	return nil
}

type SegmentActivityType string

const (
	RideSegment SegmentActivityType = "Ride"
	RunSegment  SegmentActivityType = "Run"
)

type SubRouteType string

const (
	Road         SubRouteType = "road"
	MountainBike SubRouteType = "mountain_bike"
	Cross        SubRouteType = "cross"
	Trail        SubRouteType = "train"
	Mixed        SubRouteType = "mixed"
)

func (rt *SubRouteType) UnmarshalJSON(data []byte) error {
	var subRouteType int
	err := json.Unmarshal(data, &subRouteType)
	if err != nil {
		return err
	}

	switch subRouteType {
	case 1:
		*rt = Road
	case 2:
		*rt = MountainBike
	case 3:
		*rt = Cross
	case 4:
		*rt = Trail
	case 5:
		*rt = Mixed
	default:
		return errors.New("invalid sub-route type")
	}

	return nil
}

type SportType string

const (
	AlpineSkiSport                     SportType = "AlpineSki"
	BackcountrySkiSport                SportType = "BackcountrySki"
	BadmintonSport                     SportType = "Badminton"
	CanoeingSport                      SportType = "Canoeing"
	CrossfitSport                      SportType = "Crossfit"
	EBikeRideSport                     SportType = "EBikeRide"
	EllipticalSport                    SportType = "Elliptical"
	EMountainBikeRideSport             SportType = "EMountainBikeRide"
	GolfSport                          SportType = "Golf"
	GravelRideSport                    SportType = "GravelRide"
	HandcycleSport                     SportType = "Handcycle"
	HighIntensityIntervalTrainingSport SportType = "HighIntensityIntervalTraining"
	HikeSport                          SportType = "Hike"
	IceSkateSport                      SportType = "IceSkate"
	InlineSkateSport                   SportType = "InlineSkate"
	KayakingSport                      SportType = "Kayaking"
	KitesurfSport                      SportType = "Kitesurf"
	MountainBikeRideSport              SportType = "MountainBikeRide"
	NordicSkiSport                     SportType = "NordicSki"
	PickleballSport                    SportType = "Pickleball"
	PilatesSport                       SportType = "Pilates"
	RacquetballSport                   SportType = "Racquetball"
	RideSport                          SportType = "Ride"
	RockClimbingSport                  SportType = "RockClimbing"
	RollerSkiSport                     SportType = "RollerSki"
	RowingSport                        SportType = "Rowing"
	RunSport                           SportType = "Run"
	SailSport                          SportType = "Sail"
	Skateboard                         SportType = "Skateboard"
	SnowboardSport                     SportType = "Snowboard"
	SnowshoeSport                      SportType = "Snowshoe"
	SoccerSport                        SportType = "Soccer"
	SquashSport                        SportType = "Squash"
	StairStepperSport                  SportType = "StairStepper"
	StandUpPaddlingSport               SportType = "StandUpPaddling"
	SurfingSport                       SportType = "Surfing"
	SwimSport                          SportType = "Swim"
	TableTennisSport                   SportType = "TableTennis"
	TennisSport                        SportType = "Tennis"
	TrailRunSport                      SportType = "TrailRun"
	VelomobileSport                    SportType = "Velomobile"
	VirtualRideSport                   SportType = "VirtualRide"
	VirtualRowSport                    SportType = "VirtualRow"
	VirtualRunSport                    SportType = "VirtualRun"
	WalkSportType                      SportType = "Walk"
	WeightTrainingSport                SportType = "WeightTraining"
	WheelchairSport                    SportType = "Wheelchair"
	WindsurfSport                      SportType = "Windsurf"
	WorkoutSport                       SportType = "Workout"
	YogaSport                          SportType = "Yoga"
)
