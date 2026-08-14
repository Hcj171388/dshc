package tools

import (
	"encoding/json"
	"sort"
	"strings"
)

// SortTool represents a sort tool
type SortTool struct{}

func (t *SortTool) Name() string { return "sort" }

func (t *SortTool) Description() string {
	return "Sort lines or values"
}

func (t *SortTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"input": map[string]interface{}{
				"type":        "string",
				"description": "Input text",
			},
			"order": map[string]interface{}{
				"type":        "string",
				"description": "Order: asc, desc",
			},
		},
		"required": []string{"input"},
	}
}

func (t *SortTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Input string `json:"input"`
		Order string `json:"order"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	lines := strings.Split(params.Input, "\n")
	sort.Strings(lines)
	if params.Order == "desc" {
		for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
			lines[i], lines[j] = lines[j], lines[i]
		}
	}
	return strings.Join(lines, "\n"), nil
}
