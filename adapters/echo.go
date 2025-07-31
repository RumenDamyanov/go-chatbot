package adapters

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	gochatbot "github.com/RumenDamyanov/go-chatbot"
)

// EchoAdapter provides Echo framework integration for go-chatbot.
type EchoAdapter struct {
	chatbot *gochatbot.Chatbot
	timeout time.Duration
}

// NewEchoAdapter creates a new Echo adapter with the provided chatbot instance.
func NewEchoAdapter(bot *gochatbot.Chatbot) *EchoAdapter {
	return &EchoAdapter{
		chatbot: bot,
		timeout: 30 * time.Second,
	}
}

// WithTimeout sets the request timeout for the adapter.
func (a *EchoAdapter) WithTimeout(timeout time.Duration) *EchoAdapter {
	a.timeout = timeout
	return a
}

// ChatHandler returns an Echo handler function for chat endpoints.
func (a *EchoAdapter) ChatHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), a.timeout)
		defer cancel()

		var req ChatRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, ChatResponse{
				Success: false,
				Error:   "Invalid request format: " + err.Error(),
			})
		}

		// Validate required fields
		if req.Message == "" {
			return c.JSON(http.StatusBadRequest, ChatResponse{
				Success: false,
				Error:   "Message is required",
			})
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

			return c.JSON(statusCode, ChatResponse{
				Success: false,
				Error:   err.Error(),
			})
		}

		return c.JSON(http.StatusOK, ChatResponse{
			Response: response,
			Success:  true,
		})
	}
}

// HealthHandler returns an Echo handler function for health check endpoints.
func (a *EchoAdapter) HealthHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
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
			return c.JSON(http.StatusServiceUnavailable, response)
		}

		return c.JSON(http.StatusOK, response)
	}
}

// StreamChatHandler returns an Echo handler function for streaming chat endpoints.
// This is a placeholder for future streaming implementation.
func (a *EchoAdapter) StreamChatHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, map[string]string{
			"error": "Streaming chat not yet implemented",
		})
	}
}

// SetupRoutes sets up the standard chatbot routes on an Echo router.
func (a *EchoAdapter) SetupRoutes(e *echo.Echo) {
	chatGroup := e.Group("/chat")
	chatGroup.POST("/", a.ChatHandler())
	chatGroup.POST("/stream", a.StreamChatHandler())
	chatGroup.GET("/health", a.HealthHandler())
}

// SetupRoutesWithPrefix sets up the chatbot routes with a custom prefix.
func (a *EchoAdapter) SetupRoutesWithPrefix(e *echo.Echo, prefix string) {
	chatGroup := e.Group(prefix)
	chatGroup.POST("/", a.ChatHandler())
	chatGroup.POST("/stream", a.StreamChatHandler())
	chatGroup.GET("/health", a.HealthHandler())
}

// Middleware returns an Echo middleware that adds chatbot functionality to the context.
func (a *EchoAdapter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("chatbot", a.chatbot)
			return next(c)
		}
	}
}

// GetChatbotFromEchoContext extracts the chatbot instance from Echo context.
func GetChatbotFromEchoContext(c echo.Context) (*gochatbot.Chatbot, bool) {
	bot := c.Get("chatbot")
	if bot == nil {
		return nil, false
	}

	chatbotInstance, ok := bot.(*gochatbot.Chatbot)
	return chatbotInstance, ok
}
