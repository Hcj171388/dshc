package tools

import (
	"encoding/json"
	"fmt"
)

// VideoTool represents a video processing tool
type VideoTool struct{}

func (t *VideoTool) Name() string { return "video_process" }

func (t *VideoTool) Description() string {
	return "Process videos"
}

func (t *VideoTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action: trim, convert, extract_audio",
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "Video path",
			},
			"start": map[string]interface{}{
				"type":        "number",
				"description": "Start time",
			},
			"duration": map[string]interface{}{
				"type":        "number",
				"description": "Duration",
			},
		},
		"required": []string{"action", "path"},
	}
}

func (t *VideoTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Action   string  `json:"action"`
		Path     string  `json:"path"`
		Start    float64 `json:"start"`
		Duration float64 `json:"duration"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	return fmt.Sprintf("Processed video %s: %s", params.Path, params.Action), nil
}
