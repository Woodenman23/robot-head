package shared

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMessageSerialization(t *testing.T) {
	// Create a test message
	msg := Message{
		Type:      MessageTypeUserInput,
		Timestamp: time.Now().Unix(),
		Data:      "Hello, robot!",
	}
	
	// Serialize to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	
	// Deserialize back
	var decoded Message
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}
	
	// Verify fields match
	if decoded.Type != msg.Type {
		t.Errorf("Expected type %v, got %v", msg.Type, decoded.Type)
	}
	
	if decoded.Data != msg.Data {
		t.Errorf("Expected data %v, got %v", msg.Data, decoded.Data)
	}
}

func TestAllMessageTypes(t *testing.T) {
	messageTypes := []MessageType{
		MessageTypeUserInput,
		MessageTypeAIResponse,
		MessageTypeStatus,
		MessageTypeError,
	}
	
	for _, msgType := range messageTypes {
		msg := Message{
			Type:      msgType,
			Timestamp: time.Now().Unix(),
			Data:      "test data",
		}
		
		jsonData, err := json.Marshal(msg)
		if err != nil {
			t.Errorf("Failed to marshal %v: %v", msgType, err)
			continue
		}
		
		var decoded Message
		err = json.Unmarshal(jsonData, &decoded)
		if err != nil {
			t.Errorf("Failed to unmarshal %v: %v", msgType, err)
			continue
		}
		
		if decoded.Type != msgType {
			t.Errorf("Type mismatch for %v: got %v", msgType, decoded.Type)
		}
	}
}

func TestMessageWithDifferentDataTypes(t *testing.T) {
	testCases := []struct {
		name string
		data interface{}
	}{
		{"string", "hello world"},
		{"number", 42},
		{"boolean", true},
		{"map", map[string]string{"key": "value"}},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := Message{
				Type:      MessageTypeStatus,
				Timestamp: time.Now().Unix(),
				Data:      tc.data,
			}
			
			jsonData, err := json.Marshal(msg)
			if err != nil {
				t.Fatalf("Failed to marshal %v: %v", tc.name, err)
			}
			
			var decoded Message
			err = json.Unmarshal(jsonData, &decoded)
			if err != nil {
				t.Fatalf("Failed to unmarshal %v: %v", tc.name, err)
			}
			
			if decoded.Type != msg.Type {
				t.Errorf("Type mismatch for %v", tc.name)
			}
		})
	}
}