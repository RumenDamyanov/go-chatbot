package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEchoAdapter(t *testing.T) {
	bot := setupTestBot()
	adapter := NewEchoAdapter(bot)

	assert.NotNil(t, adapter)
	assert.Equal(t, bot, adapter.chatbot)
	assert.Equal(t, 30*time.Second, adapter.timeout)
}

func TestEchoAdapter_WithTimeout(t *testing.T) {
	bot := setupTestBot()
	adapter := NewEchoAdapter(bot).WithTimeout(10 * time.Second)

	assert.Equal(t, 10*time.Second, adapter.timeout)
}

func TestEchoAdapter_ChatHandler(t *testing.T) {
	bot := setupTestBot()
	adapter := NewEchoAdapter(bot)

	e := echo.New()
	e.POST("/chat", adapter.ChatHandler())

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

			req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			e.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response ChatResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
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

func TestEchoAdapter_HealthHandler(t *testing.T) {
	bot := setupTestBot()
	adapter := NewEchoAdapter(bot)

	e := echo.New()
	e.GET("/health", adapter.HealthHandler())

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response.Status)
	assert.Equal(t, "local", response.Provider)
	assert.Equal(t, "free-model", response.Model)
	assert.Greater(t, response.Timestamp, int64(0))
	assert.Empty(t, response.Error)
}

func TestEchoAdapter_StreamChatHandler(t *testing.T) {
	bot := setupTestBot()
	adapter := NewEchoAdapter(bot)

	e := echo.New()
	e.POST("/stream", adapter.StreamChatHandler())

	req := httptest.NewRequest("POST", "/stream", nil)
	w := httptest.NewRecorder()

	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

func TestEchoAdapter_SetupRoutes(t *testing.T) {
	bot := setupTestBot()
	adapter := NewEchoAdapter(bot)

	e := echo.New()
	adapter.SetupRoutes(e)

	// Test that routes are properly set up
	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	// Test POST /chat/
	req := httptest.NewRequest("POST", "/chat/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test GET /chat/health
	req = httptest.NewRequest("GET", "/chat/health", nil)
	w = httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test POST /chat/stream
	req = httptest.NewRequest("POST", "/chat/stream", nil)
	w = httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

func TestEchoAdapter_SetupRoutesWithPrefix(t *testing.T) {
	bot := setupTestBot()
	adapter := NewEchoAdapter(bot)

	e := echo.New()
	adapter.SetupRoutesWithPrefix(e, "/api/v1/chatbot")

	// Test that routes are properly set up with prefix
	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	req := httptest.NewRequest("POST", "/api/v1/chatbot/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEchoAdapter_Middleware(t *testing.T) {
	bot := setupTestBot()
	adapter := NewEchoAdapter(bot)

	e := echo.New()
	e.Use(adapter.Middleware())
	e.GET("/test", func(c echo.Context) error {
		retrievedBot, exists := GetChatbotFromEchoContext(c)
		assert.True(t, exists)
		assert.Equal(t, bot, retrievedBot)
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetChatbotFromEchoContext(t *testing.T) {
	bot := setupTestBot()

	// Test with chatbot in context
	e := echo.New()
	c := e.NewContext(httptest.NewRequest("GET", "/", nil), httptest.NewRecorder())
	c.Set("chatbot", bot)

	retrievedBot, exists := GetChatbotFromEchoContext(c)
	assert.True(t, exists)
	assert.Equal(t, bot, retrievedBot)

	// Test without chatbot in context
	c2 := e.NewContext(httptest.NewRequest("GET", "/", nil), httptest.NewRecorder())

	retrievedBot2, exists2 := GetChatbotFromEchoContext(c2)
	assert.False(t, exists2)
	assert.Nil(t, retrievedBot2)

	// Test with wrong type in context
	c3 := e.NewContext(httptest.NewRequest("GET", "/", nil), httptest.NewRecorder())
	c3.Set("chatbot", "not a chatbot")

	retrievedBot3, exists3 := GetChatbotFromEchoContext(c3)
	assert.False(t, exists3)
	assert.Nil(t, retrievedBot3)
}

func TestEchoAdapter_ChatHandler_ContextTimeout(t *testing.T) {
	bot := setupTestBot()
	adapter := NewEchoAdapter(bot).WithTimeout(1 * time.Millisecond) // Very short timeout

	e := echo.New()
	e.POST("/chat", adapter.ChatHandler())

	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Add a context that might timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	// The status could be 200 (if fast enough) or 408 (if timeout)
	// We just verify it doesn't crash
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusRequestTimeout)
}
