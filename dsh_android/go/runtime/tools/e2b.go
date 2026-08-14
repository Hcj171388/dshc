package tools

import (
	"encoding/json"
	"fmt"
)

// E2BTool represents an E2B (Escape to the Box) tool
type E2BTool struct{}

func (t *E2BTool) Name() string { return "e2b_run" }

func (t *E2BTool) Description() string {
	return "Run code in an ephemeral sandbox"
}

func (t *E2BTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"code": map[string]interface{}{
				"type":        "string",
				"description": "Code to execute",
			},
			"language": map[string]interface{}{
				"type":        "string",
				"description": "Programming language",
			},
			"timeout_ms": map[string]interface{}{
				"type":        "integer",
				"description": "Timeout in milliseconds",
			},
		},
		"required": []string{"code"},
	}
}

func (t *E2BTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Code      string `json:"code"`
		Language  string `json:"language"`
		TimeoutMs int    `json:"timeout_ms"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	return fmt.Sprintf("Executed %s code: %s", params.Language, params.Code), nil
}
