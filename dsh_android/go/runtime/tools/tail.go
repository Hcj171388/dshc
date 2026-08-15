package tools

import (
	"encoding/json"
	"strings"
)

// TailTool represents a tail command tool
type TailTool struct{}

func (t *TailTool) Name() string { return "tail" }

func (t *TailTool) Description() string {
	return "Show last lines of text"
}

func (t *TailTool) Schema() map[string]interface{} {
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

func (t *TailTool) Execute(input json.RawMessage) (string, error) {
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
		lines = lines[len(lines)-params.Lines:]
	}
	return strings.Join(lines, "\n"), nil
}
