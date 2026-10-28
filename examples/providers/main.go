package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.rumenx.com/chatbot/config"
	"go.rumenx.com/chatbot/models"
)

func main() {
	fmt.Println("🤖 Go Chatbot - AI Provider Examples")
	fmt.Println("=====================================")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test message
	message := "Hello! Can you tell me a brief fact about Go programming language?"

	// 1. Free Model (always available)
	fmt.Println("\n1. 🆓 Free Model:")
	testFreeModel(ctx, message)

	// 2. OpenAI (if API key is available)
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		fmt.Println("\n2. 🤖 OpenAI:")
		testOpenAI(ctx, message, apiKey)
	} else {
		fmt.Println("\n2. 🤖 OpenAI: Skipped (OPENAI_API_KEY not set)")
	}

	// 3. Anthropic (if API key is available)
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		fmt.Println("\n3. 🧠 Anthropic Claude:")
		testAnthropic(ctx, message, apiKey)
	} else {
		fmt.Println("\n3. 🧠 Anthropic Claude: Skipped (ANTHROPIC_API_KEY not set)")
	}

	// 4. Google Gemini (if API key is available)
	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		fmt.Println("\n4. 💎 Google Gemini:")
		testGemini(ctx, message, apiKey)
	} else {
		fmt.Println("\n4. 💎 Google Gemini: Skipped (GEMINI_API_KEY not set)")
	}

	// 5. xAI Grok (if API key is available)
	if apiKey := os.Getenv("XAI_API_KEY"); apiKey != "" {
		fmt.Println("\n5. 🚀 xAI Grok:")
		testXAI(ctx, message, apiKey)
	} else {
		fmt.Println("\n5. 🚀 xAI Grok: Skipped (XAI_API_KEY not set)")
	}

	// 6. Meta LLaMA (if API key is available)
	if apiKey := os.Getenv("META_API_KEY"); apiKey != "" {
		fmt.Println("\n6. 🦙 Meta LLaMA:")
		testMeta(ctx, message, apiKey)
	} else {
		fmt.Println("\n6. 🦙 Meta LLaMA: Skipped (META_API_KEY not set)")
	}

	// 7. Ollama (if running locally)
	fmt.Println("\n7. 🏠 Ollama (Local):")
	testOllama(ctx, message)

	fmt.Println("\n✅ AI Provider examples completed!")
	fmt.Println("\nTo test with real APIs, set the following environment variables:")
	fmt.Println("  export OPENAI_API_KEY=\"your-openai-key\"")
	fmt.Println("  export ANTHROPIC_API_KEY=\"your-anthropic-key\"")
	fmt.Println("  export GEMINI_API_KEY=\"your-gemini-key\"")
	fmt.Println("  export XAI_API_KEY=\"your-xai-key\"")
	fmt.Println("  export META_API_KEY=\"your-meta-key\"")
	fmt.Println("\nFor Ollama, install and run locally: https://ollama.ai")
}

func testFreeModel(ctx context.Context, message string) {
	model := models.NewFreeModel()

	fmt.Printf("   Provider: %s\n", model.Provider())
	fmt.Printf("   Model: %s\n", model.Name())

	response, err := model.Ask(ctx, message, nil)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Response: %s\n", response)
}

func testOpenAI(ctx context.Context, message, apiKey string) {
	model, err := models.NewOpenAIModel(config.OpenAIConfig{
		APIKey: apiKey,
		Model:  "gpt-3.5-turbo",
	})
	if err != nil {
		fmt.Printf("   ❌ Failed to create model: %v\n", err)
		return
	}

	fmt.Printf("   Provider: %s\n", model.Provider())
	fmt.Printf("   Model: %s\n", model.Name())

	response, err := model.Ask(ctx, message, nil)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Response: %s\n", response)
}

func testAnthropic(ctx context.Context, message, apiKey string) {
	model, err := models.NewAnthropicModel(config.AnthropicConfig{
		APIKey: apiKey,
		Model:  "claude-sonnet-5",
	})
	if err != nil {
		fmt.Printf("   ❌ Failed to create model: %v\n", err)
		return
	}

	fmt.Printf("   Provider: %s\n", model.Provider())
	fmt.Printf("   Model: %s\n", model.Name())

	response, err := model.Ask(ctx, message, nil)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Response: %s\n", response)
}

func testGemini(ctx context.Context, message, apiKey string) {
	model, err := models.NewGeminiModel(config.GeminiConfig{
		APIKey: apiKey,
		Model:  "gemini-3.8-flash",
	})
	if err != nil {
		fmt.Printf("   ❌ Failed to create model: %v\n", err)
		return
	}

	fmt.Printf("   Provider: %s\n", model.Provider())
	fmt.Printf("   Model: %s\n", model.Name())

	response, err := model.Ask(ctx, message, nil)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Response: %s\n", response)
}

func testXAI(ctx context.Context, message, apiKey string) {
	model, err := models.NewXAIModel(config.XAIConfig{
		APIKey: apiKey,
		Model:  "grok-4.6",
	})
	if err != nil {
		fmt.Printf("   ❌ Failed to create model: %v\n", err)
		return
	}

	fmt.Printf("   Provider: %s\n", model.Provider())
	fmt.Printf("   Model: %s\n", model.Name())

	response, err := model.Ask(ctx, message, nil)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Response: %s\n", response)
}

func testMeta(ctx context.Context, message, apiKey string) {
	model, err := models.NewMetaModel(config.MetaConfig{
		APIKey: apiKey,
		Model:  "llama-3.3-70b-instruct",
	})
	if err != nil {
		fmt.Printf("   ❌ Failed to create model: %v\n", err)
		return
	}

	fmt.Printf("   Provider: %s\n", model.Provider())
	fmt.Printf("   Model: %s\n", model.Name())

	response, err := model.Ask(ctx, message, nil)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Response: %s\n", response)
}

func testOllama(ctx context.Context, message string) {
	model, err := models.NewOllamaModel(config.OllamaConfig{
		Model: "llama3.3", // Default model
	})
	if err != nil {
		fmt.Printf("   ❌ Failed to create model: %v\n", err)
		return
	}

	fmt.Printf("   Provider: %s\n", model.Provider())
	fmt.Printf("   Model: %s\n", model.Name())

	// First check if Ollama is running
	if err := model.Health(ctx); err != nil {
		fmt.Printf("   ⚠️  Ollama not available: %v\n", err)
		fmt.Printf("   💡 To use Ollama:\n")
		fmt.Printf("      1. Install: https://ollama.ai\n")
		fmt.Printf("      2. Run: ollama pull llama3.3\n")
		fmt.Printf("      3. Start: ollama serve\n")
		return
	}

	response, err := model.Ask(ctx, message, nil)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Response: %s\n", response)
}

func init() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
