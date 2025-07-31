package models

import (
	"context"
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
