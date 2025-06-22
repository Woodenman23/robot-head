package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ElevenlabsAPIURL = "https://api.elevenlabs.io/v1/text-to-speech"
	Model = "eleven_turbo_v2"
	VoiceID   = "Oe8Lhg3t63j9BsrTQBjx" // Yowz - South London Bloke
)

type TTSRequest struct {
	Text	string		`json:"text"`
	ModelID	string		`json:"model_id"`
	VoiceSettings VoiceSettings `json:"voice_settings"`
}

type VoiceSettings struct {
	Stability 		float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
}

func getElevenLabsAPIKey() (string, error) {
	if key := os.Getenv("ELEVENLABS_API_KEY"); key != "" {
		return key, nil
	}

	  	homeDir, err := os.UserHomeDir()
  	if err != nil {
  		return "", fmt.Errorf("could not get home directory: %v", err)
  	}

  	keyPath := filepath.Join(homeDir, ".api_keys", "elevenlabs_key")
  	keyBytes, err := os.ReadFile(keyPath)
  	if err != nil {
  		return "", fmt.Errorf("could not read ElevenLabs API key: %v", err)
  	}

  	return strings.TrimSpace(string(keyBytes)), nil
}

func generateSpeech(text string) ([]byte, error) {
	apiKey, err := getElevenLabsAPIKey()
	if err != nil {
		return nil, err
	}

	request := TTSRequest {
		Text : text,
		ModelID: Model,
		VoiceSettings : VoiceSettings{
			Stability: 0.5,
			SimilarityBoost: 0.5,
		},
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s", ElevenlabsAPIURL, VoiceID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "audio/mpeg")
  	req.Header.Set("Content-Type", "application/json")
  	req.Header.Set("xi-api-key", apiKey)

  	client := &http.Client{Timeout: 30 * time.Second}
  	resp, err := client.Do(req)
  	if err != nil {
  		return nil, err
  	}
  	defer resp.Body.Close()

  	if resp.StatusCode != http.StatusOK {
  	body, err := io.ReadAll(resp.Body)
  	if err != nil {
  		return nil, fmt.Errorf("TTS API error %d, couldn't read error response", resp.StatusCode)
  	}
  	return nil, fmt.Errorf("TTS API error %d: %s", resp.StatusCode, string(body))
  }


  	return io.ReadAll(resp.Body)

}

