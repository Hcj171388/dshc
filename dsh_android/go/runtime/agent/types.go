package agent

import "encoding/json"

// AgentEvent represents an event from the agent loop
type AgentEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
