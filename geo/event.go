package geo

import (
	"time"
)

const (
	XRAY_FLUX_CHANGED = iota
	SOLAR_WIND_SPEED_CHANGED
	GROUND_FIELD_CHANGED
)

const (
	XRAY_FLARE_A = iota
	XRAY_FLARE_B
	XRAY_FLARE_C
	XRAY_FLARE_M
	XRAY_FLARE_X
)

type GeoEvent struct {
	Time        time.Time
	TimeStart   time.Time
	TimeEnd     time.Time
	Event       int8
	Cat         int8
	Value       float32
	Processed   bool
	Description string
}

type GeoEventReport struct {
	Events []GeoEvent
}
