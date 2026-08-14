package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Message represents a chat message
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// ToolCall represents a tool call from the model
type ToolCall struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Function `json:"function"`
}

// Function represents a function call
type Function struct {
	Name      string         `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ToolResult represents the result of a tool call
type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
}

// Response represents the LLM response
type Response struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a choice in the response
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Config represents LLM configuration
type Config struct {
	APIKey       string
	BaseURL      string
	Model        string
	MaxTokens    int
	Timeout      time.Duration
}

// Client is the LLM client
type Client struct {
	config Config
	http   *http.Client
}

// NewClient creates a new LLM client
func NewClient(cfg Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.deepseek.com"
	}
	if cfg.Model == "" {
		cfg.Model = "deepseek-chat"
	}
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = 4096
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 60 * time.Second
	}
	return &Client{
		config: cfg,
		http:   &http.Client{Timeout: cfg.Timeout},
	}
}

// CallOptions represents options for a call
type CallOptions struct {
	Messages   []Message
	Tools      []Tool
	Stream     bool
	Temperature float64
	MaxTokens   int
}

// Tool represents a tool available to the model
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Call sends a request to the LLM and returns the response
func (c *Client) Call(ctx context.Context, opts CallOptions) (*Response, error) {
	reqBody := map[string]interface{}{
		"model":       c.config.Model,
		"messages":    opts.Messages,
		"max_tokens":  c.config.MaxTokens,
		"temperature": opts.Temperature,
	}

	if len(opts.Tools) > 0 {
		reqBody["tools"] = opts.Tools
		reqBody["tool_choice"] = "auto"
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.BaseURL+"/v1/chat/completions", bytes.NewReader(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// CallWithTools is a convenience method that handles tool calls
func (c *Client) CallWithTools(ctx context.Context, messages []Message, tools []Tool, executeTool func(string, json.RawMessage) (string, error)) ([]Message, error) {
	opts := CallOptions{
		Messages:  messages,
		Tools:     tools,
		Temperature: 0.7,
	}

	resp, err := c.Call(ctx, opts)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return messages, nil
	}

	assistantMsg := resp.Choices[0].Message
	messages = append(messages, assistantMsg)

	// Handle tool calls
	for _, toolCall := range assistantMsg.Content.(map[string]interface{})["tool_calls"].([]interface{}) {
		tc := toolCall.(map[string]interface{})
		funcCall := tc["function"].(map[string]interface{})
		name := funcCall["name"].(string)
		args := json.RawMessage(funcCall["arguments"].(string))

		result, err := executeTool(name, args)
		if err != nil {
			result = fmt.Sprintf("Error: %v", err)
		}

		messages = append(messages, Message{
			Role: "tool",
			Content: map[string]interface{}{
				"tool_call_id": tc["id"],
				"content":      result,
			},
		})
	}

	return messages, nil
}
