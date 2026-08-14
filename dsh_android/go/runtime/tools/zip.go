package tools

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ZipTool represents a zip archive tool
type ZipTool struct{}

func (t *ZipTool) Name() string { return "zip" }

func (t *ZipTool) Description() string {
	return "Create ZIP archives"
}

func (t *ZipTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"output": map[string]interface{}{
				"type":        "string",
				"description": "Output ZIP file path",
			},
			"files": map[string]interface{}{
				"type":        "array",
				"description": "Files to include",
			},
		},
		"required": []string{"output", "files"},
	}
}

func (t *ZipTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Output string   `json:"output"`
		Files  []string `json:"files"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	f, err := os.Create(params.Output)
	if err != nil {
		return "", err
	}
	defer f.Close()
	
	w := zip.NewWriter(f)
	for _, file := range params.Files {
		err := addFileToZip(w, file)
		if err != nil {
			w.Close()
			return "", err
		}
	}
	w.Close()
	return fmt.Sprintf("Created zip: %s", params.Output), nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	info, err := file.Stat()
	if err != nil {
		return err
	}
	
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate
	
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	
	_, err = io.Copy(writer, file)
	return err
}
