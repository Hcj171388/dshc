package tools

import (
	"encoding/json"
	"fmt"
)

type Tool interface {
	Name() string
	Description() string
	Schema() map[string]interface{}
	Execute(input json.RawMessage) (string, error)
}

// LLMTool represents a tool for the LLM API
type LLMTool struct {
	Type        string                 `json:"type"`
	Function    LLMFunction            `json:"function"`
}

// LLMFunction represents a function for the LLM API
type LLMFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{tools: make(map[string]Tool)}
}

func (r *Registry) Register(t Tool) {
	r.tools[t.Name()] = t
}

func (r *Registry) Get(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

func (r *Registry) Call(name string, input json.RawMessage) (string, error) {
	t, ok := r.tools[name]
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", name)
	}
	return t.Execute(input)
}

// ToLLMTools converts registered tools to LLM tool format
func (r *Registry) ToLLMTools() []LLMTool {
	tools := make([]LLMTool, 0, len(r.tools))
	for _, t := range r.tools {
		tools = append(tools, LLMTool{
			Type: "function",
			Function: LLMFunction{
				Name:        t.Name(),
				Description: t.Description(),
				Parameters:  t.Schema(),
			},
		})
	}
	return tools
}
