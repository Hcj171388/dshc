package tools

import (
	"encoding/json"
	"fmt"
)

// JsonQueryTool represents a JSON query tool
type JsonQueryTool struct{}

func (t *JsonQueryTool) Name() string { return "json_query" }

func (t *JsonQueryTool) Description() string {
	return "Query JSON data"
}

func (t *JsonQueryTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"data": map[string]interface{}{
				"type":        "string",
				"description": "JSON string",
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "JSONPath expression",
			},
		},
		"required": []string{"data", "path"},
	}
}

func (t *JsonQueryTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Data string `json:"data"`
		Path string `json:"path"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	var result interface{}
	if err := json.Unmarshal([]byte(params.Data), &result); err != nil {
		return "", err
	}
	
	// Simple path parsing (e.g., "foo.bar[0]")
	parts := []string{}
	current := ""
	for _, ch := range params.Path {
		switch ch {
		case '.':
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		case '[':
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		case ']':
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		default:
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	
	v := result
	for _, p := range parts {
		switch val := v.(type) {
		case map[string]interface{}:
			v = val[p]
		case []interface{}:
			idx, _ := fmt.Sscanf(p, "%d", new(int))
			if idx < len(val) {
				v = val[idx]
			} else {
				return "", fmt.Errorf("index %d out of bounds", idx)
			}
		default:
			return "", fmt.Errorf("cannot index %T", v)
		}
	}
	
	out, _ := json.Marshal(v)
	return string(out), nil
}
