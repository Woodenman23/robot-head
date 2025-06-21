package main

import (
	"robot-head/shared"
	"testing"
)

func TestCreateUserMessage(t *testing.T) {
	text := "Hello robot!"
	msg := createUserMessage(text)
	
	if msg.Type != shared.MessageTypeUserInput {
		t.Errorf("Expected type %v, got %v", shared.MessageTypeUserInput, msg.Type)
	}
	
	if msg.Data != text {
		t.Errorf("Expected data %v, got %v", text, msg.Data)
	}
	
	if msg.Timestamp == 0 {
		t.Error("Timestamp should not be zero")
	}
}

func TestCreateUserMessageEmptyString(t *testing.T) {
	text := ""
	msg := createUserMessage(text)
	
	if msg.Data != "" {
		t.Errorf("Expected empty data, got %v", msg.Data)
	}
	
	if msg.Type != shared.MessageTypeUserInput {
		t.Error("Type should still be MessageTypeUserInput")
	}
}