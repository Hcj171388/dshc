package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// UniqTool represents a unique values tool
type UniqTool struct{}

func (t *UniqTool) Name() string { return "uniq" }

func (t *UniqTool) Description() string {
	return "Remove duplicate lines"
}

func (t *UniqTool) Schema() map[string]interface{} {
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

func (t *UniqTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Input string `json:"input"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	lines := strings.Split(params.Input, "\n")
	seen := make(map[string]bool)
	var result []string
	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}
	return fmt.Sprintf("%s", strings.Join(result, "\n")), nil
}
