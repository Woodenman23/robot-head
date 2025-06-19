package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"robot-head/shared"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// Handle WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected via WebSocket")

	// Keep connection alive and handle messages
	for {
		var msg shared.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		
		fmt.Printf("Received: %+v\n", msg)
		// TODO: Process message and send response
	}
}

func main() {
	portNum := ":9001"

	// health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is healthy!")
	})

	// Test endpoint
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You made a %s request to /test", r.Method)
	})

	// landing page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Robot Head Server!")
	})

	http.HandleFunc("/ws", handleWebSocket)

	fmt.Printf("Server started on localhost%s\n", portNum)
	log.Fatal(http.ListenAndServe(portNum, nil))
}
