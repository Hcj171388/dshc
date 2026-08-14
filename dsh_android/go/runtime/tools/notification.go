package tools

import (
	"encoding/json"
	"fmt"
)

// NotificationTool represents a notification tool
type NotificationTool struct{}

func (t *NotificationTool) Name() string { return "notify" }

func (t *NotificationTool) Description() string {
	return "Send notifications"
}

func (t *NotificationTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"message": map[string]interface{}{
				"type":        "string",
				"description": "Notification message",
			},
			"priority": map[string]interface{}{
				"type":        "string",
				"description": "Priority: low, medium, high",
			},
		},
		"required": []string{"message"},
	}
}

func (t *NotificationTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Message string `json:"message"`
		Priority string `json:"priority"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	return fmt.Sprintf("Notification sent: %s (priority: %s)", params.Message, params.Priority), nil
}
