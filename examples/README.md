# Framework Examples

This directory contains practical examples demonstrating how to integrate the go-chatbot package with different Go web frameworks.

## Available Examples

### Gin Framework (`gin_example.go`)

Shows how to integrate with [Gin](https://github.com/gin-gonic/gin):
- Basic chatbot setup with free model
- Route configuration using `SetupRoutes()`
- Middleware integration for context injection
- Custom route accessing chatbot from context

**Run:**
```bash
go mod tidy
go run gin_example.go
```

### Echo Framework (`echo_example.go`)

Shows how to integrate with [Echo](https://github.com/labstack/echo):
- Basic chatbot setup with free model
- Route configuration using `SetupRoutes()`
- Middleware integration for context injection
- Custom route accessing chatbot from context

**Run:**
```bash
go mod tidy
go run echo_example.go
```

### Fiber Framework (`fiber_example.go`)

Shows how to integrate with [Fiber](https://github.com/gofiber/fiber):
- Basic chatbot setup with free model
- Route configuration using `SetupRoutes()`
- Middleware integration for context injection
- Custom route accessing chatbot from context

**Run:**
```bash
go mod tidy
go run fiber_example.go
```

### Chi Framework (`chi_example.go`)

Shows how to integrate with [Chi](https://github.com/go-chi/chi):
- Basic chatbot setup with free model
- Route configuration using `SetupRoutes()`
- Middleware integration for context injection
- Custom route accessing chatbot from context

**Run:**
```bash
go mod tidy
go run chi_example.go
```

## Common Endpoints

All examples expose the same endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/chat/` | POST | Send chat messages |
| `/chat/health` | GET | Health check |
| `/status` | GET | Custom status with chatbot info |

## Testing the Examples

1. **Start any example:**
   ```bash
   cd examples
   go run gin_example.go  # or echo_example.go, fiber_example.go, chi_example.go
   ```

2. **Test chat endpoint:**
   ```bash
   curl -X POST http://localhost:8080/chat/ \
     -H "Content-Type: application/json" \
     -d '{"message": "Hello!"}'
   ```

3. **Test health endpoint:**
   ```bash
   curl http://localhost:8080/chat/health
   ```

4. **Test status endpoint:**
   ```bash
   curl http://localhost:8080/status
   ```

## Expected Responses

**Chat Response:**
```json
{
  "success": true,
  "response": "Hello! I'm a simple free chatbot. How can I help you today?"
}
```

**Health Response:**
```json
{
  "status": "healthy",
  "provider": "local",
  "model": "free-model",
  "timestamp": 1640995200
}
```

**Status Response:**
```json
{
  "status": "ok",
  "model": "free-model",
  "provider": "local"
}
```

## Configuration

All examples use the free model by default for easy testing. To use other providers:

1. **Set environment variables:**
   ```bash
   export OPENAI_API_KEY="your-api-key"
   ```

2. **Modify the configuration in the example:**
   ```go
   cfg := config.Default()
   cfg.Model = "openai"  // or "anthropic", "gemini", etc.
   ```

3. **Run the example as usual**

## Production Considerations

These examples are for demonstration purposes. For production use:

- Add proper error handling and logging
- Implement authentication and authorization
- Configure rate limiting and abuse prevention
- Use proper configuration management (not hardcoded values)
- Add monitoring and health checks
- Configure CORS and security headers as needed
