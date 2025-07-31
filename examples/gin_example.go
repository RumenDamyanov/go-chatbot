//go:build ignore
// +build ignore

package main

import (
	"log"

	gochatbot "github.com/RumenDamyanov/go-chatbot"
	"github.com/RumenDamyanov/go-chatbot/adapters"
	"github.com/RumenDamyanov/go-chatbot/config"
	"github.com/gin-gonic/gin"
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

	// Create Gin router
	r := gin.Default()

	// Create Gin adapter and setup routes
	adapter := adapters.NewGinAdapter(chatbot)
	adapter.SetupRoutes(r)

	// Add middleware for context injection (optional)
	r.Use(adapter.Middleware())

	// Add a custom route that uses the chatbot from context
	r.GET("/status", func(c *gin.Context) {
		if bot, exists := adapters.GetChatbotFromContext(c); exists {
			c.JSON(200, gin.H{
				"status":   "ok",
				"model":    bot.GetModel().Name(),
				"provider": bot.GetModel().Provider(),
			})
		} else {
			c.JSON(500, gin.H{"error": "chatbot not found in context"})
		}
	})

	log.Println("Starting Gin server on :8080")
	log.Println("Chat endpoint: http://localhost:8080/chat/")
	log.Println("Health endpoint: http://localhost:8080/chat/health")
	log.Println("Status endpoint: http://localhost:8080/status")

	r.Run(":8080")
}
