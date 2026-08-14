package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type BashTool struct {
	DefaultTimeoutMs int
}

func (b *BashTool) Name() string { return "bash" }

func (b *BashTool) Description() string {
	return "Execute a bash command and return the output"
}

func (b *BashTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "The bash command to execute",
			},
			"timeout_ms": map[string]interface{}{
				"type":        "integer",
				"description": "Timeout in milliseconds",
			},
		},
		"required": []string{"command"},
	}
}

func (b *BashTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Command   string `json:"command"`
		TimeoutMs int    `json:"timeout_ms"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	if params.Command == "" {
		return "", fmt.Errorf("command is required")
	}
	if params.TimeoutMs == 0 {
		params.TimeoutMs = b.DefaultTimeoutMs
	}
	if params.TimeoutMs == 0 {
		params.TimeoutMs = 30000
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(params.TimeoutMs)*time.Millisecond)
	defer cancel()
	out, err := exec.CommandContext(ctx, "bash", "-c", params.Command).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w, output: %s", err, string(out))
	}
	return string(out), nil
}
