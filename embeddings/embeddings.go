// Package embeddings provides text embedding functionality for the go-chatbot package.
package embeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/RumenDamyanov/go-chatbot/config"
)

// Vector represents an embedding vector.
type Vector []float64

// EmbeddingResponse represents the response from an embedding API.
type EmbeddingResponse struct {
	Embeddings []Vector `json:"embeddings"`
	Error      string   `json:"error,omitempty"`
}

// EmbeddingProvider defines the interface for embedding providers.
type EmbeddingProvider interface {
	// Embed generates embeddings for the given texts.
	Embed(ctx context.Context, texts []string) ([]Vector, error)

	// EmbedSingle generates an embedding for a single text.
	EmbedSingle(ctx context.Context, text string) (Vector, error)

	// Dimensions returns the dimensionality of the embeddings.
	Dimensions() int

	// Model returns the model name.
	Model() string

	// Provider returns the provider name.
	Provider() string
}

// OpenAIEmbeddingProvider implements embedding using OpenAI's API.
type OpenAIEmbeddingProvider struct {
	config     config.OpenAIConfig
	httpClient *http.Client
	model      string
	dimensions int
}

// OpenAIEmbeddingRequest represents a request to OpenAI's embedding API.
type OpenAIEmbeddingRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

// OpenAIEmbeddingResponse represents OpenAI's embedding API response.
type OpenAIEmbeddingResponse struct {
	Data []struct {
		Embedding Vector `json:"embedding"`
		Index     int    `json:"index"`
	} `json:"data"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// NewOpenAIEmbeddingProvider creates a new OpenAI embedding provider.
func NewOpenAIEmbeddingProvider(cfg config.OpenAIConfig, model string) *OpenAIEmbeddingProvider {
	if model == "" {
		model = "text-embedding-3-small" // Default OpenAI embedding model
	}

	dimensions := 1536 // Default for text-embedding-3-small
	if model == "text-embedding-3-large" {
		dimensions = 3072
	}

	return &OpenAIEmbeddingProvider{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		model:      model,
		dimensions: dimensions,
	}
}

// Embed generates embeddings for multiple texts.
func (p *OpenAIEmbeddingProvider) Embed(ctx context.Context, texts []string) ([]Vector, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts provided")
	}

	// OpenAI has a limit of 2048 texts per request, batch if necessary
	const batchSize = 2048
	var allEmbeddings []Vector

	for i := 0; i < len(texts); i += batchSize {
		end := i + batchSize
		if end > len(texts) {
			end = len(texts)
		}

		batch := texts[i:end]
		embeddings, err := p.embedBatch(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("failed to embed batch %d-%d: %w", i, end, err)
		}

		allEmbeddings = append(allEmbeddings, embeddings...)
	}

	return allEmbeddings, nil
}

// EmbedSingle generates an embedding for a single text.
func (p *OpenAIEmbeddingProvider) EmbedSingle(ctx context.Context, text string) (Vector, error) {
	embeddings, err := p.Embed(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	return embeddings[0], nil
}

// embedBatch processes a batch of texts.
func (p *OpenAIEmbeddingProvider) embedBatch(ctx context.Context, texts []string) ([]Vector, error) {
	// Build request
	request := OpenAIEmbeddingRequest{
		Input: texts,
		Model: p.model,
	}

	// Marshal request
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	endpoint := "https://api.openai.com/v1/embeddings"
	if p.config.Endpoint != "" && p.config.Endpoint != "https://api.openai.com/v1/chat/completions" {
		// Use custom endpoint if provided and not the chat endpoint
		endpoint = p.config.Endpoint + "/embeddings"
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)

	// Send request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var openaiResp OpenAIEmbeddingResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if openaiResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", openaiResp.Error.Message)
	}

	// Extract embeddings
	embeddings := make([]Vector, len(openaiResp.Data))
	for _, data := range openaiResp.Data {
		if data.Index < len(embeddings) {
			embeddings[data.Index] = data.Embedding
		}
	}

	return embeddings, nil
}

// Dimensions returns the dimensionality of the embeddings.
func (p *OpenAIEmbeddingProvider) Dimensions() int {
	return p.dimensions
}

// Model returns the model name.
func (p *OpenAIEmbeddingProvider) Model() string {
	return p.model
}

// Provider returns the provider name.
func (p *OpenAIEmbeddingProvider) Provider() string {
	return "openai"
}

// VectorStore provides vector storage and similarity search functionality.
type VectorStore struct {
	vectors   []Vector
	metadata  []map[string]interface{}
	provider  EmbeddingProvider
	threshold float64
}

// NewVectorStore creates a new vector store.
func NewVectorStore(provider EmbeddingProvider) *VectorStore {
	return &VectorStore{
		provider:  provider,
		threshold: 0.7, // Default similarity threshold
	}
}

// AddTexts adds texts to the vector store.
func (vs *VectorStore) AddTexts(ctx context.Context, texts []string, metadata []map[string]interface{}) error {
	if len(texts) != len(metadata) {
		return fmt.Errorf("texts and metadata length mismatch: %d vs %d", len(texts), len(metadata))
	}

	// Generate embeddings
	embeddings, err := vs.provider.Embed(ctx, texts)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Add to store
	vs.vectors = append(vs.vectors, embeddings...)
	vs.metadata = append(vs.metadata, metadata...)

	return nil
}

// AddText adds a single text to the vector store.
func (vs *VectorStore) AddText(ctx context.Context, text string, metadata map[string]interface{}) error {
	return vs.AddTexts(ctx, []string{text}, []map[string]interface{}{metadata})
}

// Search finds similar texts in the vector store.
func (vs *VectorStore) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	if len(vs.vectors) == 0 {
		return nil, fmt.Errorf("vector store is empty")
	}

	// Generate query embedding
	queryVector, err := vs.provider.EmbedSingle(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// Calculate similarities
	similarities := make([]SearchResult, len(vs.vectors))
	for i, vector := range vs.vectors {
		similarity := CosineSimilarity(queryVector, vector)
		similarities[i] = SearchResult{
			Index:      i,
			Similarity: similarity,
			Metadata:   vs.metadata[i],
		}
	}

	// Sort by similarity (descending)
	for i := 0; i < len(similarities)-1; i++ {
		for j := i + 1; j < len(similarities); j++ {
			if similarities[i].Similarity < similarities[j].Similarity {
				similarities[i], similarities[j] = similarities[j], similarities[i]
			}
		}
	}

	// Apply threshold and limit
	var results []SearchResult
	for _, result := range similarities {
		if result.Similarity >= vs.threshold && len(results) < limit {
			results = append(results, result)
		}
	}

	return results, nil
}

// SearchResult represents a search result from the vector store.
type SearchResult struct {
	Index      int                    `json:"index"`
	Similarity float64                `json:"similarity"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// SetThreshold sets the similarity threshold for search results.
func (vs *VectorStore) SetThreshold(threshold float64) {
	vs.threshold = threshold
}

// Count returns the number of vectors in the store.
func (vs *VectorStore) Count() int {
	return len(vs.vectors)
}

// Clear removes all vectors from the store.
func (vs *VectorStore) Clear() {
	vs.vectors = nil
	vs.metadata = nil
}

// CosineSimilarity calculates the cosine similarity between two vectors.
func CosineSimilarity(a, b Vector) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// EuclideanDistance calculates the Euclidean distance between two vectors.
func EuclideanDistance(a, b Vector) float64 {
	if len(a) != len(b) {
		return math.Inf(1)
	}

	var sum float64
	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

// DotProduct calculates the dot product of two vectors.
func DotProduct(a, b Vector) float64 {
	if len(a) != len(b) {
		return 0
	}

	var product float64
	for i := 0; i < len(a); i++ {
		product += a[i] * b[i]
	}

	return product
}

// Normalize normalizes a vector to unit length.
func Normalize(v Vector) Vector {
	var norm float64
	for _, val := range v {
		norm += val * val
	}
	norm = math.Sqrt(norm)

	if norm == 0 {
		return v
	}

	normalized := make(Vector, len(v))
	for i, val := range v {
		normalized[i] = val / norm
	}

	return normalized
}
