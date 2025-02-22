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
	goesxraylist, err := getGoesXrayData(url)
	if err != nil {
		logError.Println("Network error, failed to fetch Goes data")
		return err
	}

	e, err := processGeosArray(goesxraylist, 10) // last 5 minutes

	logInfo.Println(fmt.Sprintf("Event detected: %s, %d, %d, %f ", e.Time, e.Event, e.Cat, e.Value))
	if e.Event == XRAY_FLUX_CHANGED && e.Cat > 2 {
		var currentEvents []GeoEvent
		currentEvents, err = readFile("events.json")
		if err != nil {
			logError.Println("Failed to read events file")
			return err
		}
		currentEvents = append(currentEvents, e)
		err = writeFile(currentEvents, "events.json")

		if err != nil {
			logError.Println("Failed to write events file")
			return err
		}
	}
	return nil
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

func processGeosArray(geoslist []GoesXray, index_stop int) (GeoEvent, error) {
	var goes_xray_max_flux GoesXray
	idxlast := len(geoslist) - 1
	for _, v := range geoslist[(idxlast - index_stop):idxlast] {
		if v.Flux > goes_xray_max_flux.Flux {
			goes_xray_max_flux = v
		}
	}
	if goes_xray_max_flux.Flux < 1.0e-7 {
		return GeoEvent{Time: goes_xray_max_flux.TimeTag, Event: XRAY_FLUX_CHANGED, Description: "A Xray event", Value: float32(goes_xray_max_flux.Flux), Cat: XRAY_FLARE_A}, nil
	}
	if goes_xray_max_flux.Flux > 1.0e-7 && goes_xray_max_flux.Flux < 1.0e-6 {
		return GeoEvent{Time: goes_xray_max_flux.TimeTag, Event: XRAY_FLUX_CHANGED, Description: "B Xray event", Value: float32(goes_xray_max_flux.Flux), Cat: XRAY_FLARE_B}, nil
	}
	if goes_xray_max_flux.Flux > 1.0e-6 && goes_xray_max_flux.Flux < 1.0e-5 {
		return GeoEvent{Time: goes_xray_max_flux.TimeTag, Event: XRAY_FLUX_CHANGED, Description: "C Xray event", Value: float32(goes_xray_max_flux.Flux), Cat: XRAY_FLARE_C}, nil
	}
	if goes_xray_max_flux.Flux > 1.0e-5 && goes_xray_max_flux.Flux < 1.0e-4 {
		return GeoEvent{Time: goes_xray_max_flux.TimeTag, Event: XRAY_FLUX_CHANGED, Description: "M Xray event", Value: float32(goes_xray_max_flux.Flux), Cat: XRAY_FLARE_M}, nil
	}
	if goes_xray_max_flux.Flux > 1.0e-4 {
		return GeoEvent{Time: goes_xray_max_flux.TimeTag, Event: XRAY_FLUX_CHANGED, Description: "X Xray event", Value: float32(goes_xray_max_flux.Flux), Cat: XRAY_FLARE_X}, nil
	}
	return GeoEvent{}, errors.New(fmt.Sprintf("Unknown category: %f", goes_xray_max_flux.Flux))

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
