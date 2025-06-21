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

var model string = "gpt-4o"

type LLMRequest struct {
	Model  string   `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role  string  `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	Choices  []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func getAPIKey() (string, error) {
	// try env, if not, look in file
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
			return key, nil
		}
	homeDir, err := os.UserHomeDir()
  	if err != nil {
  		return "", fmt.Errorf("could not get home directory: %v", err)
  	}

  	keyPath := filepath.Join(homeDir, ".api_keys", "openai_key")
  	keyBytes, err := os.ReadFile(keyPath)
  	if err != nil {
  		return "", fmt.Errorf("could not read API key from %s: %v", keyPath, err)
  	}

  	key := strings.TrimSpace(string(keyBytes))
  	if key == "" {
  		return "", fmt.Errorf("API key file is empty")
  	}

  	return key, nil

}

func callLLM(userMessage string) (string, error) {
	apiKey, err := getAPIKey()
	if err != nil {
		return "", fmt.Errorf("failed to get API key: %w", err)
	}

	// Create request structure
	request := LLMRequest{
		Model: model,
		Messages: []Message{
				{
					Role:	"system",
					Content: 
						`You are a helpful, voice-based assistant.
						Speak naturally, like you are talking to a friend.
						Keep your answers short and to the point.
						Use conversational language, contractions, and 
						occasionally check in like 'Want to hear more?'`,
            	},
				
				{
					Role:	"user",
					Content: userMessage,
				},
		},
	}

	// Encode in JSON
	jsonData, err := json.Marshal(request)
  	if err != nil {
  		return "", fmt.Errorf("failed to marshal request: %w", err)
  	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
  		return "", fmt.Errorf("failed to create request: %w", err)
  	}

  	// Set headers
  	req.Header.Set("Content-Type", "application/json")
  	req.Header.Set("Authorization", "Bearer "+apiKey)

  	// Make request
  	client := &http.Client{Timeout: 30 * time.Second}
  	resp, err := client.Do(req)
  	if err != nil {
  		return "", fmt.Errorf("API request failed: %w", err)
  	}
  	defer resp.Body.Close()

  	// Read response
  	body, err := io.ReadAll(resp.Body)
  	if err != nil {
  		return "", fmt.Errorf("failed to read response: %w", err)
  	}

  	// Parse JSON response
  	var llmResponse LLMResponse
  	err = json.Unmarshal(body, &llmResponse)
  	if err != nil {
  		return "", fmt.Errorf("failed to parse response: %w", err)
  	}

  	// Extract message
  	if len(llmResponse.Choices) == 0 {
  		return "", fmt.Errorf("no response choices returned")
  	}

  	return llmResponse.Choices[0].Message.Content, nil
}