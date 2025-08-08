// Package models provides AI model implementations for the go-chatbot package.
package models

import (
	"context"
	"errors"
	"fmt"

	"go.rumenx.com/chatbot/config"
)

// Model represents an AI model that can process chat messages.
type Model interface {
	// Ask sends a message to the AI model and returns the response.
	Ask(ctx context.Context, message string, context map[string]interface{}) (string, error)

	// Name returns the name of the model.
	Name() string

	// Provider returns the provider of the model.
	Provider() string
}

// HealthChecker is an optional interface that models can implement
// to provide health check functionality.
type HealthChecker interface {
	Health(ctx context.Context) error
}

// StreamingModel is an optional interface for models that support streaming responses.
type StreamingModel interface {
	AskStream(ctx context.Context, message string, context map[string]interface{}) (<-chan string, error)
}

// ModelFactory creates AI models based on configuration.
type ModelFactory struct{}

// NewFromConfig creates a new AI model based on the configuration.
func NewFromConfig(cfg *config.Config) (Model, error) {
	if cfg == nil {
		return nil, errors.New("config cannot be nil")
	}

	switch cfg.Model {
	case "openai":
		return NewOpenAIModel(cfg.OpenAI)
	case "anthropic":
		return NewAnthropicModel(cfg.Anthropic)
	case "gemini":
		return NewGeminiModel(cfg.Gemini)
	case "xai":
		return NewXAIModel(cfg.XAI)
	case "meta":
		return NewMetaModel(cfg.Meta)
	case "ollama":
		return NewOllamaModel(cfg.Ollama)
	case "free":
		return NewFreeModel(), nil
	default:
		return nil, fmt.Errorf("unsupported model: %s", cfg.Model)
	}
}

// Registry holds available model constructors.
type Registry struct {
	constructors map[string]func(interface{}) (Model, error)
}

// NewRegistry creates a new model registry.
func NewRegistry() *Registry {
	return &Registry{
		constructors: make(map[string]func(interface{}) (Model, error)),
	}
}

// Register registers a model constructor.
func (r *Registry) Register(name string, constructor func(interface{}) (Model, error)) {
	r.constructors[name] = constructor
}

// Create creates a model instance by name.
func (r *Registry) Create(name string, config interface{}) (Model, error) {
	constructor, exists := r.constructors[name]
	if !exists {
		return nil, fmt.Errorf("unknown model: %s", name)
	}
	return constructor(config)
}

// ListAvailable returns a list of available model names.
func (r *Registry) ListAvailable() []string {
	names := make([]string, 0, len(r.constructors))
	for name := range r.constructors {
		names = append(names, name)
	}
	return names
}

// Default registry instance.
var DefaultRegistry = NewRegistry()

func init() {
	// Register default models
	DefaultRegistry.Register("openai", func(cfg interface{}) (Model, error) {
		if openaiCfg, ok := cfg.(config.OpenAIConfig); ok {
			return NewOpenAIModel(openaiCfg)
		}
		return nil, errors.New("invalid OpenAI config")
	})

	DefaultRegistry.Register("anthropic", func(cfg interface{}) (Model, error) {
		if anthropicCfg, ok := cfg.(config.AnthropicConfig); ok {
			return NewAnthropicModel(anthropicCfg)
		}
		return nil, errors.New("invalid Anthropic config")
	})

	DefaultRegistry.Register("gemini", func(cfg interface{}) (Model, error) {
		if geminiCfg, ok := cfg.(config.GeminiConfig); ok {
			return NewGeminiModel(geminiCfg)
		}
		return nil, errors.New("invalid Gemini config")
	})

	DefaultRegistry.Register("xai", func(cfg interface{}) (Model, error) {
		if xaiCfg, ok := cfg.(config.XAIConfig); ok {
			return NewXAIModel(xaiCfg)
		}
		return nil, errors.New("invalid xAI config")
	})

	DefaultRegistry.Register("meta", func(cfg interface{}) (Model, error) {
		if metaCfg, ok := cfg.(config.MetaConfig); ok {
			return NewMetaModel(metaCfg)
		}
		return nil, errors.New("invalid Meta config")
	})

	DefaultRegistry.Register("ollama", func(cfg interface{}) (Model, error) {
		if ollamaCfg, ok := cfg.(config.OllamaConfig); ok {
			return NewOllamaModel(ollamaCfg)
		}
		return nil, errors.New("invalid Ollama config")
	})

	DefaultRegistry.Register("free", func(cfg interface{}) (Model, error) {
		return NewFreeModel(), nil
	})
}
