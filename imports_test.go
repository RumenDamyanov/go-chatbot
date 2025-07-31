package gochatbot

import (
	"testing"

	"github.com/RumenDamyanov/go-chatbot/config"
)

// TestBasicChatbotCreation tests that we can create a basic chatbot
func TestBasicChatbotCreation(t *testing.T) {
	cfg := &config.Config{
		Model: "free",
	}

	chatbot, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	if chatbot == nil {
		t.Fatal("Chatbot is nil")
	}

	t.Log("✅ Basic chatbot creation successful")
}
