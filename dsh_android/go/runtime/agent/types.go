package agent

import "encoding/json"

type AgentEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Config struct {
	MaxTurns      int `json:"max_turns"`
	ToolTimeoutMs int `json:"tool_timeout_ms"`
	MaxParallel   int `json:"max_parallel"`
}

func DefaultConfig() Config {
	return Config{MaxTurns: 20, ToolTimeoutMs: 30000, MaxParallel: 5}
}
