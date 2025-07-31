package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// TestRequest represents a test API request
type TestRequest struct {
	Method   string            `json:"method"`
	URL      string            `json:"url"`
	Headers  map[string]string `json:"headers"`
	Body     interface{}       `json:"body"`
	Expected int               `json:"expected_status"`
}

// ChatRequest matches the server's ChatRequest struct
type ChatRequest struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
	UseEmbeddings  bool   `json:"use_embeddings"`
	Stream         bool   `json:"stream"`
}

func main() {
	baseURL := "http://localhost:8080"

	fmt.Println("🧪 Advanced Chatbot API Test Suite")
	fmt.Println("===================================")
	fmt.Println("")

	// Wait for server to be ready
	fmt.Print("🔄 Checking if server is running... ")
	if !waitForServer(baseURL, 5*time.Second) {
		fmt.Println("❌")
		fmt.Println("Error: Server is not running on", baseURL)
		fmt.Println("Please start the server first: go run examples/advanced_demo.go")
		os.Exit(1)
	}
	fmt.Println("✅")

	// Test cases
	tests := []TestRequest{
		{
			Method:   "GET",
			URL:      baseURL + "/status",
			Expected: 200,
		},
		{
			Method:   "GET",
			URL:      baseURL + "/conversations",
			Expected: 200,
		},
		{
			Method:   "POST",
			URL:      baseURL + "/conversations",
			Headers:  map[string]string{"Content-Type": "application/json"},
			Body:     map[string]string{"title": "Test Conversation"},
			Expected: 200,
		},
		{
			Method:  "POST",
			URL:     baseURL + "/knowledge",
			Headers: map[string]string{"Content-Type": "application/json"},
			Body: map[string]string{
				"content": "Go is a programming language developed at Google",
				"id":      "go_facts_test",
			},
			Expected: 201,
		},
		{
			Method:  "POST",
			URL:     baseURL + "/chat",
			Headers: map[string]string{"Content-Type": "application/json"},
			Body: ChatRequest{
				Message:       "Hello, this is a test message",
				Stream:        false,
				UseEmbeddings: false,
			},
			Expected: 200,
		},
	}

	// Run tests
	passed := 0
	total := len(tests)

	for i, test := range tests {
		fmt.Printf("🧪 Test %d/%d: %s %s\n", i+1, total, test.Method, test.URL)

		if runTest(test) {
			fmt.Println("   ✅ PASS")
			passed++
		} else {
			fmt.Println("   ❌ FAIL")
		}
		fmt.Println("")
	}

	// Summary
	fmt.Println("📊 Test Results")
	fmt.Println("===============")
	fmt.Printf("✅ Passed: %d/%d\n", passed, total)
	if passed == total {
		fmt.Println("🎉 All tests passed! The advanced features are working correctly.")
	} else {
		fmt.Printf("⚠️  %d tests failed. Check the server logs for details.\n", total-passed)
	}
}

func waitForServer(baseURL string, timeout time.Duration) bool {
	client := &http.Client{Timeout: 1 * time.Second}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/status")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return true
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func runTest(test TestRequest) bool {
	client := &http.Client{Timeout: 10 * time.Second}

	var body io.Reader
	if test.Body != nil {
		jsonBody, err := json.Marshal(test.Body)
		if err != nil {
			fmt.Printf("   Error marshaling request body: %v\n", err)
			return false
		}
		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(test.Method, test.URL, body)
	if err != nil {
		fmt.Printf("   Error creating request: %v\n", err)
		return false
	}

	// Add headers
	for key, value := range test.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("   Error making request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Read response body for debugging
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != test.Expected {
		fmt.Printf("   Expected status %d, got %d\n", test.Expected, resp.StatusCode)
		if len(respBody) > 0 && len(respBody) < 200 {
			fmt.Printf("   Response: %s\n", string(respBody))
		}
		return false
	}

	fmt.Printf("   Status: %d, Response length: %d bytes\n", resp.StatusCode, len(respBody))
	return true
}
