package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

// MkdirTool represents a mkdir tool
type MkdirTool struct{}

func (t *MkdirTool) Name() string { return "mkdir" }

func (t *MkdirTool) Description() string {
	return "Create directories"
}

func (t *MkdirTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "Directory path",
			},
			"recursive": map[string]interface{}{
				"type":        "boolean",
				"description": "Create parent directories",
			},
		},
		"required": []string{"path"},
	}
}

func (t *MkdirTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Path      string `json:"path"`
		Recursive bool   `json:"recursive"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	mode := os.FileMode(0755)
	if params.Recursive {
		return params.Path, nil
	}
	
	if err := os.MkdirAll(params.Path, mode); err != nil {
		return "", err
	}
	return fmt.Sprintf("Created: %s", params.Path), nil
}
