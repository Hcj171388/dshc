package tools

import (
	"encoding/json"
	"os"
)

// TeeTool represents a tee command tool
type TeeTool struct{}

func (t *TeeTool) Name() string { return "tee" }

func (t *TeeTool) Description() string {
	return "Write to file and stdout"
}

func (t *TeeTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"input": map[string]interface{}{
				"type":        "string",
				"description": "Input text",
			},
			"file": map[string]interface{}{
				"type":        "string",
				"description": "Output file",
			},
		},
		"required": []string{"input", "file"},
	}
}

func (t *TeeTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Input string `json:"input"`
		File  string `json:"file"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	f, err := os.Create(params.File)
	if err != nil {
		return "", err
	}
	defer f.Close()
	
	f.WriteString(params.Input)
	return params.Input, nil
}
