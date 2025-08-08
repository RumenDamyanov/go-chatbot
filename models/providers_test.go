package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.rumenx.com/chatbot/config"
)

func TestNewXAIModel(t *testing.T) {
	tests := []struct {
		name      string
		config    config.XAIConfig
		expectErr bool
	}{
		{
			name: "valid config",
			config: config.XAIConfig{
				APIKey: "test-key",
				Model:  "grok-beta",
			},
			expectErr: false,
		},
		{
			name: "missing API key",
			config: config.XAIConfig{
				Model: "grok-beta",
			},
			expectErr: true,
		},
		{
			name: "default model",
			config: config.XAIConfig{
				APIKey: "test-key",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewXAIModel(tt.config)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, model)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, model)
				assert.Equal(t, "xai", model.Provider())
				if tt.config.Model != "" {
					assert.Equal(t, tt.config.Model, model.Name())
				} else {
					assert.Equal(t, "grok-beta", model.Name())
				}
			}
		})
	}
}

func TestNewMetaModel(t *testing.T) {
	tests := []struct {
		name      string
		config    config.MetaConfig
		expectErr bool
	}{
		{
			name: "valid config",
			config: config.MetaConfig{
				APIKey: "test-key",
				Model:  "llama-3.2-3b-instruct",
			},
			expectErr: false,
		},
		{
			name: "missing API key",
			config: config.MetaConfig{
				Model: "llama-3.2-3b-instruct",
			},
			expectErr: true,
		},
		{
			name: "default model",
			config: config.MetaConfig{
				APIKey: "test-key",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewMetaModel(tt.config)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, model)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, model)
				assert.Equal(t, "meta", model.Provider())
				if tt.config.Model != "" {
					assert.Equal(t, tt.config.Model, model.Name())
				} else {
					assert.Equal(t, "llama-3.2-3b-instruct", model.Name())
				}
			}
		})
	}
}

func TestNewOllamaModel(t *testing.T) {
	tests := []struct {
		name      string
		config    config.OllamaConfig
		expectErr bool
	}{
		{
			name: "valid config",
			config: config.OllamaConfig{
				Model: "llama3.2",
			},
			expectErr: false,
		},
		{
			name:      "default model",
			config:    config.OllamaConfig{},
			expectErr: false,
		},
		{
			name: "with endpoint",
			config: config.OllamaConfig{
				Model:    "mistral",
				Endpoint: "http://localhost:11434",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewOllamaModel(tt.config)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, model)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, model)
				assert.Equal(t, "ollama", model.Provider())
				if tt.config.Model != "" {
					assert.Equal(t, tt.config.Model, model.Name())
				} else {
					assert.Equal(t, "llama3.2", model.Name())
				}
			}
		})
	}
}

func TestXAIModel_ContextCancellation(t *testing.T) {
	model, err := NewXAIModel(config.XAIConfig{
		APIKey: "test-key",
		Model:  "grok-beta",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	response, err := model.Ask(ctx, "Hello", nil)
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestMetaModel_ContextCancellation(t *testing.T) {
	model, err := NewMetaModel(config.MetaConfig{
		APIKey: "test-key",
		Model:  "llama-3.2-3b-instruct",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	response, err := model.Ask(ctx, "Hello", nil)
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestOllamaModel_ContextCancellation(t *testing.T) {
	model, err := NewOllamaModel(config.OllamaConfig{
		Model: "llama3.2",
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	response, err := model.Ask(ctx, "Hello", nil)
	assert.Error(t, err)
	assert.Empty(t, response)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestOllamaModel_Health_NotRunning(t *testing.T) {
	model, err := NewOllamaModel(config.OllamaConfig{
		Model:    "llama3.2",
		Endpoint: "http://localhost:9999", // Non-existent endpoint
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = model.Health(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "health check failed")
}
