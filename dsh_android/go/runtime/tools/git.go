package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// GitTool represents a Git tool
type GitTool struct{}

func (t *GitTool) Name() string { return "git_command" }

func (t *GitTool) Description() string {
	return "Execute Git commands"
}

func (t *GitTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "Git command to execute",
			},
		},
		"required": []string{"command"},
	}
}

func (t *GitTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	cmd := exec.Command("bash", "-c", fmt.Sprintf("git %s", params.Command))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v\nOutput: %s", err, string(out)), nil
	}
	return string(out), nil
}
