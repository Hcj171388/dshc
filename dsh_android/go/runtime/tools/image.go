package tools

import (
	"encoding/json"
	"fmt"
)

// ImageTool represents an image processing tool
type ImageTool struct{}

func (t *ImageTool) Name() string { return "image_process" }

func (t *ImageTool) Description() string {
	return "Process images"
}

func (t *ImageTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action: resize, rotate, convert",
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "Image path",
			},
			"width": map[string]interface{}{
				"type":        "integer",
				"description": "Width",
			},
			"height": map[string]interface{}{
				"type":        "integer",
				"description": "Height",
			},
		},
		"required": []string{"action", "path"},
	}
}

func (t *ImageTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Action string `json:"action"`
		Path   string `json:"path"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	return fmt.Sprintf("Processed image %s: %s (%dx%d)", params.Path, params.Action, params.Width, params.Height), nil
}
