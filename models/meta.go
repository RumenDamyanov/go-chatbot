package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RumenDamyanov/go-chatbot/config"
)

// MetaModel implements the Model interface for Meta's LLaMA API.
type MetaModel struct {
	config     config.MetaConfig
	httpClient *http.Client
}

// NewMetaModel creates a new Meta model instance.
func NewMetaModel(cfg config.MetaConfig) (*MetaModel, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("Meta API key is required")
	}
	if cfg.Model == "" {
		cfg.Model = "llama-3.2-3b-instruct" // Default model
	}

	return &MetaModel{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// metaRequest represents the request structure for Meta's API.
// Meta LLaMA uses OpenAI-compatible API format
type metaRequest struct {
	Model       string        `json:"model"`
	Messages    []metaMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	Stop        []string      `json:"stop,omitempty"`
}

// metaMessage represents a message in the conversation.
type metaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// metaResponse represents the response from Meta's API.
type metaResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []metaChoice `json:"choices"`
	Usage   metaUsage    `json:"usage"`
}

// metaChoice represents a choice in the response.
type metaChoice struct {
	Index        int         `json:"index"`
	Message      metaMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// metaUsage represents token usage information.
type metaUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// metaError represents an error response from the API.
type metaError struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Ask sends a message to Meta LLaMA and returns the response.
func (m *MetaModel) Ask(ctx context.Context, message string, context map[string]interface{}) (string, error) {
	// Prepare the request
	req := metaRequest{
		Model: m.config.Model,
		Messages: []metaMessage{
			{
				Role:    "user",
				Content: message,
			},
		},
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        1.0,
		Stream:      false,
	}

	// Add conversation history if provided
	if history, ok := context["history"]; ok {
		if hist, ok := history.([]map[string]interface{}); ok {
			var messages []metaMessage
			for _, msg := range hist {
				if role, roleOk := msg["role"].(string); roleOk {
					if content, contentOk := msg["content"].(string); contentOk {
						// Meta uses OpenAI-compatible roles: "user", "assistant", "system"
						messages = append(messages, metaMessage{
							Role:    role,
							Content: content,
						})
					}
				}
			}
			// Add current message at the end
			messages = append(messages, metaMessage{
				Role:    "user",
				Content: message,
			})
			req.Messages = messages
		}
	}

	// Add system message if provided
	if systemMsg, ok := context["system"]; ok {
		if sys, ok := systemMsg.(string); ok {
			// Prepend system message
			req.Messages = append([]metaMessage{
				{
					Role:    "system",
					Content: sys,
				},
			}, req.Messages...)
		}
	}

	// Override parameters from context if provided
	if temp, ok := context["temperature"]; ok {
		if temperature, ok := temp.(float64); ok {
			req.Temperature = temperature
		}
	}
	if maxTokens, ok := context["max_tokens"]; ok {
		if tokens, ok := maxTokens.(int); ok {
			req.MaxTokens = tokens
		}
	}
	if topP, ok := context["top_p"]; ok {
		if tp, ok := topP.(float64); ok {
			req.TopP = tp
		}
	}
	if stop, ok := context["stop"]; ok {
		if stopSequences, ok := stop.([]string); ok {
			req.Stop = stopSequences
		}
	}

	// Marshal the request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Construct URL - Meta LLaMA is often accessed through platforms like Replicate or Together AI
	endpoint := "https://api.llama-api.com" // Default endpoint (hypothetical)
	if m.config.Endpoint != "" {
		endpoint = m.config.Endpoint
	}
	url := fmt.Sprintf("%s/v1/chat/completions", endpoint)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+m.config.APIKey)

	// Send the request
	resp, err := m.httpClient.Do(httpReq)
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
		var errResp metaError
		if err := json.Unmarshal(body, &errResp); err == nil {
			return "", fmt.Errorf("Meta API error: %s", errResp.Error.Message)
		}
		return "", fmt.Errorf("Meta API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var metaResp metaResponse
	if err := json.Unmarshal(body, &metaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract the text content
	if len(metaResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	choice := metaResp.Choices[0]
	if choice.Message.Content == "" {
		return "", fmt.Errorf("no content in response message")
	}

	return choice.Message.Content, nil
}

// Name returns the name of the model.
func (m *MetaModel) Name() string {
	return m.config.Model
}

// Provider returns the provider name.
func (m *MetaModel) Provider() string {
	return "meta"
}

// Health checks if the Meta API is accessible.
func (m *MetaModel) Health(ctx context.Context) error {
	// Create a simple test request
	req := metaRequest{
		Model: m.config.Model,
		Messages: []metaMessage{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		MaxTokens: 10,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal health check request: %w", err)
	}

	endpoint := "https://api.llama-api.com"
	if m.config.Endpoint != "" {
		endpoint = m.config.Endpoint
	}
	url := fmt.Sprintf("%s/v1/chat/completions", endpoint)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+m.config.APIKey)

	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid API key")
	}

	if resp.StatusCode >= 500 {
		return fmt.Errorf("Meta API server error: %d", resp.StatusCode)
	}

	return nil
}
