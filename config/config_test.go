package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	assert.NotNil(t, cfg)
	assert.Equal(t, "free", cfg.Model)
	assert.Equal(t, "You are a helpful, friendly chatbot.", cfg.Prompt)
	assert.Equal(t, "en", cfg.Language)
	assert.Equal(t, "neutral", cfg.Tone)
	assert.True(t, cfg.Emojis)
	assert.True(t, cfg.Deescalate)
	assert.False(t, cfg.Funny)
	assert.Equal(t, 30*time.Second, cfg.Timeout)
	assert.Equal(t, 256, cfg.MaxTokens)
	assert.Equal(t, 0.7, cfg.Temperature)
}

func TestDefaultWithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("CHATBOT_MODEL", "openai")
	os.Setenv("CHATBOT_PROMPT", "Test prompt")
	os.Setenv("OPENAI_API_KEY", "test-key")
	defer func() {
		os.Unsetenv("CHATBOT_MODEL")
		os.Unsetenv("CHATBOT_PROMPT")
		os.Unsetenv("OPENAI_API_KEY")
	}()

	cfg := Default()

	assert.Equal(t, "openai", cfg.Model)
	assert.Equal(t, "Test prompt", cfg.Prompt)
	assert.Equal(t, "test-key", cfg.OpenAI.APIKey)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errType error
	}{
		{
			name:    "valid free model config",
			config:  Default(),
			wantErr: false,
		},
		{
			name: "invalid empty model",
			config: &Config{
				Model:       "",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 0.7,
			},
			wantErr: true,
			errType: ErrInvalidModel,
		},
		{
			name: "invalid timeout",
			config: &Config{
				Model:       "free",
				Timeout:     0,
				MaxTokens:   256,
				Temperature: 0.7,
			},
			wantErr: true,
			errType: ErrInvalidTimeout,
		},
		{
			name: "invalid max tokens",
			config: &Config{
				Model:       "free",
				Timeout:     30 * time.Second,
				MaxTokens:   0,
				Temperature: 0.7,
			},
			wantErr: true,
			errType: ErrInvalidMaxTokens,
		},
		{
			name: "invalid temperature",
			config: &Config{
				Model:       "free",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 3.0,
			},
			wantErr: true,
			errType: ErrInvalidTemperature,
		},
		{
			name: "openai without api key",
			config: &Config{
				Model:       "openai",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 0.7,
				OpenAI: OpenAIConfig{
					APIKey: "",
				},
			},
			wantErr: true,
			errType: ErrMissingAPIKey,
		},
		{
			name: "anthropic without api key",
			config: &Config{
				Model:       "anthropic",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 0.7,
				Anthropic: AnthropicConfig{
					APIKey: "",
				},
			},
			wantErr: true,
			errType: ErrMissingAPIKey,
		},
		{
			name: "gemini without api key",
			config: &Config{
				Model:       "gemini",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 0.7,
				Gemini: GeminiConfig{
					APIKey: "",
				},
			},
			wantErr: true,
			errType: ErrMissingAPIKey,
		},
		{
			name: "xai without api key",
			config: &Config{
				Model:       "xai",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 0.7,
				XAI: XAIConfig{
					APIKey: "",
				},
			},
			wantErr: true,
			errType: ErrMissingAPIKey,
		},
		{
			name: "meta without api key",
			config: &Config{
				Model:       "meta",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 0.7,
				Meta: MetaConfig{
					APIKey: "",
				},
			},
			wantErr: true,
			errType: ErrMissingAPIKey,
		},
		{
			name: "ollama without endpoint",
			config: &Config{
				Model:       "ollama",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 0.7,
				Ollama: OllamaConfig{
					Endpoint: "",
				},
			},
			wantErr: true,
			errType: ErrMissingEndpoint,
		},
		{
			name: "unsupported model",
			config: &Config{
				Model:       "unknown-model",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 0.7,
			},
			wantErr: true,
			errType: ErrUnsupportedModel,
		},
		{
			name: "valid anthropic config",
			config: &Config{
				Model:       "anthropic",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 0.7,
				Anthropic: AnthropicConfig{
					APIKey: "test-key",
				},
			},
			wantErr: false,
		},
		{
			name: "valid temperature boundary low",
			config: &Config{
				Model:       "free",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 0.0,
			},
			wantErr: false,
		},
		{
			name: "valid temperature boundary high",
			config: &Config{
				Model:       "free",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: 2.0,
			},
			wantErr: false,
		},
		{
			name: "invalid temperature too low",
			config: &Config{
				Model:       "free",
				Timeout:     30 * time.Second,
				MaxTokens:   256,
				Temperature: -0.1,
			},
			wantErr: true,
			errType: ErrInvalidTemperature,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetEnvHelpers(t *testing.T) {
	t.Run("getEnv", func(t *testing.T) {
		os.Setenv("TEST_STRING", "test_value")
		defer os.Unsetenv("TEST_STRING")

		result := getEnv("TEST_STRING", "default")
		assert.Equal(t, "test_value", result)

		result = getEnv("NON_EXISTENT", "default")
		assert.Equal(t, "default", result)
	})

	t.Run("getIntEnv", func(t *testing.T) {
		os.Setenv("TEST_INT", "42")
		defer os.Unsetenv("TEST_INT")

		result := getIntEnv("TEST_INT", 0)
		assert.Equal(t, 42, result)

		result = getIntEnv("NON_EXISTENT", 10)
		assert.Equal(t, 10, result)
	})

	t.Run("getBoolEnv", func(t *testing.T) {
		os.Setenv("TEST_BOOL", "true")
		defer os.Unsetenv("TEST_BOOL")

		result := getBoolEnv("TEST_BOOL", false)
		assert.True(t, result)

		result = getBoolEnv("NON_EXISTENT", false)
		assert.False(t, result)
	})

	t.Run("getFloatEnv", func(t *testing.T) {
		os.Setenv("TEST_FLOAT", "3.14")
		defer os.Unsetenv("TEST_FLOAT")

		result := getFloatEnv("TEST_FLOAT", 0.0)
		assert.Equal(t, 3.14, result)

		result = getFloatEnv("NON_EXISTENT", 2.5)
		assert.Equal(t, 2.5, result)

		// Test invalid float value
		os.Setenv("INVALID_FLOAT", "not-a-number")
		defer os.Unsetenv("INVALID_FLOAT")
		result = getFloatEnv("INVALID_FLOAT", 1.0)
		assert.Equal(t, 1.0, result) // Should return default on parse error
	})

	t.Run("getDurationEnv", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "5s")
		defer os.Unsetenv("TEST_DURATION")

		result := getDurationEnv("TEST_DURATION", time.Minute)
		assert.Equal(t, 5*time.Second, result)

		result = getDurationEnv("NON_EXISTENT", time.Minute)
		assert.Equal(t, time.Minute, result)
	})
}
