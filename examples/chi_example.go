//go:build ignore
// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	// Create Chi router
	r := chi.NewRouter()

	// Create Chi adapter and setup routes
	adapter := adapters.NewChiAdapter(chatbot)
	adapter.SetupRoutes(r)

	// Add middleware for context injection (optional)
	r.Use(adapter.Middleware())

	// Add a custom route that uses the chatbot from context
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		if bot, exists := adapters.GetChatbotFromChiContext(r); exists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"status": "ok",
				"model": "` + bot.GetModel().Name() + `",
				"provider": "` + bot.GetModel().Provider() + `"
			}`
			w.Write([]byte(response))
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "chatbot not found in context"}`))
		}
	})

	log.Println("Starting Chi server on :8080")
	log.Println("Chat endpoint: http://localhost:8080/chat/")
	log.Println("Health endpoint: http://localhost:8080/chat/health")
	log.Println("Status endpoint: http://localhost:8080/status")

	log.Fatal(http.ListenAndServe(":8080", r))
}
