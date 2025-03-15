package geo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type GoesXray struct {
	TimeTag              time.Time `json:"time_tag"`
	Satellite            int       `json:"satellite"`
	Flux                 float64   `json:"flux"`
	ObservedFlux         float64   `json:"observed_flux"`
	ElectronCorrection   float64   `json:"electron_correction"`
	ElectronContaminaton bool      `json:"electron_contaminaton"`
	Energy               string    `json:"energy"`
}

func Run(url string, ctx context.Context) error {
	logError := ctx.Value("logerror").(*log.Logger)
	logInfo := ctx.Value("loginfo").(*log.Logger)

	events, readEventsEror := readFile("events.json")
	if readEventsEror != nil {
		logError.Println("Failed to read events file")
		return readEventsEror
	}
	goesXrayArray, getDataError := getGoesXrayData(url)
	if getDataError != nil {
		logError.Println("Network error, failed to fetch Goes data")
		return getDataError
	}
	if len(events) == 0 {
		events, msg, detectEventError := DetectEvent(goesXrayArray, time.Now().Add(time.Duration(-100)*time.Hour)) // 48 hours from now as initial value
		if detectEventError != nil {
			logError.Printf("Failed to detect events: %s", detectEventError)
			return detectEventError
		}
		if msg != "" {
			logInfo.Print(msg)
		}
		logInfo.Printf("No prev. events. New events detected: %d", len(events))
		writeError := writeFile(events, "events.json")
		if writeError != nil {
			logError.Println("Failed to write events file")
			return writeError
		}
	} else {
		logInfo.Printf("Previous events found: %d", len(events))
		updatedEvents := []GeoEvent{}
		oldEvents := []GeoEvent{}
		var timePeak time.Time
		var timeStart time.Time
		timePeak = time.Now().Add(time.Duration(-100) * time.Hour)
		var reventEvent GeoEvent

		for _, event := range events {
			if event.Time.After(timePeak) {
				timePeak = event.Time
				reventEvent = event
			}
		}
		// TODO: remove element that not processed ?
		if reventEvent.Processed {
			timeStart = reventEvent.TimeEnd
			oldEvents = events
		}
		if !reventEvent.Processed {
			timeStart = reventEvent.TimeStart
			oldEvents = events[0 : len(events)-1]
		}
		newEvents, msg, detectEventError := DetectEvent(goesXrayArray, timeStart)
		if detectEventError != nil {
			logError.Printf("Failed to detect events: %s", detectEventError)
			updatedEvents = events
		}
		if msg != "" {
			logInfo.Print(msg)
		}
		updatedEvents = append(oldEvents, newEvents...)
		writeError := writeFile(updatedEvents, "events.json")
		if writeError != nil {
			logError.Println("Failed to write events file")
			return writeError
		}
	}
	return nil
}

// IndexAt returns the index of the element in the array that is closest to the value
func IndexAt(array []GoesXray, value time.Time, tolerance int) (int, int) {
	length := len(array)
	if length == 0 || value.Before(array[0].TimeTag) || value.After(array[length-1].TimeTag) {
		return -1, 0
	}
	low, high := 0, length-1
	checks := 0
	tolDur := time.Duration(tolerance) * time.Minute
	for low <= high {
		checks++
		mid := (low + high) / 2
		midTime := array[mid].TimeTag
		if value.After(midTime.Add(-tolDur)) && value.Before(midTime.Add(tolDur)) {
			return mid, checks
		}
		if value.After(midTime) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1, checks
}

// Estimation of Full width half maximum
func GetFwhfIndices(array []GoesXray, peakIdx int) (int, int) {
	length := len(array)
	if peakIdx < 0 || peakIdx >= length {
		return -1, -1
	}
	peak := array[peakIdx].Flux
	half := peak / 2
	low, high := peakIdx, peakIdx
	for low > 0 && array[low].Flux > half {
		low--
	}
	for high < length-1 && array[high].Flux > half {
		high++
	}
	halfCutoff := half + (0.1 * half)
	if array[low].Flux > halfCutoff {
		low = -1
	}
	if array[high].Flux > halfCutoff {
		high = -1
	}
	return low, high
}

func DetectEvent(array []GoesXray, timeFrom time.Time) ([]GeoEvent, string, error) {
	cutoff := XRAY_FLARE_M
	startTime := timeFrom
	msg := ""
	if array[0].TimeTag.After(timeFrom) {
		startTime = array[0].TimeTag
	}
	indx, _ := IndexAt(array, startTime, 1)
	if indx == -1 {
		return nil, "", errors.New("Failed to find index")
	}
	var prev GoesXray
	rate := 0.0
	events := []GeoEvent{}
	for _, v := range array[indx:] {
		if prev.Flux > v.Flux && rate > 0 {
			var newEvent GeoEvent
			if prev.Flux < 1.0e-7 {
				newEvent = GeoEvent{Time: prev.TimeTag, Event: XRAY_FLUX_CHANGED, Description: "A Xray event", Value: float32(prev.Flux), Cat: XRAY_FLARE_A}
			}
			if prev.Flux > 1.0e-7 && prev.Flux < 1.0e-6 {
				newEvent = GeoEvent{Time: prev.TimeTag, Event: XRAY_FLUX_CHANGED, Description: "B Xray event", Value: float32(prev.Flux), Cat: XRAY_FLARE_B}
			}
			if prev.Flux > 1.0e-6 && prev.Flux < 1.0e-5 {
				newEvent = GeoEvent{Time: prev.TimeTag, Event: XRAY_FLUX_CHANGED, Description: "C Xray event", Value: float32(prev.Flux), Cat: XRAY_FLARE_C}
			}
			if prev.Flux > 1.0e-5 && prev.Flux < 1.0e-4 {
				newEvent = GeoEvent{Time: prev.TimeTag, Event: XRAY_FLUX_CHANGED, Description: "M Xray event", Value: float32(prev.Flux), Cat: XRAY_FLARE_M}
			}
			if prev.Flux > 1.0e-4 {
				newEvent = GeoEvent{Time: prev.TimeTag, Event: XRAY_FLUX_CHANGED, Description: "X Xray event", Value: float32(prev.Flux), Cat: XRAY_FLARE_X}
			}

			if newEvent.Cat >= int8(cutoff) {
				events = append(events, newEvent)
			}
		}
		rate = v.Flux - prev.Flux
		prev = v
	}
	msg = fmt.Sprintf("New events detected: %d", len(events))
	updatedEvents := []GeoEvent{}
	for _, event := range events {
		updatedEvent, _, updateErr := updateEvent(array, event)
		if updateErr != nil {
			updatedEvents = append(updatedEvents, event)
		} else {
			updatedEvents = append(updatedEvents, updatedEvent)
		}
	}
	return updatedEvents, msg, nil
}

func updateEvent(array []GoesXray, event GeoEvent) (GeoEvent, string, error) {
	index, _ := IndexAt(array, event.Time, 2)
	fwhmLower, fwhmUpper := GetFwhfIndices(array, index)

	if fwhmLower == -1 && fwhmUpper == -1 {
		return GeoEvent{}, "", errors.New("Failed to find fwhm")
	} else if fwhmLower != -1 && fwhmUpper != -1 {

		return GeoEvent{Event: event.Event, Time: event.Time, TimeStart: array[fwhmLower].TimeTag, TimeEnd: array[fwhmUpper].TimeTag, Cat: event.Cat, Value: event.Value, Processed: true, Description: fmt.Sprintf("Duration (mins): %d", (fwhmUpper-fwhmLower)/2)}, "Ongoing event closed", nil
	} else {
		// Unprocessed event still ongoing
		return GeoEvent{Event: event.Event, Time: event.Time, Cat: event.Cat, Value: event.Value, Processed: false, Description: "Ongoing event"}, "Ongoing event", nil
	}
}

func getGoesXrayData(url string) ([]GoesXray, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, err
	}
	body, _ := io.ReadAll(res.Body)
	var jsonData []GoesXray
	errunmarshal := json.Unmarshal(body, &jsonData)
	if errunmarshal != nil {
		return nil, err
	}
	return jsonData, nil

}

func writeFile(events []GeoEvent, name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(events)
	if err != nil {
		return err
	}
	return nil
}

func readFile(name string) ([]GeoEvent, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	jsonData := []GeoEvent{}
	err = json.NewDecoder(file).Decode(&jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
