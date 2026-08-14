package tools

import (
	"encoding/json"
	"fmt"
)

// CodeRuntimeTool represents a code runtime tool
type CodeRuntimeTool struct{}

func (t *CodeRuntimeTool) Name() string { return "code_runtime" }

func (t *CodeRuntimeTool) Description() string {
	return "Run code in various languages"
}

func (t *CodeRuntimeTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"language": map[string]interface{}{
				"type":        "string",
				"description": "Programming language",
			},
			"code": map[string]interface{}{
				"type":        "string",
				"description": "Code to execute",
			},
			"args": map[string]interface{}{
				"type":        "array",
				"description": "Arguments",
			},
		},
		"required": []string{"language", "code"},
	}
}

func (t *CodeRuntimeTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Language string   `json:"language"`
		Code     string   `json:"code"`
		Args     []string `json:"args"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	return fmt.Sprintf("Executed %s code: %s with args %v", params.Language, params.Code, params.Args), nil
}
