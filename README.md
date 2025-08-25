# go-chatbot

[![CI](https://github.com/rumendamyanov/go-chatbot/actions/workflows/ci.yml/badge.svg)](https://github.com/rumendamyanov/go-chatbot/actions/workflows/ci.yml)
![CodeQL](https://github.com/rumendamyanov/go-chatbot/actions/workflows/github-code-scanning/codeql/badge.svg)
![Dependabot](https://github.com/rumendamyanov/go-chatbot/actions/workflows/dependabot/dependabot-updates/badge.svg)
[![codecov](https://codecov.io/gh/rumendamyanov/go-chatbot/branch/master/graph/badge.svg)](https://codecov.io/gh/rumendamyanov/go-chatbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/rumendamyanov/go-chatbot?)](https://goreportcard.com/report/github.com/rumendamyanov/go-chatbot)
[![Go Reference](https://pkg.go.dev/badge/go.rumenx.com/chatbot.svg)](https://pkg.go.dev/go.rumenx.com/chatbot)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/rumendamyanov/go-chatbot/blob/master/LICENSE.md)

> 📖 **Documentation**: [Contributing](CONTRIBUTING.md) · [Security](SECURITY.md) · [Changelog](CHANGELOG.md) · [Funding](FUNDING.md)

**go-chatbot** is a modern, framework-agnostic Go package for integrating AI-powered chat functionality into any web application. It features out-of-the-box support for popular Go web frameworks, a flexible model abstraction for using OpenAI, Anthropic, xAI, Google Gemini, Meta, and more, and is designed for easy customization and extension. Build your own UI or use the provided minimal frontend as a starting point. High test coverage, static analysis, and coding standards are included.

## About

A modern, framework-agnostic Go package for adding customizable, AI-powered chat functionality to any web application. Includes adapters for popular Go web frameworks, supports multiple AI providers, and is easy to extend.

### Project Inspiration

This project is inspired by and related to the [php-chatbot](https://github.com/RumenDamyanov/php-chatbot) project, bringing similar AI-powered chat functionality to the Go ecosystem with enhanced enterprise features and modern architecture patterns.

## ✨ Key Features

### Core Functionality

- Plug-and-play chat API and UI components (minimal frontend dependencies)
- Go web framework support via adapters/middleware (Gin, Echo, Fiber, Chi, etc.)
- AI model abstraction via interfaces (swap models easily)
- Customizable prompts, tone, language, and scope
- Emoji and multi-script support (Cyrillic, Greek, Armenian, Asian, etc.)
- Security safeguards against abuse
- High test coverage (Go standard testing + testify)
- Static analysis (golangci-lint, staticcheck) & Go coding standards (gofmt, goimports)
- Structured logging with slog
- Context-aware operations with cancellation support
- Graceful error handling and recovery

### 🚀 Advanced Features

- **Real-time Streaming**: Server-Sent Events (SSE) for live response streaming
- **Vector Embeddings**: OpenAI embeddings integration with similarity search capabilities
- **Database Persistence**: Full conversation history with SQLite/PostgreSQL support
- **Knowledge Base**: Context-aware responses using embedded knowledge
- **Session Management**: Multi-user conversation tracking and persistence

## 📚 Documentation & Wiki

Comprehensive documentation and guides are available in our [GitHub Wiki](https://github.com/rumendamyanov/go-chatbot/wiki):

### 🚀 Getting Started
- **[Installation Guide](https://github.com/rumendamyanov/go-chatbot/wiki/Installation-Guide)** - Step-by-step installation for all environments
- **[Quick Start Guide](https://github.com/rumendamyanov/go-chatbot/wiki/Quick-Start-Guide)** - Get up and running in minutes
- **[Configuration](https://github.com/rumendamyanov/go-chatbot/wiki/Configuration)** - Complete configuration reference

### 🔧 Implementation Guides
- **[Framework Integration](https://github.com/rumendamyanov/go-chatbot/wiki/Framework-Integration)** - Gin, Echo, Fiber, Chi, and vanilla net/http setup
- **[Frontend Integration](https://github.com/rumendamyanov/go-chatbot/wiki/Frontend-Integration)** - React, Vue, Angular components and examples
- **[AI Models](https://github.com/rumendamyanov/go-chatbot/wiki/AI-Models)** - Provider comparison and configuration

### 📖 Examples & Best Practices
- **[Examples](https://github.com/rumendamyanov/go-chatbot/wiki/Examples)** - Real-world implementations and use cases
- **[Best Practices](https://github.com/rumendamyanov/go-chatbot/wiki/Best-Practices)** - Production deployment and security guidelines
- **[Security & Filtering](https://github.com/rumendamyanov/go-chatbot/wiki/Security-and-Filtering)** - Content filtering and abuse prevention

### 🛠️ Development & Support
- **[API Reference](https://pkg.go.dev/go.rumenx.com/chatbot)** - Complete API documentation
- **[Troubleshooting](https://github.com/rumendamyanov/go-chatbot/wiki/Troubleshooting)** - Common issues and solutions
- **[Contributing](https://github.com/rumendamyanov/go-chatbot/wiki/Contributing)** - How to contribute to the project
- **[FAQ](https://github.com/rumendamyanov/go-chatbot/wiki/FAQ)** - Frequently asked questions

> 💡 **Tip**: The wiki contains production-ready examples, troubleshooting guides, and comprehensive API documentation that goes beyond this README.

### 📋 Technical Documentation

For detailed technical summaries and architectural overviews, see the [`docs/`](docs/) directory:

- **[Project Summary](docs/PROJECT_SUMMARY.md)** - Complete project overview and capabilities
- **[Advanced Features](docs/ADVANCED_FEATURES_SUMMARY.md)** - Streaming, embeddings, and database implementation details
- **[AI Providers Guide](docs/AI_PROVIDERS_SUMMARY.md)** - Provider comparison and integration guide

## Supported AI Providers & Models

| Provider | Models | API Key Required | Type |
|----------|--------|------------------|------|
| OpenAI | gpt-4.1, gpt-4o, gpt-4o-mini, gpt-3.5-turbo, etc. | Yes | Remote |
| Anthropic | Claude 3 Sonnet, 3.7, 4, etc. | Yes | Remote |
| xAI | Grok-1, Grok-1.5, etc. | Yes | Remote |
| Google | Gemini 1.5 Pro, Gemini 1.5 Flash, etc. | Yes | Remote |
| Meta | Llama 3 (8B, 70B), etc. | Yes | Remote |
| Ollama | llama2, mistral, phi3, and any local Ollama model | No (local) / Opt | Local/Remote |
| Free model | Simple fallback, no API key required | No | Local |

## Framework Adapters

The package includes built-in adapters for popular Go web frameworks, providing consistent API patterns and easy integration:

### Available Adapters

| Framework | Package | Adapter Type | Features |
|-----------|---------|--------------|----------|
| **Gin** | `github.com/gin-gonic/gin` | `adapters.GinAdapter` | Route setup, middleware, context extraction |
| **Echo** | `github.com/labstack/echo/v4` | `adapters.EchoAdapter` | Handler functions, middleware, context extraction |
| **Fiber** | `github.com/gofiber/fiber/v2` | `adapters.FiberAdapter` | Fast HTTP handlers, middleware, context extraction |
| **Chi** | `github.com/go-chi/chi/v5` | `adapters.ChiAdapter` | Standard net/http compatible, middleware, context extraction |

### Common Adapter Features

All adapters provide:
- **Chat Handler**: `POST /chat/` - Process chat messages
- **Health Handler**: `GET /chat/health` - Health check endpoint
- **Stream Handler**: `POST /chat/stream` - Streaming responses (placeholder)
- **Middleware**: Inject chatbot instance into request context
- **Route Setup**: Easy route configuration with optional custom prefixes
- **Timeout Support**: Configurable request timeouts
- **Error Handling**: Consistent JSON error responses

### Basic Usage Pattern

```go
// 1. Create chatbot instance
chatbot, err := gochatbot.New(config)

// 2. Create framework-specific adapter
adapter := adapters.NewGinAdapter(chatbot)        // or NewEchoAdapter, NewFiberAdapter, NewChiAdapter

// 3. Setup routes
adapter.SetupRoutes(router)                       // Default: /chat/*, /chat/health, /chat/stream
adapter.SetupRoutesWithPrefix(router, "/api/v1") // Custom prefix: /api/v1/*, /api/v1/health, /api/v1/stream

// 4. Optional: Use middleware for context injection
router.Use(adapter.Middleware())
```

## Installation

```bash
go get go.rumenx.com/chatbot
```

### Gin Framework

```go
package main

import (
    "github.com/gin-gonic/gin"
    "go.rumenx.com/chatbot"
    "go.rumenx.com/chatbot/adapters"
    "go.rumenx.com/chatbot/config"
)

func main() {
    // Create configuration and chatbot
    cfg := config.Default()
    chatbot, err := gochatbot.New(cfg)
    if err != nil {
        panic(err)
    }

    // Create Gin router and adapter
    r := gin.Default()
    adapter := adapters.NewGinAdapter(chatbot)

    // Setup chatbot routes
    adapter.SetupRoutes(r)
    // Or with custom prefix: adapter.SetupRoutesWithPrefix(r, "/api/v1/bot")

    r.Run(":8080")
}
```

### Echo Framework

```go
package main

import (
    "github.com/labstack/echo/v4"
    "go.rumenx.com/chatbot"
    "go.rumenx.com/chatbot/adapters"
    "go.rumenx.com/chatbot/config"
)

func main() {
    // Create configuration and chatbot
    cfg := config.Default()
    chatbot, err := gochatbot.New(cfg)
    if err != nil {
        panic(err)
    }

    // Create Echo server and adapter
    e := echo.New()
    adapter := adapters.NewEchoAdapter(chatbot)

    // Setup chatbot routes
    adapter.SetupRoutes(e)
    // Or with custom prefix: adapter.SetupRoutesWithPrefix(e, "/api/v1/bot")

    e.Logger.Fatal(e.Start(":8080"))
}
```

### Fiber Framework

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "go.rumenx.com/chatbot"
    "go.rumenx.com/chatbot/adapters"
    "go.rumenx.com/chatbot/config"
)

func main() {
    // Create configuration and chatbot
    cfg := config.Default()
    chatbot, err := gochatbot.New(cfg)
    if err != nil {
        panic(err)
    }

    // Create Fiber app and adapter
    app := fiber.New()
    adapter := adapters.NewFiberAdapter(chatbot)

    // Setup chatbot routes
    adapter.SetupRoutes(app)
    // Or with custom prefix: adapter.SetupRoutesWithPrefix(app, "/api/v1/bot")

    app.Listen(":8080")
}
```

### Chi Framework

```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "go.rumenx.com/chatbot"
    "go.rumenx.com/chatbot/adapters"
    "go.rumenx.com/chatbot/config"
)

func main() {
    // Create configuration and chatbot
    cfg := config.Default()
    chatbot, err := gochatbot.New(cfg)
    if err != nil {
        panic(err)
    }

    // Create Chi router and adapter
    r := chi.NewRouter()
    adapter := adapters.NewChiAdapter(chatbot)

    // Setup chatbot routes
    adapter.SetupRoutes(r)
    // Or with custom prefix: adapter.SetupRoutesWithPrefix(r, "/api/v1/bot")

    http.ListenAndServe(":8080", r)
}
```

### Vanilla net/http

```go
import (
    "net/http"
    "go.rumenx.com/chatbot"
    "go.rumenx.com/chatbot/config"
)

func main() {
    cfg := config.Default()
    chatbot := gochatbot.New(cfg)

    http.HandleFunc("/api/chat", chatbot.HandleHTTP)
    http.ListenAndServe(":8080", nil)
}
```

## Usage

1. Add the chat endpoint to your web application using one of the framework adapters.
2. Configure your preferred AI model and prompts in the config.
3. Optionally, customize the frontend (CSS/JS) in `/web`.

### API Keys & Credentials

**Never hardcode API keys or secrets in your codebase.**

- Use environment variables or your infrastructure's secret management.
- The config will check for environment variables (e.g. `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, etc.) first.
- See `.env.example` for reference.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "go.rumenx.com/chatbot"
    "go.rumenx.com/chatbot/config"
    "go.rumenx.com/chatbot/models"
)

func main() {
    // Create configuration
    cfg := config.Default()
    cfg.Model = "openai" // or "anthropic", "gemini", etc.
    cfg.OpenAI.APIKey = "your-api-key"

    // Create model and chatbot
    model, err := models.NewFromConfig(cfg)
    if err != nil {
        log.Fatal(err)
    }

    chatbot := gochatbot.New(cfg, gochatbot.WithModel(model))

    // Ask a question
    ctx := context.Background()
    reply, err := chatbot.Ask(ctx, "Hello!")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Reply:", reply)
}
```

## �️ Advanced Features Implementation

The package includes three enterprise-level features that elevate it beyond basic chatbot functionality:

### Real-time Streaming Responses

Server-Sent Events (SSE) streaming for live response delivery:

```go
import "go.rumenx.com/chatbot/streaming"

// Create stream handler
streamHandler := streaming.NewStreamHandler()

// Start streaming
err := chatbot.AskStream(ctx, "Tell me a story", func(chunk string) error {
    return streamHandler.WriteChunk(chunk)
})
```

**Features:**

- Real-time token streaming with SSE
- Automatic chunk processing and error handling
- Context cancellation support
- Browser and curl compatible

### Vector Embeddings & Knowledge Base

OpenAI embeddings integration with semantic search:

```go
import "go.rumenx.com/chatbot/embeddings"

// Create embedding provider
provider := embeddings.NewOpenAIEmbeddingProvider(apiKey, "text-embedding-3-small")

// Create knowledge base
vectorStore := embeddings.NewInMemoryVectorStore()

// Add knowledge
err := vectorStore.Add(ctx, provider, "id1", "Go is a programming language...")

// Search for relevant context
results, err := vectorStore.Search(ctx, provider, "What is Go?", 5)
```

**Features:**

- OpenAI text-embedding-3-small/large support
- Vector similarity search with cosine distance
- In-memory vector store with search capabilities
- Context enhancement for intelligent responses

### Database Persistence & Conversation History

Full SQL-based conversation management:

```go
import "go.rumenx.com/chatbot/database"

// Create conversation store
store, err := database.NewSQLConversationStore("sqlite", "./chat.db")

// Create conversation
convID, err := store.CreateConversation(ctx, "user123", "Chat Title")

// Save messages
err = store.SaveMessage(ctx, convID, "user", "Hello!")
err = store.SaveMessage(ctx, convID, "assistant", "Hi there!")

// Get conversation history
history, err := store.GetConversation(ctx, convID)
```

**Features:**

- SQLite and PostgreSQL support
- Complete conversation and message CRUD operations
- User session management
- Conversation history and search functionality

### Complete Integration Example

See `examples/advanced_demo.go` for a comprehensive implementation that combines all three features:

- Streaming responses with real-time UI updates
- Knowledge base integration for context-aware answers
- Persistent conversation history across sessions
- Full REST API with web interface

```bash
# Run the advanced demo
go run examples/advanced_demo.go
# Visit http://localhost:8080 for web interface
```

## Example .env

```env
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=...
GEMINI_API_KEY=...
XAI_API_KEY=...
```

## API Endpoint Contract

- Endpoint: `/api/chat`
- Method: `POST`
- Request:

  ```json
  { "message": "Hello" }
  ```

- Response:

  ```json
  { "reply": "Hi! How can I help you?" }
  ```

## Testing

```bash
go test ./...
```

## Running Tests with Coverage

```bash
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Static Analysis

```bash
golangci-lint run
```

## Security

- Input validation and abuse prevention built-in.
- Rate limiting and content filtering available via config.
- Context-aware operations with timeout support.

## Rate Limiting & Abuse Prevention

You can implement rate limiting and abuse prevention using middleware:

```go
// Example with Gin
import "github.com/gin-contrib/ratelimit"

r.Use(ratelimit.RateLimiter(store, &ratelimit.Options{
    ErrorHandler: ratelimit.DefaultErrorHandler,
    KeyFunc:      ratelimit.KeyByIP,
}))
```

## JavaScript Framework Components

You can use the provided chat popup as a plain HTML/JS snippet, or integrate a modern component for Vue, React, or Angular:

- **Vue**: `web/components/GoChatbotVue.vue`
- **React**: `web/components/GoChatbotReact.jsx`
- **Angular**: `web/components/GoChatbotAngular.html` (with TypeScript logic)

### JS Component Usage

1. **Copy the component** into your app's source tree.
2. **Import and register** it in your app (see your framework's docs).
3. **Customize the backend endpoint** (`/api/chat`) as needed.
4. **Build your assets** (e.g. with Vite, Webpack, or your framework's CLI).

#### Example (Vue)

```js
// main.js
import GoChatbotVue from './components/GoChatbotVue.vue';
app.component('GoChatbotVue', GoChatbotVue);
```

#### Example (React)

```js
import GoChatbotReact from './components/GoChatbotReact.jsx';
<GoChatbotReact />
```

#### Example (Angular)

- Copy the HTML, TypeScript, and CSS into your Angular component files.
- Register and use `<go-chatbot-angular></go-chatbot-angular>` in your template.

### JS Component Customization

- You can style or extend the components as needed.
- You may add your own framework component in a similar way.
- The backend endpoint is framework-agnostic; you can point it to any Go route/handler.

## TypeScript Example (React)

A modern TypeScript React chatbot component is provided as an example in `web/components/GoChatbotTs.tsx`.

How to use:

1. Copy `GoChatbotTs.tsx` into your React app's components folder.
2. Import and use it in your app:
   ```typescript
   import GoChatbotTs from './components/GoChatbotTs';
   // ...
   <GoChatbotTs />
   ```
3. Make sure your backend endpoint `/api/chat` is set up as described above.
4. Style as needed (the component uses inline styles for demo purposes, but you can use your own CSS/SCSS).

> This component is a minimal, framework-agnostic starting point. You are encouraged to extend or restyle it to fit your app.

## Backend Integration Example

To handle chat requests, add a route/handler in your backend that receives POST requests at `/api/chat` and returns a JSON response.

Example for Gin:

```go
import (
    "net/http"
    "github.com/gin-gonic/gin"
    "go.rumenx.com/chatbot"
    "go.rumenx.com/chatbot/config"
    "go.rumenx.com/chatbot/models"
)

func ChatHandler(cfg *config.Config) gin.HandlerFunc {
    model, _ := models.NewFromConfig(cfg)
    chatbot := gochatbot.New(cfg, gochatbot.WithModel(model))

    return func(c *gin.Context) {
        var req struct {
            Message string `json:"message"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        reply, err := chatbot.Ask(c.Request.Context(), req.Message)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"reply": reply})
    }
}
```

## Frontend Styles (SCSS)

The chat popup styles are written in modern SCSS for maintainability. You can find the source in `web/css/chatbot.scss`.

To compile SCSS to CSS:

1. **Install Sass** (if you haven't already):
   ```sh
   npm install --save-dev sass
   ```

2. **Add this script** to your `package.json`:
   ```json
   "scripts": {
     "scss": "sass web/css/chatbot.scss web/css/chatbot.css --no-source-map"
   }
   ```

3. **Compile SCSS**:
   ```sh
   npm run scss
   ```

Or use the Sass CLI directly:

```sh
sass web/css/chatbot.scss web/css/chatbot.css --no-source-map
```

To watch for changes automatically:

```sh
sass --watch web/css/chatbot.scss:web/css/chatbot.css --no-source-map
```

> Only commit the compiled `chatbot.css` for production/deployment. For more options, see [Sass documentation](https://sass-lang.com/documentation/cli/dart-sass).

## Customizing Frontend Styles & Views (Best Practice)

The provided CSS/SCSS and view files (`web/css/`, `web/views/`) are **optional and basic**. They are meant as a starting point for your own implementation. **You are encouraged to build your own UI and styles on top of or instead of these defaults.**

### How to Safely Override Views and Styles

- **Copy the provided view files** (e.g. `web/views/popup.html`) into your own application's views directory. Update your app to use your custom view.
- **Copy and modify the SCSS/CSS** (`web/css/chatbot.scss` or `chatbot.css`) into your own asset pipeline. Import, extend, or replace as needed.
- **For Go applications:** Embed the views using Go's `embed` package or copy them to your static files directory.

### Principle

- **Treat the package's frontend as a reference implementation.**
- **Override or extend in your own codebase for all customizations.**
- **Never edit vendor files directly.**

## Configuration Best Practices

- **Always use environment variables for secrets and API keys.**
- **Use Go's `embed` package for embedding static assets.**
- **Keep your customizations outside the package directory for upgrade safety.**

**Example usage:**

```go
package main

import (
    "os"
    "go.rumenx.com/chatbot/config"
)

func main() {
    cfg := &config.Config{
        Model: "openai",
        OpenAI: config.OpenAIConfig{
            APIKey: os.Getenv("OPENAI_API_KEY"),
            Model:  "gpt-4o",
        },
        Prompt:   "You are a helpful, friendly chatbot.",
        Language: "en",
        Tone:     "neutral",
    }

    // Use configuration...
}
```

> This approach ensures your configuration is safe and easy to manage in any deployment environment.

## Best Practices

- Use environment variables for secrets
- Leverage Go's context package for cancellation and timeouts
- Use the provided JS/TS components as a starting point, not as production code
- Keep your customizations outside the package directory for upgrade safety
- Follow Go coding standards (gofmt, goimports, golint)
- Use structured logging with slog
- Handle errors gracefully and provide meaningful error messages

## What's Included

| Feature                | Location/Example                                 | Optional/Required |
|------------------------|--------------------------------------------------|-------------------|
| Go Core Package        | `./`                                             | Required          |
| **Framework Adapters** |                                                  |                   |
| Gin Adapter            | `adapters/gin/`                                  | Optional          |
| Echo Adapter           | `adapters/echo/`                                 | Optional          |
| Fiber Adapter          | `adapters/fiber/`                                | Optional          |
| Chi Adapter            | `adapters/chi/`                                  | Optional          |
| **Core Packages**      |                                                  |                   |
| Config Package         | `config/`                                        | Required          |
| Models Package         | `models/`                                        | Required          |
| Middleware Package     | `middleware/`                                    | Required          |
| **Advanced Features**  |                                                  |                   |
| Streaming Package      | `streaming/`                                     | Optional          |
| Embeddings Package     | `embeddings/`                                    | Optional          |
| Database Package       | `database/`                                      | Optional          |
| **Frontend Assets**    |                                                  |                   |
| HTML Views             | `web/views/`                                     | Optional          |
| CSS/SCSS               | `web/css/chatbot.scss`/`chatbot.css`             | Optional          |
| JS/TS Components       | `web/components/`                                | Optional          |
| **Examples & Config**  |                                                  |                   |
| Advanced Demo          | `examples/advanced_demo.go`                      | Optional          |
| Basic Examples         | `examples/basic/`, `examples/providers/`         | Optional          |
| Example .env           | `.env.example`                                   | Optional          |
| Tests                  | `*_test.go`                                      | Optional          |

## Chat Message Filtering Middleware

The package includes a configurable chat message filtering middleware to help ensure safe, appropriate, and guideline-aligned AI responses. This middleware:

- Filters and optionally rephrases user-submitted messages before they reach the AI model.
- Appends hidden system instructions (not visible in chat history) to the AI context, enforcing safety and communication guidelines.
- All filtering rules (profanities, aggression patterns, link regex) and system instructions are fully configurable.

Example configuration:

```go
cfg := &config.Config{
    MessageFiltering: config.MessageFilteringConfig{
        Instructions: []string{
            "Avoid sharing external links.",
            "Refrain from quoting controversial sources.",
            "Use appropriate language.",
            "Reject harmful or dangerous requests.",
            "De-escalate potential conflicts and calm aggressive or rude users.",
        },
        Profanities:        []string{"badword1", "badword2"},
        AggressionPatterns: []string{"hate", "kill", "stupid", "idiot"},
        LinkPattern:        `https?://[\w\.-]+`,
    },
}
```

Example usage:

```go
import "go.rumenx.com/chatbot/middleware"

filter := middleware.NewChatMessageFilter(cfg.MessageFiltering)

// Before sending to the AI model:
filtered, err := filter.Handle(ctx, userMessage)
if err != nil {
    // Handle error
}

reply, err := chatbot.Ask(ctx, filtered.Message, filtered.Context)
```

Purpose:

- Promotes safe, respectful, and effective communication.
- Prevents misuse, abuse, and unsafe outputs.
- All rules are transparent and configurable—no hidden censorship or manipulation.

## Questions & Support

For questions, issues, or feature requests, please use the [GitHub Issues](https://github.com/rumendamyanov/go-chatbot/issues) page.

For security vulnerabilities, please see our [Security Policy](SECURITY.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute.

## Code of Conduct

This project adheres to a Code of Conduct to ensure a welcoming and inclusive environment for all contributors. See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for details.

## Support the Project

If you find this project helpful, consider supporting its development. See [FUNDING.md](FUNDING.md) for ways to contribute.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes and releases.

## License

MIT. See [LICENSE.md](LICENSE.md).
