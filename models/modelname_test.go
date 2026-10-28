package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.rumenx.com/chatbot/config"
)

func TestConstructorsRejectUnsafeModelName(t *testing.T) {
	const unsafe = "../other"

	tests := []struct {
		name string
		new  func() error
	}{
		{
			name: "openai",
			new: func() error {
				_, err := NewOpenAIModel(config.OpenAIConfig{APIKey: "test-key", Model: unsafe})
				return err
			},
		},
		{
			name: "anthropic",
			new: func() error {
				_, err := NewAnthropicModel(config.AnthropicConfig{APIKey: "test-key", Model: unsafe})
				return err
			},
		},
		{
			name: "meta",
			new: func() error {
				_, err := NewMetaModel(config.MetaConfig{APIKey: "test-key", Model: unsafe})
				return err
			},
		},
		{
			name: "xai",
			new: func() error {
				_, err := NewXAIModel(config.XAIConfig{APIKey: "test-key", Model: unsafe})
				return err
			},
		},
		{
			name: "ollama",
			new: func() error {
				_, err := NewOllamaModel(config.OllamaConfig{Model: unsafe})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.new()
			require.Error(t, err)
			assert.ErrorContains(t, err, "not allowed")
		})
	}
}
