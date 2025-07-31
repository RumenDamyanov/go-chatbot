package models

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"
	"time"
)

// FreeModel is a simple fallback model that provides basic responses
// without requiring API keys or external services.
type FreeModel struct {
	responses []string
}

// NewFreeModel creates a new free model instance.
func NewFreeModel() *FreeModel {
	return &FreeModel{
		responses: []string{
			"I'm a simple chatbot. How can I help you today?",
			"Thanks for your message! I'm here to assist you.",
			"Hello! I'm a basic AI assistant. What would you like to know?",
			"I appreciate you reaching out. How can I be of service?",
			"Hi there! I'm ready to help with any questions you might have.",
			"Thank you for your question. I'll do my best to help!",
			"Hello! I'm an AI assistant. Feel free to ask me anything.",
			"I'm here to help! What can I assist you with today?",
		},
	}
}

// Ask processes a message and returns a response.
func (f *FreeModel) Ask(ctx context.Context, message string, context map[string]interface{}) (string, error) {
	// Simulate processing time with context awareness
	select {
	case <-time.After(100 * time.Millisecond):
		// Continue processing
	case <-ctx.Done():
		return "", ctx.Err()
	}

	// Simple response based on message content
	message = strings.ToLower(strings.TrimSpace(message))

	switch {
	case strings.Contains(message, "hello") || strings.Contains(message, "hi"):
		return "Hello! Nice to meet you. How can I help you today?", nil
	case strings.Contains(message, "how are you"):
		return "I'm doing well, thank you for asking! How are you?", nil
	case strings.Contains(message, "thank"):
		return "You're welcome! I'm happy to help.", nil
	case strings.Contains(message, "bye") || strings.Contains(message, "goodbye"):
		return "Goodbye! Have a great day!", nil
	case strings.Contains(message, "help"):
		return "I'm here to help! Feel free to ask me any questions.", nil
	case strings.Contains(message, "name"):
		return "I'm a simple AI chatbot. You can call me Bot!", nil
	case strings.Contains(message, "?"):
		return "That's an interesting question! While I'm a basic model, I'll do my best to help.", nil
	default:
		// Return a random response for other messages using crypto/rand
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(f.responses))))
		if err != nil {
			// Fallback to first response if random generation fails
			return f.responses[0], nil
		}
		return f.responses[n.Int64()], nil
	}
}

// Name returns the name of the model.
func (f *FreeModel) Name() string {
	return "free-model"
}

// Provider returns the provider of the model.
func (f *FreeModel) Provider() string {
	return "local"
}

// Health checks if the model is healthy.
func (f *FreeModel) Health(ctx context.Context) error {
	// Free model is always healthy since it's local
	return nil
}
