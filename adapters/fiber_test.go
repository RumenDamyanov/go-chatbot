package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFiberAdapter(t *testing.T) {
	bot := setupTestBot()
	adapter := NewFiberAdapter(bot)

	assert.NotNil(t, adapter)
	assert.Equal(t, bot, adapter.chatbot)
	assert.Equal(t, 30*time.Second, adapter.timeout)
}

func TestFiberAdapter_WithTimeout(t *testing.T) {
	bot := setupTestBot()
	adapter := NewFiberAdapter(bot).WithTimeout(10 * time.Second)

	assert.Equal(t, 10*time.Second, adapter.timeout)
}

func TestFiberAdapter_ChatHandler(t *testing.T) {
	bot := setupTestBot()
	adapter := NewFiberAdapter(bot)

	app := fiber.New()
	app.Post("/chat", adapter.ChatHandler())

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

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			responseBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var response ChatResponse
			err = json.Unmarshal(responseBody, &response)
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

func TestFiberAdapter_HealthHandler(t *testing.T) {
	bot := setupTestBot()
	adapter := NewFiberAdapter(bot)

	app := fiber.New()
	app.Get("/health", adapter.HealthHandler())

	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response HealthResponse
	err = json.Unmarshal(responseBody, &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response.Status)
	assert.Equal(t, "local", response.Provider)
	assert.Equal(t, "free-model", response.Model)
	assert.Greater(t, response.Timestamp, int64(0))
	assert.Empty(t, response.Error)
}

func TestFiberAdapter_StreamChatHandler(t *testing.T) {
	bot := setupTestBot()
	adapter := NewFiberAdapter(bot)

	app := fiber.New()
	app.Post("/stream", adapter.StreamChatHandler())

	req, err := http.NewRequest("POST", "/stream", nil)
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
}

func TestFiberAdapter_SetupRoutes(t *testing.T) {
	bot := setupTestBot()
	adapter := NewFiberAdapter(bot)

	app := fiber.New()
	adapter.SetupRoutes(app)

	// Test that routes are properly set up
	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	// Test POST /chat/
	req, err := http.NewRequest("POST", "/chat/", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test GET /chat/health
	req, err = http.NewRequest("GET", "/chat/health", nil)
	require.NoError(t, err)
	resp, err = app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test POST /chat/stream
	req, err = http.NewRequest("POST", "/chat/stream", nil)
	require.NoError(t, err)
	resp, err = app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
}

func TestFiberAdapter_SetupRoutesWithPrefix(t *testing.T) {
	bot := setupTestBot()
	adapter := NewFiberAdapter(bot)

	app := fiber.New()
	adapter.SetupRoutesWithPrefix(app, "/api/v1/chatbot")

	// Test that routes are properly set up with prefix
	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	req, err := http.NewRequest("POST", "/api/v1/chatbot/", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberAdapter_Middleware(t *testing.T) {
	bot := setupTestBot()
	adapter := NewFiberAdapter(bot)

	app := fiber.New()
	app.Use(adapter.Middleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		retrievedBot, exists := GetChatbotFromFiberContext(c)
		assert.True(t, exists)
		assert.Equal(t, bot, retrievedBot)
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetChatbotFromFiberContext(t *testing.T) {
	bot := setupTestBot()

	app := fiber.New()

	// Test with chatbot in context via middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("chatbot", bot)
		retrievedBot, exists := GetChatbotFromFiberContext(c)
		assert.True(t, exists)
		assert.Equal(t, bot, retrievedBot)
		return c.JSON(fiber.Map{"test": "with_chatbot"})
	})

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test without chatbot in context
	app2 := fiber.New()
	app2.Get("/", func(c *fiber.Ctx) error {
		retrievedBot, exists := GetChatbotFromFiberContext(c)
		assert.False(t, exists)
		assert.Nil(t, retrievedBot)
		return c.JSON(fiber.Map{"test": "no_chatbot"})
	})

	req2, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	resp2, err := app2.Test(req2)
	require.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	// Test with wrong type in context
	app3 := fiber.New()
	app3.Use(func(c *fiber.Ctx) error {
		c.Locals("chatbot", "not a chatbot")
		retrievedBot, exists := GetChatbotFromFiberContext(c)
		assert.False(t, exists)
		assert.Nil(t, retrievedBot)
		return c.JSON(fiber.Map{"test": "wrong_type"})
	})

	req3, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	resp3, err := app3.Test(req3)
	require.NoError(t, err)
	defer resp3.Body.Close()
	assert.Equal(t, http.StatusOK, resp3.StatusCode)
}

func TestFiberAdapter_ChatHandler_ContextTimeout(t *testing.T) {
	bot := setupTestBot()
	adapter := NewFiberAdapter(bot).WithTimeout(1 * time.Millisecond) // Very short timeout

	app := fiber.New()
	app.Post("/chat", adapter.ChatHandler())

	chatReq := ChatRequest{Message: "Hello"}
	body, _ := json.Marshal(chatReq)

	req, err := http.NewRequest("POST", "/chat", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Add a context that might timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// The status could be 200 (if fast enough) or 408 (if timeout)
	// We just verify it doesn't crash
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusRequestTimeout)
}
