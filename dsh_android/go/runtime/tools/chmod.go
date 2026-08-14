package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

// ChmodTool represents a chmod tool
type ChmodTool struct{}

func (t *ChmodTool) Name() string { return "chmod" }

func (t *ChmodTool) Description() string {
	return "Change file permissions"
}

func (t *ChmodTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "File path",
			},
			"mode": map[string]interface{}{
				"type":        "string",
				"description": "Permission mode (e.g., 755)",
			},
		},
		"required": []string{"path", "mode"},
	}
}

func (t *ChmodTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Path string `json:"path"`
		Mode string `json:"mode"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	var perm os.FileMode
	fmt.Sscanf(params.Mode, "%o", &perm)
	
	if err := os.Chmod(params.Path, perm); err != nil {
		return "", err
	}
	return fmt.Sprintf("Changed permissions of %s to %s", params.Path, params.Mode), nil
}
