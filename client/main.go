package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"robot-head/shared"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

func openWebsocket() (*websocket.Conn, error) {
	// Websocket server URL
	serverURL := url.URL{Scheme: "ws", Host: "localhost:9001", Path: "/ws"}
	fmt.Printf("Connecting to %s\n", serverURL.String())
	// Connect to URL
	conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to server!")
	return conn, nil
}

func connectWithRetry() (*websocket.Conn, error) {
	maxRetries := 5
	baseDelay := 1 * time.Second
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		conn, err := openWebsocket()
		if err == nil {
			return conn, nil
		}
		
		if attempt < maxRetries-1 {
			delay := baseDelay * time.Duration(1<<attempt) // Exponential backoff
			fmt.Printf("Connection failed, retrying in %v... (attempt %d/%d)\n", delay, attempt+1, maxRetries)
			time.Sleep(delay)
		}
	}
	
	return nil, fmt.Errorf("failed to connect after %d attempts", maxRetries)
}

func createUserMessage(text string) shared.Message {
	return shared.Message{
		Type:      shared.MessageTypeUserInput,
		Timestamp: time.Now().Unix(),
		Data:      text,
	}
}

func listenForMessages(conn *websocket.Conn) {
	for {
		var response shared.Message
		err := conn.ReadJSON(&response)
		if err != nil {
			log.Println("Connection closed", err)
			break
		}
	switch  response.Type {
	case shared.MessageTypeAIResponse:
			fmt.Printf("\nRobot: %v\n\n", response.Data)
	case shared.MessageTypeAudio:
		// Parse audio data
		audioDataJSON, err := json.Marshal(response.Data)
		if err != nil {
			log.Printf("Failed to marshal audio data: %v\n", err)
			continue
		}

		var audioData shared.AudioData
		err = json.Unmarshal(audioDataJSON, &audioData)
		if err != nil {
			log.Printf("Failed to unmarshal audio data: %v\n", err)
			continue
		}

		fmt.Printf("\nPlaying audio for: %s\n", audioData.Text)
		go func() {
			err := playAudio(audioData.AudioData)
			if err != nil {
				log.Printf("Failed to play audio: %v\n", err)
			}
		}()		
		default:
			// Other message types (status, error, etc.)
			fmt.Printf("Server: %v\n", response.Data)
		}	
	}
}


func sendMessages(conn *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type messages (Ctrl+C to quit): ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		

		userMessage := createUserMessage(text)
		err := conn.WriteJSON(userMessage)
		if err != nil {
			log.Println("Failed to send message:", err)
			break
		}
	}
}

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

	// Resample if neccessary
	resampled := beep.Resample(4, format.SampleRate, sr, streamer)

	done := make(chan bool)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		done <- true	
	})))

	<-done
	return nil
}

func main() {
	fmt.Println("Robot Head Client starting...")

	conn, err := connectWithRetry()
	if err != nil {
		log.Fatal("Websocket connection failed:", err)
	}
	defer conn.Close()

	// Test server connection
	testMessage := shared.Message{
		Type:      shared.MessageTypeStatus,
		Timestamp: time.Now().Unix(),
		Data:      "Robot head client connected",
	}

	err = conn.WriteJSON(testMessage)
	if err != nil {
		log.Fatal("Failed to send message:", err)
	}
	log.Println("Client connected and ready to send/recieve messages.")
	fmt.Println("Sent connection message to server.")

	go listenForMessages(conn)

	sendMessages(conn)
}
