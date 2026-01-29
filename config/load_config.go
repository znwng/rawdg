package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type TimerConfig struct {
	SoundFilePath string `json:"sound_file_path"`
	Repetition    int    `json:"repetition"`
}

func LoadConfig(config_file_path string) TimerConfig {
	// If config file doesn't exist, generate it
	if _, err := os.Stat(config_file_path); os.IsNotExist(err) {
		if err := GenerateConfig(); err != nil {
			panic(fmt.Sprintf("Failed to generate config: %v", err))
		}
	}
	f, err := os.Open(config_file_path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	var cfg TimerConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	return cfg
}

