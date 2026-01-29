package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"rawdg/config" // Custom package to load config
	"rawdg/sound"  // Custom package to play sound
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./rdg <duration>\nExample: ./rdg 1h2m3s")
		os.Exit(1)
	}

	// Load config from ~/.config/rdg/rdgconf.json
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(home, ".config", "rdg", "rdgconf.json")
	conf := config.LoadConfig(configPath)

	dur, err := time.ParseDuration(os.Args[1])
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	remaining := dur

	for remaining >= 0 {
		h := int(remaining.Hours())
		m := int(remaining.Minutes()) % 60
		s := int(remaining.Seconds()) % 60

		fmt.Print("\033[3J\033[H\033[2J")
		fmt.Printf("\r%02d:%02d:%02d", h, m, s)
		<-ticker.C
		remaining -= time.Second
	}

	fmt.Print("\033[3J\033[H\033[2J")
	fmt.Println("\rDone")

	if conf.SoundFilePath != "" {
		repetition := conf.Repetition
		if repetition <= 0 {
			repetition = 3 // temporary default
		}

		for range repetition {
			sound.PlaySound(conf.SoundFilePath)
		}
	}
}

