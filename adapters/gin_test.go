package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gochatbot "go.rumenx.com/chatbot"
	"go.rumenx.com/chatbot/config"
)

func setupTestBot() *gochatbot.Chatbot {
	cfg := &config.Config{
		Model:   "free",
		Timeout: 5 * time.Second,
		RateLimit: config.RateLimitConfig{
			RequestsPerMinute: 60, // Allow 60 requests per minute
			BurstSize:         10, // Allow bursts of 10 requests
		},
		MessageFiltering: config.MessageFilteringConfig{
			Enabled: false, // Disable filtering for tests
		},
	}

	bot, _ := gochatbot.New(cfg)
	return bot
}

func TestNewGinAdapter(t *testing.T) {
	bot := setupTestBot()
	adapter := NewGinAdapter(bot)

	assert.NotNil(t, adapter)
	assert.Equal(t, bot, adapter.chatbot)
	assert.Equal(t, 30*time.Second, adapter.timeout)
}

func TestGinAdapter_WithTimeout(t *testing.T) {
	bot := setupTestBot()
	adapter := NewGinAdapter(bot).WithTimeout(10 * time.Second)

	assert.Equal(t, 10*time.Second, adapter.timeout)
}

func TestGinAdapter_ChatHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bot := setupTestBot()
	adapter := NewGinAdapter(bot)

	router := gin.New()
	router.POST("/chat", adapter.ChatHandler())

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
			router.ServeHTTP(w, req)

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

func TestGinAdapter_HealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bot := setupTestBot()
	adapter := NewGinAdapter(bot)

	router := gin.New()
	router.GET("/health", adapter.HealthHandler())

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

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

func TestGinAdapter_StreamChatHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bot := setupTestBot()
	adapter := NewGinAdapter(bot)

	router := gin.New()
	router.POST("/stream", adapter.StreamChatHandler())

	req := httptest.NewRequest("POST", "/stream", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

func TestGinAdapter_SetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bot := setupTestBot()
	adapter := NewGinAdapter(bot)

	router := gin.New()
	adapter.SetupRoutes(router)

	// Test that routes are properly set up
	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	// Test POST /chat/
	req := httptest.NewRequest("POST", "/chat/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test GET /chat/health
	req = httptest.NewRequest("GET", "/chat/health", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test POST /chat/stream
	req = httptest.NewRequest("POST", "/chat/stream", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

func TestGinAdapter_SetupRoutesWithPrefix(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bot := setupTestBot()
	adapter := NewGinAdapter(bot)

	router := gin.New()
	adapter.SetupRoutesWithPrefix(router, "/api/v1/chatbot")

	// Test that routes are properly set up with prefix
	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	req := httptest.NewRequest("POST", "/api/v1/chatbot/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGinAdapter_Middleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bot := setupTestBot()
	adapter := NewGinAdapter(bot)

	router := gin.New()
	router.Use(adapter.Middleware())
	router.GET("/test", func(c *gin.Context) {
		retrievedBot, exists := GetChatbotFromContext(c)
		assert.True(t, exists)
		assert.Equal(t, bot, retrievedBot)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetChatbotFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bot := setupTestBot()

	// Test with chatbot in context
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("chatbot", bot)

	retrievedBot, exists := GetChatbotFromContext(c)
	assert.True(t, exists)
	assert.Equal(t, bot, retrievedBot)

	// Test without chatbot in context
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())

	retrievedBot2, exists2 := GetChatbotFromContext(c2)
	assert.False(t, exists2)
	assert.Nil(t, retrievedBot2)

	// Test with wrong type in context
	c3, _ := gin.CreateTestContext(httptest.NewRecorder())
	c3.Set("chatbot", "not a chatbot")

	retrievedBot3, exists3 := GetChatbotFromContext(c3)
	assert.False(t, exists3)
	assert.Nil(t, retrievedBot3)
}

func TestGinAdapter_ChatHandler_ContextTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bot := setupTestBot()
	adapter := NewGinAdapter(bot).WithTimeout(1 * time.Millisecond) // Very short timeout

	router := gin.New()
	router.POST("/chat", adapter.ChatHandler())

	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	req := httptest.NewRequest("POST", "/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Add a context that might timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// The status could be 200 (if fast enough) or 408 (if timeout)
	// We just verify it doesn't crash
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusRequestTimeout)
}
