# Basic Example

This is a basic example demonstrating how to use the go-chatbot package.

## Running the Example

1. Navigate to this directory:
   ```bash
   cd examples/basic
   ```

2. Run the example:
   ```bash
   go run main.go
   ```

3. The server will start on port 8080. You can:
   - Send chat requests to `http://localhost:8080/api/chat`
   - Check health at `http://localhost:8080/health`
   - Access the web interface at `http://localhost:8080`

## Testing the API

You can test the chat API using curl:

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, how are you?"}'
```

Expected response:
```json
{
  "reply": "Hello! Nice to meet you. How can I help you today?"
}
```

## Configuration

This example uses the default configuration with the free model. To use other AI providers, set the appropriate environment variables:

```bash
export CHATBOT_MODEL=openai
export OPENAI_API_KEY=your-api-key-here
go run main.go
```
