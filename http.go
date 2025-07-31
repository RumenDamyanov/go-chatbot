package gochatbot

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	clientIPContextKey contextKey = "client_ip"
)

// Health status constants
const (
	healthStatusHealthy   = "healthy"
	healthStatusUnhealthy = "unhealthy"
)

// ChatRequest represents an incoming chat request.
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse represents a chat response.
type ChatResponse struct {
	Reply string `json:"reply"`
	Error string `json:"error,omitempty"`
}

// HTTPHandler provides HTTP handling functionality for the chatbot.
type HTTPHandler struct {
	chatbot *Chatbot
}

// NewHTTPHandler creates a new HTTP handler for the chatbot.
func NewHTTPHandler(chatbot *Chatbot) *HTTPHandler {
	return &HTTPHandler{
		chatbot: chatbot,
	}
}

// HandleHTTP handles HTTP requests for chat functionality.
func (h *HTTPHandler) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle OPTIONS requests for CORS
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON request")
		return
	}

	// Validate request
	if strings.TrimSpace(req.Message) == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Message cannot be empty")
		return
	}

	// Create context with client information
	ctx := context.WithValue(r.Context(), clientIPContextKey, h.getClientIP(r))

	// Add timeout if not already set
	if h.chatbot.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.chatbot.timeout)
		defer cancel()
	}

	// Process chat request
	reply, err := h.chatbot.Ask(ctx, req.Message)
	if err != nil {
		// Check for specific error types
		if ctx.Err() == context.DeadlineExceeded {
			h.writeErrorResponse(w, http.StatusRequestTimeout, "Request timeout")
			return
		}
		if strings.Contains(err.Error(), "rate limit") {
			h.writeErrorResponse(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	// Send response
	response := ChatResponse{
		Reply: reply,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Error encoding response, but headers already sent
		return
	}
}

// writeErrorResponse writes an error response to the client.
func (h *HTTPHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	response := ChatResponse{
		Error: message,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Error encoding response, but headers already sent
		return
	}
}

// getClientIP extracts the client IP address from the request.
func (h *HTTPHandler) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}

// Health handles health check requests.
func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.chatbot.Health(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		response := map[string]interface{}{
			"status": healthStatusUnhealthy,
			"error":  err.Error(),
		}
		if encErr := json.NewEncoder(w).Encode(response); encErr != nil {
			// Error encoding response, but headers already sent
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status": healthStatusHealthy,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Error encoding response, but headers already sent
		return
	}
}

// HandleHTTP is a convenience method to create and handle HTTP requests.
func (c *Chatbot) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	handler := NewHTTPHandler(c)
	handler.HandleHTTP(w, r)
}

// HandleStreamHTTP handles streaming HTTP requests for chat functionality.
func (h *HTTPHandler) HandleStreamHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle OPTIONS requests for CORS
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST requests
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate request
	if strings.TrimSpace(req.Message) == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Message cannot be empty")
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Add client IP to context
	clientIP := h.getClientIP(r)

	// Process streaming request
	err := h.chatbot.AskStream(ctx, w, req.Message, WithContext("client_ip", clientIP))
	if err != nil {
		// If we couldn't set up streaming, fall back to error response
		h.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleStreamHTTP is a convenience method to create and handle streaming HTTP requests.
func (c *Chatbot) HandleStreamHTTP(w http.ResponseWriter, r *http.Request) {
	handler := NewHTTPHandler(c)
	handler.HandleStreamHTTP(w, r)
}
