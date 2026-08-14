package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DiffTool represents a diff tool
type DiffTool struct{}

func (t *DiffTool) Name() string { return "diff" }

func (t *DiffTool) Description() string {
	return "Compare two strings"
}

func (t *DiffTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"a": map[string]interface{}{
				"type":        "string",
				"description": "First string",
			},
			"b": map[string]interface{}{
				"type":        "string",
				"description": "Second string",
			},
		},
		"required": []string{"a", "b"},
	}
}

func (t *DiffTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		A string `json:"a"`
		B string `json:"b"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	if params.A == params.B {
		return "Strings are identical", nil
	}
	
	linesA := strings.Split(params.A, "\n")
	linesB := strings.Split(params.B, "\n")
	
	var result strings.Builder
	for i := 0; i < len(linesA) || i < len(linesB); i++ {
		lineA := ""
		lineB := ""
		if i < len(linesA) {
			lineA = linesA[i]
		}
		if i < len(linesB) {
			lineB = linesB[i]
		}
		if lineA != lineB {
			if lineA != "" {
				result.WriteString(fmt.Sprintf("- %s\n", lineA))
			}
			if lineB != "" {
				result.WriteString(fmt.Sprintf("+ %s\n", lineB))
			}
		}
	}
	return result.String(), nil
}
