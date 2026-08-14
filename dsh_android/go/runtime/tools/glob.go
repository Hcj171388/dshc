package tools

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

// GlobTool represents a glob search tool
type GlobTool struct{}

func (t *GlobTool) Name() string { return "glob" }

func (t *GlobTool) Description() string {
	return "Find files matching glob pattern"
}

func (t *GlobTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"pattern": map[string]interface{}{
				"type":        "string",
				"description": "Glob pattern",
			},
		},
		"required": []string{"pattern"},
	}
}

func (t *GlobTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Pattern string `json:"pattern"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	matches, err := filepath.Glob(params.Pattern)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", matches), nil
}
