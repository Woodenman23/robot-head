package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"robot-head/shared"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func createResponse(msg shared.Message) shared.Message {
	// Handle voice messages (audio input from client)
	if msg.Type == shared.MessageTypeAudio {
		// Parse audio data from client
		audioDataJSON, err := json.Marshal(msg.Data)
		if err != nil {
			log.Printf("Failed to marshal audio data: %v", err)
			return shared.Message{
				Type:      shared.MessageTypeError,
				Timestamp: time.Now().Unix(),
				Data:      "Failed to process audio data",
			}
		}

		var audioData shared.AudioData
		err = json.Unmarshal(audioDataJSON, &audioData)
		if err != nil {
			log.Printf("Failed to unmarshal audio data: %v", err)
			return shared.Message{
				Type:      shared.MessageTypeError,
				Timestamp: time.Now().Unix(),
				Data:      "Failed to parse audio data",
			}
		}

		// Transcribe audio to text using Whisper
		transcript, err := transcribeAudio(audioData.AudioData)
		if err != nil {
			log.Printf("Speech-to-text error: %v", err)
			return shared.Message{
				Type:      shared.MessageTypeError,
				Timestamp: time.Now().Unix(),
				Data:      "Sorry, I couldn't understand what you said.",
			}
		}

		// Skip processing if no speech detected - return empty message
		if transcript == "" || transcript == "[BLANK_AUDIO]" {
			return shared.Message{} // Empty message - won't be sent
		}

		fmt.Printf("User: %s\n", transcript)

		// Process transcript with OpenAI
		aiResponse, err := callLLM(transcript)
		if err != nil {
			log.Printf("OpenAI API error: %v", err)
			return shared.Message{
				Type:      shared.MessageTypeError,
				Timestamp: time.Now().Unix(),
				Data:      "Sorry, I'm having trouble thinking right now.",
			}
		}

		fmt.Printf("Robot: %s\n", aiResponse)

		// Generate speech from AI response
		audioBytes, err := generateSpeech(aiResponse)
		if err != nil {
			log.Printf("TTS error: %v", err)
			// Fallback to text response
			return shared.Message{
				Type:      shared.MessageTypeAIResponse,
				Timestamp: time.Now().Unix(),
				Data:      aiResponse,
			}
		}

		// Return audio response
		responseAudioData := shared.AudioData{
			Text:      aiResponse,
			AudioData: audioBytes,
			MimeType:  "audio/mpeg",
		}

		return shared.Message{
			Type:      shared.MessageTypeAudio,
			Timestamp: time.Now().Unix(),
			Data:      responseAudioData,
		}
	}

	// Handle text input messages (fallback)
	if msg.Type == shared.MessageTypeUserInput {
		if userText, ok := msg.Data.(string); ok {
			aiResponse, err := callLLM(userText)
			if err != nil {
				log.Printf("OpenAI API error: %v", err)
				return shared.Message{
					Type:      shared.MessageTypeError,
					Timestamp: time.Now().Unix(),
					Data:      "Sorry, I'm having trouble thinking right now.",
				}
			}

			audioBytes, err := generateSpeech(aiResponse)
			if err != nil {
				log.Printf("TTS error: %v", err)
				// Fallback to text response
				return shared.Message{
					Type:      shared.MessageTypeAIResponse,
					Timestamp: time.Now().Unix(),
					Data:      aiResponse,
				}
			}

			// Return audio message
			audioData := shared.AudioData{
				Text:      aiResponse,
				AudioData: audioBytes,
				MimeType:  "audio/mpeg",
			}

			return shared.Message{
				Type:      shared.MessageTypeAudio,
				Timestamp: time.Now().Unix(),
				Data:      audioData,
			}
		}
	}

	// Fallback for other message types
	return shared.Message{
		Type:      shared.MessageTypeStatus,
		Timestamp: time.Now().Unix(),
		Data:      fmt.Sprintf("Received: %v", msg.Data),
	}
}


func handleMessageExchange(conn *websocket.Conn) {
	for {
		var msg shared.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		
		// Log message type without printing binary data
		if msg.Type == shared.MessageTypeAudio {
			fmt.Printf("Received audio message\n")
		} else {
			fmt.Printf("Received: %+v\n", msg)
		}
		
		response := createResponse(msg)
		// Only send response if it has content (not empty message)
		if response.Type != "" {
			err = conn.WriteJSON(response)
			if err != nil {
				log.Println("Failed to send response:", err)
				break
			}
		}
	}
}

// Handle WebSocket connections
func establishWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected via WebSocket")
	handleMessageExchange(conn)
}


func main() {
	// Initialize Whisper model
	fmt.Println("Loading Whisper model...")
	if err := initWhisper(); err != nil {
		log.Fatal("Failed to initialize Whisper:", err)
	}

	port := getEnv("PORT", "9001")
	host := getEnv("HOST", "0.0.0.0")
	portNum := ":" + port

	server := &http.Server{
		Addr: host + portNum,
	}

	// health check endpoint - used for Docker health checks
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is healthy!")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Robot Head Server!")
	})

	http.HandleFunc("/ws", establishWebsocketConnection)
	fmt.Printf("Server running on localhost%s\n", portNum)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
			return
		}
	}()

	// wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")

  	// Graceful shutdown with 30 second timeout
  	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
  	defer cancel()

  	if err := server.Shutdown(ctx); err != nil {
  		log.Fatal("Server forced to shutdown:", err)
  	}

  	fmt.Println("Server exited")

}
