package main

import (
	"fmt"
	"log"
	"net/url"
	"github.com/gorilla/websocket"
	//"robot-head/shared"
)

func main() {
	fmt.Println("Robot Head Client starting...")
	
	// Websocket server URL
	serverURL := url.URL{Scheme: "ws", Host: "localhost:9001", Path: "/ws"}
	fmt.Printf("Connecting to %s\n", serverURL.String())
	// Connect to URL
	conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		log.Fatal("Websocket connection failed:", err)
	}
	defer conn.Close()
	fmt.Println("Connected to server!")

	// TODO: Handle messages from server
	// TODO: Send messages to server
	
	log.Println("Client ready")
}