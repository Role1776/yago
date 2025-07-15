package yago

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// requestPayload is the structure for the request body sent to the YandexGPT API.
type requestPayload struct {
	ModelURI          string            `json:"modelUri"`
	CompletionOptions CompletionOptions `json:"completionOptions"`
	Messages          []Message         `json:"messages"`
}

// Response is the top-level structure of the API response.
type Response struct {
	Result Result `json:"result"`
}

// Result contains the alternatives and usage statistics.
type Result struct {
	Alternatives []Alternative `json:"alternatives"`
	Usage        Usage         `json:"usage"`
	ModelVersion string        `json:"modelVersion"`
}

// Alternative represents a single possible completion.
type Alternative struct {
	Message Message `json:"message"`
	Status  string  `json:"status"`
}

// Usage provides token usage statistics for the request.
type Usage struct {
	InputTextTokens  string `json:"inputTextTokens"`
	CompletionTokens string `json:"completionTokens"`
	TotalTokens      string `json:"totalTokens"`
}

// Generate sends the provided messages to the model and returns the response.
func (g *GenerativeModel) Generate(ctx context.Context, messages []Message) (*Response, error) {
	finalMessages := make([]Message, 0, len(messages)+1)
	if g.SystemInstruction != "" {
		finalMessages = append(finalMessages, Message{Role: RoleSystem, Text: g.SystemInstruction})
	}
	finalMessages = append(finalMessages, messages...)
	if len(finalMessages) == 0 {
		return nil, fmt.Errorf("empty message list")
	}

	payload := requestPayload{
		ModelURI:          g.modelURI,
		CompletionOptions: g.CompletionOptions,
		Messages:          finalMessages,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", g.c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Api-Key "+g.c.apiKey)

	resp, err := g.c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
