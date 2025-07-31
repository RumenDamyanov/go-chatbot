package adapters

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	gochatbot "github.com/RumenDamyanov/go-chatbot"
)

// FiberAdapter provides Fiber framework integration for go-chatbot.
type FiberAdapter struct {
	chatbot *gochatbot.Chatbot
	timeout time.Duration
}

// NewFiberAdapter creates a new Fiber adapter with the provided chatbot instance.
func NewFiberAdapter(bot *gochatbot.Chatbot) *FiberAdapter {
	return &FiberAdapter{
		chatbot: bot,
		timeout: 30 * time.Second,
	}
}

// WithTimeout sets the request timeout for the adapter.
func (a *FiberAdapter) WithTimeout(timeout time.Duration) *FiberAdapter {
	a.timeout = timeout
	return a
}

// ChatHandler returns a Fiber handler function for chat endpoints.
func (a *FiberAdapter) ChatHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), a.timeout)
		defer cancel()

		var req ChatRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ChatResponse{
				Success: false,
				Error:   "Invalid request format: " + err.Error(),
			})
		}

		// Validate required fields
		if req.Message == "" {
			return c.Status(fiber.StatusBadRequest).JSON(ChatResponse{
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
			statusCode := fiber.StatusInternalServerError
			// Check for specific error types
			if ctx.Err() == context.DeadlineExceeded {
				statusCode = fiber.StatusRequestTimeout
			}

			return c.Status(statusCode).JSON(ChatResponse{
				Success: false,
				Error:   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(ChatResponse{
			Response: response,
			Success:  true,
		})
	}
}

// HealthHandler returns a Fiber handler function for health check endpoints.
func (a *FiberAdapter) HealthHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
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
			return c.Status(fiber.StatusServiceUnavailable).JSON(response)
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

// StreamChatHandler returns a Fiber handler function for streaming chat endpoints.
// This is a placeholder for future streaming implementation.
func (a *FiberAdapter) StreamChatHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"error": "Streaming chat not yet implemented",
		})
	}
}

// SetupRoutes sets up the standard chatbot routes on a Fiber app.
func (a *FiberAdapter) SetupRoutes(app *fiber.App) {
	chatGroup := app.Group("/chat")
	chatGroup.Post("/", a.ChatHandler())
	chatGroup.Post("/stream", a.StreamChatHandler())
	chatGroup.Get("/health", a.HealthHandler())
}

// SetupRoutesWithPrefix sets up the chatbot routes with a custom prefix.
func (a *FiberAdapter) SetupRoutesWithPrefix(app *fiber.App, prefix string) {
	chatGroup := app.Group(prefix)
	chatGroup.Post("/", a.ChatHandler())
	chatGroup.Post("/stream", a.StreamChatHandler())
	chatGroup.Get("/health", a.HealthHandler())
}

// Middleware returns a Fiber middleware that adds chatbot functionality to the context.
func (a *FiberAdapter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("chatbot", a.chatbot)
		return c.Next()
	}
}

// GetChatbotFromFiberContext extracts the chatbot instance from Fiber context.
func GetChatbotFromFiberContext(c *fiber.Ctx) (*gochatbot.Chatbot, bool) {
	bot := c.Locals("chatbot")
	if bot == nil {
		return nil, false
	}

	chatbotInstance, ok := bot.(*gochatbot.Chatbot)
	return chatbotInstance, ok
}
