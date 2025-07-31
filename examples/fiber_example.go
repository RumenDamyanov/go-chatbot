//go:build ignore
// +build ignore

package main

import (
	"log"

	gochatbot "github.com/RumenDamyanov/go-chatbot"
	"github.com/RumenDamyanov/go-chatbot/adapters"
	"github.com/RumenDamyanov/go-chatbot/config"
	"github.com/gofiber/fiber/v2"
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

	// Create Fiber app
	app := fiber.New()

	// Create Fiber adapter and setup routes
	adapter := adapters.NewFiberAdapter(chatbot)
	adapter.SetupRoutes(app)

	// Add middleware for context injection (optional)
	app.Use(adapter.Middleware())

	// Add a custom route that uses the chatbot from context
	app.Get("/status", func(c *fiber.Ctx) error {
		if bot, exists := adapters.GetChatbotFromFiberContext(c); exists {
			return c.JSON(fiber.Map{
				"status":   "ok",
				"model":    bot.GetModel().Name(),
				"provider": bot.GetModel().Provider(),
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "chatbot not found in context",
		})
	})

	log.Println("Starting Fiber server on :8080")
	log.Println("Chat endpoint: http://localhost:8080/chat/")
	log.Println("Health endpoint: http://localhost:8080/chat/health")
	log.Println("Status endpoint: http://localhost:8080/status")

	log.Fatal(app.Listen(":8080"))
}
