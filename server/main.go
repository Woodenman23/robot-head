package main

import (
	"context"
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
        // Call OpenAI for user input messages
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
		
		fmt.Printf("Received: %+v\n", msg)
		
		response := createResponse(msg)
		err = conn.WriteJSON(response)
		if err != nil {
			log.Println("Failed to send response:", err)
			break
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

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
			return
		}
		fmt.Printf("Server started on localhost%s\n", portNum)
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
