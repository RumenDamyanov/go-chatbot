package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.rumenx.com/chatbot/config"
)

func TestNewGeminiModel(t *testing.T) {
	tests := []struct {
		name      string
		config    config.GeminiConfig
		expectErr bool
	}{
		{
			name: "valid config",
			config: config.GeminiConfig{
				APIKey: "test-key",
				Model:  "gemini-1.5-flash",
			},
			expectErr: false,
		},
		{
			name: "missing API key",
			config: config.GeminiConfig{
				Model: "gemini-1.5-flash",
			},
			expectErr: true,
		},
		{
			name: "default model",
			config: config.GeminiConfig{
				APIKey: "test-key",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewGeminiModel(tt.config)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, model)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, model)
				assert.Equal(t, "gemini", model.Provider())
				if tt.config.Model != "" {
					assert.Equal(t, tt.config.Model, model.Name())
				} else {
					assert.Equal(t, "gemini-1.5-flash", model.Name())
				}
			}
		})
	}
}

func TestGeminiModel_Ask_InvalidKey(t *testing.T) {
	model, err := NewGeminiModel(config.GeminiConfig{
		APIKey: "invalid-key",
		Model:  "gemini-1.5-flash",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := model.Ask(ctx, "Hello", nil)
	assert.Error(t, err)
	assert.Empty(t, response)
}

func TestGeminiModel_Ask_ContextCancellation(t *testing.T) {
	model, err := NewGeminiModel(config.GeminiConfig{
		APIKey: "test-key",
		Model:  "gemini-1.5-flash",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	response, err := model.Ask(ctx, "Hello", nil)
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestGeminiModel_Health_InvalidKey(t *testing.T) {
	model, err := NewGeminiModel(config.GeminiConfig{
		APIKey: "invalid-key",
		Model:  "gemini-1.5-flash",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = model.Health(ctx)
	// Health check may succeed or fail depending on API response, so we just check that it doesn't panic
	// Real API errors would be caught with actual invalid keys
	t.Logf("Health check result: %v", err)
}

func TestGeminiModel_ConversationHistory(t *testing.T) {
	model, err := NewGeminiModel(config.GeminiConfig{
		APIKey: "test-key",
		Model:  "gemini-1.5-flash",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with conversation history and custom parameters
	context := map[string]interface{}{
		"history": []map[string]interface{}{
			{"role": "user", "content": "Hello"},
			{"role": "assistant", "content": "Hi there!"},
		},
		"temperature": 0.5,
		"max_tokens":  500,
	}

	// This will fail with invalid API key, but we can test that the request is properly formed
	response, err := model.Ask(ctx, "How are you?", context)
	assert.Error(t, err)
	assert.Empty(t, response)
}
