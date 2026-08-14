package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type FsTools struct {
	ReadLimit     int
	ReadMaxBytes  int
	ReadMaxLineLen int
}

func NewFsTools(readLimit, readMaxBytes, readMaxLineLen int) *FsTools {
	return &FsTools{
		ReadLimit:      readLimit,
		ReadMaxBytes:   readMaxBytes,
		ReadMaxLineLen: readMaxLineLen,
	}
}

func (f *FsTools) RegisterAll(reg *Registry) {
	reg.Register(&ReadFileTool{fs: f})
	reg.Register(&WriteFileTool{fs: f})
	reg.Register(&ListDirTool{fs: f})
	reg.Register(&SearchFilesTool{fs: f})
}

type ReadFileTool struct{ fs *FsTools }

func (t *ReadFileTool) Name() string { return "read_file" }
func (t *ReadFileTool) Description() string {
	return "Read the contents of a text file"
}
func (t *ReadFileTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{"type": "string", "description": "File path"},
		},
		"required": []string{"path"},
	}
}
func (t *ReadFileTool) Execute(input json.RawMessage) (string, error) {
	var params struct{ Path string `json:"path"` }
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	data, err := os.ReadFile(params.Path)
	if err != nil {
		return "", err
	}
	if len(data) > t.fs.ReadMaxBytes {
		data = data[:t.fs.ReadMaxBytes]
	}
	return string(data), nil
}

type WriteFileTool struct{ fs *FsTools }

func (t *WriteFileTool) Name() string { return "write_file" }
func (t *WriteFileTool) Description() string {
	return "Write content to a file"
}
func (t *WriteFileTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path":    map[string]interface{}{"type": "string"},
			"content": map[string]interface{}{"type": "string"},
		},
		"required": []string{"path", "content"},
	}
}
func (t *WriteFileTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(params.Path), 0755); err != nil {
		return "", err
	}
	return "", os.WriteFile(params.Path, []byte(params.Content), 0644)
}

type ListDirTool struct{ fs *FsTools }

func (t *ListDirTool) Name() string { return "list_dir" }
func (t *ListDirTool) Description() string {
	return "List files in a directory"
}
func (t *ListDirTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{"type": "string"},
		},
		"required": []string{"path"},
	}
}
func (t *ListDirTool) Execute(input json.RawMessage) (string, error) {
	var params struct{ Path string `json:"path"` }
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	entries, err := os.ReadDir(params.Path)
	if err != nil {
		return "", err
	}
	var lines []string
	for _, e := range entries {
		lines = append(lines, e.Name())
	}
	return strings.Join(lines, "\n"), nil
}

type SearchFilesTool struct{ fs *FsTools }

func (t *SearchFilesTool) Name() string { return "search_files" }
func (t *SearchFilesTool) Description() string {
	return "Search for files matching a pattern"
}
func (t *SearchFilesTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path":    map[string]interface{}{"type": "string"},
			"pattern": map[string]interface{}{"type": "string"},
		},
		"required": []string{"path", "pattern"},
	}
}
func (t *SearchFilesTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Path    string `json:"path"`
		Pattern string `json:"pattern"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	var matches []string
	filepath.Walk(params.Path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		matched, _ := filepath.Match(params.Pattern, info.Name())
		if matched {
			matches = append(matches, p)
		}
		return nil
	})
	return strings.Join(matches, "\n"), nil
}
