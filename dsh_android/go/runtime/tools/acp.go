package tools

import (
	"encoding/json"
	"fmt"
)

// ACPTool represents an ACP (Agent Client Protocol) tool
type ACPTool struct{}

func (t *ACPTool) Name() string { return "acp_call" }

func (t *ACPTool) Description() string {
	return "Call an ACP server for automation"
}

func (t *ACPTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"server": map[string]interface{}{
				"type":        "string",
				"description": "ACP server URL",
			},
			"method": map[string]interface{}{
				"type":        "string",
				"description": "RPC method name",
			},
			"params": map[string]interface{}{
				"type":        "object",
				"description": "RPC parameters",
			},
		},
		"required": []string{"server", "method"},
	}
}

func (t *ACPTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Server string                 `json:"server"`
		Method string                 `json:"method"`
		Params map[string]interface{} `json:"params"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	return fmt.Sprintf("ACP call to %s/%s", params.Server, params.Method), nil
}
