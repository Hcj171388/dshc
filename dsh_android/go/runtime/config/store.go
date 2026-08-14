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
	MaxTurns       int     `json:"max_turns"`
	ToolTimeoutMs  int     `json:"tool_timeout_ms"`
	Model          string  `json:"model"`
	APIKey         string  `json:"api_key"`
	BaseURL        string  `json:"base_url"`
	Temperature    float64 `json:"temperature"`
	SystemPrompt   string  `json:"system_prompt"`
}

type ToolsConfig struct {
	Bash BashConfig `json:"bash"`
	Fs   FsConfig   `json:"fs"`
	Web  WebConfig  `json:"web"`
}

type BashConfig struct {
	TimeoutMs int `json:"timeout_ms"`
}

type FsConfig struct {
	ReadLimit      int `json:"read_limit"`
	ReadMaxBytes   int `json:"read_max_bytes"`
	ReadMaxLineLen int `json:"read_max_line_len"`
}

type WebConfig struct {
	TimeoutMs int `json:"timeout_ms"`
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

func defaultConfig() AppConfig {
	return AppConfig{
		Agent: AgentConfig{
			MaxTurns:      20,
			ToolTimeoutMs: 30000,
			Model:         "deepseek-chat",
			Temperature:   0.7,
			SystemPrompt:  "You are a helpful assistant that can execute commands and manipulate files.",
		},
		Tools: ToolsConfig{
			Bash: BashConfig{TimeoutMs: 30000},
			Fs:   FsConfig{ReadLimit: 1000, ReadMaxBytes: 102400, ReadMaxLineLen: 10000},
			Web:  WebConfig{TimeoutMs: 30000},
		},
	}
}
