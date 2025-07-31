// Package gochatbot provides a framework-agnostic Go package for integrating
// AI-powered chat functionality into web applications.
//
// It supports multiple AI providers (OpenAI, Anthropic, Google Gemini, etc.)
// and includes built-in security features, message filtering, and framework adapters.
package gochatbot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/RumenDamyanov/go-chatbot/config"
	"github.com/RumenDamyanov/go-chatbot/middleware"
	"github.com/RumenDamyanov/go-chatbot/models"
	"github.com/RumenDamyanov/go-chatbot/streaming"
)

// Chatbot represents the main chatbot instance.
type Chatbot struct {
	config    *config.Config
	model     models.Model
	filter    *middleware.ChatMessageFilter
	rateLimit *middleware.RateLimiter
	timeout   time.Duration
}

// Option represents a configuration option for the Chatbot.
type Option func(*Chatbot)

// WithModel sets a custom AI model for the chatbot.
func WithModel(model models.Model) Option {
	return func(c *Chatbot) {
		c.model = model
	}
}

// WithTimeout sets a custom timeout for AI requests.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Chatbot) {
		c.timeout = timeout
	}
}

// WithFilter sets a custom message filter.
func WithFilter(filter *middleware.ChatMessageFilter) Option {
	return func(c *Chatbot) {
		c.filter = filter
	}
}

// WithRateLimit sets a custom rate limiter.
func WithRateLimit(limiter *middleware.RateLimiter) Option {
	return func(c *Chatbot) {
		c.rateLimit = limiter
	}
}

// New creates a new Chatbot instance with the given configuration and options.
func New(cfg *config.Config, opts ...Option) (*Chatbot, error) {
	if cfg == nil {
		return nil, errors.New("config cannot be nil")
	}

	// Create default model if none provided
	var model models.Model
	var err error

	chatbot := &Chatbot{
		config:  cfg,
		timeout: cfg.Timeout,
	}

	// Apply options
	for _, opt := range opts {
		opt(chatbot)
	}

	// Create model if not provided via options
	if chatbot.model == nil {
		model, err = models.NewFromConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create model: %w", err)
		}
		chatbot.model = model
	}

	// Create message filter
	if chatbot.filter == nil {
		chatbot.filter = middleware.NewChatMessageFilter(cfg.MessageFiltering)
	}

	// Create rate limiter
	if chatbot.rateLimit == nil {
		chatbot.rateLimit = middleware.NewRateLimiter(cfg.RateLimit)
	}

	return chatbot, nil
}

// Ask sends a message to the AI model and returns the response.
// It applies message filtering and rate limiting before processing.
func (c *Chatbot) Ask(ctx context.Context, message string, options ...AskOption) (string, error) {
	if message == "" {
		return "", errors.New("message cannot be empty")
	}

	// Create context with timeout
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	// Apply rate limiting
	if c.rateLimit != nil {
		if err := c.rateLimit.Allow(ctx); err != nil {
			return "", fmt.Errorf("rate limit exceeded: %w", err)
		}
	}

	// Apply message filtering
	filtered, err := c.filter.Handle(ctx, message)
	if err != nil {
		return "", fmt.Errorf("message filtering failed: %w", err)
	}

	// Parse options
	askOpts := &askOptions{
		context: filtered.Context,
	}
	for _, opt := range options {
		opt(askOpts)
	}

	// Send to AI model
	response, err := c.model.Ask(ctx, filtered.Message, askOpts.context)
	if err != nil {
		return "", fmt.Errorf("AI model request failed: %w", err)
	}

	return response, nil
}

// AskOption represents an option for the Ask method.
type AskOption func(*askOptions)

type askOptions struct {
	context map[string]interface{}
}

// WithContext adds additional context to the AI request.
func WithContext(key string, value interface{}) AskOption {
	return func(opts *askOptions) {
		if opts.context == nil {
			opts.context = make(map[string]interface{})
		}
		opts.context[key] = value
	}
}

// GetConfig returns the chatbot's configuration.
func (c *Chatbot) GetConfig() *config.Config {
	return c.config
}

// GetModel returns the chatbot's AI model.
func (c *Chatbot) GetModel() models.Model {
	return c.model
}

// Health checks if the chatbot and its dependencies are healthy.
func (c *Chatbot) Health(ctx context.Context) error {
	// Check if model is available
	if c.model == nil {
		return errors.New("AI model is not initialized")
	}

	// Check model health if supported
	if healthChecker, ok := c.model.(models.HealthChecker); ok {
		if err := healthChecker.Health(ctx); err != nil {
			return fmt.Errorf("AI model health check failed: %w", err)
		}
	}

	return nil
}

// AskStream sends a message to the AI model and returns a streaming response.
// It applies message filtering and rate limiting before processing.
func (c *Chatbot) AskStream(ctx context.Context, w http.ResponseWriter, message string, options ...AskOption) error {
	if message == "" {
		return errors.New("message cannot be empty")
	}

	// Create streaming handler
	streamHandler, err := streaming.NewStreamHandler(w)
	if err != nil {
		return fmt.Errorf("failed to create stream handler: %w", err)
	}
	defer streamHandler.Close()

	// Create context with timeout
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	// Apply rate limiting
	if c.rateLimit != nil {
		if err := c.rateLimit.Allow(ctx); err != nil {
			return streamHandler.WriteError("", fmt.Sprintf("rate limit exceeded: %v", err))
		}
	}

	// Apply message filtering
	filtered, err := c.filter.Handle(ctx, message)
	if err != nil {
		return streamHandler.WriteError("", fmt.Sprintf("message filtering failed: %v", err))
	}

	// Parse options
	askOpts := &askOptions{
		context: filtered.Context,
	}
	for _, opt := range options {
		opt(askOpts)
	}

	// Check if model supports streaming
	streamingModel, isStreaming := c.model.(models.StreamingModel)
	if !isStreaming {
		// Fallback to regular Ask and send as single chunk
		response, err := c.model.Ask(ctx, filtered.Message, askOpts.context)
		if err != nil {
			return streamHandler.WriteError("", fmt.Sprintf("AI model request failed: %v", err))
		}

		// Send as single chunk
		err = streamHandler.WriteChunk(streaming.StreamResponse{
			ID:      "single-chunk",
			Content: response,
			Done:    false,
		})
		if err != nil {
			return err
		}

		return streamHandler.WriteDone("single-chunk")
	}

	// Get streaming response
	responseCh, err := streamingModel.AskStream(ctx, filtered.Message, askOpts.context)
	if err != nil {
		return streamHandler.WriteError("", fmt.Sprintf("streaming request failed: %v", err))
	}

	// Process streaming response
	processor := streaming.NewStreamProcessor("stream", streamHandler)
	return processor.ProcessChannel(ctx, responseCh)
}
