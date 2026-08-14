package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// CurlTool represents a curl tool
type CurlTool struct{}

func (t *CurlTool) Name() string { return "curl" }

func (t *CurlTool) Description() string {
	return "Execute curl command"
}

func (t *CurlTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "Curl command",
			},
		},
		"required": []string{"command"},
	}
}

func (t *CurlTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	cmd := exec.Command("bash", "-c", fmt.Sprintf("curl %s", params.Command))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v\nOutput: %s", err, string(out)), nil
	}
	return string(out), nil
}
