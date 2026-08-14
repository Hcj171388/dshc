package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ContextTool manages context information
type ContextTool struct {
	ContextDir string
}

func (t *ContextTool) Name() string { return "context_manage" }

func (t *ContextTool) Description() string {
	return "Manage context and knowledge base"
}

func (t *ContextTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"load", "save", "query"},
				"description": "Action to perform",
			},
			"key": map[string]interface{}{
				"type":        "string",
				"description": "Context key",
			},
			"value": map[string]interface{}{
				"type":        "string",
				"description": "Context value",
			},
		},
		"required": []string{"action"},
	}
}

func (t *ContextTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Action string `json:"action"`
		Key    string `json:"key"`
		Value  string `json:"value"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	if t.ContextDir == "" {
		t.ContextDir = "./context"
	}
	
	switch params.Action {
	case "load":
		return t.loadContext(params.Key)
	case "save":
		return t.saveContext(params.Key, params.Value)
	case "query":
		return t.queryContext(params.Key)
	default:
		return "", fmt.Errorf("unknown action: %s", params.Action)
	}
}

func (t *ContextTool) loadContext(key string) (string, error) {
	path := filepath.Join(t.ContextDir, key+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Context not found: %s", key), nil
	}
	return string(data), nil
}

func (t *ContextTool) saveContext(key, value string) (string, error) {
	os.MkdirAll(t.ContextDir, 0755)
	path := filepath.Join(t.ContextDir, key+".json")
	return "", os.WriteFile(path, []byte(value), 0644)
}

func (t *ContextTool) queryContext(key string) (string, error) {
	path := filepath.Join(t.ContextDir, key+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return "Not found", nil
	}
	return string(data), nil
}
