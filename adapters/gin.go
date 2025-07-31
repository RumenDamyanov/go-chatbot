// Package adapters provides framework-specific integrations for the go-chatbot package.
package adapters

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	gochatbot "github.com/RumenDamyanov/go-chatbot"
)

// GinAdapter provides Gin framework integration for go-chatbot.
type GinAdapter struct {
	chatbot *gochatbot.Chatbot
	timeout time.Duration
}

// NewGinAdapter creates a new Gin adapter with the provided chatbot instance.
func NewGinAdapter(bot *gochatbot.Chatbot) *GinAdapter {
	return &GinAdapter{
		chatbot: bot,
		timeout: 30 * time.Second,
	}
}

// WithTimeout sets the request timeout for the adapter.
func (a *GinAdapter) WithTimeout(timeout time.Duration) *GinAdapter {
	a.timeout = timeout
	return a
}

// ChatRequest represents the expected request format for chat endpoints.
type ChatRequest struct {
	Message string                 `json:"message" binding:"required"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// ChatResponse represents the response format for chat endpoints.
type ChatResponse struct {
	Response string `json:"response"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
}

// HealthResponse represents the response format for health check endpoints.
type HealthResponse struct {
	Status    string `json:"status"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	Timestamp int64  `json:"timestamp"`
	Error     string `json:"error,omitempty"`
}

// ChatHandler returns a Gin handler function for chat endpoints.
func (a *GinAdapter) ChatHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), a.timeout)
		defer cancel()

		var req ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ChatResponse{
				Success: false,
				Error:   "Invalid request format: " + err.Error(),
			})
			return
		}

		// Convert context map to AskOptions
		var askOptions []gochatbot.AskOption
		if req.Context != nil {
			for key, value := range req.Context {
				askOptions = append(askOptions, gochatbot.WithContext(key, value))
			}
		}

		response, err := a.chatbot.Ask(ctx, req.Message, askOptions...)
		if err != nil {
			statusCode := http.StatusInternalServerError
			// Check for specific error types
			if ctx.Err() == context.DeadlineExceeded {
				statusCode = http.StatusRequestTimeout
			}

			c.JSON(statusCode, ChatResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, ChatResponse{
			Response: response,
			Success:  true,
		})
	}
}

// HealthHandler returns a Gin handler function for health check endpoints.
func (a *GinAdapter) HealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		response := HealthResponse{
			Status:    "healthy",
			Provider:  a.chatbot.GetModel().Provider(),
			Model:     a.chatbot.GetModel().Name(),
			Timestamp: time.Now().Unix(),
		}

		// Use the chatbot's health check method
		if err := a.chatbot.Health(ctx); err != nil {
			response.Status = "unhealthy"
			response.Error = err.Error()
			c.JSON(http.StatusServiceUnavailable, response)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// StreamChatHandler returns a Gin handler function for streaming chat endpoints.
// This is a placeholder for future streaming implementation.
func (a *GinAdapter) StreamChatHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "Streaming chat not yet implemented",
		})
	}
}

// SetupRoutes sets up the standard chatbot routes on a Gin router.
func (a *GinAdapter) SetupRoutes(router gin.IRouter) {
	chatGroup := router.Group("/chat")
	{
		chatGroup.POST("/", a.ChatHandler())
		chatGroup.POST("/stream", a.StreamChatHandler())
		chatGroup.GET("/health", a.HealthHandler())
	}
}

// SetupRoutesWithPrefix sets up the chatbot routes with a custom prefix.
func (a *GinAdapter) SetupRoutesWithPrefix(router gin.IRouter, prefix string) {
	chatGroup := router.Group(prefix)
	{
		chatGroup.POST("/", a.ChatHandler())
		chatGroup.POST("/stream", a.StreamChatHandler())
		chatGroup.GET("/health", a.HealthHandler())
	}
}

// Middleware returns a Gin middleware that adds chatbot functionality to the context.
func (a *GinAdapter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("chatbot", a.chatbot)
		c.Next()
	}
}

// GetChatbotFromContext extracts the chatbot instance from Gin context.
func GetChatbotFromContext(c *gin.Context) (*gochatbot.Chatbot, bool) {
	bot, exists := c.Get("chatbot")
	if !exists {
		return nil, false
	}

	chatbotInstance, ok := bot.(*gochatbot.Chatbot)
	return chatbotInstance, ok
}
