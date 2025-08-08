//go:build ignore
// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	gochatbot "go.rumenx.com/chatbot"
	"go.rumenx.com/chatbot/adapters"
	"go.rumenx.com/chatbot/config"
)

func main() {
	// Create configuration with free model for demo
	cfg := config.Default()
	cfg.Model = "free"

	// Create chatbot instance
	chatbot, err := gochatbot.New(cfg)
	if err != nil {
		log.Fatal("Failed to create chatbot:", err)
	}

	// Create Echo server
	e := echo.New()

	// Create Echo adapter and setup routes
	adapter := adapters.NewEchoAdapter(chatbot)
	adapter.SetupRoutes(e)

	// Add middleware for context injection (optional)
	e.Use(adapter.Middleware())

	// Add a custom route that uses the chatbot from context
	e.GET("/status", func(c echo.Context) error {
		if bot, exists := adapters.GetChatbotFromEchoContext(c); exists {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status":   "ok",
				"model":    bot.GetModel().Name(),
				"provider": bot.GetModel().Provider(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "chatbot not found in context",
		})
	})

	log.Println("Starting Echo server on :8080")
	log.Println("Chat endpoint: http://localhost:8080/chat/")
	log.Println("Health endpoint: http://localhost:8080/chat/health")
	log.Println("Status endpoint: http://localhost:8080/status")

	e.Logger.Fatal(e.Start(":8080"))
}
