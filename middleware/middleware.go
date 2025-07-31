// Package middleware provides message filtering and rate limiting functionality.
package middleware

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/RumenDamyanov/go-chatbot/config"
)

// ChatMessageFilter provides message filtering capabilities.
type ChatMessageFilter struct {
	config          config.MessageFilteringConfig
	profanityRegex  *regexp.Regexp
	aggressionRegex *regexp.Regexp
	linkRegex       *regexp.Regexp
	mutex           sync.RWMutex
}

// FilteredMessage represents a filtered message with additional context.
type FilteredMessage struct {
	Message string
	Context map[string]interface{}
}

// NewChatMessageFilter creates a new message filter.
func NewChatMessageFilter(cfg config.MessageFilteringConfig) *ChatMessageFilter {
	filter := &ChatMessageFilter{
		config: cfg,
	}

	// Compile regex patterns
	if len(cfg.Profanities) > 0 {
		pattern := strings.Join(cfg.Profanities, "|")
		filter.profanityRegex = regexp.MustCompile(`(?i)\b(` + pattern + `)\b`)
	}

	if len(cfg.AggressionPatterns) > 0 {
		pattern := strings.Join(cfg.AggressionPatterns, "|")
		filter.aggressionRegex = regexp.MustCompile(`(?i)\b(` + pattern + `)\b`)
	}

	if cfg.LinkPattern != "" {
		filter.linkRegex = regexp.MustCompile(cfg.LinkPattern)
	}

	return filter
}

// Handle processes and filters a message.
func (f *ChatMessageFilter) Handle(ctx context.Context, message string) (*FilteredMessage, error) {
	if !f.config.Enabled {
		return &FilteredMessage{
			Message: message,
			Context: make(map[string]interface{}),
		}, nil
	}

	f.mutex.RLock()
	defer f.mutex.RUnlock()

	filtered := message
	context := make(map[string]interface{})

	// Filter profanities
	if f.profanityRegex != nil {
		filtered = f.profanityRegex.ReplaceAllString(filtered, "***")
	}

	// Filter aggression patterns
	if f.aggressionRegex != nil {
		if f.aggressionRegex.MatchString(filtered) {
			context["aggression_detected"] = true
		}
	}

	// Filter links
	if f.linkRegex != nil {
		if f.linkRegex.MatchString(filtered) {
			filtered = f.linkRegex.ReplaceAllString(filtered, "[link removed]")
			context["links_filtered"] = true
		}
	}

	// Add system instructions to context
	if len(f.config.Instructions) > 0 {
		context["system_instructions"] = f.config.Instructions
	}

	return &FilteredMessage{
		Message: filtered,
		Context: context,
	}, nil
}

// UpdateConfig updates the filter configuration.
func (f *ChatMessageFilter) UpdateConfig(cfg config.MessageFilteringConfig) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.config = cfg

	// Recompile regex patterns
	if len(cfg.Profanities) > 0 {
		pattern := strings.Join(cfg.Profanities, "|")
		f.profanityRegex = regexp.MustCompile(`(?i)\b(` + pattern + `)\b`)
	} else {
		f.profanityRegex = nil
	}

	if len(cfg.AggressionPatterns) > 0 {
		pattern := strings.Join(cfg.AggressionPatterns, "|")
		f.aggressionRegex = regexp.MustCompile(`(?i)\b(` + pattern + `)\b`)
	} else {
		f.aggressionRegex = nil
	}

	if cfg.LinkPattern != "" {
		f.linkRegex = regexp.MustCompile(cfg.LinkPattern)
	} else {
		f.linkRegex = nil
	}
}

// RateLimiter provides rate limiting functionality.
type RateLimiter struct {
	config   config.RateLimitConfig
	requests map[string][]time.Time
	mutex    sync.RWMutex
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:   cfg,
		requests: make(map[string][]time.Time),
	}
}

// Allow checks if a request is allowed based on rate limiting rules.
func (r *RateLimiter) Allow(ctx context.Context) error {
	// Extract client identifier from context (IP, user ID, etc.)
	clientID := r.getClientID(ctx)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.config.Window)

	// Clean old requests
	if requests, exists := r.requests[clientID]; exists {
		validRequests := make([]time.Time, 0, len(requests))
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}
		r.requests[clientID] = validRequests
	}

	// Check if within limit
	requestCount := len(r.requests[clientID])
	if requestCount >= r.config.RequestsPerMinute {
		return fmt.Errorf("rate limit exceeded: %d requests in %v", requestCount, r.config.Window)
	}

	// Add current request
	r.requests[clientID] = append(r.requests[clientID], now)

	return nil
}

// getClientID extracts a client identifier from the context.
func (r *RateLimiter) getClientID(ctx context.Context) string {
	// Try to get IP address from context
	if ip, ok := ctx.Value("client_ip").(string); ok {
		return ip
	}

	// Try to get user ID from context
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}

	// Fallback to a default identifier
	return "default"
}

// Cleanup removes old request records to prevent memory leaks.
func (r *RateLimiter) Cleanup() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.config.Window)

	for clientID, requests := range r.requests {
		validRequests := make([]time.Time, 0, len(requests))
		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}

		if len(validRequests) == 0 {
			delete(r.requests, clientID)
		} else {
			r.requests[clientID] = validRequests
		}
	}
}

// StartCleanupRoutine starts a background routine to clean up old records.
func (r *RateLimiter) StartCleanupRoutine(ctx context.Context) {
	ticker := time.NewTicker(r.config.Window)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.Cleanup()
		}
	}
}
