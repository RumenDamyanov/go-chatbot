package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.rumenx.com/chatbot/config"
)

// XAIModel implements the Model interface for xAI's Grok API.
type XAIModel struct {
	config     config.XAIConfig
	httpClient *http.Client
}

// NewXAIModel creates a new xAI model instance.
func NewXAIModel(cfg config.XAIConfig) (*XAIModel, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("xAI API key is required")
	}
	if cfg.Model == "" {
		cfg.Model = "grok-beta" // Default model
	}

	return &XAIModel{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// xaiRequest represents the request structure for xAI's API.
// xAI uses OpenAI-compatible API format
type xaiRequest struct {
	Model       string       `json:"model"`
	Messages    []xaiMessage `json:"messages"`
	MaxTokens   int          `json:"max_tokens,omitempty"`
	Temperature float64      `json:"temperature,omitempty"`
	TopP        float64      `json:"top_p,omitempty"`
	Stream      bool         `json:"stream,omitempty"`
}

// xaiMessage represents a message in the conversation.
type xaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// xaiResponse represents the response from xAI's API.
type xaiResponse struct {
	ID      string      `json:"id"`
	Object  string      `json:"object"`
	Created int64       `json:"created"`
	Model   string      `json:"model"`
	Choices []xaiChoice `json:"choices"`
	Usage   xaiUsage    `json:"usage"`
}

// xaiChoice represents a choice in the response.
type xaiChoice struct {
	Index        int        `json:"index"`
	Message      xaiMessage `json:"message"`
	FinishReason string     `json:"finish_reason"`
}

// xaiUsage represents token usage information.
type xaiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// xaiError represents an error response from the API.
type xaiError struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Ask sends a message to xAI Grok and returns the response.
func (x *XAIModel) Ask(ctx context.Context, message string, context map[string]interface{}) (string, error) {
	// Prepare the request
	req := xaiRequest{
		Model: x.config.Model,
		Messages: []xaiMessage{
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
			var messages []xaiMessage
			for _, msg := range hist {
				if role, roleOk := msg["role"].(string); roleOk {
					if content, contentOk := msg["content"].(string); contentOk {
						// xAI uses OpenAI-compatible roles: "user", "assistant", "system"
						messages = append(messages, xaiMessage{
							Role:    role,
							Content: content,
						})
					}
				}
			}
			// Add current message at the end
			messages = append(messages, xaiMessage{
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
			req.Messages = append([]xaiMessage{
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

	// Marshal the request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Construct URL
	endpoint := "https://api.x.ai"
	if x.config.Endpoint != "" {
		endpoint = x.config.Endpoint
	}
	url := fmt.Sprintf("%s/v1/chat/completions", endpoint)

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+x.config.APIKey)

	// Send the request
	resp, err := x.httpClient.Do(httpReq)
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
		var errResp xaiError
		if err := json.Unmarshal(body, &errResp); err == nil {
			return "", fmt.Errorf("xAI API error: %s", errResp.Error.Message)
		}
		return "", fmt.Errorf("xAI API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var xaiResp xaiResponse
	if err := json.Unmarshal(body, &xaiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract the text content
	if len(xaiResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	choice := xaiResp.Choices[0]
	if choice.Message.Content == "" {
		return "", fmt.Errorf("no content in response message")
	}

	return choice.Message.Content, nil
}

// Name returns the name of the model.
func (x *XAIModel) Name() string {
	return x.config.Model
}

// Provider returns the provider name.
func (x *XAIModel) Provider() string {
	return "xai"
}

// Health checks if the xAI API is accessible.
func (x *XAIModel) Health(ctx context.Context) error {
	// Create a simple test request
	req := xaiRequest{
		Model: x.config.Model,
		Messages: []xaiMessage{
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

	endpoint := "https://api.x.ai"
	if x.config.Endpoint != "" {
		endpoint = x.config.Endpoint
	}
	url := fmt.Sprintf("%s/v1/chat/completions", endpoint)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+x.config.APIKey)

	resp, err := x.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid API key")
	}

	if resp.StatusCode >= 500 {
		return fmt.Errorf("xAI API server error: %d", resp.StatusCode)
	}

	return nil
}
