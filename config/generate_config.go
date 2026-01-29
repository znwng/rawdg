package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func GenerateConfig() error {
	// Determine user config directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get user home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "rdg")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	soundFilesDir := filepath.Join(home, ".config", "rdg", "sound_files")
	if err := os.MkdirAll(soundFilesDir, 0755); err != nil {
		return fmt.Errorf("cannot create soudn_files directory: %w", err)
	}

	// Default config with empty sound file path
	cfg := TimerConfig{
		SoundFilePath: "",
		Repetition:    1,
	}

	// Write config file
	configPath := filepath.Join(configDir, "rdgconf.json")
	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("cannot create config file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ") // pretty print
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("cannot encode config: %w", err)
	}

	return nil
}

