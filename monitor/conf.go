package monitor

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	GoesServiceUrl string `json:"goesServiceUrl"`
	TaskInterval   int    `json:"taskInterval"`
}

func ReadConfigFile(name string) (Configuration, error) {
	file, err := os.Open(name)
	if err != nil {
		return Configuration{}, err
	}
	defer file.Close()
	jsonData := Configuration{}
	err = json.NewDecoder(file).Decode(&jsonData)
	if err != nil {
		return Configuration{}, err
	}
	return jsonData, nil
}
