package sound

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func PlaySound(sound_file_path string) {
	f, err := os.Open(sound_file_path)
	if err != nil {
		fmt.Println("No sound file found at the specified souend_file_path in ~/.config/rdg/rdgconf.json")
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := streamer.Close(); err != nil {
			panic(err)
		}
	}()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/20))
	done := make(chan bool)

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

