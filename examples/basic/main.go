package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	gochatbot "github.com/RumenDamyanov/go-chatbot"
	"github.com/RumenDamyanov/go-chatbot/config"
)

func main() {
	// Create configuration
	cfg := config.Default()

	// You can override settings programmatically
	cfg.Prompt = "You are a helpful assistant for a Go chatbot demo."

	// Create chatbot instance
	chatbot, err := gochatbot.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create chatbot: %v", err)
	}

	// Test the chatbot directly
	ctx := context.Background()
	reply, err := chatbot.Ask(ctx, "Hello! How are you?")
	if err != nil {
		log.Printf("Error asking chatbot: %v", err)
	} else {
		fmt.Printf("Chatbot reply: %s\n", reply)
	}

	// Set up HTTP server
	http.HandleFunc("/api/chat", chatbot.HandleHTTP)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		handler := gochatbot.NewHTTPHandler(chatbot)
		handler.Health(w, r)
	})

	// Serve static files (if you want to include frontend)
	http.Handle("/", http.FileServer(http.Dir("./web/")))

	fmt.Println("Starting server on :8080")
	fmt.Println("Chat endpoint: http://localhost:8080/api/chat")
	fmt.Println("Health check: http://localhost:8080/health")
	fmt.Println("Web interface: http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
