package tools

import (
	"encoding/json"
)

// TruncateTool represents a truncate tool
type TruncateTool struct{}

func (t *TruncateTool) Name() string { return "truncate" }

func (t *TruncateTool) Description() string {
	return "Truncate text to specified length"
}

func (t *TruncateTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"input": map[string]interface{}{
				"type":        "string",
				"description": "Input text",
			},
			"length": map[string]interface{}{
				"type":        "integer",
				"description": "Max length",
			},
		},
		"required": []string{"input", "length"},
	}
}

func (t *TruncateTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Input  string `json:"input"`
		Length int    `json:"length"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	runes := []rune(params.Input)
	if len(runes) > params.Length {
		return string(runes[:params.Length]) + "...", nil
	}
	return params.Input, nil
}
