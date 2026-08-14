package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	Agent AgentConfig `json:"agent"`
	Tools ToolsConfig `json:"tools"`
}

type AgentConfig struct {
	MaxTurns      int `json:"max_turns"`
	ToolTimeoutMs int `json:"tool_timeout_ms"`
	MaxParallel   int `json:"max_parallel"`
}

type ToolsConfig struct {
	Bash BashConfig `json:"bash"`
	Fs   FsConfig   `json:"fs"`
}

type BashConfig struct {
	TimeoutMs int `json:"timeout_ms"`
}

type FsConfig struct {
	ReadLimit     int `json:"read_limit"`
	ReadMaxBytes  int `json:"read_max_bytes"`
	ReadMaxLineLen int `json:"read_max_line_len"`
}

type Store struct {
	path string
	data AppConfig
}

func NewStore(dataDir string) (*Store, error) {
	path := filepath.Join(dataDir, "config.json")
	cfg := defaultConfig()
	if data, err := os.ReadFile(path); err == nil {
		json.Unmarshal(data, &cfg)
	}
	return &Store{path: path, data: cfg}, nil
}

func (s *Store) Get() *AppConfig {
	return &s.data
}

func (s *Store) Save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (a *AgentConfig) ToAgentConfig() map[string]interface{} {
	return map[string]interface{}{
		"max_turns":      a.MaxTurns,
		"tool_timeout_ms": a.ToolTimeoutMs,
		"max_parallel":   a.MaxParallel,
	}
}

func defaultConfig() AppConfig {
	return AppConfig{
		Agent: AgentConfig{MaxTurns: 20, ToolTimeoutMs: 30000, MaxParallel: 5},
		Tools: ToolsConfig{
			Bash: BashConfig{TimeoutMs: 30000},
			Fs:   FsConfig{ReadLimit: 1000, ReadMaxBytes: 102400, ReadMaxLineLen: 10000},
		},
	}
}
