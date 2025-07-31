#!/bin/bash

# Advanced Go Chatbot Demo Runner
# This script demonstrates the advanced features integration

set -e

echo "🚀 Advanced Go Chatbot Demo"
echo "==============================================="
echo ""

# Check if OpenAI API key is set
if [ -z "$OPENAI_API_KEY" ]; then
    echo "⚠️  Warning: OPENAI_API_KEY environment variable not set"
    echo "   The demo will use a placeholder API key and may not work with real OpenAI API"
    echo ""
fi

# Build the demo
echo "🔨 Building advanced demo..."
go build -o bin/advanced_demo examples/advanced_demo.go

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed!"
    exit 1
fi

echo ""
echo "🎯 Features included in this demo:"
echo "   ✅ Streaming responses with Server-Sent Events"
echo "   ✅ OpenAI embeddings for enhanced context"
echo "   ✅ SQLite database for conversation persistence"
echo "   ✅ Vector similarity search capabilities"
echo "   ✅ RESTful API for all operations"
echo ""

echo "🌐 Starting server on http://localhost:8080..."
echo "   Press Ctrl+C to stop the server"
echo "   Open http://localhost:8080 in your browser for interactive documentation"
echo ""

# Create bin directory if it doesn't exist
mkdir -p bin

# Run the demo
./bin/advanced_demo
