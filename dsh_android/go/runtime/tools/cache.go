package tools

import (
	"encoding/json"
	"sync"
	"fmt"
)

// CacheTool represents a cache tool
type CacheTool struct {
	mu    sync.Mutex
	store map[string]interface{}
}

func (t *CacheTool) Name() string { return "cache" }

func (t *CacheTool) Description() string {
	return "Cache key-value pairs"
}

func (t *CacheTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action: set, get, delete, clear",
			},
			"key": map[string]interface{}{
				"type":        "string",
				"description": "Cache key",
			},
			"value": map[string]interface{}{
				"type":        "string",
				"description": "Cache value",
			},
		},
		"required": []string{"action"},
	}
}

func (t *CacheTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Action string `json:"action"`
		Key    string `json:"key"`
		Value  string `json:"value"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if t.store == nil {
		t.store = make(map[string]interface{})
	}
	
	switch params.Action {
	case "set":
		t.store[params.Key] = params.Value
		return "Cached", nil
	case "get":
		v, ok := t.store[params.Key]
		if !ok {
			return "Not found", nil
		}
		return fmt.Sprintf("%v", v), nil
	case "delete":
		delete(t.store, params.Key)
		return "Deleted", nil
	case "clear":
		t.store = make(map[string]interface{})
		return "Cleared", nil
	default:
		return "Unknown action", nil
	}
}
