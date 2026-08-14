package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

// RenameTool represents a file rename tool
type RenameTool struct{}

func (t *RenameTool) Name() string { return "rename" }

func (t *RenameTool) Description() string {
	return "Rename files or directories"
}

func (t *RenameTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"old_path": map[string]interface{}{
				"type":        "string",
				"description": "Old path",
			},
			"new_path": map[string]interface{}{
				"type":        "string",
				"description": "New path",
			},
		},
		"required": []string{"old_path", "new_path"},
	}
}

func (t *RenameTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		OldPath string `json:"old_path"`
		NewPath string `json:"new_path"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	if err := os.Rename(params.OldPath, params.NewPath); err != nil {
		return "", err
	}
	return fmt.Sprintf("Renamed: %s -> %s", params.OldPath, params.NewPath), nil
}
