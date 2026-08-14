package mcp

import "encoding/json"

// Tool represents an MCP tool for the registry
type MCPTool struct {
	Client *Client
	Name   string
}

func (t *MCPTool) GetMCPClient() *Client {
	return t.Client
}

func (t *MCPTool) GetMCPName() string {
	return t.Name
}

// Execute calls the MCP tool
func (t *MCPTool) Execute(input json.RawMessage) (string, error) {
	return t.Client.CallTool(t.Name, input)
}
