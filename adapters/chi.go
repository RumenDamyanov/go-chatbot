package adapters

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	gochatbot "go.rumenx.com/chatbot"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	chatbotContextKey contextKey = "chatbot"
)

// ChiAdapter wraps a chatbot for use with the Chi framework
type ChiAdapter struct {
	chatbot *gochatbot.Chatbot
	timeout time.Duration
}

// NewChiAdapter creates a new Chi adapter for the chatbot
func NewChiAdapter(chatbot *gochatbot.Chatbot) *ChiAdapter {
	return &ChiAdapter{
		chatbot: chatbot,
		timeout: 30 * time.Second,
	}
}

// WithTimeout sets the timeout for chat operations
func (adapter *ChiAdapter) WithTimeout(timeout time.Duration) *ChiAdapter {
	adapter.timeout = timeout
	return adapter
}

// ChatHandler returns a Chi handler for chat requests
func (adapter *ChiAdapter) ChatHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), adapter.timeout)
		defer cancel()

		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response := ChatResponse{
				Success: false,
				Error:   "Invalid JSON",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if req.Message == "" {
			response := ChatResponse{
				Success: false,
				Error:   "Message is required",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		chatResponse, err := adapter.chatbot.Ask(ctx, req.Message)
		if err != nil {
			// Check if it's a timeout error
			if ctx.Err() == context.DeadlineExceeded {
				response := ChatResponse{
					Success: false,
					Error:   "Request timeout",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusRequestTimeout)
				json.NewEncoder(w).Encode(response)
				return
			}

			response := ChatResponse{
				Success: false,
				Error:   err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := ChatResponse{
			Success:  true,
			Response: chatResponse,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// HealthHandler returns a Chi handler for health checks
func (adapter *ChiAdapter) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		model := adapter.chatbot.GetModel()

		response := HealthResponse{
			Status:    "healthy",
			Provider:  model.Provider(),
			Model:     model.Name(),
			Timestamp: time.Now().Unix(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// StreamChatHandler returns a Chi handler for streaming chat (placeholder)
func (adapter *ChiAdapter) StreamChatHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Streaming not implemented", http.StatusNotImplemented)
	}
}

// SetupRoutes sets up the default routes on a Chi router
func (adapter *ChiAdapter) SetupRoutes(r chi.Router) {
	r.Route("/chat", func(r chi.Router) {
		r.Post("/", adapter.ChatHandler())
		r.Get("/health", adapter.HealthHandler())
		r.Post("/stream", adapter.StreamChatHandler())
	})
}

// SetupRoutesWithPrefix sets up routes with a custom prefix
func (adapter *ChiAdapter) SetupRoutesWithPrefix(r chi.Router, prefix string) {
	r.Route(prefix, func(r chi.Router) {
		r.Post("/", adapter.ChatHandler())
		r.Get("/health", adapter.HealthHandler())
		r.Post("/stream", adapter.StreamChatHandler())
	})
}

// Middleware returns Chi middleware that adds the chatbot to the request context
func (adapter *ChiAdapter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), chatbotContextKey, adapter.chatbot)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetChatbotFromChiContext retrieves the chatbot from a Chi request context
func GetChatbotFromChiContext(r *http.Request) (*gochatbot.Chatbot, bool) {
	chatbot, ok := r.Context().Value(chatbotContextKey).(*gochatbot.Chatbot)
	return chatbot, ok
}
