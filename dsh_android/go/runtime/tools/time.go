package tools

import (
	"encoding/json"
	"time"
)

// TimeTool represents a time tool
type TimeTool struct{}

func (t *TimeTool) Name() string { return "time" }

func (t *TimeTool) Description() string {
	return "Get current time or parse dates"
}

func (t *TimeTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action: now, format, parse",
			},
			"format": map[string]interface{}{
				"type":        "string",
				"description": "Time format",
			},
			"input": map[string]interface{}{
				"type":        "string",
				"description": "Input date string",
			},
		},
		"required": []string{"action"},
	}
}

func (t *TimeTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Action string `json:"action"`
		Format string `json:"format"`
		Input  string `json:"input"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	switch params.Action {
	case "now":
		return time.Now().Format("2006-01-02 15:04:05"), nil
	case "format":
		if params.Format == "" {
			params.Format = "2006-01-02 15:04:05"
		}
		return time.Now().Format(params.Format), nil
	case "parse":
		t, err := time.Parse("2006-01-02 15:04:05", params.Input)
		if err != nil {
			return "", err
		}
		return t.Format(params.Format), nil
	default:
		return time.Now().Format("2006-01-02 15:04:05"), nil
	}
}
