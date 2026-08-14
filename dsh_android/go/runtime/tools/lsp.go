package tools

import (
	"encoding/json"
	"fmt"
)

// LSPTool represents a Language Server Protocol tool
type LSPTool struct {
	ServerName string
}

func (t *LSPTool) Name() string { return "lsp_call" }

func (t *LSPTool) Description() string {
	return "Call a language server for code intelligence"
}

func (t *LSPTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"server": map[string]interface{}{
				"type":        "string",
				"description": "Language server name",
			},
			"method": map[string]interface{}{
				"type":        "string",
				"description": "LSP method (e.g., textDocument/definition)",
			},
			"params": map[string]interface{}{
				"type":        "object",
				"description": "LSP parameters",
			},
		},
		"required": []string{"server", "method"},
	}
}

func (t *LSPTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Server string      `json:"server"`
		Method string      `json:"method"`
		Params json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	return fmt.Sprintf("LSP call to %s: %s", params.Server, params.Method), nil
}
