package main

import (
	"fmt"
	"log"
	"robot-head/shared"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/gorilla/websocket"
)

func recordAudio(duration time.Duration) ([]byte, error) {
	err := portaudio.Initialize()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PortAudio: %v", err)
	}
	defer portaudio.Terminate()

	sampleRate := 16000
	channels := 1
	framesPerBuffer := 1024

	// Audio data buffer
	totalSamples := int(float64(sampleRate) * duration.Seconds())
	audioBuffer := make([]float32, totalSamples)
	bufferIndex := 0

	// Input stream parameters
	inputParams := portaudio.HighLatencyParameters(nil, nil)
	inputParams.Input.Device, _ = portaudio.DefaultInputDevice()
	inputParams.Input.Channels = channels
	inputParams.SampleRate = float64(sampleRate)
	inputParams.FramesPerBuffer = framesPerBuffer

	stream, err := portaudio.OpenStream(inputParams, func(in []float32) {
		// Prevent buffer overflow
		remainingSpace := len(audioBuffer) - bufferIndex
		if remainingSpace > 0 {
			copyLen := len(in)
			if copyLen > remainingSpace {
				copyLen = remainingSpace
			}
			copy(audioBuffer[bufferIndex:bufferIndex+copyLen], in[:copyLen])
			bufferIndex += copyLen
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open audio stream: %v", err)
	}
	defer stream.Close()

	// Start recording
	err = stream.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start recording: %v", err)
	}

	// Record for specified duration
	time.Sleep(duration)

	// Stop recording
	err = stream.Stop()
	if err != nil {
		return nil, fmt.Errorf("failed to stop recording: %v", err)
	}

	// Convert float32 samples to bytes (16-bit PCM)
	audioBytes := make([]byte, bufferIndex*2)
	for i := 0; i < bufferIndex; i++ {
		sample := int16(audioBuffer[i] * 32767)
		audioBytes[i*2] = byte(sample)
		audioBytes[i*2+1] = byte(sample >> 8)
	}

	return audioBytes, nil
}

func sendVoiceMessage(conn *websocket.Conn, audioData []byte) error {
	voiceData := shared.AudioData{
		AudioData: audioData,
		MimeType:  "audio/pcm",
	}

	voiceMessage := shared.Message{
		Type:      shared.MessageTypeAudio,
		Timestamp: time.Now().Unix(),
		Data:      voiceData,
	}

	return conn.WriteJSON(voiceMessage)
}

func sendVoiceMessages(conn *websocket.Conn) {
	fmt.Println("Say something")

	for {
		audioData, err := recordAudio(3 * time.Second)
		if err != nil {
			log.Printf("Failed to record audio: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		err = sendVoiceMessage(conn, audioData)
		if err != nil {
			log.Println("Failed to send voice message:", err)
			break
		}
	}
}