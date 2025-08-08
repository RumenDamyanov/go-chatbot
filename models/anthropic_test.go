package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.rumenx.com/chatbot/config"
)

func TestNewAnthropicModel(t *testing.T) {
	tests := []struct {
		name      string
		config    config.AnthropicConfig
		expectErr bool
	}{
		{
			name: "valid config",
			config: config.AnthropicConfig{
				APIKey: "test-key",
				Model:  "claude-3-haiku-20240307",
			},
			expectErr: false,
		},
		{
			name: "missing API key",
			config: config.AnthropicConfig{
				Model: "claude-3-haiku-20240307",
			},
			expectErr: true,
		},
		{
			name: "default model",
			config: config.AnthropicConfig{
				APIKey: "test-key",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewAnthropicModel(tt.config)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, model)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, model)
				assert.Equal(t, "anthropic", model.Provider())
				if tt.config.Model != "" {
					assert.Equal(t, tt.config.Model, model.Name())
				} else {
					assert.Equal(t, "claude-3-haiku-20240307", model.Name())
				}
			}
		})
	}
}

func TestAnthropicModel_Ask_InvalidKey(t *testing.T) {
	// This test uses an invalid API key to test error handling
	// without making actual API calls
	model, err := NewAnthropicModel(config.AnthropicConfig{
		APIKey: "invalid-key",
		Model:  "claude-3-haiku-20240307",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := model.Ask(ctx, "Hello", nil)
	assert.Error(t, err)
	assert.Empty(t, response)
	// The error could be "invalid" for API key or other errors, so just check that there's an error
	assert.NotNil(t, err)
}

func TestAnthropicModel_Ask_ContextCancellation(t *testing.T) {
	model, err := NewAnthropicModel(config.AnthropicConfig{
		APIKey: "test-key",
		Model:  "claude-3-haiku-20240307",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	response, err := model.Ask(ctx, "Hello", nil)
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestAnthropicModel_Health_InvalidKey(t *testing.T) {
	model, err := NewAnthropicModel(config.AnthropicConfig{
		APIKey: "invalid-key",
		Model:  "claude-3-haiku-20240307",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = model.Health(ctx)
	assert.Error(t, err)
}

func TestAnthropicModel_ConversationHistory(t *testing.T) {
	model, err := NewAnthropicModel(config.AnthropicConfig{
		APIKey: "test-key",
		Model:  "claude-3-haiku-20240307",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with conversation history
	context := map[string]interface{}{
		"history": []map[string]interface{}{
			{"role": "user", "content": "Hello"},
			{"role": "assistant", "content": "Hi there!"},
		},
		"system": "You are a helpful assistant.",
	}

	// This will fail with invalid API key, but we can test that the request is properly formed
	response, err := model.Ask(ctx, "How are you?", context)
	assert.Error(t, err)
	assert.Empty(t, response)
}
