// Package config provides configuration management for the go-chatbot package.
package config

import (
	"os"
	"strconv"
	"time"
)

// Config represents the main configuration for the chatbot.
type Config struct {
	// AI Model Configuration
	Model string `json:"model" yaml:"model"`

	// OpenAI Configuration
	OpenAI OpenAIConfig `json:"openai" yaml:"openai"`

	// Anthropic Configuration
	Anthropic AnthropicConfig `json:"anthropic" yaml:"anthropic"`

	// Google Gemini Configuration
	Gemini GeminiConfig `json:"gemini" yaml:"gemini"`

	// xAI Configuration
	XAI XAIConfig `json:"xai" yaml:"xai"`

	// Meta Configuration
	Meta MetaConfig `json:"meta" yaml:"meta"`

	// Ollama Configuration
	Ollama OllamaConfig `json:"ollama" yaml:"ollama"`

	// Chatbot Behavior
	Prompt   string `json:"prompt" yaml:"prompt"`
	Language string `json:"language" yaml:"language"`
	Tone     string `json:"tone" yaml:"tone"`

	// Security and Rate Limiting
	RateLimit        RateLimitConfig        `json:"rate_limit" yaml:"rate_limit"`
	MessageFiltering MessageFilteringConfig `json:"message_filtering" yaml:"message_filtering"`

	// Request Configuration
	Timeout     time.Duration `json:"timeout" yaml:"timeout"`
	MaxTokens   int           `json:"max_tokens" yaml:"max_tokens"`
	Temperature float64       `json:"temperature" yaml:"temperature"`

	// Feature Flags
	Emojis     bool `json:"emojis" yaml:"emojis"`
	Deescalate bool `json:"deescalate" yaml:"deescalate"`
	Funny      bool `json:"funny" yaml:"funny"`

	// Allowed Scripts
	AllowedScripts []string `json:"allowed_scripts" yaml:"allowed_scripts"`
}

// OpenAIConfig contains OpenAI-specific configuration.
type OpenAIConfig struct {
	APIKey   string `json:"api_key" yaml:"api_key"`
	Model    string `json:"model" yaml:"model"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// AnthropicConfig contains Anthropic-specific configuration.
type AnthropicConfig struct {
	APIKey   string `json:"api_key" yaml:"api_key"`
	Model    string `json:"model" yaml:"model"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// GeminiConfig contains Google Gemini-specific configuration.
type GeminiConfig struct {
	APIKey   string `json:"api_key" yaml:"api_key"`
	Model    string `json:"model" yaml:"model"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// XAIConfig contains xAI-specific configuration.
type XAIConfig struct {
	APIKey   string `json:"api_key" yaml:"api_key"`
	Model    string `json:"model" yaml:"model"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// MetaConfig contains Meta-specific configuration.
type MetaConfig struct {
	APIKey   string `json:"api_key" yaml:"api_key"`
	Model    string `json:"model" yaml:"model"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// OllamaConfig contains Ollama-specific configuration.
type OllamaConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Model    string `json:"model" yaml:"model"`
}

// RateLimitConfig contains rate limiting configuration.
type RateLimitConfig struct {
	RequestsPerMinute int           `json:"requests_per_minute" yaml:"requests_per_minute"`
	BurstSize         int           `json:"burst_size" yaml:"burst_size"`
	Window            time.Duration `json:"window" yaml:"window"`
}

// MessageFilteringConfig contains message filtering configuration.
type MessageFilteringConfig struct {
	Instructions       []string `json:"instructions" yaml:"instructions"`
	Profanities        []string `json:"profanities" yaml:"profanities"`
	AggressionPatterns []string `json:"aggression_patterns" yaml:"aggression_patterns"`
	LinkPattern        string   `json:"link_pattern" yaml:"link_pattern"`
	Enabled            bool     `json:"enabled" yaml:"enabled"`
}

// Default returns a default configuration with environment variable overrides.
func Default() *Config {
	return &Config{
		Model: getEnv("CHATBOT_MODEL", "free"),
		OpenAI: OpenAIConfig{
			APIKey:   getEnv("OPENAI_API_KEY", ""),
			Model:    getEnv("OPENAI_MODEL", "gpt-4o"),
			Endpoint: getEnv("OPENAI_ENDPOINT", "https://api.openai.com/v1/chat/completions"),
		},
		Anthropic: AnthropicConfig{
			APIKey:   getEnv("ANTHROPIC_API_KEY", ""),
			Model:    getEnv("ANTHROPIC_MODEL", "claude-3-sonnet-20240229"),
			Endpoint: getEnv("ANTHROPIC_ENDPOINT", "https://api.anthropic.com/v1/messages"),
		},
		Gemini: GeminiConfig{
			APIKey:   getEnv("GEMINI_API_KEY", ""),
			Model:    getEnv("GEMINI_MODEL", "gemini-1.5-pro"),
			Endpoint: getEnv("GEMINI_ENDPOINT", "https://generativelanguage.googleapis.com/v1beta/models"),
		},
		XAI: XAIConfig{
			APIKey:   getEnv("XAI_API_KEY", ""),
			Model:    getEnv("XAI_MODEL", "grok-1"),
			Endpoint: getEnv("XAI_ENDPOINT", "https://api.x.ai/v1/chat/completions"),
		},
		Meta: MetaConfig{
			APIKey:   getEnv("META_API_KEY", ""),
			Model:    getEnv("META_MODEL", "llama-3-70b"),
			Endpoint: getEnv("META_ENDPOINT", "https://api.meta.ai/v1/chat/completions"),
		},
		Ollama: OllamaConfig{
			Endpoint: getEnv("OLLAMA_ENDPOINT", "http://localhost:11434/api/chat"),
			Model:    getEnv("OLLAMA_MODEL", "llama2"),
		},
		Prompt:      getEnv("CHATBOT_PROMPT", "You are a helpful, friendly chatbot."),
		Language:    getEnv("CHATBOT_LANGUAGE", "en"),
		Tone:        getEnv("CHATBOT_TONE", "neutral"),
		Timeout:     getDurationEnv("CHATBOT_TIMEOUT", 30*time.Second),
		MaxTokens:   getIntEnv("CHATBOT_MAX_TOKENS", 256),
		Temperature: getFloatEnv("CHATBOT_TEMPERATURE", 0.7),
		Emojis:      getBoolEnv("CHATBOT_EMOJIS", true),
		Deescalate:  getBoolEnv("CHATBOT_DEESCALATE", true),
		Funny:       getBoolEnv("CHATBOT_FUNNY", false),
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getIntEnv("RATE_LIMIT_REQUESTS", 10),
			BurstSize:         getIntEnv("RATE_LIMIT_BURST", 5),
			Window:            getDurationEnv("RATE_LIMIT_WINDOW", time.Minute),
		},
		MessageFiltering: MessageFilteringConfig{
			Instructions: []string{
				"Avoid sharing external links.",
				"Refrain from quoting controversial sources.",
				"Use appropriate language.",
				"Reject harmful or dangerous requests.",
				"De-escalate potential conflicts and calm aggressive or rude users.",
			},
			Profanities:        []string{},
			AggressionPatterns: []string{"hate", "kill", "stupid", "idiot"},
			LinkPattern:        `https?://[\w\.-]+`,
			Enabled:            getBoolEnv("FILTER_ENABLED", true),
		},
		AllowedScripts: []string{"Latin", "Cyrillic", "Greek", "Armenian", "Han", "Kana", "Hangul"},
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Model == "" {
		return ErrInvalidModel
	}

	if c.Timeout <= 0 {
		return ErrInvalidTimeout
	}

	if c.MaxTokens <= 0 {
		return ErrInvalidMaxTokens
	}

	if c.Temperature < 0 || c.Temperature > 2 {
		return ErrInvalidTemperature
	}

	// Validate model-specific configuration
	switch c.Model {
	case "openai":
		if c.OpenAI.APIKey == "" {
			return ErrMissingAPIKey
		}
	case "anthropic":
		if c.Anthropic.APIKey == "" {
			return ErrMissingAPIKey
		}
	case "gemini":
		if c.Gemini.APIKey == "" {
			return ErrMissingAPIKey
		}
	case "xai":
		if c.XAI.APIKey == "" {
			return ErrMissingAPIKey
		}
	case "meta":
		if c.Meta.APIKey == "" {
			return ErrMissingAPIKey
		}
	case "ollama":
		if c.Ollama.Endpoint == "" {
			return ErrMissingEndpoint
		}
	case "free":
		// No validation needed for free model
	default:
		return ErrUnsupportedModel
	}

	return nil
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getFloatEnv(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
