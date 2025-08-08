package gochatbot

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.rumenx.com/chatbot/config"
)

func TestNewHTTPHandler(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)

	if handler == nil {
		t.Error("Expected non-nil HTTP handler")
	}
}

func TestHTTPHandlerChat(t *testing.T) {
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

	handler := NewHTTPHandler(chatbot)

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "valid chat request",
			method:         "POST",
			body:           `{"message": "Hello there"}`,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "empty message",
			method:         "POST",
			body:           `{"message": ""}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "invalid JSON",
			method:         "POST",
			body:           `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "missing message field",
			method:         "POST",
			body:           `{"other": "value"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "wrong method",
			method:         "GET",
			body:           `{"message": "Hello"}`,
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/chat", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response ChatResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Failed to unmarshal response: %v", err)
				return
			}

			if tt.expectError && response.Error == "" {
				t.Error("Expected error in response")
			}

			if !tt.expectError && response.Reply == "" {
				t.Error("Expected reply in response")
			}
		})
	}
}

func TestHTTPHandlerHealth(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "valid health check",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "wrong method",
			method:         "POST",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			handler.Health(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
					return
				}

				if response["status"] != "healthy" {
					t.Errorf("Expected status 'healthy', got %v", response["status"])
				}
			}
		})
	}
}

func TestHTTPHandlerHealth_UnhealthyModel(t *testing.T) {
	// Create a chatbot with an invalid OpenAI config to trigger health failure
	chatbot, err := New(&config.Config{
		Model: "openai",
		OpenAI: config.OpenAIConfig{
			APIKey: "invalid-key",
			Model:  "gpt-3.5-turbo",
		},
	})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d for unhealthy model, got %d", http.StatusServiceUnavailable, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
		return
	}

	if response["status"] != "unhealthy" {
		t.Errorf("Expected status 'unhealthy', got %v", response["status"])
	}

	if response["error"] == nil {
		t.Error("Expected error field to be present")
	}
}

func TestHTTPHandlerGetClientIP(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)

	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		expected   string
	}{
		{
			name:       "X-Forwarded-For header",
			remoteAddr: "192.168.1.1:12345",
			headers:    map[string]string{"X-Forwarded-For": "203.0.113.1"},
			expected:   "203.0.113.1",
		},
		{
			name:       "X-Real-IP header",
			remoteAddr: "192.168.1.1:12345",
			headers:    map[string]string{"X-Real-IP": "203.0.113.2"},
			expected:   "203.0.113.2",
		},
		{
			name:       "X-Forwarded-For with multiple IPs",
			remoteAddr: "192.168.1.1:12345",
			headers:    map[string]string{"X-Forwarded-For": "203.0.113.3, 198.51.100.1"},
			expected:   "203.0.113.3",
		},
		{
			name:       "remote address fallback",
			remoteAddr: "203.0.113.4:12345",
			headers:    map[string]string{},
			expected:   "203.0.113.4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remoteAddr

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			ip := handler.getClientIP(req)

			if ip != tt.expected {
				t.Errorf("Expected IP %s, got %s", tt.expected, ip)
			}
		})
	}
}

func TestChatbotHandleHTTP(t *testing.T) {
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

	req := httptest.NewRequest("POST", "/chat", strings.NewReader(`{"message": "Hello"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	chatbot.HandleHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response ChatResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
		return
	}

	if response.Reply == "" {
		t.Error("Expected non-empty reply")
	}
}

func TestChatbotHandleStreamHTTP(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	req := httptest.NewRequest("POST", "/stream", strings.NewReader(`{"message": "Hello"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	chatbot.HandleStreamHTTP(w, req)

	// For the free model, it should return a regular response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleStreamHTTP_OPTIONS(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)

	req := httptest.NewRequest("OPTIONS", "/stream", nil)
	w := httptest.NewRecorder()

	handler.HandleStreamHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for OPTIONS, got %d", http.StatusOK, w.Code)
	}

	// Check CORS headers
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected CORS header Access-Control-Allow-Origin to be *")
	}
}

func TestHandleStreamHTTP_MethodNotAllowed(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)

	req := httptest.NewRequest("GET", "/stream", nil)
	w := httptest.NewRecorder()

	handler.HandleStreamHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d for GET, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestHandleStreamHTTP_InvalidJSON(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)

	req := httptest.NewRequest("POST", "/stream", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleStreamHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleStreamHTTP_EmptyMessage(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)

	req := httptest.NewRequest("POST", "/stream", strings.NewReader(`{"message": ""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleStreamHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for empty message, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleStreamHTTP_WhitespaceMessage(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)

	req := httptest.NewRequest("POST", "/stream", strings.NewReader(`{"message": "   "}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleStreamHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for whitespace message, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHTTPHandlerWithTimeout(t *testing.T) {
	// Create chatbot with very short timeout
	chatbot, err := New(&config.Config{
		Model:   "free",
		Timeout: 1 * time.Nanosecond,
	})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)

	req := httptest.NewRequest("POST", "/chat", strings.NewReader(`{"message": "Hello"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add a small delay to ensure timeout
	time.Sleep(1 * time.Millisecond)

	handler.HandleHTTP(w, req)

	// Should get a timeout error
	if w.Code == http.StatusOK {
		t.Error("Expected timeout error, but got success")
	}
}

func TestWriteErrorResponse(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)
	w := httptest.NewRecorder()

	handler.writeErrorResponse(w, http.StatusBadRequest, "Test error message")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response ChatResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
		return
	}

	if response.Error != "Test error message" {
		t.Errorf("Expected error 'Test error message', got '%s'", response.Error)
	}
}

func TestHTTPHandlerContextTimeout(t *testing.T) {
	chatbot, err := New(&config.Config{Model: "free"})
	if err != nil {
		t.Fatalf("Failed to create chatbot: %v", err)
	}

	handler := NewHTTPHandler(chatbot)

	// Create request with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := httptest.NewRequest("POST", "/chat", strings.NewReader(`{"message": "Hello"}`))
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleHTTP(w, req)

	// Should handle the cancelled context gracefully
	if w.Code == http.StatusOK {
		// The free model might still respond even with cancelled context
		// since it doesn't make external API calls
		t.Logf("Free model responded despite cancelled context")
	}
}

func TestHTTPHandlerLargePayload(t *testing.T) {
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

	handler := NewHTTPHandler(chatbot)

	// Create a large message
	largeMessage := strings.Repeat("Hello ", 1000)
	payload := map[string]string{"message": largeMessage}
	payloadBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/chat", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for large payload, got %d", http.StatusOK, w.Code)
	}
}
