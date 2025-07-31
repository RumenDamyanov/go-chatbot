package streaming

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewStreamHandler(t *testing.T) {
	w := httptest.NewRecorder()

	handler, err := NewStreamHandler(w)
	if err != nil {
		t.Errorf("Failed to create stream handler: %v", err)
	}

	if handler == nil {
		t.Error("Expected non-nil stream handler")
	}

	// Check that headers were set
	headers := w.Header()
	if contentType := headers.Get("Content-Type"); contentType != "text/event-stream" {
		t.Errorf("Expected Content-Type 'text/event-stream', got '%s'", contentType)
	}
}

func TestNewStreamHandlerNonFlusher(t *testing.T) {
	// Use a ResponseWriter that doesn't implement Flusher
	w := &nonFlusherWriter{
		header: make(http.Header),
	}

	handler, err := NewStreamHandler(w)
	if err == nil {
		t.Error("Expected error when ResponseWriter doesn't implement Flusher")
	}

	if handler != nil {
		t.Error("Expected nil handler when ResponseWriter doesn't implement Flusher")
	}
}

func TestStreamHandler_WriteChunk(t *testing.T) {
	w := httptest.NewRecorder()
	handler, err := NewStreamHandler(w)
	if err != nil {
		t.Fatalf("Failed to create stream handler: %v", err)
	}

	chunk := StreamResponse{
		ID:      "test-1",
		Content: "Hello, world!",
		Done:    false,
	}

	err = handler.WriteChunk(chunk)
	if err != nil {
		t.Errorf("Failed to write chunk: %v", err)
	}

	response := w.Body.String()
	if !strings.Contains(response, "Hello, world!") {
		t.Errorf("Expected response to contain 'Hello, world!', got: %s", response)
	}

	// Check SSE format
	if !strings.HasPrefix(response, "data: ") {
		t.Error("Expected SSE format with 'data: ' prefix")
	}

	if !strings.HasSuffix(response, "\n\n") {
		t.Error("Expected SSE format with '\\n\\n' suffix")
	}
}

func TestStreamHandler_WriteError(t *testing.T) {
	w := httptest.NewRecorder()
	handler, err := NewStreamHandler(w)
	if err != nil {
		t.Fatalf("Failed to create stream handler: %v", err)
	}

	err = handler.WriteError("test-1", "Something went wrong")
	if err != nil {
		t.Errorf("Failed to write error: %v", err)
	}

	response := w.Body.String()
	if !strings.Contains(response, "Something went wrong") {
		t.Errorf("Expected response to contain error message, got: %s", response)
	}

	if !strings.Contains(response, `"done":true`) {
		t.Error("Expected error chunk to have done:true")
	}
}

func TestStreamHandler_WriteDone(t *testing.T) {
	w := httptest.NewRecorder()
	handler, err := NewStreamHandler(w)
	if err != nil {
		t.Fatalf("Failed to create stream handler: %v", err)
	}

	err = handler.WriteDone("test-1")
	if err != nil {
		t.Errorf("Failed to write done: %v", err)
	}

	response := w.Body.String()
	if !strings.Contains(response, `"done":true`) {
		t.Error("Expected done chunk to have done:true")
	}

	if !strings.Contains(response, `"id":"test-1"`) {
		t.Error("Expected done chunk to have correct ID")
	}
}

func TestStreamHandler_Close(t *testing.T) {
	w := httptest.NewRecorder()
	handler, err := NewStreamHandler(w)
	if err != nil {
		t.Fatalf("Failed to create stream handler: %v", err)
	}

	// Should not panic
	handler.Close()
}

func TestNewStreamProcessor(t *testing.T) {
	w := httptest.NewRecorder()
	handler, err := NewStreamHandler(w)
	if err != nil {
		t.Fatalf("Failed to create stream handler: %v", err)
	}

	processor := NewStreamProcessor("test-request", handler)
	if processor == nil {
		t.Error("Expected non-nil stream processor")
	}
}

func TestStreamProcessor_ProcessChannel(t *testing.T) {
	w := httptest.NewRecorder()
	handler, err := NewStreamHandler(w)
	if err != nil {
		t.Fatalf("Failed to create stream handler: %v", err)
	}

	processor := NewStreamProcessor("test-request", handler)

	// Create a channel and send some data
	ch := make(chan string, 3)
	ch <- "First"
	ch <- "Second"
	ch <- "Third"
	close(ch)

	ctx := context.Background()
	err = processor.ProcessChannel(ctx, ch)
	if err != nil {
		t.Errorf("Failed to process channel: %v", err)
	}

	response := w.Body.String()

	// Verify all messages were processed
	expectedMessages := []string{"First", "Second", "Third"}
	for _, msg := range expectedMessages {
		if !strings.Contains(response, msg) {
			t.Errorf("Expected response to contain '%s'", msg)
		}
	}

	// Should have a done marker at the end
	if !strings.Contains(response, `"done":true`) {
		t.Error("Expected response to have done marker")
	}
}

func TestStreamProcessor_ProcessChannelCancellation(t *testing.T) {
	w := httptest.NewRecorder()
	handler, err := NewStreamHandler(w)
	if err != nil {
		t.Fatalf("Failed to create stream handler: %v", err)
	}

	processor := NewStreamProcessor("test-request", handler)

	// Create a channel that blocks
	ch := make(chan string)

	// Create context and cancel immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = processor.ProcessChannel(ctx, ch)

	// Should handle cancellation gracefully
	if err != nil && err != context.Canceled {
		t.Errorf("Expected context.Canceled error or nil, got: %v", err)
	}
}

func TestStreamProcessor_ProcessOpenAIStream(t *testing.T) {
	w := httptest.NewRecorder()
	handler, err := NewStreamHandler(w)
	if err != nil {
		t.Fatalf("Failed to create stream handler: %v", err)
	}

	processor := NewStreamProcessor("test-request", handler)

	// Create a mock HTTP response with SSE format
	sseData := `data: {"choices":[{"delta":{"content":"Hello"}}]}

data: {"choices":[{"delta":{"content":" World"}}]}

data: [DONE]

`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte(sseData))
	}))
	defer server.Close()

	// Make a request to get the response
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	ctx := context.Background()
	err = processor.ProcessOpenAIStream(ctx, resp)
	if err != nil {
		t.Errorf("Failed to process OpenAI stream: %v", err)
	}

	response := w.Body.String()

	// Should contain the content from delta messages
	if !strings.Contains(response, "Hello") {
		t.Error("Expected response to contain 'Hello'")
	}
	if !strings.Contains(response, " World") {
		t.Error("Expected response to contain ' World'")
	}
}

func TestStreamProcessor_ProcessOpenAIStreamContextCancel(t *testing.T) {
	w := httptest.NewRecorder()
	handler, err := NewStreamHandler(w)
	if err != nil {
		t.Fatalf("Failed to create stream handler: %v", err)
	}

	processor := NewStreamProcessor("test-request", handler)

	// Create a slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte("data: slow response\n\n"))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Cancel context immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = processor.ProcessOpenAIStream(ctx, resp)

	// Should handle cancellation gracefully
	if err != nil && err != context.Canceled {
		t.Errorf("Expected context.Canceled error or nil, got: %v", err)
	}
}

func TestNewStreamingClient(t *testing.T) {
	client := NewStreamingClient(30 * time.Second)
	if client == nil {
		t.Error("Expected non-nil streaming client")
	}
}

func TestStreamingClient_MakeStreamingRequest(t *testing.T) {
	client := NewStreamingClient(30 * time.Second)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: test\n\n"))
	}))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	ctx := context.Background()
	resp, err := client.MakeStreamingRequest(ctx, req)
	if err != nil {
		t.Errorf("Failed to make streaming request: %v", err)
	}

	if resp == nil {
		t.Error("Expected non-nil response")
	} else {
		resp.Body.Close()
	}
}

func TestStreamingClient_MakeStreamingRequestTimeout(t *testing.T) {
	client := NewStreamingClient(1 * time.Millisecond) // Very short timeout

	// Create a slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Longer than timeout
		w.Write([]byte("data: test\n\n"))
	}))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	ctx := context.Background()
	resp, err := client.MakeStreamingRequest(ctx, req)

	// Should timeout
	if err == nil {
		t.Error("Expected timeout error")
		if resp != nil {
			resp.Body.Close()
		}
	}
}

// nonFlusherWriter is a test helper that doesn't implement http.Flusher
type nonFlusherWriter struct {
	header  http.Header
	written []byte
}

func (w *nonFlusherWriter) Header() http.Header {
	return w.header
}

func (w *nonFlusherWriter) Write(data []byte) (int, error) {
	w.written = append(w.written, data...)
	return len(data), nil
}

func (w *nonFlusherWriter) WriteHeader(statusCode int) {
	// No-op for testing
}

func TestExtractAnthropicContent(t *testing.T) {
	tests := []struct {
		name     string
		chunk    map[string]interface{}
		expected string
	}{
		{
			name: "content_block_delta with text",
			chunk: map[string]interface{}{
				"type": "content_block_delta",
				"delta": map[string]interface{}{
					"text": "Hello world",
				},
			},
			expected: "Hello world",
		},
		{
			name: "content_block_start with text",
			chunk: map[string]interface{}{
				"type": "content_block_start",
				"content_block": map[string]interface{}{
					"text": "Starting message",
				},
			},
			expected: "Starting message",
		},
		{
			name: "unknown event type",
			chunk: map[string]interface{}{
				"type": "unknown_event",
				"data": "some data",
			},
			expected: "",
		},
		{
			name: "content_block_delta without text",
			chunk: map[string]interface{}{
				"type": "content_block_delta",
				"delta": map[string]interface{}{
					"other": "field",
				},
			},
			expected: "",
		},
		{
			name: "content_block_delta with invalid delta",
			chunk: map[string]interface{}{
				"type":  "content_block_delta",
				"delta": "not a map",
			},
			expected: "",
		},
		{
			name: "content_block_start with invalid content_block",
			chunk: map[string]interface{}{
				"type":          "content_block_start",
				"content_block": "not a map",
			},
			expected: "",
		},
		{
			name: "missing type field",
			chunk: map[string]interface{}{
				"data": "some data",
			},
			expected: "",
		},
		{
			name: "type is not string",
			chunk: map[string]interface{}{
				"type": 123,
				"data": "some data",
			},
			expected: "",
		},
		{
			name:     "empty chunk",
			chunk:    map[string]interface{}{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractAnthropicContent(tt.chunk)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestStreamProcessor_ProcessAnthropicStream(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   string
		expectedChunks []string
		expectError    bool
	}{
		{
			name: "valid anthropic stream",
			responseBody: `data: {"type": "content_block_delta", "delta": {"text": "Hello"}}
data: {"type": "content_block_delta", "delta": {"text": " world"}}
data: {"type": "content_block_delta", "delta": {"text": "!"}}`,
			expectedChunks: []string{"Hello", " world", "!"},
			expectError:    false,
		},
		{
			name: "content_block_start event",
			responseBody: `data: {"type": "content_block_start", "content_block": {"text": "Starting"}}
data: {"type": "content_block_delta", "delta": {"text": " message"}}`,
			expectedChunks: []string{"Starting", " message"},
			expectError:    false,
		},
		{
			name: "mixed valid and invalid chunks",
			responseBody: `data: {"type": "content_block_delta", "delta": {"text": "Valid"}}
data: {"invalid": "json"
data: {"type": "content_block_delta", "delta": {"text": " chunk"}}`,
			expectedChunks: []string{"Valid", " chunk"},
			expectError:    false,
		},
		{
			name:           "empty stream",
			responseBody:   "",
			expectedChunks: []string{},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test HTTP response writer and handler
			w := httptest.NewRecorder()
			handler, err := NewStreamHandler(w)
			if err != nil {
				t.Fatalf("failed to create stream handler: %v", err)
			}

			processor := NewStreamProcessor("test-request", handler)

			// Create mock HTTP response
			response := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(tt.responseBody)),
			}

			ctx := context.Background()
			err = processor.ProcessAnthropicStream(ctx, response)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify the response was written to the recorder
			responseBody := w.Body.String()

			// Count data: lines to verify chunks were written
			dataLines := strings.Count(responseBody, "data:")
			expectedDataLines := len(tt.expectedChunks)

			if dataLines < expectedDataLines {
				t.Errorf("expected at least %d data lines, got %d", expectedDataLines, dataLines)
			}

			// Verify content appears in response
			for _, expectedContent := range tt.expectedChunks {
				if !strings.Contains(responseBody, expectedContent) {
					t.Errorf("expected content '%s' not found in response", expectedContent)
				}
			}
		})
	}
}

func TestStreamProcessor_ProcessAnthropicStream_ContextCancellation(t *testing.T) {
	w := httptest.NewRecorder()
	handler, err := NewStreamHandler(w)
	if err != nil {
		t.Fatalf("failed to create stream handler: %v", err)
	}

	processor := NewStreamProcessor("test-request", handler)

	// Create a simple response for testing context cancellation
	responseBody := `data: {"type": "content_block_delta", "delta": {"text": "chunk1"}}
data: {"type": "content_block_delta", "delta": {"text": "chunk2"}}
data: {"type": "content_block_delta", "delta": {"text": "chunk3"}}`

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}

	// Create a context that cancels immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err = processor.ProcessAnthropicStream(ctx, response)

	// The function may return an error from WriteError, or may complete without error
	// but should write the cancellation message to the response
	responseOutput := w.Body.String()
	if err != nil || strings.Contains(responseOutput, "Request cancelled") {
		// Test passes if either we get an error OR the cancellation message appears
		t.Logf("Context cancellation handled correctly - error: %v, response: %s", err, responseOutput)
	} else {
		t.Errorf("expected either error or cancellation message in response. Error: %v, Response: %s", err, responseOutput)
	}
}

func TestStreamingClient_MakeStreamingRequest_WithValidRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
		w.Write([]byte("data: test chunk\n"))
	}))
	defer server.Close()

	client := NewStreamingClient(30 * time.Second)

	// Create a proper HTTP request
	req, err := http.NewRequest("POST", server.URL, strings.NewReader(`{"message":"test"}`))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	ctx := context.Background()
	resp, err := client.MakeStreamingRequest(ctx, req)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp == nil {
		t.Error("expected non-nil response")
	}

	if resp != nil {
		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	}
}
