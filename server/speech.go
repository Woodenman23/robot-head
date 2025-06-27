package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

var whisperModel whisper.Model

func initWhisper() error {
	model, err := whisper.New("./models/ggml-base.en.bin")
	if err != nil {
		return fmt.Errorf("failed to load whisper model: %v", err)
	}
	whisperModel = model
	log.Println("Whisper model loaded successfully")
	return nil
}

func transcribeAudio(audioData []byte) (string, error) {
	if whisperModel == nil {
		return "", fmt.Errorf("whisper model not initialized")
	}

	// Convert bytes to float32 samples
	// audioData is 16-bit PCM, need to convert to float32
	samples := make([]float32, len(audioData)/2)
	for i := 0; i < len(samples); i++ {
		// Convert 16-bit PCM to float32 (-1.0 to 1.0)
		sample := int16(audioData[i*2]) | int16(audioData[i*2+1])<<8
		samples[i] = float32(sample) / 32768.0
	}

	// Create processing context
	context, err := whisperModel.NewContext()
	if err != nil {
		return "", fmt.Errorf("failed to create whisper context: %v", err)
	}

	// Process audio
	if err := context.Process(samples, nil, nil, nil); err != nil {
		return "", fmt.Errorf("failed to process audio: %v", err)
	}

	// Extract transcript using NextSegment
	var transcript string
	for {
		segment, err := context.NextSegment()
		if err != nil {
			break // EOF or error, we're done
		}
		transcript += segment.Text + " "
	}

	// Clean up transcript
	transcript = strings.TrimSpace(transcript)
	return transcript, nil
}