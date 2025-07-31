// Package streaming provides streaming response functionality for the go-chatbot package.
package streaming

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// StreamResponse represents a streaming response chunk.
type StreamResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
	Error   string `json:"error,omitempty"`
}

// StreamHandler handles Server-Sent Events (SSE) streaming.
type StreamHandler struct {
	writer  http.ResponseWriter
	flusher http.Flusher
	done    chan bool
}

// NewStreamHandler creates a new streaming handler.
func NewStreamHandler(w http.ResponseWriter) (*StreamHandler, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming unsupported: ResponseWriter does not implement http.Flusher")
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	return &StreamHandler{
		writer:  w,
		flusher: flusher,
		done:    make(chan bool),
	}, nil
}

// WriteChunk writes a streaming chunk to the response.
func (s *StreamHandler) WriteChunk(chunk StreamResponse) error {
	data, err := json.Marshal(chunk)
	if err != nil {
		return fmt.Errorf("failed to marshal chunk: %w", err)
	}

	// Write SSE format
	_, err = fmt.Fprintf(s.writer, "data: %s\n\n", data)
	if err != nil {
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	s.flusher.Flush()
	return nil
}

// WriteError writes an error chunk to the response.
func (s *StreamHandler) WriteError(id, errorMsg string) error {
	return s.WriteChunk(StreamResponse{
		ID:    id,
		Error: errorMsg,
		Done:  true,
	})
}

// WriteDone writes a completion chunk to the response.
func (s *StreamHandler) WriteDone(id string) error {
	return s.WriteChunk(StreamResponse{
		ID:   id,
		Done: true,
	})
}

// Close closes the stream.
func (s *StreamHandler) Close() {
	close(s.done)
}

// StreamProcessor processes streaming data from various sources.
type StreamProcessor struct {
	requestID string
	handler   *StreamHandler
}

// NewStreamProcessor creates a new stream processor.
func NewStreamProcessor(requestID string, handler *StreamHandler) *StreamProcessor {
	return &StreamProcessor{
		requestID: requestID,
		handler:   handler,
	}
}

// ProcessChannel processes a channel of strings and streams them.
func (sp *StreamProcessor) ProcessChannel(ctx context.Context, ch <-chan string) error {
	defer sp.handler.WriteDone(sp.requestID)

	for {
		select {
		case <-ctx.Done():
			return sp.handler.WriteError(sp.requestID, "Request cancelled")
		case content, ok := <-ch:
			if !ok {
				// Channel closed, we're done
				return nil
			}

			err := sp.handler.WriteChunk(StreamResponse{
				ID:      sp.requestID,
				Content: content,
				Done:    false,
			})
			if err != nil {
				return fmt.Errorf("failed to write chunk: %w", err)
			}
		}
	}
}

// ProcessOpenAIStream processes OpenAI streaming response format.
func (sp *StreamProcessor) ProcessOpenAIStream(ctx context.Context, response *http.Response) error {
	defer sp.handler.WriteDone(sp.requestID)
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return sp.handler.WriteError(sp.requestID, "Request cancelled")
		default:
		}

		line := scanner.Text()

		// Skip empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, ":") {
			continue
		}

		// Parse SSE format
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// Check for end of stream
			if data == "[DONE]" {
				return nil
			}

			// Parse JSON data
			var chunk map[string]interface{}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue // Skip malformed chunks
			}

			// Extract content from OpenAI format
			content := extractOpenAIContent(chunk)
			if content != "" {
				err := sp.handler.WriteChunk(StreamResponse{
					ID:      sp.requestID,
					Content: content,
					Done:    false,
				})
				if err != nil {
					return fmt.Errorf("failed to write chunk: %w", err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return sp.handler.WriteError(sp.requestID, fmt.Sprintf("Stream reading error: %v", err))
	}

	return nil
}

// ProcessAnthropicStream processes Anthropic streaming response format.
func (sp *StreamProcessor) ProcessAnthropicStream(ctx context.Context, response *http.Response) error {
	defer sp.handler.WriteDone(sp.requestID)
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return sp.handler.WriteError(sp.requestID, "Request cancelled")
		default:
		}

		line := scanner.Text()

		// Parse event: lines and data: lines
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			var chunk map[string]interface{}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			// Extract content from Anthropic format
			content := extractAnthropicContent(chunk)
			if content != "" {
				err := sp.handler.WriteChunk(StreamResponse{
					ID:      sp.requestID,
					Content: content,
					Done:    false,
				})
				if err != nil {
					return fmt.Errorf("failed to write chunk: %w", err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return sp.handler.WriteError(sp.requestID, fmt.Sprintf("Stream reading error: %v", err))
	}

	return nil
}

// extractOpenAIContent extracts content from OpenAI streaming format.
func extractOpenAIContent(chunk map[string]interface{}) string {
	choices, ok := chunk["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return ""
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return ""
	}

	delta, ok := choice["delta"].(map[string]interface{})
	if !ok {
		return ""
	}

	content, ok := delta["content"].(string)
	if !ok {
		return ""
	}

	return content
}

// extractAnthropicContent extracts content from Anthropic streaming format.
func extractAnthropicContent(chunk map[string]interface{}) string {
	eventType, ok := chunk["type"].(string)
	if !ok {
		return ""
	}

	switch eventType {
	case "content_block_delta":
		if delta, ok := chunk["delta"].(map[string]interface{}); ok {
			if text, ok := delta["text"].(string); ok {
				return text
			}
		}
	case "content_block_start":
		if contentBlock, ok := chunk["content_block"].(map[string]interface{}); ok {
			if text, ok := contentBlock["text"].(string); ok {
				return text
			}
		}
	}

	return ""
}

// StreamingClient provides utilities for making streaming requests.
type StreamingClient struct {
	client  *http.Client
	timeout time.Duration
}

// NewStreamingClient creates a new streaming client.
func NewStreamingClient(timeout time.Duration) *StreamingClient {
	return &StreamingClient{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// MakeStreamingRequest makes a streaming HTTP request.
func (sc *StreamingClient) MakeStreamingRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Set streaming headers
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	// Create a context with timeout
	reqCtx, cancel := context.WithTimeout(ctx, sc.timeout)
	defer cancel()

	req = req.WithContext(reqCtx)

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("streaming request failed: %w", err)
	}

	// Check for successful status
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("streaming request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
