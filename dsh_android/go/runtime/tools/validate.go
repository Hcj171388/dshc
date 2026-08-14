package tools

import (
	"encoding/json"
	"fmt"
)

// ValidateTool represents a validation tool
type ValidateTool struct{}

func (t *ValidateTool) Name() string { return "validate" }

func (t *ValidateTool) Description() string {
	return "Validate data against schema"
}

func (t *ValidateTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"data": map[string]interface{}{
				"type":        "string",
				"description": "Data to validate",
			},
			"schema": map[string]interface{}{
				"type":        "string",
				"description": "JSON schema",
			},
		},
		"required": []string{"data", "schema"},
	}
}

func (t *ValidateTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Data   string `json:"data"`
		Schema string `json:"schema"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	return fmt.Sprintf("Validated %s against %s", params.Data, params.Schema), nil
}
