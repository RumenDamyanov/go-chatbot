# AI Provider Documentation

This document provides detailed information about all supported AI providers in the go-chatbot package.

## Supported Providers

### 1. 🆓 Free Model (Local)
A built-in, offline chatbot that works without any API keys or external dependencies.

**Configuration:**
```go
model := models.NewFreeModel()
```

**Features:**
- No API key required
- Offline operation
- Context-aware responses
- Pattern-based conversation handling
- Instant responses

**Use Case:** Development, testing, or basic chatbot functionality without external dependencies.

---

### 2. 🤖 OpenAI (GPT Models)
Integration with OpenAI's GPT models including GPT-3.5, GPT-4, and newer variants.

**Configuration:**
```go
model, err := models.NewOpenAIModel(config.OpenAIConfig{
    APIKey: "your-api-key",
    Model:  "gpt-3.5-turbo", // or "gpt-4", "gpt-4-turbo", etc.
})
```

**Environment Variables:**
```bash
export OPENAI_API_KEY="your-openai-api-key"
```

**Features:**
- Multiple model variants (GPT-3.5, GPT-4, etc.)
- Streaming responses support
- Conversation history
- Function calling support
- Token usage tracking

**Use Case:** Production applications requiring high-quality conversational AI.

---

### 3. 🧠 Anthropic Claude
Integration with Anthropic's Claude models, known for helpful, harmless, and honest AI.

**Configuration:**
```go
model, err := models.NewAnthropicModel(config.AnthropicConfig{
    APIKey: "your-api-key",
    Model:  "claude-3-haiku-20240307", // or "claude-3-sonnet-20240229", etc.
})
```

**Environment Variables:**
```bash
export ANTHROPIC_API_KEY="your-anthropic-api-key"
```

**Features:**
- Claude 3 model family support
- Large context windows
- Strong reasoning capabilities
- Built-in safety features
- Conversation history support

**Use Case:** Applications requiring nuanced understanding and safe AI responses.

---

### 4. 💎 Google Gemini
Integration with Google's Gemini models, offering multimodal AI capabilities.

**Configuration:**
```go
model, err := models.NewGeminiModel(config.GeminiConfig{
    APIKey: "your-api-key",
    Model:  "gemini-1.5-flash", // or "gemini-1.5-pro", etc.
})
```

**Environment Variables:**
```bash
export GEMINI_API_KEY="your-gemini-api-key"
```

**Features:**
- Gemini 1.5 model support
- Multimodal capabilities (text, images)
- Large context windows
- Safety settings configuration
- Fast response times

**Use Case:** Applications requiring multimodal AI or integration with Google's ecosystem.

---

### 5. 🚀 xAI Grok
Integration with xAI's Grok models, designed for wit and real-time information.

**Configuration:**
```go
model, err := models.NewXAIModel(config.XAIConfig{
    APIKey: "your-api-key",
    Model:  "grok-beta",
})
```

**Environment Variables:**
```bash
export XAI_API_KEY="your-xai-api-key"
```

**Features:**
- OpenAI-compatible API
- Real-time information access
- Conversational and witty responses
- Standard chat completion format

**Use Case:** Applications requiring up-to-date information and engaging conversational style.

---

### 6. 🦙 Meta LLaMA
Integration with Meta's LLaMA models through various providers.

**Configuration:**
```go
model, err := models.NewMetaModel(config.MetaConfig{
    APIKey: "your-api-key",
    Model:  "llama-3.2-3b-instruct",
    Endpoint: "https://your-llama-provider.com", // e.g., Together AI, Replicate
})
```

**Environment Variables:**
```bash
export META_API_KEY="your-meta-api-key"
```

**Features:**
- Various LLaMA model sizes
- Open-source foundation
- Customizable endpoints
- Cost-effective for high volume

**Use Case:** Applications requiring open-source models or custom deployments.

---

### 7. 🏠 Ollama (Local Models)
Integration with Ollama for running large language models locally.

**Configuration:**
```go
model, err := models.NewOllamaModel(config.OllamaConfig{
    Model:    "llama3.2",
    Endpoint: "http://localhost:11434", // default
})
```

**Setup:**
```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Pull a model
ollama pull llama3.2

# Start Ollama server
ollama serve
```

**Features:**
- Local model execution
- No API keys required
- Privacy-focused (no data leaves your machine)
- Multiple model support
- Both chat and completion APIs
- Custom model parameters

**Use Case:** Privacy-sensitive applications, offline operation, or development without API costs.

## Usage Examples

### Basic Usage
```go
package main

import (
    "context"
    "fmt"
    "go.rumenx.com/chatbot/models"
    "go.rumenx.com/chatbot/config"
)

func main() {
    // Using OpenAI
    model, err := models.NewOpenAIModel(config.OpenAIConfig{
        APIKey: "your-api-key",
        Model:  "gpt-3.5-turbo",
    })
    if err != nil {
        panic(err)
    }

    response, err := model.Ask(context.Background(), "Hello!", nil)
    if err != nil {
        panic(err)
    }

    fmt.Println(response)
}
```

### With Conversation History
```go
context := map[string]interface{}{
    "history": []map[string]interface{}{
        {"role": "user", "content": "What is Go?"},
        {"role": "assistant", "content": "Go is a programming language..."},
    },
    "system": "You are a helpful programming assistant.",
}

response, err := model.Ask(ctx, "Tell me more about Go's concurrency", context)
```

### Model Factory
```go
cfg := &config.Config{
    Model: "openai",
    OpenAI: config.OpenAIConfig{
        APIKey: "your-api-key",
        Model:  "gpt-3.5-turbo",
    },
}

model, err := models.NewFromConfig(cfg)
if err != nil {
    panic(err)
}

response, err := model.Ask(context.Background(), "Hello!", nil)
```

## Health Checks

All models support health checking:

```go
// Check if the model/API is accessible
if healthChecker, ok := model.(models.HealthChecker); ok {
    if err := healthChecker.Health(context.Background()); err != nil {
        fmt.Printf("Model health check failed: %v\n", err)
    }
}
```

## Error Handling

All models return detailed error information:

```go
response, err := model.Ask(ctx, "Hello", nil)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "invalid"):
        // Handle API key issues
    case strings.Contains(err.Error(), "timeout"):
        // Handle timeout issues
    case strings.Contains(err.Error(), "rate limit"):
        // Handle rate limiting
    default:
        // Handle other errors
    }
}
```

## Testing

Run the provider examples to test all integrations:

```bash
cd examples/providers
go run main.go
```

Set environment variables for the providers you want to test:

```bash
export OPENAI_API_KEY="your-key"
export ANTHROPIC_API_KEY="your-key"
export GEMINI_API_KEY="your-key"
export XAI_API_KEY="your-key"
export META_API_KEY="your-key"
```

## Cost Considerations

| Provider | Cost Model | Notes |
|----------|------------|-------|
| Free Model | Free | No costs, local processing |
| OpenAI | Per token | Higher quality, various pricing tiers |
| Anthropic | Per token | Competitive pricing, large context |
| Gemini | Per token | Generous free tier, competitive pricing |
| xAI | Per token | Newer provider, competitive rates |
| Meta | Varies | Depends on hosting provider |
| Ollama | Hardware only | One-time hardware cost, no ongoing fees |

## Best Practices

1. **Use Free Model for Development**: Start with the free model during development
2. **Environment Variables**: Store API keys securely in environment variables
3. **Error Handling**: Always handle API errors gracefully
4. **Context Management**: Use context for timeouts and cancellation
5. **Health Checks**: Implement health checks for production systems
6. **Rate Limiting**: Respect API rate limits and implement retry logic
7. **Local Models**: Consider Ollama for privacy-sensitive applications

## Contributing

To add a new AI provider:

1. Create a new file `models/newprovider.go`
2. Implement the `Model` interface
3. Add configuration to `config/config.go`
4. Update the model factory in `models/models.go`
5. Add tests in `models/newprovider_test.go`
6. Update this documentation
