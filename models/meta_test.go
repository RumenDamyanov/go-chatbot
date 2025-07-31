package models

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/RumenDamyanov/go-chatbot/config"
)

func TestNewMetaModel_Detailed(t *testing.T) {
	cfg := config.MetaConfig{
		APIKey: "test-key",
		Model:  "llama-2-70b-chat",
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if model == nil {
		t.Fatal("Expected model to be created")
	}

	if model.Name() != "llama-2-70b-chat" {
		t.Errorf("Expected model name 'llama-2-70b-chat', got '%s'", model.Name())
	}

	if model.Provider() != "meta" {
		t.Errorf("Expected provider 'meta', got '%s'", model.Provider())
	}
}

func TestNewMetaModel_InvalidConfig(t *testing.T) {
	cfg := config.MetaConfig{
		APIKey: "", // Empty API key
		Model:  "llama-2-70b-chat",
	}

	_, err := NewMetaModel(cfg)
	if err == nil {
		t.Fatal("Expected error for empty API key")
	}

	if !strings.Contains(err.Error(), "API key is required") {
		t.Errorf("Expected 'API key is required' error, got: %v", err)
	}
}

func TestMetaModel_Ask_Success_Detailed(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json")
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got '%s'", authHeader)
		}

		// Mock successful response
		response := metaResponse{
			Choices: []metaChoice{
				{
					Message: metaMessage{
						Role:    "assistant",
						Content: "Hello! How can I help you?",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := config.MetaConfig{
		APIKey:   "test-key",
		Model:    "llama-2-70b-chat",
		Endpoint: server.URL,
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx := context.Background()
	response, err := model.Ask(ctx, "Hello", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "Hello! How can I help you?"
	if response != expected {
		t.Errorf("Expected response '%s', got '%s'", expected, response)
	}
}

func TestMetaModel_Ask_APIError_Detailed(t *testing.T) {
	// Create mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	cfg := config.MetaConfig{
		APIKey:   "test-key",
		Model:    "llama-2-70b-chat",
		Endpoint: server.URL,
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx := context.Background()
	_, err = model.Ask(ctx, "Hello", nil)
	if err == nil {
		t.Fatal("Expected error for API failure")
	}

	if !strings.Contains(err.Error(), "API error") {
		t.Errorf("Expected API error, got: %v", err)
	}
}

func TestMetaModel_Ask_InvalidJSON(t *testing.T) {
	// Create mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	cfg := config.MetaConfig{
		APIKey:   "test-key",
		Model:    "llama-2-70b-chat",
		Endpoint: server.URL,
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx := context.Background()
	_, err = model.Ask(ctx, "Hello", nil)
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "unmarshal response") {
		t.Errorf("Expected unmarshal error, got: %v", err)
	}
}

func TestMetaModel_Ask_EmptyResponse(t *testing.T) {
	// Create mock server that returns empty choices
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := metaResponse{
			Choices: []metaChoice{}, // Empty choices
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := config.MetaConfig{
		APIKey:   "test-key",
		Model:    "llama-2-70b-chat",
		Endpoint: server.URL,
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx := context.Background()
	_, err = model.Ask(ctx, "Hello", nil)
	if err == nil {
		t.Fatal("Expected error for empty response")
	}

	if !strings.Contains(err.Error(), "no choices") {
		t.Errorf("Expected 'no choices' error, got: %v", err)
	}
}

func TestMetaModel_Ask_ContextCancellation(t *testing.T) {
	// Create mock server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		response := metaResponse{
			Choices: []metaChoice{
				{
					Message: metaMessage{
						Role:    "assistant",
						Content: "Response",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := config.MetaConfig{
		APIKey:   "test-key",
		Model:    "llama-2-70b-chat",
		Endpoint: server.URL,
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = model.Ask(ctx, "Hello", nil)
	if err == nil {
		t.Fatal("Expected error for context cancellation")
	}

	if !strings.Contains(err.Error(), "context") {
		t.Errorf("Expected context cancellation error, got: %v", err)
	}
}

func TestMetaModel_Health_Success(t *testing.T) {
	// Create mock server for health check
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got '%s'", authHeader)
		}

		// Mock successful health check response
		response := metaResponse{
			Choices: []metaChoice{
				{
					Message: metaMessage{
						Role:    "assistant",
						Content: "OK",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := config.MetaConfig{
		APIKey:   "test-key",
		Model:    "llama-2-70b-chat",
		Endpoint: server.URL,
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx := context.Background()
	err = model.Health(ctx)
	if err != nil {
		t.Errorf("Expected no error for health check, got %v", err)
	}
}

func TestMetaModel_Health_Unauthorized(t *testing.T) {
	// Create mock server that returns unauthorized
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}))
	defer server.Close()

	cfg := config.MetaConfig{
		APIKey:   "invalid-key",
		Model:    "llama-2-70b-chat",
		Endpoint: server.URL,
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx := context.Background()
	err = model.Health(ctx)
	if err == nil {
		t.Fatal("Expected error for unauthorized access")
	}

	if !strings.Contains(err.Error(), "invalid API key") {
		t.Errorf("Expected 'invalid API key' error, got: %v", err)
	}
}

func TestMetaModel_Health_ServerError(t *testing.T) {
	// Create mock server that returns server error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	cfg := config.MetaConfig{
		APIKey:   "test-key",
		Model:    "llama-2-70b-chat",
		Endpoint: server.URL,
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx := context.Background()
	err = model.Health(ctx)
	if err == nil {
		t.Fatal("Expected error for server error")
	}

	if !strings.Contains(err.Error(), "meta API server error") {
		t.Errorf("Expected 'meta API server error' error, got: %v", err)
	}
}

func TestMetaModel_Health_NetworkError(t *testing.T) {
	cfg := config.MetaConfig{
		APIKey:   "test-key",
		Model:    "llama-2-70b-chat",
		Endpoint: "http://invalid-endpoint.local",
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx := context.Background()
	err = model.Health(ctx)
	if err == nil {
		t.Fatal("Expected error for network failure")
	}

	if !strings.Contains(err.Error(), "health check failed") {
		t.Errorf("Expected 'health check failed' error, got: %v", err)
	}
}

func TestMetaModel_Health_ContextCancellation(t *testing.T) {
	// Create mock server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.MetaConfig{
		APIKey:   "test-key",
		Model:    "llama-2-70b-chat",
		Endpoint: server.URL,
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err = model.Health(ctx)
	if err == nil {
		t.Fatal("Expected error for context cancellation")
	}

	if !strings.Contains(err.Error(), "health check failed") {
		t.Errorf("Expected 'health check failed' error, got: %v", err)
	}
}

func TestMetaModel_CustomEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := metaResponse{
			Choices: []metaChoice{
				{
					Message: metaMessage{
						Role:    "assistant",
						Content: "Custom endpoint response",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := config.MetaConfig{
		APIKey:   "test-key",
		Model:    "llama-2-70b-chat",
		Endpoint: server.URL,
	}

	model, err := NewMetaModel(cfg)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	ctx := context.Background()
	response, err := model.Ask(ctx, "Hello", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "Custom endpoint response"
	if response != expected {
		t.Errorf("Expected response '%s', got '%s'", expected, response)
	}
}
