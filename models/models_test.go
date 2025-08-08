package models

import (
	"testing"

	"go.rumenx.com/chatbot/config"
)

func TestNewFromConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      config.Config
		expectError bool
		expectType  string
	}{
		{
			name: "free model",
			config: config.Config{
				Model: "free",
			},
			expectError: false,
			expectType:  "free-model",
		},
		{
			name: "openai model with key",
			config: config.Config{
				Model: "openai",
				OpenAI: config.OpenAIConfig{
					APIKey: "test-key",
					Model:  "gpt-3.5-turbo",
				},
			},
			expectError: false,
			expectType:  "gpt-3.5-turbo",
		},
		{
			name: "openai model without key",
			config: config.Config{
				Model: "openai",
				OpenAI: config.OpenAIConfig{
					Model: "gpt-3.5-turbo",
				},
			},
			expectError: true,
		},
		{
			name: "unknown model",
			config: config.Config{
				Model: "unknown",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := NewFromConfig(&tt.config)
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
				return
			}
			if model.Name() != tt.expectType {
				t.Errorf("expected model type %s, got %s", tt.expectType, model.Name())
			}
		})
	}
}

func TestRegistry_Create(t *testing.T) {
	registry := NewRegistry()

	// Register a test model
	registry.Register("test", func(cfg interface{}) (Model, error) {
		return NewFreeModel(), nil
	})

	tests := []struct {
		name        string
		modelName   string
		config      interface{}
		expectError bool
	}{
		{
			name:        "existing model",
			modelName:   "test",
			config:      nil,
			expectError: false,
		},
		{
			name:        "non-existing model",
			modelName:   "nonexistent",
			config:      nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := registry.Create(tt.modelName, tt.config)
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

func TestRegistry_ListAvailable(t *testing.T) {
	registry := NewRegistry()

	// Should be empty initially
	models := registry.ListAvailable()
	if len(models) != 0 {
		t.Errorf("expected empty list, got %d models", len(models))
	}

	// Register some models
	registry.Register("model1", func(cfg interface{}) (Model, error) {
		return NewFreeModel(), nil
	})
	registry.Register("model2", func(cfg interface{}) (Model, error) {
		return NewFreeModel(), nil
	})

	models = registry.ListAvailable()
	if len(models) != 2 {
		t.Errorf("expected 2 models, got %d", len(models))
	}

	// Check if both models are present
	found1, found2 := false, false
	for _, model := range models {
		if model == "model1" {
			found1 = true
		}
		if model == "model2" {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Error("not all registered models found in list")
	}
}

func TestDefaultRegistry(t *testing.T) {
	// Test that default models are registered
	availableModels := DefaultRegistry.ListAvailable()

	expectedModels := []string{"openai", "anthropic", "gemini", "xai", "meta", "ollama", "free"}

	if len(availableModels) < len(expectedModels) {
		t.Errorf("expected at least %d models, got %d", len(expectedModels), len(availableModels))
	}

	// Check that each expected model is available
	for _, expected := range expectedModels {
		found := false
		for _, available := range availableModels {
			if available == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected model %s not found in default registry", expected)
		}
	}
}

func TestDefaultRegistry_CreateFreeModel(t *testing.T) {
	model, err := DefaultRegistry.Create("free", nil)
	if err != nil {
		t.Errorf("unexpected error creating free model: %v", err)
		return
	}
	if model == nil {
		t.Error("expected model but got nil")
		return
	}
	if model.Name() != "free-model" {
		t.Errorf("expected model name 'free-model', got '%s'", model.Name())
	}
}

func TestDefaultRegistry_CreateOpenAIModel(t *testing.T) {
	config := config.OpenAIConfig{
		APIKey: "test-key",
		Model:  "gpt-3.5-turbo",
	}

	model, err := DefaultRegistry.Create("openai", config)
	if err != nil {
		t.Errorf("unexpected error creating openai model: %v", err)
		return
	}
	if model == nil {
		t.Error("expected model but got nil")
		return
	}
	if model.Provider() != "openai" {
		t.Errorf("expected provider 'openai', got '%s'", model.Provider())
	}
}

func TestDefaultRegistry_CreateInvalidConfig(t *testing.T) {
	// Try to create openai model with wrong config type
	_, err := DefaultRegistry.Create("openai", "invalid config")
	if err == nil {
		t.Error("expected error with invalid config type")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
