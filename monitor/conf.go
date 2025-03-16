package monitor

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	GoesServiceUrl    string   `json:"goesServiceUrl"`
	TaskInterval      int      `json:"taskInterval"`
	NotifySender      string   `json:"notifySender"`
	NotifyRecipients  []string `json:"notifyRecipients"`
	NotifySmtpAddress string   `json:"notifySmtpAddress"`
	NotifySmtpPort    string   `json:"notifySmtpPort"`
	NotifySmtpPass    string   `json:"notifySmtpPass"`
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
