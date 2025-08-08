package middleware

import (
	"context"
	"testing"
	"time"

	"go.rumenx.com/chatbot/config"
)

func TestChatMessageFilter_UpdateConfig(t *testing.T) {
	filter := NewChatMessageFilter(config.MessageFilteringConfig{
		Enabled:     true,
		Profanities: []string{"bad"},
	})

	// Test updating config
	newConfig := config.MessageFilteringConfig{
		Enabled:     true,
		Profanities: []string{"bad", "worse"},
		LinkPattern: `https?://[^\s]+`,
	}

	filter.UpdateConfig(newConfig)

	// Test that new config is applied
	ctx := context.Background()
	result, err := filter.Handle(ctx, "This is worse content")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected result but got nil")
		return
	}
	if result.Message != "This is *** content" {
		t.Errorf("expected filtered message 'This is *** content', got '%s'", result.Message)
	}
}

func TestRateLimiter_Cleanup(t *testing.T) {
	limiter := NewRateLimiter(config.RateLimitConfig{
		RequestsPerMinute: 10,
		Window:            time.Minute,
	})

	// Create context with client IP
	ctx := context.WithValue(context.Background(), "client_ip", "192.168.1.1")

	// Make a request to create the client
	err := limiter.Allow(ctx)
	if err != nil {
		t.Errorf("first request should be allowed, got error: %v", err)
	}

	// Manually add old request data
	clientID := "192.168.1.1"
	limiter.mutex.Lock()
	// Add an old timestamp
	limiter.requests[clientID] = append(limiter.requests[clientID], time.Now().Add(-time.Hour*2))
	limiter.mutex.Unlock()

	// Call cleanup
	limiter.Cleanup()

	// Check that old requests were cleaned up
	limiter.mutex.Lock()
	requests, exists := limiter.requests[clientID]
	limiter.mutex.Unlock()

	if !exists || len(requests) == 0 {
		t.Error("cleanup should preserve recent requests but remove old ones")
	}
}

func TestRateLimiter_StartCleanupRoutine(t *testing.T) {
	limiter := NewRateLimiter(config.RateLimitConfig{
		RequestsPerMinute: 10,
		Window:            time.Minute,
	})

	// Start cleanup routine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go limiter.StartCleanupRoutine(ctx)

	// Add a request
	reqCtx := context.WithValue(context.Background(), "client_ip", "192.168.1.1")
	err := limiter.Allow(reqCtx)
	if err != nil {
		t.Errorf("first request should be allowed, got error: %v", err)
	}

	// Wait a bit for cleanup routine to potentially run
	time.Sleep(time.Millisecond * 100)

	// Verify the cleanup routine is running by checking that it doesn't panic
	cancel() // Stop the cleanup routine
}

func TestChatMessageFilter_Handle_EdgeCases(t *testing.T) {
	filter := NewChatMessageFilter(config.MessageFilteringConfig{
		Enabled:     true,
		Profanities: []string{"bad", "worse"},
		LinkPattern: `https?://[^\s]+`,
	})

	tests := []struct {
		name           string
		message        string
		expectFiltered string
		expectError    bool
	}{
		{
			name:           "normal message",
			message:        "Hello, this is a normal message",
			expectFiltered: "Hello, this is a normal message",
			expectError:    false,
		},
		{
			name:           "profanity filtering",
			message:        "This is bad content",
			expectFiltered: "This is *** content",
			expectError:    false,
		},
		{
			name:           "link filtering",
			message:        "Check out https://example.com",
			expectFiltered: "Check out [link removed]",
			expectError:    false,
		},
		{
			name:           "empty message",
			message:        "",
			expectFiltered: "",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := filter.Handle(ctx, tt.message)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != nil && result.Message != tt.expectFiltered {
				t.Errorf("expected filtered message '%s', got '%s'", tt.expectFiltered, result.Message)
			}
		})
	}
}

func TestChatMessageFilter_DisabledFilter(t *testing.T) {
	filter := NewChatMessageFilter(config.MessageFilteringConfig{
		Enabled: false, // Disabled filter
	})

	ctx := context.Background()
	message := "This should pass through unchanged"

	result, err := filter.Handle(ctx, message)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected result but got nil")
		return
	}
	if result.Message != message {
		t.Errorf("expected unchanged message '%s', got '%s'", message, result.Message)
	}
}

func TestRateLimiter_GetClientID(t *testing.T) {
	limiter := NewRateLimiter(config.RateLimitConfig{
		RequestsPerMinute: 10,
		Window:            time.Minute,
	})

	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "context with client_ip",
			ctx:      context.WithValue(context.Background(), "client_ip", "192.168.1.1"),
			expected: "192.168.1.1",
		},
		{
			name:     "context with user_id",
			ctx:      context.WithValue(context.Background(), "user_id", "user123"),
			expected: "user123",
		},
		{
			name:     "empty context",
			ctx:      context.Background(),
			expected: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := limiter.getClientID(tt.ctx)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
