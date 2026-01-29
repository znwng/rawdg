package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"rawdg/config"
	"rawdg/sound"

	"golang.org/x/term"
)

// Function to get terminal width
func getTerminalWidth() int {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return 80
	}
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

// Get progress bar width
func progressBarWidth() int {
	reserved := len("00:00:00 []")
	w := getTerminalWidth() - reserved
	if w < 10 {
		return 10
	}
	return w
}

func renderProgressBar(elapsed, total time.Duration, width int) string {
	if total <= 0 {
		return strings.Repeat("█", width)
	}
	progress := float64(elapsed) / float64(total)
	if progress > 1 {
		progress = 1
	}
	filled := int(progress * float64(width))
	return strings.Repeat("█", filled) + strings.Repeat(" ", width-filled)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./rdg <duration>")
		fmt.Println("Example: ./rdg 1h2m3s")
		os.Exit(1)
	}

	total, err := time.ParseDuration(os.Args[1])
	if err != nil {
		panic(err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(home, ".config", "rdg", "rdgconf.json")
	conf := config.LoadConfig(configPath)

	fmt.Print("\033[?1049h")
	fmt.Print("\033[2J\033[H")
	defer func() {
		fmt.Print("\033[?25h")
		fmt.Print("\033[?1049l")
	}()

	start := time.Now()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		elapsed := time.Since(start)
		remaining := total - elapsed

		if remaining <= 0 {
			break
		}

		secs := int(math.Ceil(remaining.Seconds()))
		h := secs / 3600
		m := (secs % 3600) / 60
		s := secs % 60

		barWidth := progressBarWidth()
		bar := renderProgressBar(elapsed, total, barWidth)

		fmt.Printf("\r%02d:%02d:%02d [%s]", h, m, s, bar)

		<-ticker.C
	}

	barWidth := progressBarWidth()
	fmt.Printf("\r00:00:00 [%s]\n", renderProgressBar(total, total, barWidth))

	if conf.SoundFilePath != "" {
		repetition := conf.Repetition
		if repetition <= 0 {
			repetition = 3
		}
		for range repetition {
			sound.PlaySound(conf.SoundFilePath)
		}
	}
}

