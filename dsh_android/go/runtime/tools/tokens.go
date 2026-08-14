package tools

import (
	"encoding/json"
	"fmt"
)

// TokenTool represents a token counter tool
type TokenTool struct{}

func (t *TokenTool) Name() string { return "token_count" }

func (t *TokenTool) Description() string {
	return "Count tokens in text"
}

func (t *TokenTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"text": map[string]interface{}{
				"type":        "string",
				"description": "Text to count tokens",
			},
		},
		"required": []string{"text"},
	}
}

func (t *TokenTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	// Simple word-based token counting
	tokens := len([]rune(params.Text))
	return fmt.Sprintf("%d tokens", tokens), nil
}
