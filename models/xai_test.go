package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RumenDamyanov/go-chatbot/config"
)

func TestXAIModel_Health(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectError    bool
		errorContains  string
	}{
		{
			name: "successful health check",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"choices": [{"message": {"content": "Hi!"}}]}`))
			},
			expectError: false,
		},
		{
			name: "unauthorized error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Invalid API key"}`))
			},
			expectError:   true,
			errorContains: "invalid API key",
		},
		{
			name: "server error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal server error"}`))
			},
			expectError:   true,
			errorContains: "server error: 500",
		},
		{
			name: "bad request but not auth error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "Bad request"}`))
			},
			expectError: false, // Health check passes for non-auth, non-server errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			// Create model with test server endpoint
			config := config.XAIConfig{
				APIKey:   "test-key",
				Model:    "grok-beta",
				Endpoint: server.URL,
			}
			model, err := NewXAIModel(config)
			if err != nil {
				t.Fatalf("failed to create model: %v", err)
			}

			// Test health check
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err = model.Health(ctx)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
					return
				}
				if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestXAIModel_Health_ContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Delay longer than context timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := config.XAIConfig{
		APIKey:   "test-key",
		Model:    "grok-beta",
		Endpoint: server.URL,
	}
	model, err := NewXAIModel(config)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = model.Health(ctx)
	if err == nil {
		t.Error("expected timeout error but got none")
	}
}

func TestXAIModel_Ask_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"choices": [
				{
					"message": {
						"content": "Hello! How can I help you today?"
					}
				}
			]
		}`))
	}))
	defer server.Close()

	config := config.XAIConfig{
		APIKey:   "test-api-key",
		Model:    "grok-beta",
		Endpoint: server.URL,
	}
	model, err := NewXAIModel(config)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx := context.Background()
	response, err := model.Ask(ctx, "Hello", nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if response != "Hello! How can I help you today?" {
		t.Errorf("expected specific response, got: %s", response)
	}
}

func TestXAIModel_Ask_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"message": "Invalid request parameters"
			}
		}`))
	}))
	defer server.Close()

	config := config.XAIConfig{
		APIKey:   "test-api-key",
		Model:    "grok-beta",
		Endpoint: server.URL,
	}
	model, err := NewXAIModel(config)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx := context.Background()
	_, err = model.Ask(ctx, "Hello", nil)

	if err == nil {
		t.Error("expected error for API error response")
	}
}
