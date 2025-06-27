package main

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

func playAudio(audioData []byte) error {
	// Initialize speaker
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second))

	// create a reader from the MP3 data
	reader := bytes.NewReader(audioData)
	streamer, format, err := mp3.Decode(io.NopCloser(reader))
	if err != nil {
		return fmt.Errorf("failed to decode MP3: %v", err)
	}
	defer streamer.Close()

	// Resample if necessary
	resampled := beep.Resample(4, format.SampleRate, sr, streamer)

	done := make(chan bool)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}