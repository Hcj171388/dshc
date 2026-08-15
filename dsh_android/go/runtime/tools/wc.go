package tools

import (
	"encoding/json"
	"strings"
	"fmt"
)

// WcTool represents a word count tool
type WcTool struct{}

func (t *WcTool) Name() string { return "wc" }

func (t *WcTool) Description() string {
	return "Count words, lines, characters"
}

func (t *WcTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"input": map[string]interface{}{
				"type":        "string",
				"description": "Input text",
			},
		},
		"required": []string{"input"},
	}
}

func (t *WcTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Input string `json:"input"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	lines := len(strings.Split(params.Input, "\n"))
	words := len(strings.Fields(params.Input))
	chars := len([]rune(params.Input))
	return fmt.Sprintf("%d lines, %d words, %d characters", lines, words, chars), nil
}
