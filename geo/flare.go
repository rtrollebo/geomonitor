package geo

import (
	"time"
)

type SolarFlare struct {
	TimeFwhm time.Time
}

// Get time of full width half max entry
func getTimeFwhmEntry(array []float64, backgroundFlux float64) {
	// to be implemented
}
