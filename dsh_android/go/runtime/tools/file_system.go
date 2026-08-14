package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// FileSystemTool represents a file system tool
type FileSystemTool struct{}

func (t *FileSystemTool) Name() string { return "fs_operation" }

func (t *FileSystemTool) Description() string {
	return "Perform file system operations"
}

func (t *FileSystemTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action: mkdir, touch, rm, copy, move",
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "Target path",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "Content for write/touch",
			},
		},
		"required": []string{"action", "path"},
	}
}

func (t *FileSystemTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Action  string `json:"action"`
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	switch params.Action {
	case "mkdir":
		if err := os.MkdirAll(params.Path, 0755); err != nil {
			return "", err
		}
		return "Created directory: " + params.Path, nil
	case "touch":
		f, err := os.Create(params.Path)
		if err != nil {
			return "", err
		}
		f.Close()
		return "Created file: " + params.Path, nil
	case "rm":
		if err := os.RemoveAll(params.Path); err != nil {
			return "", err
		}
		return "Removed: " + params.Path, nil
	case "copy":
		src, err := os.ReadFile(params.Path)
		if err != nil {
			return "", err
		}
		dst := filepath.Join(filepath.Dir(params.Path), "copy_" + filepath.Base(params.Path))
		if err := os.WriteFile(dst, src, 0644); err != nil {
			return "", err
		}
		return "Copied to: " + dst, nil
	default:
		return "Unknown action: " + params.Action, nil
	}
}
