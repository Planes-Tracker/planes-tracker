package types

import (
	"encoding/json"
	"os"
)

type Config struct {
	DatabaseName    string      `json:"databaseName"`
	PollRate        int         `json:"pollRate"`
	Debug           bool        `json:"debug"`

	Location        Coordinates `json:"location"`
	Radius          Radius      `json:"radius"`
}

func NewConfigFromFile(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}



	return &config, nil
}
