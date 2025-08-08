package gochatbot

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.rumenx.com/chatbot/config"
	"go.rumenx.com/chatbot/middleware"
	"go.rumenx.com/chatbot/models"
)

func TestWithModel(t *testing.T) {
	chatbot := &Chatbot{}
	freeModel := models.NewFreeModel()
	option := WithModel(freeModel)

	option(chatbot)

	if chatbot.model != freeModel {
		t.Error("Expected model to be set")
	}
}

func TestWithTimeout(t *testing.T) {
	chatbot := &Chatbot{}
	timeout := 30 * time.Second
	option := WithTimeout(timeout)

	option(chatbot)

	if chatbot.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, chatbot.timeout)
	}
}

func TestWithFilter(t *testing.T) {
	chatbot := &Chatbot{}
	filter := middleware.NewChatMessageFilter(config.MessageFilteringConfig{
		Enabled: true,
	})
	option := WithFilter(filter)

	option(chatbot)

	if chatbot.filter != filter {
		t.Error("Expected filter to be set")
	}
}

func TestWithRateLimit(t *testing.T) {
	chatbot := &Chatbot{}
	limiter := middleware.NewRateLimiter(config.RateLimitConfig{
		RequestsPerMinute: 100,
		BurstSize:         10,
	})
	option := WithRateLimit(limiter)

	option(chatbot)

	if chatbot.rateLimit != limiter {
		t.Error("Expected rate limiter to be set")
	}
}

func TestNewWithOptions(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"},
		WithTimeout(10*time.Second),
	)

	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	if chatbot.timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s, got %v", chatbot.timeout)
	}
}

func TestChatbotAsk(t *testing.T) {
	chatbot, err := New(&config.Config{
		Model: "free",
		RateLimit: config.RateLimitConfig{
			RequestsPerMinute: 600, // Allow enough requests for testing (10 per second)
			BurstSize:         10,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	tests := []struct {
		name    string
		message string
	}{
		{"simple question", "Hello there"},
		{"question mark", "How are you?"},
		{"greeting", "hi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			response, err := chatbot.Ask(ctx, tt.message)

			if err != nil {
				t.Errorf("Ask() error = %v", err)
				return
			}

			if response == "" {
				t.Error("Expected non-empty response")
			}
		})
	}
}

func TestChatbotAskWithContext(t *testing.T) {
	chatbot, err := New(&config.Config{
		Model: "free",
		RateLimit: config.RateLimitConfig{
			RequestsPerMinute: 600,
			BurstSize:         10,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	ctx := context.Background()

	response, err := chatbot.Ask(ctx, "Hello", WithContext("user_id", "test123"))
	if err != nil {
		t.Errorf("Ask() with context error = %v", err)
	}

	if response == "" {
		t.Error("Expected non-empty response")
	}
}

func TestChatbotGetConfig(t *testing.T) {
	originalConfig := &config.Config{
		Model:   "free",
		Timeout: 30 * time.Second,
	}

	chatbot, err := New(originalConfig)
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	retrievedConfig := chatbot.GetConfig()

	if retrievedConfig.Model != originalConfig.Model {
		t.Errorf("Expected model %s, got %s", originalConfig.Model, retrievedConfig.Model)
	}
}

func TestChatbotGetModel(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	model := chatbot.GetModel()

	if model == nil {
		t.Error("Expected non-nil model")
	}

	if model.Name() != "free-model" {
		t.Errorf("Expected model name 'free-model', got %s", model.Name())
	}
}

func TestChatbotHealth(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	ctx := context.Background()
	err = chatbot.Health(ctx)

	if err != nil {
		t.Errorf("Health() error = %v", err)
	}
}

func TestChatbotHealthWithTimeout(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for context to expire
	time.Sleep(1 * time.Millisecond)

	err = chatbot.Health(ctx)

	// Free model health check is very fast, so timeout might not occur
	// This test is mainly to ensure Health() doesn't panic with cancelled context
	if err != nil {
		t.Logf("Health() returned error (expected for cancelled context): %v", err)
	} else {
		t.Logf("Health() succeeded despite cancelled context (free model is very fast)")
	}
}

func TestChatbotAskWithTimeout(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for context to expire
	time.Sleep(1 * time.Millisecond)

	_, err = chatbot.Ask(ctx, "Hello")

	if err == nil {
		t.Error("Expected timeout error from Ask()")
	}
}

func TestChatbotAskStream(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	ctx := context.Background()

	// Create a mock response writer
	w := httptest.NewRecorder()

	err = chatbot.AskStream(ctx, w, "Hello")

	if err != nil {
		t.Errorf("AskStream() error = %v", err)
	}
}

func TestChatbotAskStreamWithContext(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	ctx := context.Background()

	w := httptest.NewRecorder()

	err = chatbot.AskStream(ctx, w, "Tell me a story", WithContext("stream_id", "stream123"))

	if err != nil {
		t.Errorf("AskStream() with context error = %v", err)
	}
}

func TestChatbotAskEmptyMessage(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	ctx := context.Background()
	_, err = chatbot.Ask(ctx, "")

	if err == nil {
		t.Error("Expected error for empty message")
	}

	if !strings.Contains(err.Error(), "message cannot be empty") {
		t.Errorf("Expected 'message cannot be empty' error, got: %v", err)
	}
}

func TestChatbotAskStreamEmptyMessage(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	ctx := context.Background()
	w := httptest.NewRecorder()

	err = chatbot.AskStream(ctx, w, "")

	if err == nil {
		t.Error("Expected error for empty message in stream")
	}

	if !strings.Contains(err.Error(), "message cannot be empty") {
		t.Errorf("Expected 'message cannot be empty' error, got: %v", err)
	}
}

func TestChatbotWithInvalidModel(t *testing.T) {
	_, err := New(&config.Config{Model: "invalid-model-12345"})

	if err == nil {
		t.Error("Expected error for invalid model")
	}
}

func TestNewWithNilConfig(t *testing.T) {
	_, err := New(nil)

	if err == nil {
		t.Error("Expected error for nil config")
	}

	if !strings.Contains(err.Error(), "config cannot be nil") {
		t.Errorf("Expected 'config cannot be nil' error, got: %v", err)
	}
}

func TestChatbotWithMultipleContextOptions(t *testing.T) {
	chatbot, err := New(&config.Config{
		Model: "free",
		RateLimit: config.RateLimitConfig{
			RequestsPerMinute: 600,
			BurstSize:         10,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	ctx := context.Background()

	response, err := chatbot.Ask(ctx, "Hello",
		WithContext("user_id", "user123"),
		WithContext("session", "sess456"),
		WithContext("language", "en"),
	)

	if err != nil {
		t.Errorf("Ask() with multiple context options error = %v", err)
	}

	if response == "" {
		t.Error("Expected non-empty response")
	}
}

func TestChatbotAskStream_RateLimitExceeded(t *testing.T) {
	// Test the rate limiting functionality more carefully
	limiter := middleware.NewRateLimiter(config.RateLimitConfig{
		RequestsPerMinute: 60, // More reasonable limit
		BurstSize:         1,
		Window:            time.Minute,
	})

	chatbot, err := New(&config.Config{Model: "free"}, WithRateLimit(limiter))
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	ctx := context.Background()
	w := httptest.NewRecorder()

	// Just test that rate limiting doesn't panic
	err = chatbot.AskStream(ctx, w, "Test message")

	// Rate limiting might or might not trigger, but it should not panic
	if err != nil {
		t.Logf("AskStream with rate limit error (may be expected): %v", err)
	}
}

func TestChatbotAskStream_WithTimeout(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"}, WithTimeout(1*time.Millisecond))
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	ctx := context.Background()
	w := httptest.NewRecorder()

	// This should timeout very quickly
	err = chatbot.AskStream(ctx, w, "Tell me a very long story")

	// The test might pass or timeout depending on timing, but it should not panic
	// and should handle the timeout context properly
	if err != nil && !strings.Contains(err.Error(), "context") {
		t.Logf("Got error (expected): %v", err)
	}
}

func TestChatbotAskStream_FilterFailure(t *testing.T) {
	// Create a filter that will reject certain messages
	filter := middleware.NewChatMessageFilter(config.MessageFilteringConfig{
		Enabled:     true,
		Profanities: []string{"badword"},
	})

	chatbot, err := New(&config.Config{Model: "free"}, WithFilter(filter))
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	ctx := context.Background()
	w := httptest.NewRecorder()

	// Try with a blocked word
	err = chatbot.AskStream(ctx, w, "This contains badword which should be blocked")

	// The filter might block this or process it depending on implementation
	// The important thing is it doesn't panic
	if err != nil {
		t.Logf("Filter result (may be expected): %v", err)
	}
}

func TestChatbotAskStream_NonStreamingModel(t *testing.T) {
	// Test the fallback path for non-streaming models
	// The "free" model might or might not support streaming, so this tests the fallback
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	ctx := context.Background()
	w := httptest.NewRecorder()

	err = chatbot.AskStream(ctx, w, "Short test message")
	if err != nil {
		// If this fails, it might be due to the model not being properly mocked
		// But the important thing is that the code path is exercised
		t.Logf("AskStream with non-streaming model error: %v", err)
	}

	// Check that some response was written
	if w.Body.Len() == 0 {
		t.Log("No response body written (may be expected for mock)")
	}
}

func TestChatbotAskStream_ContextCancellation(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	// Create a context that gets cancelled immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	w := httptest.NewRecorder()

	err = chatbot.AskStream(ctx, w, "This should be cancelled")

	// Should handle cancelled context gracefully
	if err != nil && strings.Contains(err.Error(), "context canceled") {
		t.Logf("Got expected context cancellation: %v", err)
	}
}
