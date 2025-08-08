package models

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.rumenx.com/chatbot/config"
)

// OpenAIModel implements the Model interface for OpenAI's API.
type OpenAIModel struct {
	config     config.OpenAIConfig
	httpClient *http.Client
}

// OpenAIRequest represents a request to the OpenAI API.
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// OpenAIResponse represents a response from the OpenAI API.
type OpenAIResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Choice represents a response choice.
type Choice struct {
	Message Message `json:"message"`
}

// APIError represents an API error.
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// NewOpenAIModel creates a new OpenAI model instance.
func NewOpenAIModel(cfg config.OpenAIConfig) (*OpenAIModel, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://api.openai.com/v1/chat/completions"
	}

	if cfg.Model == "" {
		cfg.Model = "gpt-4o"
	}

	return &OpenAIModel{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Ask sends a message to the OpenAI API and returns the response.
func (o *OpenAIModel) Ask(ctx context.Context, message string, context map[string]interface{}) (string, error) {
	// Prepare system prompt
	systemPrompt := "You are a helpful chatbot."
	if prompt, ok := context["prompt"].(string); ok && prompt != "" {
		systemPrompt = prompt
	}

	// Prepare request
	request := OpenAIRequest{
		Model: o.config.Model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: message},
		},
	}

	// Add optional parameters from context
	if temp, ok := context["temperature"].(float64); ok {
		request.Temperature = temp
	}
	if maxTokens, ok := context["max_tokens"].(int); ok {
		request.MaxTokens = maxTokens
	}

	// Marshal request
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", o.config.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.config.APIKey)

	// Send request
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if openaiResp.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", openaiResp.Error.Message)
	}

	// Check for choices
	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return openaiResp.Choices[0].Message.Content, nil
}

// Name returns the name of the model.
func (o *OpenAIModel) Name() string {
	return o.config.Model
}

// Provider returns the provider of the model.
func (o *OpenAIModel) Provider() string {
	return "openai"
}

// Health checks if the OpenAI API is accessible.
func (o *OpenAIModel) Health(ctx context.Context) error {
	// Simple health check by making a minimal request
	_, err := o.Ask(ctx, "Hi", map[string]interface{}{
		"max_tokens": 1,
	})
	return err
}

// AskStream sends a streaming request to OpenAI and returns a channel of responses.
func (o *OpenAIModel) AskStream(ctx context.Context, message string, context map[string]interface{}) (<-chan string, error) {
	// Prepare messages
	messages := []Message{
		{Role: "user", Content: message},
	}

	// Build request
	request := OpenAIRequest{
		Model:    o.config.Model,
		Messages: messages,
		Stream:   true,
	}

	// Apply context parameters
	if temp, ok := context["temperature"].(float64); ok {
		request.Temperature = temp
	}
	if maxTokens, ok := context["max_tokens"].(int); ok {
		request.MaxTokens = maxTokens
	}

	// Marshal request
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", o.config.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.config.APIKey)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	// Send request
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Check status
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Create response channel
	responseCh := make(chan string, 10)

	// Start goroutine to read streaming response
	go func() {
		defer close(responseCh)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			// Skip empty lines and comments
			if len(line) == 0 || strings.HasPrefix(line, ":") {
				continue
			}

			// Parse SSE format
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")

				// Check for end of stream
				if data == "[DONE]" {
					return
				}

				// Parse JSON data
				var chunk map[string]interface{}
				if err := json.Unmarshal([]byte(data), &chunk); err != nil {
					continue // Skip malformed chunks
				}

				// Extract content
				content := extractOpenAIStreamContent(chunk)
				if content != "" {
					select {
					case responseCh <- content:
					case <-ctx.Done():
						return
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			// Log error but don't panic
			select {
			case responseCh <- fmt.Sprintf("[ERROR: %v]", err):
			case <-ctx.Done():
			}
		}
	}()

	return responseCh, nil
}

// extractOpenAIStreamContent extracts content from OpenAI streaming format.
func extractOpenAIStreamContent(chunk map[string]interface{}) string {
	choices, ok := chunk["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return ""
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return ""
	}

	delta, ok := choice["delta"].(map[string]interface{})
	if !ok {
		return ""
	}

	content, ok := delta["content"].(string)
	if !ok {
		return ""
	}

	return content
}
