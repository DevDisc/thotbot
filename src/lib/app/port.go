package app

import (
	"encoding/json"
	"io/ioutil"
)

type Port struct {
	Holdings map[string]int `json:"holdings"`
}

func LoadPort(path string) (*Port, error) {
	jsonBits, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var port Port
	err = json.Unmarshal(jsonBits, &port)
	if err != nil {
		return nil, err
	}
	return &port, nil
}

func SavePort(port *Port, path string) error {
	jsonData, err := json.Marshal(port)
	if err != nil {
		return err
	}

	// Overwrite existing
	err = ioutil.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
