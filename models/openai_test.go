package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/RumenDamyanov/go-chatbot/config"
)

func TestNewOpenAIModel(t *testing.T) {
	tests := []struct {
		name        string
		config      config.OpenAIConfig
		expectError bool
	}{
		{
			name: "valid config",
			config: config.OpenAIConfig{
				APIKey: "test-key",
				Model:  "gpt-3.5-turbo",
			},
			expectError: false,
		},
		{
			name: "missing API key",
			config: config.OpenAIConfig{
				Model: "gpt-3.5-turbo",
			},
			expectError: true,
		},
		{
			name: "default model",
			config: config.OpenAIConfig{
				APIKey: "test-key",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewOpenAIModel(tt.config)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if model == nil {
				t.Error("expected model but got nil")
			}
		})
	}
}

func TestOpenAIModel_Ask_InvalidKey(t *testing.T) {
	model, err := NewOpenAIModel(config.OpenAIConfig{
		APIKey: "invalid-key",
		Model:  "gpt-3.5-turbo",
	})
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = model.Ask(ctx, "test message", nil)
	if err == nil {
		t.Error("expected error with invalid API key")
	}
}

func TestOpenAIModel_Ask_ContextCancellation(t *testing.T) {
	model, err := NewOpenAIModel(config.OpenAIConfig{
		APIKey: "test-key",
		Model:  "gpt-3.5-turbo",
	})
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err = model.Ask(ctx, "test message", nil)
	if err == nil {
		t.Error("expected error with cancelled context")
	}
}

func TestOpenAIModel_Health_InvalidKey(t *testing.T) {
	model, err := NewOpenAIModel(config.OpenAIConfig{
		APIKey: "invalid-key",
		Model:  "gpt-3.5-turbo",
	})
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = model.Health(ctx)
	if err == nil {
		t.Error("expected error with invalid API key")
	}
	t.Logf("Health check result: %v", err)
}

func TestOpenAIModel_AskStream_InvalidKey(t *testing.T) {
	model, err := NewOpenAIModel(config.OpenAIConfig{
		APIKey: "invalid-key",
		Model:  "gpt-3.5-turbo",
	})
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch, err := model.AskStream(ctx, "test message", nil)
	if err == nil {
		t.Error("expected error with invalid API key")
	}
	if ch != nil {
		t.Error("expected nil channel with error")
	}
}

func TestOpenAIModel_AskStream_Success(t *testing.T) {
	// Create a mock server for streaming
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json")
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got '%s'", authHeader)
		}

		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Send streaming response chunks
		chunks := []string{
			`data: {"choices":[{"delta":{"content":"Hello"}}],"id":"chatcmpl-test"}`,
			`data: {"choices":[{"delta":{"content":" world"}}],"id":"chatcmpl-test"}`,
			`data: {"choices":[{"delta":{"content":"!"}}],"id":"chatcmpl-test"}`,
			`data: [DONE]`,
		}

		for _, chunk := range chunks {
			w.Write([]byte(chunk + "\n\n"))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			time.Sleep(10 * time.Millisecond)
		}
	}))
	defer server.Close()

	cfg := config.OpenAIConfig{
		APIKey:   "test-key",
		Model:    "gpt-3.5-turbo",
		Endpoint: server.URL,
	}

	model, err := NewOpenAIModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx := context.Background()
	ch, err := model.AskStream(ctx, "Hello", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ch == nil {
		t.Fatal("Expected channel, got nil")
	}

	// Collect all chunks
	var result strings.Builder
	for chunk := range ch {
		result.WriteString(chunk)
	}

	expected := "Hello world!"
	if result.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result.String())
	}
}

func TestOpenAIModel_AskStream_ContextCancellation(t *testing.T) {
	// Create a mock server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte(`data: {"choices":[{"delta":{"content":"test"}}]}`))
	}))
	defer server.Close()

	cfg := config.OpenAIConfig{
		APIKey:   "test-key",
		Model:    "gpt-3.5-turbo",
		Endpoint: server.URL,
	}

	model, err := NewOpenAIModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	ch, err := model.AskStream(ctx, "Hello", nil)
	if err == nil {
		t.Fatal("Expected error for context cancellation")
	}

	if ch != nil {
		t.Error("Expected nil channel with error")
	}
}

func TestOpenAIModel_AskStream_HTTPError(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	cfg := config.OpenAIConfig{
		APIKey:   "test-key",
		Model:    "gpt-3.5-turbo",
		Endpoint: server.URL,
	}

	model, err := NewOpenAIModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx := context.Background()
	ch, err := model.AskStream(ctx, "Hello", nil)
	if err == nil {
		t.Fatal("Expected error for HTTP failure")
	}

	if ch != nil {
		t.Error("Expected nil channel with error")
	}
}

func TestOpenAIModel_ConversationHistory(t *testing.T) {
	model, err := NewOpenAIModel(config.OpenAIConfig{
		APIKey: "invalid-key", // Will fail but we test the structure
		Model:  "gpt-3.5-turbo",
	})
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	contextData := map[string]interface{}{
		"conversation_id": "test-123",
		"user_id":         "user-456",
	}

	_, err = model.Ask(ctx, "test message", contextData)
	// Should fail due to invalid key, but tests the context handling
	if err == nil {
		t.Error("expected error with invalid API key")
	}
}

func TestOpenAIModel_Name(t *testing.T) {
	model, err := NewOpenAIModel(config.OpenAIConfig{
		APIKey: "test-key",
		Model:  "gpt-4",
	})
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	if model.Name() != "gpt-4" {
		t.Errorf("expected model name 'gpt-4', got '%s'", model.Name())
	}
}

func TestOpenAIModel_Provider(t *testing.T) {
	model, err := NewOpenAIModel(config.OpenAIConfig{
		APIKey: "test-key",
		Model:  "gpt-3.5-turbo",
	})
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	if model.Provider() != "openai" {
		t.Errorf("expected provider 'openai', got '%s'", model.Provider())
	}
}

func TestExtractOpenAIStreamContent(t *testing.T) {
	tests := []struct {
		name     string
		chunk    map[string]interface{}
		expected string
	}{
		{
			name: "valid OpenAI stream chunk",
			chunk: map[string]interface{}{
				"choices": []interface{}{
					map[string]interface{}{
						"delta": map[string]interface{}{
							"content": "Hello world",
						},
					},
				},
			},
			expected: "Hello world",
		},
		{
			name: "chunk with empty content",
			chunk: map[string]interface{}{
				"choices": []interface{}{
					map[string]interface{}{
						"delta": map[string]interface{}{
							"content": "",
						},
					},
				},
			},
			expected: "",
		},
		{
			name: "chunk missing choices",
			chunk: map[string]interface{}{
				"id": "test-id",
			},
			expected: "",
		},
		{
			name: "chunk with empty choices array",
			chunk: map[string]interface{}{
				"choices": []interface{}{},
			},
			expected: "",
		},
		{
			name: "chunk with invalid choices type",
			chunk: map[string]interface{}{
				"choices": "not an array",
			},
			expected: "",
		},
		{
			name: "chunk with invalid choice type",
			chunk: map[string]interface{}{
				"choices": []interface{}{
					"not a map",
				},
			},
			expected: "",
		},
		{
			name: "chunk missing delta",
			chunk: map[string]interface{}{
				"choices": []interface{}{
					map[string]interface{}{
						"index": 0,
					},
				},
			},
			expected: "",
		},
		{
			name: "chunk with invalid delta type",
			chunk: map[string]interface{}{
				"choices": []interface{}{
					map[string]interface{}{
						"delta": "not a map",
					},
				},
			},
			expected: "",
		},
		{
			name: "chunk delta missing content",
			chunk: map[string]interface{}{
				"choices": []interface{}{
					map[string]interface{}{
						"delta": map[string]interface{}{
							"role": "assistant",
						},
					},
				},
			},
			expected: "",
		},
		{
			name: "chunk with non-string content",
			chunk: map[string]interface{}{
				"choices": []interface{}{
					map[string]interface{}{
						"delta": map[string]interface{}{
							"content": 123,
						},
					},
				},
			},
			expected: "",
		},
		{
			name: "multiple choices - uses first",
			chunk: map[string]interface{}{
				"choices": []interface{}{
					map[string]interface{}{
						"delta": map[string]interface{}{
							"content": "First choice",
						},
					},
					map[string]interface{}{
						"delta": map[string]interface{}{
							"content": "Second choice",
						},
					},
				},
			},
			expected: "First choice",
		},
		{
			name:     "empty chunk",
			chunk:    map[string]interface{}{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractOpenAIStreamContent(tt.chunk)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
