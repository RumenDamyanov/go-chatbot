package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChiAdapter(t *testing.T) {
	bot := setupTestBot()
	adapter := NewChiAdapter(bot)

	assert.NotNil(t, adapter)
	assert.Equal(t, bot, adapter.chatbot)
	assert.Equal(t, 30*time.Second, adapter.timeout)
}

func TestChiAdapter_WithTimeout(t *testing.T) {
	bot := setupTestBot()
	adapter := NewChiAdapter(bot).WithTimeout(10 * time.Second)

	assert.Equal(t, 10*time.Second, adapter.timeout)
}

func TestChiAdapter_ChatHandler(t *testing.T) {
	bot := setupTestBot()
	adapter := NewChiAdapter(bot)

	r := chi.NewRouter()
	r.Post("/chat", adapter.ChatHandler())

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectSuccess  bool
	}{
		{
			name: "valid chat request",
			requestBody: ChatRequest{
				Message: "Hello",
				Context: map[string]interface{}{
					"test": "value",
				},
			},
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
		},
		{
			name:           "missing message",
			requestBody:    ChatRequest{},
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req, err := http.NewRequest("POST", "/chat", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var response ChatResponse
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.expectSuccess, response.Success)

			if tt.expectSuccess {
				assert.NotEmpty(t, response.Response)
				assert.Empty(t, response.Error)
			} else {
				assert.NotEmpty(t, response.Error)
			}
		})
	}
}

func TestChiAdapter_HealthHandler(t *testing.T) {
	bot := setupTestBot()
	adapter := NewChiAdapter(bot)

	r := chi.NewRouter()
	r.Get("/health", adapter.HealthHandler())

	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response HealthResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response.Status)
	assert.Equal(t, "local", response.Provider)
	assert.Equal(t, "free-model", response.Model)
	assert.Greater(t, response.Timestamp, int64(0))
	assert.Empty(t, response.Error)
}

func TestChiAdapter_StreamChatHandler(t *testing.T) {
	bot := setupTestBot()
	adapter := NewChiAdapter(bot)

	r := chi.NewRouter()
	r.Post("/stream", adapter.StreamChatHandler())

	req, err := http.NewRequest("POST", "/stream", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotImplemented, rr.Code)
}

func TestChiAdapter_SetupRoutes(t *testing.T) {
	bot := setupTestBot()
	adapter := NewChiAdapter(bot)

	r := chi.NewRouter()
	adapter.SetupRoutes(r)

	// Test that routes are properly set up
	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	// Test POST /chat/
	req, err := http.NewRequest("POST", "/chat/", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test GET /chat/health
	req, err = http.NewRequest("GET", "/chat/health", nil)
	require.NoError(t, err)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test POST /chat/stream
	req, err = http.NewRequest("POST", "/chat/stream", nil)
	require.NoError(t, err)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotImplemented, rr.Code)
}

func TestChiAdapter_SetupRoutesWithPrefix(t *testing.T) {
	bot := setupTestBot()
	adapter := NewChiAdapter(bot)

	r := chi.NewRouter()
	adapter.SetupRoutesWithPrefix(r, "/api/v1/chatbot")

	// Test that routes are properly set up with prefix
	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	req, err := http.NewRequest("POST", "/api/v1/chatbot/", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestChiAdapter_Middleware(t *testing.T) {
	bot := setupTestBot()
	adapter := NewChiAdapter(bot)

	r := chi.NewRouter()
	r.Use(adapter.Middleware())
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		retrievedBot, exists := GetChatbotFromChiContext(r)
		assert.True(t, exists)
		assert.Equal(t, bot, retrievedBot)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetChatbotFromChiContext(t *testing.T) {
	bot := setupTestBot()

	// Test with chatbot in context
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "chatbot", bot)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		retrievedBot, exists := GetChatbotFromChiContext(r)
		assert.True(t, exists)
		assert.Equal(t, bot, retrievedBot)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"test": "with_chatbot"})
	})

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test without chatbot in context
	r2 := chi.NewRouter()
	r2.Get("/", func(w http.ResponseWriter, r *http.Request) {
		retrievedBot, exists := GetChatbotFromChiContext(r)
		assert.False(t, exists)
		assert.Nil(t, retrievedBot)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"test": "no_chatbot"})
	})

	req2, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	rr2 := httptest.NewRecorder()
	r2.ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusOK, rr2.Code)

	// Test with wrong type in context
	r3 := chi.NewRouter()
	r3.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "chatbot", "not a chatbot")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r3.Get("/", func(w http.ResponseWriter, r *http.Request) {
		retrievedBot, exists := GetChatbotFromChiContext(r)
		assert.False(t, exists)
		assert.Nil(t, retrievedBot)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"test": "wrong_type"})
	})

	req3, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	rr3 := httptest.NewRecorder()
	r3.ServeHTTP(rr3, req3)
	assert.Equal(t, http.StatusOK, rr3.Code)
}

func TestChiAdapter_ChatHandler_ContextTimeout(t *testing.T) {
	bot := setupTestBot()
	adapter := NewChiAdapter(bot).WithTimeout(1 * time.Millisecond) // Very short timeout

	r := chi.NewRouter()
	r.Post("/chat", adapter.ChatHandler())

	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	req, err := http.NewRequest("POST", "/chat", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Add a context that might timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// The status could be 200 (if fast enough) or 408 (if timeout)
	// We just verify it doesn't crash
	assert.True(t, rr.Code == http.StatusOK || rr.Code == http.StatusRequestTimeout)
}
