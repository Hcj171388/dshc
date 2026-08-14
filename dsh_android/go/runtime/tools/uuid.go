package tools

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

// UUIDTool represents a UUID generator tool
type UUIDTool struct{}

func (t *UUIDTool) Name() string { return "uuid" }

func (t *UUIDTool) Description() string {
	return "Generate UUIDs"
}

func (t *UUIDTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"count": map[string]interface{}{
				"type":        "integer",
				"description": "Number of UUIDs to generate",
			},
		},
	}
}

func (t *UUIDTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	if params.Count <= 0 {
		params.Count = 1
	}
	
	var result []string
	for i := 0; i < params.Count; i++ {
		result = append(result, uuid.New().String())
	}
	return fmt.Sprintf("%v", result), nil
}
