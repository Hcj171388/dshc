package tools

import (
	"encoding/json"
	"strings"
)

// HeadTool represents a head command tool
type HeadTool struct{}

func (t *HeadTool) Name() string { return "head" }

func (t *HeadTool) Description() string {
	return "Show first lines of text"
}

func (t *HeadTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"input": map[string]interface{}{
				"type":        "string",
				"description": "Input text",
			},
			"lines": map[string]interface{}{
				"type":        "integer",
				"description": "Number of lines",
			},
		},
		"required": []string{"input"},
	}
}

func (t *HeadTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Input string `json:"input"`
		Lines int    `json:"lines"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	if params.Lines <= 0 {
		params.Lines = 10
	}
	
	lines := strings.Split(params.Input, "\n")
	if len(lines) > params.Lines {
		lines = lines[:params.Lines]
	}
	return strings.Join(lines, "\n"), nil
}
