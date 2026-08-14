package tools

import (
	"encoding/json"
	"fmt"
)

// GoalTool manages goals and objectives
type GoalTool struct{}

func (t *GoalTool) Name() string { return "goal_set" }

func (t *GoalTool) Description() string {
	return "Set or update a goal for the agent"
}

func (t *GoalTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"goal": map[string]interface{}{
				"type":        "string",
				"description": "The goal to set",
			},
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Optional goal ID",
			},
		},
		"required": []string{"goal"},
	}
}

func (t *GoalTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Goal string `json:"goal"`
		ID   string `json:"id"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	if params.ID == "" {
		params.ID = "default"
	}
	return fmt.Sprintf("Goal set: %s (ID: %s)", params.Goal, params.ID), nil
}
