package tools

import (
	"encoding/json"
	"fmt"
)

// ConvertTool represents a data conversion tool
type ConvertTool struct{}

func (t *ConvertTool) Name() string { return "convert" }

func (t *ConvertTool) Description() string {
	return "Convert between data formats"
}

func (t *ConvertTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"input": map[string]interface{}{
				"type":        "string",
				"description": "Input data",
			},
			"from": map[string]interface{}{
				"type":        "string",
				"description": "Source format: json, xml, csv, yaml",
			},
			"to": map[string]interface{}{
				"type":        "string",
				"description": "Target format: json, xml, csv, yaml",
			},
		},
		"required": []string{"input", "from", "to"},
	}
}

func (t *ConvertTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Input string `json:"input"`
		From  string `json:"from"`
		To    string `json:"to"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	return fmt.Sprintf("Converted %s to %s", params.From, params.To), nil
}
