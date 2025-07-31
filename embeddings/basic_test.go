package embeddings

import (
	"testing"

	"github.com/RumenDamyanov/go-chatbot/config"
)

func TestOpenAIEmbeddingProvider_Basic(t *testing.T) {
	config := config.OpenAIConfig{APIKey: "test-key"}
	provider := NewOpenAIEmbeddingProvider(config, "test-model")

	if provider == nil {
		t.Error("expected non-nil provider")
	}

	if provider.Model() != "test-model" {
		t.Errorf("expected model 'test-model', got '%s'", provider.Model())
	}

	if provider.Provider() != "openai" {
		t.Errorf("expected provider 'openai', got '%s'", provider.Provider())
	}

	// Test dimensions for different models
	provider1 := NewOpenAIEmbeddingProvider(config, "text-embedding-3-small")
	if provider1.Dimensions() != 1536 {
		t.Errorf("expected 1536 dimensions for text-embedding-3-small, got %d", provider1.Dimensions())
	}

	provider2 := NewOpenAIEmbeddingProvider(config, "text-embedding-3-large")
	if provider2.Dimensions() != 3072 {
		t.Errorf("expected 3072 dimensions for text-embedding-3-large, got %d", provider2.Dimensions())
	}
}
