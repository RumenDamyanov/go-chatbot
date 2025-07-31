package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/RumenDamyanov/go-chatbot/config"
)

// GeminiModel implements the Model interface for Google's Gemini API.
type GeminiModel struct {
	config     config.GeminiConfig
	httpClient *http.Client
}

// NewGeminiModel creates a new Gemini model instance.
func NewGeminiModel(cfg config.GeminiConfig) (*GeminiModel, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("gemini API key is required")
	}
	if cfg.Model == "" {
		cfg.Model = "gemini-1.5-flash" // Default model
	}

	return &GeminiModel{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// geminiRequest represents the request structure for Gemini's API.
type geminiRequest struct {
	Contents         []geminiContent         `json:"contents"`
	GenerationConfig *geminiGenerationConfig `json:"generationConfig,omitempty"`
	SafetySettings   []geminiSafetySetting   `json:"safetySettings,omitempty"`
}

// geminiContent represents content in the request.
type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

// geminiPart represents a part of the content.
type geminiPart struct {
	Text string `json:"text"`
}

// geminiGenerationConfig represents generation configuration.
type geminiGenerationConfig struct {
	Temperature     float64  `json:"temperature,omitempty"`
	TopK            int      `json:"topK,omitempty"`
	TopP            float64  `json:"topP,omitempty"`
	MaxOutputTokens int      `json:"maxOutputTokens,omitempty"`
	StopSequences   []string `json:"stopSequences,omitempty"`
}

// geminiSafetySetting represents safety settings.
type geminiSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// geminiResponse represents the response from Gemini's API.
type geminiResponse struct {
	Candidates    []geminiCandidate   `json:"candidates"`
	UsageMetadata geminiUsageMetadata `json:"usageMetadata,omitempty"`
}

// geminiCandidate represents a response candidate.
type geminiCandidate struct {
	Content       geminiContent        `json:"content"`
	FinishReason  string               `json:"finishReason"`
	Index         int                  `json:"index"`
	SafetyRatings []geminySafetyRating `json:"safetyRatings,omitempty"`
}

// geminySafetyRating represents safety rating information.
type geminySafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

// geminiUsageMetadata represents usage information.
type geminiUsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

// geminiError represents an error response from the API.
type geminiError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

// Ask sends a message to Gemini and returns the response.
func (g *GeminiModel) Ask(ctx context.Context, message string, context map[string]interface{}) (string, error) {
	// Prepare the request
	req := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: message},
				},
			},
		},
		GenerationConfig: &geminiGenerationConfig{
			Temperature:     0.7,
			TopK:            40,
			TopP:            0.8,
			MaxOutputTokens: 1000,
		},
		SafetySettings: []geminiSafetySetting{
			{
				Category:  "HARM_CATEGORY_HARASSMENT",
				Threshold: "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				Category:  "HARM_CATEGORY_HATE_SPEECH",
				Threshold: "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				Category:  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
				Threshold: "BLOCK_MEDIUM_AND_ABOVE",
			},
			{
				Category:  "HARM_CATEGORY_DANGEROUS_CONTENT",
				Threshold: "BLOCK_MEDIUM_AND_ABOVE",
			},
		},
	}

	// Add conversation history if provided
	if history, ok := context["history"]; ok {
		if hist, ok := history.([]map[string]interface{}); ok {
			var contents []geminiContent
			for _, msg := range hist {
				if role, roleOk := msg["role"].(string); roleOk {
					if content, contentOk := msg["content"].(string); contentOk {
						// Gemini uses "user" and "model" roles
						geminiRole := role
						if role == "assistant" {
							geminiRole = "model"
						}
						contents = append(contents, geminiContent{
							Role:  geminiRole,
							Parts: []geminiPart{{Text: content}},
						})
					}
				}
			}
			// Add current message at the end
			contents = append(contents, geminiContent{
				Parts: []geminiPart{{Text: message}},
			})
			req.Contents = contents
		}
	}

	// Override generation config from context if provided
	if temp, ok := context["temperature"]; ok {
		if temperature, ok := temp.(float64); ok {
			req.GenerationConfig.Temperature = temperature
		}
	}
	if maxTokens, ok := context["max_tokens"]; ok {
		if tokens, ok := maxTokens.(int); ok {
			req.GenerationConfig.MaxOutputTokens = tokens
		}
	}

	// Marshal the request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Construct URL
	endpoint := "https://generativelanguage.googleapis.com"
	if g.config.Endpoint != "" {
		endpoint = g.config.Endpoint
	}
	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", endpoint, g.config.Model, g.config.APIKey)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		var errResp geminiError
		if err := json.Unmarshal(body, &errResp); err == nil {
			return "", fmt.Errorf("gemini API error: %s", errResp.Error.Message)
		}
		return "", fmt.Errorf("gemini API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var geminiResp geminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract the text content
	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	candidate := geminiResp.Candidates[0]
	if len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("no content parts in response")
	}

	var responseText strings.Builder
	for _, part := range candidate.Content.Parts {
		responseText.WriteString(part.Text)
	}

	if responseText.Len() == 0 {
		return "", fmt.Errorf("no text content in response")
	}

	return responseText.String(), nil
}

// Name returns the name of the model.
func (g *GeminiModel) Name() string {
	return g.config.Model
}

// Provider returns the provider name.
func (g *GeminiModel) Provider() string {
	return "gemini"
}

// Health checks if the Gemini API is accessible.
func (g *GeminiModel) Health(ctx context.Context) error {
	// Create a simple test request
	req := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: "Hello"},
				},
			},
		},
		GenerationConfig: &geminiGenerationConfig{
			MaxOutputTokens: 10,
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal health check request: %w", err)
	}

	endpoint := "https://generativelanguage.googleapis.com"
	if g.config.Endpoint != "" {
		endpoint = g.config.Endpoint
	}
	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", endpoint, g.config.Model, g.config.APIKey)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("invalid API key")
	}

	if resp.StatusCode >= 500 {
		return fmt.Errorf("gemini API server error: %d", resp.StatusCode)
	}

	return nil
}
