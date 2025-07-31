package embeddings

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RumenDamyanov/go-chatbot/config"
)

func TestEmbed_EmptyTexts(t *testing.T) {
	config := config.OpenAIConfig{APIKey: "test-key"}
	provider := NewOpenAIEmbeddingProvider(config, "text-embedding-3-small")

	ctx := context.Background()
	_, err := provider.Embed(ctx, []string{})

	if err == nil {
		t.Error("expected error for empty texts")
	}

	if err.Error() != "no texts provided" {
		t.Errorf("expected 'no texts provided' error, got: %v", err)
	}
}

func TestEmbedSingle_WithMockServer(t *testing.T) {
	// Create mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": {"message": "test error", "type": "invalid_request_error"}}`))
	}))
	defer server.Close()

	config := config.OpenAIConfig{
		APIKey:   "test-key",
		Endpoint: server.URL,
	}
	provider := NewOpenAIEmbeddingProvider(config, "text-embedding-3-small")

	ctx := context.Background()
	_, err := provider.EmbedSingle(ctx, "test")

	if err == nil {
		t.Error("expected error from mock server")
	}

	if !strings.Contains(err.Error(), "test error") {
		t.Errorf("expected error to contain 'test error', got: %v", err)
	}
}

func TestEmbed_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock response
		response := map[string]interface{}{
			"object": "list",
			"data": []map[string]interface{}{
				{
					"object":    "embedding",
					"index":     0,
					"embedding": []float64{0.1, 0.2, 0.3},
				},
				{
					"object":    "embedding",
					"index":     1,
					"embedding": []float64{0.4, 0.5, 0.6},
				},
			},
			"model": "text-embedding-3-small",
			"usage": map[string]interface{}{
				"prompt_tokens": 6,
				"total_tokens":  6,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := config.OpenAIConfig{
		APIKey:   "test-api-key",
		Endpoint: server.URL,
	}
	provider := NewOpenAIEmbeddingProvider(config, "text-embedding-3-small")

	ctx := context.Background()
	embeddings, err := provider.Embed(ctx, []string{"hello", "world"})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(embeddings) != 2 {
		t.Errorf("expected 2 embeddings, got %d", len(embeddings))
	}

	// Check first embedding
	if len(embeddings[0]) != 3 {
		t.Errorf("expected 3 dimensions, got %d", len(embeddings[0]))
	}
	if embeddings[0][0] != 0.1 || embeddings[0][1] != 0.2 || embeddings[0][2] != 0.3 {
		t.Errorf("unexpected first embedding values: %v", embeddings[0])
	}

	// Check second embedding
	if len(embeddings[1]) != 3 {
		t.Errorf("expected 3 dimensions, got %d", len(embeddings[1]))
	}
	if embeddings[1][0] != 0.4 || embeddings[1][1] != 0.5 || embeddings[1][2] != 0.6 {
		t.Errorf("unexpected second embedding values: %v", embeddings[1])
	}
}

func TestEmbedSingle_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"object": "list",
			"data": []map[string]interface{}{
				{
					"object":    "embedding",
					"index":     0,
					"embedding": []float64{0.7, 0.8, 0.9},
				},
			},
			"model": "text-embedding-3-small",
			"usage": map[string]interface{}{
				"prompt_tokens": 3,
				"total_tokens":  3,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := config.OpenAIConfig{
		APIKey:   "test-api-key",
		Endpoint: server.URL,
	}
	provider := NewOpenAIEmbeddingProvider(config, "text-embedding-3-small")

	ctx := context.Background()
	embedding, err := provider.EmbedSingle(ctx, "hello")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(embedding) != 3 {
		t.Errorf("expected 3 dimensions, got %d", len(embedding))
	}
	if embedding[0] != 0.7 || embedding[1] != 0.8 || embedding[2] != 0.9 {
		t.Errorf("unexpected embedding values: %v", embedding)
	}
}

func TestVectorStore_BasicOperations(t *testing.T) {
	config := config.OpenAIConfig{APIKey: "test-key"}
	provider := NewOpenAIEmbeddingProvider(config, "text-embedding-3-small")
	
	vectorStore := NewVectorStore(provider)
	if vectorStore == nil {
		t.Error("expected non-nil vector store")
	}

	// Test Count - should be 0 initially
	if vectorStore.Count() != 0 {
		t.Errorf("expected 0 items, got %d", vectorStore.Count())
	}

	// Test SetThreshold
	vectorStore.SetThreshold(0.5)

	// Test Clear
	vectorStore.Clear()
	if vectorStore.Count() != 0 {
		t.Errorf("expected 0 items after clear, got %d", vectorStore.Count())
	}
}

func TestVectorStore_AddTextWithMetadata(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"object": "list",
			"data": []map[string]interface{}{
				{
					"object":    "embedding",
					"index":     0,
					"embedding": []float64{0.1, 0.2, 0.3},
				},
			},
			"model": "text-embedding-3-small",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := config.OpenAIConfig{
		APIKey:   "test-api-key",
		Endpoint: server.URL,
	}
	provider := NewOpenAIEmbeddingProvider(config, "text-embedding-3-small")
	vectorStore := NewVectorStore(provider)

	ctx := context.Background()
	metadata := map[string]interface{}{"id": "1", "source": "test"}
	err := vectorStore.AddText(ctx, "test text", metadata)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if vectorStore.Count() != 1 {
		t.Errorf("expected 1 item, got %d", vectorStore.Count())
	}
}

func TestVectorStore_AddTextsWithMetadata(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"object": "list",
			"data": []map[string]interface{}{
				{
					"object":    "embedding",
					"index":     0,
					"embedding": []float64{0.1, 0.2, 0.3},
				},
				{
					"object":    "embedding",
					"index":     1,
					"embedding": []float64{0.4, 0.5, 0.6},
				},
			},
			"model": "text-embedding-3-small",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := config.OpenAIConfig{
		APIKey:   "test-api-key",
		Endpoint: server.URL,
	}
	provider := NewOpenAIEmbeddingProvider(config, "text-embedding-3-small")
	vectorStore := NewVectorStore(provider)

	ctx := context.Background()
	metadata1 := map[string]interface{}{"id": "1"}
	metadata2 := map[string]interface{}{"id": "2"}
	err := vectorStore.AddTexts(ctx, []string{"text1", "text2"}, []map[string]interface{}{metadata1, metadata2})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if vectorStore.Count() != 2 {
		t.Errorf("expected 2 items, got %d", vectorStore.Count())
	}
}

func TestMathFunctions(t *testing.T) {
	vec1 := []float64{1.0, 2.0, 3.0}
	vec2 := []float64{4.0, 5.0, 6.0}

	// Test DotProduct
	dotProduct := DotProduct(vec1, vec2)
	expected := 1.0*4.0 + 2.0*5.0 + 3.0*6.0 // 4 + 10 + 18 = 32
	if dotProduct != expected {
		t.Errorf("expected dot product %f, got %f", expected, dotProduct)
	}

	// Test CosineSimilarity
	similarity := CosineSimilarity(vec1, vec2)
	if similarity < 0 || similarity > 1 {
		t.Errorf("cosine similarity should be between 0 and 1, got %f", similarity)
	}

	// Test EuclideanDistance
	distance := EuclideanDistance(vec1, vec2)
	if distance < 0 {
		t.Errorf("euclidean distance should be non-negative, got %f", distance)
	}

	// Test Normalize
	normalized := Normalize(vec1)
	if len(normalized) != len(vec1) {
		t.Errorf("normalized vector should have same length as original")
	}
	
	// Check that normalized vector has magnitude 1
	magnitude := 0.0
	for _, v := range normalized {
		magnitude += v * v
	}
	magnitude = DotProduct(normalized, normalized) // Using our DotProduct function
	if magnitude < 0.99 || magnitude > 1.01 { // Allow small floating point error
		t.Errorf("normalized vector should have magnitude ~1, got %f", magnitude)
	}
}

func TestVectorStore_SearchFunctionality(t *testing.T) {
	// Create mock server for both adding and searching
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"object": "list",
			"data": []map[string]interface{}{
				{
					"object":    "embedding",
					"index":     0,
					"embedding": []float64{0.1, 0.2, 0.3},
				},
			},
			"model": "text-embedding-3-small",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := config.OpenAIConfig{
		APIKey:   "test-api-key",
		Endpoint: server.URL,
	}
	provider := NewOpenAIEmbeddingProvider(config, "text-embedding-3-small")
	vectorStore := NewVectorStore(provider)

	ctx := context.Background()
	
	// Add some texts first
	metadata := map[string]interface{}{"id": "1", "content": "test content"}
	err := vectorStore.AddText(ctx, "test text", metadata)
	if err != nil {
		t.Errorf("unexpected error adding text: %v", err)
	}

	// Now search
	results, err := vectorStore.Search(ctx, "query text", 1)
	if err != nil {
		t.Errorf("unexpected error searching: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestVectorStore_SearchEmptyStore(t *testing.T) {
	config := config.OpenAIConfig{APIKey: "test-key"}
	provider := NewOpenAIEmbeddingProvider(config, "text-embedding-3-small")
	vectorStore := NewVectorStore(provider)

	ctx := context.Background()
	_, err := vectorStore.Search(ctx, "query", 1)

	if err == nil {
		t.Error("expected error for empty vector store")
	}

	if err.Error() != "vector store is empty" {
		t.Errorf("expected 'vector store is empty' error, got: %v", err)
	}
}
