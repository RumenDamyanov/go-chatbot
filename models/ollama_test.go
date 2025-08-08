package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.rumenx.com/chatbot/config"
)

func TestOllamaModel_Ask_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": {
				"content": "Hello! How can I help you today?"
			}
		}`))
	}))
	defer server.Close()

	config := config.OllamaConfig{
		Model:    "llama2",
		Endpoint: server.URL,
	}
	model, err := NewOllamaModel(config)
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

func TestOllamaModel_Ask_WithHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": {
				"content": "I remember our previous conversation."
			}
		}`))
	}))
	defer server.Close()

	config := config.OllamaConfig{
		Model:    "llama2",
		Endpoint: server.URL,
	}
	model, err := NewOllamaModel(config)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx := context.Background()
	contextMap := map[string]interface{}{
		"history": []map[string]interface{}{
			{"role": "user", "content": "Previous message"},
			{"role": "assistant", "content": "Previous response"},
		},
	}

	response, err := model.Ask(ctx, "Current message", contextMap)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if response != "I remember our previous conversation." {
		t.Errorf("expected specific response, got: %s", response)
	}
}

func TestOllamaModel_Ask_WithSystemMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"message": {
				"content": "I'll follow the system instructions."
			}
		}`))
	}))
	defer server.Close()

	config := config.OllamaConfig{
		Model:    "llama2",
		Endpoint: server.URL,
	}
	model, err := NewOllamaModel(config)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx := context.Background()
	contextMap := map[string]interface{}{
		"system": "You are a helpful assistant.",
	}

	response, err := model.Ask(ctx, "Hello", contextMap)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if response != "I'll follow the system instructions." {
		t.Errorf("expected specific response, got: %s", response)
	}
}

func TestOllamaModel_Ask_RawMode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that it uses the generate endpoint in raw mode
		if r.URL.Path != "/api/generate" {
			t.Errorf("expected /api/generate path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"response": "Raw mode response"
		}`))
	}))
	defer server.Close()

	config := config.OllamaConfig{
		Model:    "llama2",
		Endpoint: server.URL,
	}
	model, err := NewOllamaModel(config)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx := context.Background()
	contextMap := map[string]interface{}{
		"raw": true,
	}

	response, err := model.Ask(ctx, "Hello", contextMap)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if response != "Raw mode response" {
		t.Errorf("expected raw mode response, got: %s", response)
	}
}

func TestOllamaModel_Ask_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": "Model not found"
		}`))
	}))
	defer server.Close()

	config := config.OllamaConfig{
		Model:    "nonexistent",
		Endpoint: server.URL,
	}
	model, err := NewOllamaModel(config)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx := context.Background()
	_, err = model.Ask(ctx, "Hello", nil)

	if err == nil {
		t.Error("expected error for API error response")
	}
}

func TestOllamaModel_Health_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"models": [
					{"name": "llama2"},
					{"name": "codellama"}
				]
			}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"message": {
					"content": "Hello!"
				}
			}`))
		}
	}))
	defer server.Close()

	config := config.OllamaConfig{
		Model:    "llama2",
		Endpoint: server.URL,
	}
	model, err := NewOllamaModel(config)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx := context.Background()
	err = model.Health(ctx)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestOllamaModel_Health_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Server error"}`))
	}))
	defer server.Close()

	config := config.OllamaConfig{
		Model:    "llama2",
		Endpoint: server.URL,
	}
	model, err := NewOllamaModel(config)
	if err != nil {
		t.Fatalf("failed to create model: %v", err)
	}

	ctx := context.Background()
	err = model.Health(ctx)

	if err == nil {
		t.Error("expected error for health check failure")
	}
}
