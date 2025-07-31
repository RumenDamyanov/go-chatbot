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

// OllamaModel implements the Model interface for Ollama local models.
type OllamaModel struct {
	config     config.OllamaConfig
	httpClient *http.Client
}

// NewOllamaModel creates a new Ollama model instance.
func NewOllamaModel(cfg config.OllamaConfig) (*OllamaModel, error) {
	if cfg.Model == "" {
		cfg.Model = "llama3.2" // Default model
	}

	return &OllamaModel{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // Longer timeout for local models
		},
	}, nil
}

// ollamaRequest represents the request structure for Ollama's API.
type ollamaRequest struct {
	Model    string                 `json:"model"`
	Prompt   string                 `json:"prompt,omitempty"`
	Messages []ollamaMessage        `json:"messages,omitempty"`
	Context  []int                  `json:"context,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
	Format   string                 `json:"format,omitempty"`
	Raw      bool                   `json:"raw,omitempty"`
	Stream   bool                   `json:"stream"`
}

// ollamaMessage represents a message in the conversation for chat API.
type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ollamaResponse represents the response from Ollama's API.
type ollamaResponse struct {
	Model              string         `json:"model"`
	CreatedAt          string         `json:"created_at"`
	Response           string         `json:"response,omitempty"`
	Message            *ollamaMessage `json:"message,omitempty"`
	Done               bool           `json:"done"`
	Context            []int          `json:"context,omitempty"`
	TotalDuration      int64          `json:"total_duration,omitempty"`
	LoadDuration       int64          `json:"load_duration,omitempty"`
	PromptEvalCount    int            `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64          `json:"prompt_eval_duration,omitempty"`
	EvalCount          int            `json:"eval_count,omitempty"`
	EvalDuration       int64          `json:"eval_duration,omitempty"`
}

// ollamaError represents an error response from the API.
type ollamaError struct {
	Error string `json:"error"`
}

// Ask sends a message to Ollama and returns the response.
func (o *OllamaModel) Ask(ctx context.Context, message string, context map[string]interface{}) (string, error) {
	// Determine endpoint
	endpoint := "http://localhost:11434"
	if o.config.Endpoint != "" {
		endpoint = o.config.Endpoint
	}

	// Use chat endpoint for conversation-style interactions
	useChatAPI := true
	if raw, ok := context["raw"]; ok {
		if rawMode, ok := raw.(bool); ok && rawMode {
			useChatAPI = false
		}
	}

	var url string
	var reqBody []byte
	var err error

	if useChatAPI {
		// Use chat API for conversation-style interactions
		url = fmt.Sprintf("%s/api/chat", endpoint)

		req := ollamaRequest{
			Model: o.config.Model,
			Messages: []ollamaMessage{
				{
					Role:    "user",
					Content: message,
				},
			},
			Stream: false,
		}

		// Add conversation history if provided
		if history, ok := context["history"]; ok {
			if hist, ok := history.([]map[string]interface{}); ok {
				var messages []ollamaMessage
				for _, msg := range hist {
					if role, roleOk := msg["role"].(string); roleOk {
						if content, contentOk := msg["content"].(string); contentOk {
							// Ollama uses "user", "assistant", "system" roles
							messages = append(messages, ollamaMessage{
								Role:    role,
								Content: content,
							})
						}
					}
				}
				// Add current message at the end
				messages = append(messages, ollamaMessage{
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
				req.Messages = append([]ollamaMessage{
					{
						Role:    "system",
						Content: sys,
					},
				}, req.Messages...)
			}
		}

		// Add options if provided
		if options := buildOllamaOptions(context); len(options) > 0 {
			req.Options = options
		}

		reqBody, err = json.Marshal(req)
	} else {
		// Use generate API for simple prompt completion
		url = fmt.Sprintf("%s/api/generate", endpoint)

		req := ollamaRequest{
			Model:  o.config.Model,
			Prompt: message,
			Stream: false,
		}

		// Add context from previous conversation
		if ctxData, ok := context["context"]; ok {
			if ctxArray, ok := ctxData.([]int); ok {
				req.Context = ctxArray
			}
		}

		// Add options if provided
		if options := buildOllamaOptions(context); len(options) > 0 {
			req.Options = options
		}

		reqBody, err = json.Marshal(req)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := o.httpClient.Do(httpReq)
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
		var errResp ollamaError
		if err := json.Unmarshal(body, &errResp); err == nil {
			return "", fmt.Errorf("Ollama API error: %s", errResp.Error)
		}
		return "", fmt.Errorf("Ollama API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract the text content based on API used
	if useChatAPI {
		if ollamaResp.Message == nil {
			return "", fmt.Errorf("no message in chat response")
		}
		if ollamaResp.Message.Content == "" {
			return "", fmt.Errorf("no content in response message")
		}
		return ollamaResp.Message.Content, nil
	} else {
		if ollamaResp.Response == "" {
			return "", fmt.Errorf("no response content")
		}
		return ollamaResp.Response, nil
	}
}

// buildOllamaOptions builds options map from context.
func buildOllamaOptions(context map[string]interface{}) map[string]interface{} {
	options := make(map[string]interface{})

	// Common parameters
	if temp, ok := context["temperature"]; ok {
		options["temperature"] = temp
	}
	if topP, ok := context["top_p"]; ok {
		options["top_p"] = topP
	}
	if topK, ok := context["top_k"]; ok {
		options["top_k"] = topK
	}
	if repeatPenalty, ok := context["repeat_penalty"]; ok {
		options["repeat_penalty"] = repeatPenalty
	}
	if seed, ok := context["seed"]; ok {
		options["seed"] = seed
	}
	if numCtx, ok := context["num_ctx"]; ok {
		options["num_ctx"] = numCtx
	}
	if numPredict, ok := context["num_predict"]; ok {
		options["num_predict"] = numPredict
	}
	if stop, ok := context["stop"]; ok {
		options["stop"] = stop
	}

	return options
}

// Name returns the name of the model.
func (o *OllamaModel) Name() string {
	return o.config.Model
}

// Provider returns the provider name.
func (o *OllamaModel) Provider() string {
	return "ollama"
}

// Health checks if the Ollama API is accessible.
func (o *OllamaModel) Health(ctx context.Context) error {
	endpoint := "http://localhost:11434"
	if o.config.Endpoint != "" {
		endpoint = o.config.Endpoint
	}

	// Check if Ollama is running by hitting the /api/tags endpoint
	url := fmt.Sprintf("%s/api/tags", endpoint)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := o.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("Ollama health check failed - is Ollama running?: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("Ollama server error: %d", resp.StatusCode)
	}

	// Check if the specific model is available
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			var tagsResp struct {
				Models []struct {
					Name string `json:"name"`
				} `json:"models"`
			}
			if json.Unmarshal(body, &tagsResp) == nil {
				// Check if our model is in the list
				for _, model := range tagsResp.Models {
					if model.Name == o.config.Model {
						return nil // Model is available
					}
				}
				return fmt.Errorf("model '%s' not found in Ollama. Available models can be seen with 'ollama list'", o.config.Model)
			}
		}
	}

	return nil
}
