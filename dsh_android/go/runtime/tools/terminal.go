package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// TerminalTool represents a persistent terminal session
type TerminalTool struct {
	SessionID string
	Process   *exec.Cmd
	Output    string
}

func (t *TerminalTool) Name() string { return "terminal" }

func (t *TerminalTool) Description() string {
	return "Execute commands in a persistent terminal session"
}

func (t *TerminalTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"session_id": map[string]interface{}{
				"type":        "string",
				"description": "Terminal session ID",
			},
			"command": map[string]interface{}{
				"type":        "string",
				"description": "Command to execute",
			},
		},
		"required": []string{"command"},
	}
}

func (t *TerminalTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		SessionID string `json:"session_id"`
		Command   string `json:"command"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	if params.Command == "" {
		return "", fmt.Errorf("command is required")
	}
	
	cmd := exec.Command("bash", "-c", params.Command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v\nOutput: %s", err, string(out)), nil
	}
	return string(out), nil
}
