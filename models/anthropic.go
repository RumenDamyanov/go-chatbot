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

// AnthropicModel implements the Model interface for Anthropic's Claude API.
type AnthropicModel struct {
	config     config.AnthropicConfig
	httpClient *http.Client
	maxTokens  int
}

// NewAnthropicModel creates a new Anthropic model instance.
func NewAnthropicModel(cfg config.AnthropicConfig) (*AnthropicModel, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("anthropic API key is required")
	}
	if cfg.Model == "" {
		cfg.Model = "claude-3-haiku-20240307" // Default model
	}

	return &AnthropicModel{
		config:    cfg,
		maxTokens: 1000, // Default max tokens
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// anthropicRequest represents the request structure for Anthropic's API.
type anthropicRequest struct {
	Model     string                 `json:"model"`
	MaxTokens int                    `json:"max_tokens"`
	Messages  []anthropicMessage     `json:"messages"`
	System    string                 `json:"system,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// anthropicMessage represents a message in the conversation.
type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// anthropicResponse represents the response from Anthropic's API.
type anthropicResponse struct {
	ID           string             `json:"id"`
	Type         string             `json:"type"`
	Role         string             `json:"role"`
	Content      []anthropicContent `json:"content"`
	Model        string             `json:"model"`
	Usage        anthropicUsage     `json:"usage"`
	StopReason   string             `json:"stop_reason"`
	StopSequence string             `json:"stop_sequence,omitempty"`
}

// anthropicContent represents content in the response.
type anthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// anthropicUsage represents token usage information.
type anthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// anthropicError represents an error response from the API.
type anthropicError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Ask sends a message to Claude and returns the response.
func (a *AnthropicModel) Ask(ctx context.Context, message string, context map[string]interface{}) (string, error) {
	// Prepare the request
	req := anthropicRequest{
		Model:     a.config.Model,
		MaxTokens: a.maxTokens,
		Messages: []anthropicMessage{
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	// Add system message if provided in context
	if systemMsg, ok := context["system"]; ok {
		if sys, ok := systemMsg.(string); ok {
			req.System = sys
		}
	}

	// Add conversation history if provided
	if history, ok := context["history"]; ok {
		if hist, ok := history.([]map[string]interface{}); ok {
			var messages []anthropicMessage
			for _, msg := range hist {
				if role, roleOk := msg["role"].(string); roleOk {
					if content, contentOk := msg["content"].(string); contentOk {
						// Claude uses "user" and "assistant" roles
						if role == "user" || role == "assistant" {
							messages = append(messages, anthropicMessage{
								Role:    role,
								Content: content,
							})
						}
					}
				}
			}
			// Add current message at the end
			messages = append(messages, anthropicMessage{
				Role:    "user",
				Content: message,
			})
			req.Messages = messages
		}
	}

	// Marshal the request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", a.config.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	// Send the request
	resp, err := a.httpClient.Do(httpReq)
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
		var errResp anthropicError
		if err := json.Unmarshal(body, &errResp); err == nil {
			return "", fmt.Errorf("anthropic API error: %s", errResp.Message)
		}
		return "", fmt.Errorf("anthropic API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var anthropicResp anthropicResponse
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract the text content
	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	var responseText strings.Builder
	for _, content := range anthropicResp.Content {
		if content.Type == "text" {
			responseText.WriteString(content.Text)
		}
	}

	if responseText.Len() == 0 {
		return "", fmt.Errorf("no text content in response")
	}

	return responseText.String(), nil
}

// Name returns the name of the model.
func (a *AnthropicModel) Name() string {
	return a.config.Model
}

// Provider returns the provider name.
func (a *AnthropicModel) Provider() string {
	return "anthropic"
}

// Health checks if the Anthropic API is accessible.
func (a *AnthropicModel) Health(ctx context.Context) error {
	// Create a simple test request
	req := anthropicRequest{
		Model:     a.config.Model,
		MaxTokens: 10,
		Messages: []anthropicMessage{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal health check request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", a.config.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid API key")
	}

	if resp.StatusCode >= 500 {
		return fmt.Errorf("anthropic API server error: %d", resp.StatusCode)
	}

	return nil
}
