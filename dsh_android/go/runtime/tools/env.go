package tools

import (
	"encoding/json"
	"os"
)

// EnvTool represents an environment variable tool
type EnvTool struct{}

func (t *EnvTool) Name() string { return "env_get" }

func (t *EnvTool) Description() string {
	return "Get environment variables"
}

func (t *EnvTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"key": map[string]interface{}{
				"type":        "string",
				"description": "Environment variable key",
			},
		},
		"required": []string{"key"},
	}
}

func (t *EnvTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	val := os.Getenv(params.Key)
	return val, nil
}
