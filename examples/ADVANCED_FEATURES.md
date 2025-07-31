# Advanced Features Integration Guide

This directory contains examples demonstrating the advanced features of the Go Chatbot framework, including streaming responses, embeddings, and database integration.

## 🚀 Advanced Features Overview

### 1. Streaming Responses
- **Technology**: Server-Sent Events (SSE)
- **Features**: Real-time streaming of AI responses
- **Use Case**: Interactive chat experiences with immediate feedback
- **Implementation**: Built-in streaming handler with automatic chunking

### 2. Text Embeddings & Vector Search
- **Provider**: OpenAI text-embedding-3-small/large
- **Features**: Semantic similarity search, enhanced context retrieval
- **Use Case**: Knowledge base integration, context-aware responses
- **Implementation**: Vector store with cosine similarity search

### 3. Database Integration
- **Database**: SQLite (easily configurable for PostgreSQL)
- **Features**: Conversation persistence, message history, user sessions
- **Use Case**: Multi-session conversations, analytics, data retention
- **Implementation**: SQL-based conversation store with full CRUD operations

## 📁 Files in this Directory

### `advanced_demo.go`
Complete integration example showcasing all three advanced features working together:

- **HTTP Server**: RESTful API with multiple endpoints
- **Streaming Chat**: `/chat` endpoint with streaming support
- **Conversation Management**: Full CRUD operations for conversations
- **Knowledge Base**: Vector store integration for enhanced context
- **Database Persistence**: SQLite backend for all data

## 🛠️ Quick Start

### Prerequisites
```bash
# Install dependencies
go mod tidy

# Set your OpenAI API key (optional for demo)
export OPENAI_API_KEY="your-api-key-here"
```

### Run the Advanced Demo
```bash
# Method 1: Using the provided script
./run_demo.sh

# Method 2: Direct execution
go run examples/advanced_demo.go

# Method 3: Build and run
go build -o bin/advanced_demo examples/advanced_demo.go
./bin/advanced_demo
```

The server will start on `http://localhost:8080` with an interactive documentation page.

## 🔌 API Endpoints

### Chat Operations
```bash
# Send a regular chat message
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello!", "stream": false}'

# Send a streaming chat message with embeddings
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Explain AI", "stream": true, "use_embeddings": true}'
```

### Conversation Management
```bash
# List all conversations
curl http://localhost:8080/conversations

# Create a new conversation
curl -X POST http://localhost:8080/conversations \
  -H "Content-Type: application/json" \
  -d '{"title": "My Chat Session"}'

# Get messages from a conversation
curl http://localhost:8080/conversations/conv_1234567890/messages
```

### Knowledge Base
```bash
# Add knowledge to the vector store
curl -X POST http://localhost:8080/knowledge \
  -H "Content-Type: application/json" \
  -d '{"content": "Go is a programming language", "id": "go_info_1"}'

# Check server status
curl http://localhost:8080/status
```

## 📊 Architecture Breakdown

### Streaming Implementation
```go
// Built-in streaming with SSE
func (s *AdvancedChatbotServer) handleStreamingResponse(w http.ResponseWriter, r *http.Request, conversationID, prompt string) {
    ctx := r.Context()
    err := s.chatbot.AskStream(ctx, w, prompt)
    // Automatic SSE headers and chunk processing
}
```

### Database Schema
```sql
-- Conversations table
CREATE TABLE conversations (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Messages table
CREATE TABLE messages (
    id VARCHAR(255) PRIMARY KEY,
    conversation_id VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Embeddings Integration
```go
// Generate embeddings for enhanced context
queryEmbedding, err := embeddingProvider.EmbedSingle(ctx, query)
results, err := vectorStore.Search(ctx, query, 3) // Top 3 similar results

// Enhance chat context with relevant knowledge
for _, result := range results {
    enhancedContext = append(enhancedContext, result.Content)
}
```

## 🎯 Key Features Demonstrated

### 1. Real-time Streaming
- Server-Sent Events for live response streaming
- Automatic chunking and progress indicators
- Context cancellation and error handling
- Compatible with web browsers and curl

### 2. Intelligent Context Enhancement
- Semantic search using OpenAI embeddings
- Vector similarity scoring with cosine distance
- Automatic context injection into prompts
- Knowledge base integration

### 3. Persistent Conversations
- Multi-session conversation tracking
- Message history with timestamps
- User session management
- Full conversation CRUD operations

### 4. Production-Ready Architecture
- Modular design with clean interfaces
- Comprehensive error handling
- Configurable database backends
- RESTful API design patterns

## 🔧 Configuration Options

### Database Configuration
```go
// SQLite (default)
db, err := sql.Open("sqlite3", "./chatbot.db")

// PostgreSQL
db, err := sql.Open("postgres", "postgres://user:password@localhost/chatbot?sslmode=disable")

conversationStore := database.NewSQLConversationStore(db, "postgres")
```

### Embedding Models
```go
// Small model (faster, less accurate)
embeddingProvider := embeddings.NewOpenAIEmbeddingProvider(config, "text-embedding-3-small")

// Large model (slower, more accurate)
embeddingProvider := embeddings.NewOpenAIEmbeddingProvider(config, "text-embedding-3-large")
```

### Chatbot Configuration
```go
chatbotConfig := &config.Config{
    Model: "gpt-4",           // AI model to use
    Timeout: 30 * time.Second, // Request timeout
    MaxTokens: 1000,          // Response length limit
    Temperature: 0.7,         // Response creativity
}
```

## 🧪 Testing the Features

### Test Streaming
```bash
# Terminal 1: Start the server
go run examples/advanced_demo.go

# Terminal 2: Test streaming response
curl -N -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Count from 1 to 10", "stream": true}'
```

### Test Database Persistence
```bash
# Create conversation
CONV_ID=$(curl -s -X POST http://localhost:8080/conversations \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Chat"}' | jq -r '.id')

# Send message in conversation
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d "{\"conversation_id\": \"$CONV_ID\", \"message\": \"Hello!\"}"

# Check conversation history
curl http://localhost:8080/conversations/$CONV_ID/messages
```

### Test Embeddings Integration
```bash
# Add knowledge to vector store
curl -X POST http://localhost:8080/knowledge \
  -H "Content-Type: application/json" \
  -d '{"content": "Go was created at Google in 2007", "id": "go_history"}'

# Query with embeddings enabled
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Tell me about Go", "use_embeddings": true}'
```

## 🚀 Next Steps

1. **Extend Vector Store**: Implement persistent vector storage with databases like Pinecone or Qdrant
2. **Add Authentication**: Integrate user authentication and authorization
3. **Implement Caching**: Add Redis caching for embeddings and frequent queries
4. **Add Monitoring**: Integrate metrics and logging for production deployment
5. **Scale Horizontally**: Implement load balancing and distributed storage

## 📚 Related Documentation

- [Main README](../README.md) - Project overview and basic usage
- [Streaming Package](../streaming/) - Detailed streaming implementation
- [Embeddings Package](../embeddings/) - Vector search capabilities
- [Database Package](../database/) - Persistence layer documentation
- [Configuration Guide](../config/) - Configuration options and examples

## 💡 Tips for Production

1. **Use Environment Variables**: Never hardcode API keys in source code
2. **Configure Timeouts**: Set appropriate timeouts for your use case
3. **Monitor Database Performance**: Index frequently queried columns
4. **Implement Rate Limiting**: Protect against abuse and manage costs
5. **Use Connection Pooling**: Configure database connection pools for better performance
6. **Enable CORS**: Configure CORS headers for web browser compatibility
7. **Add Health Checks**: Implement comprehensive health monitoring

## 🐛 Troubleshooting

### Common Issues
1. **"Stream not supported"**: Ensure your HTTP client supports SSE
2. **Database locked**: Use connection pooling and proper transaction handling
3. **Embedding API errors**: Check API key and rate limits
4. **Memory usage**: Monitor vector store size and implement cleanup

### Debug Mode
```bash
# Enable verbose logging
export CHATBOT_DEBUG=true
go run examples/advanced_demo.go
```

This example demonstrates enterprise-level features that can be deployed in production environments with proper configuration and monitoring.
