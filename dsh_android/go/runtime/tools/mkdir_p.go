package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

// MkdirpTool represents a mkdir -p tool
type MkdirpTool struct{}

func (t *MkdirpTool) Name() string { return "mkdir_p" }

func (t *MkdirpTool) Description() string {
	return "Create directories recursively (mkdir -p)"
}

func (t *MkdirpTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "Directory path",
			},
		},
		"required": []string{"path"},
	}
}

func (t *MkdirpTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	if err := os.MkdirAll(params.Path, 0755); err != nil {
		return "", err
	}
	return fmt.Sprintf("Created: %s", params.Path), nil
}
