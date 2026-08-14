package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyDirTool represents a directory copy tool
type CopyDirTool struct{}

func (t *CopyDirTool) Name() string { return "copy_dir" }

func (t *CopyDirTool) Description() string {
	return "Copy directories recursively"
}

func (t *CopyDirTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"src": map[string]interface{}{
				"type":        "string",
				"description": "Source directory",
			},
			"dst": map[string]interface{}{
				"type":        "string",
				"description": "Destination directory",
			},
		},
		"required": []string{"src", "dst"},
	}
}

func (t *CopyDirTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	return copyDir(params.Src, params.Dst)
}

func copyDir(src, dst string) (string, error) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return "", err
	}
	
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return "", err
	}
	
	entries, err := os.ReadDir(src)
	if err != nil {
		return "", err
	}
	
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		
		if entry.IsDir() {
			if _, err := copyDir(srcPath, dstPath); err != nil {
				return "", err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return "", err
			}
		}
	}
	
	return fmt.Sprintf("Copied %s to %s", src, dst), nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	
	_, err = io.Copy(out, in)
	return err
}
