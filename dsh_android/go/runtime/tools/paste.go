package tools

import (
	"encoding/json"
	"fmt"
)

// PasteTool represents a paste tool
type PasteTool struct{}

func (t *PasteTool) Name() string { return "paste" }

func (t *PasteTool) Description() string {
	return "Paste text from clipboard"
}

func (t *PasteTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"content": map[string]interface{}{
				"type":        "string",
				"description": "Content to paste",
			},
		},
		"required": []string{"content"},
	}
}

func (t *PasteTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	return fmt.Sprintf("Pasted: %s", params.Content), nil
}
