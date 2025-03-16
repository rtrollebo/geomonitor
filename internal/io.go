package internal

import (
	"encoding/json"
	"os"
)

func ReadFile[T any](name string) ([]T, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var jsonData []T
	err = json.NewDecoder(file).Decode(&jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func WriteFile[T any](data []T, name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(data)
	if err != nil {
		return err
	}
	return nil
}
