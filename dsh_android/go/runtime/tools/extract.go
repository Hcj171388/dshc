package tools

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractTool represents an archive extraction tool
type ExtractTool struct{}

func (t *ExtractTool) Name() string { return "extract" }

func (t *ExtractTool) Description() string {
	return "Extract archives (zip, tar, gz)"
}

func (t *ExtractTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"file": map[string]interface{}{
				"type":        "string",
				"description": "Archive file path",
			},
			"output": map[string]interface{}{
				"type":        "string",
				"description": "Output directory",
			},
		},
		"required": []string{"file"},
	}
}

func (t *ExtractTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		File   string `json:"file"`
		Output string `json:"output"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	if params.Output == "" {
		params.Output = filepath.Dir(params.File)
	}
	
	ext := strings.ToLower(filepath.Ext(params.File))
	switch ext {
	case ".zip":
		return extractZip(params.File, params.Output)
	case ".tar.gz", ".tgz":
		return extractTarGz(params.File, params.Output)
	case ".gz":
		return extractGz(params.File, params.Output)
	default:
		return fmt.Sprintf("Unsupported archive format: %s", ext), nil
	}
}

func extractZip(zipFile, outputDir string) (string, error) {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", err
	}
	defer r.Close()
	
	for _, f := range r.File {
		path := filepath.Join(outputDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(path), 0755)
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		w, err := os.Create(path)
		if err != nil {
			rc.Close()
			return "", err
		}
		io.Copy(w, rc)
		rc.Close()
		w.Close()
	}
	return fmt.Sprintf("Extracted %d files to %s", len(r.File), outputDir), nil
}

func extractTarGz(tarFile, outputDir string) (string, error) {
	f, err := os.Open(tarFile)
	if err != nil {
		return "", err
	}
	defer f.Close()
	
	gr, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gr.Close()
	
	tr := tar.NewReader(gr)
	count := 0
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		path := filepath.Join(outputDir, header.Name)
		if header.Typeflag == tar.TypeDir {
			os.MkdirAll(path, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(path), 0755)
		w, err := os.Create(path)
		if err != nil {
			return "", err
		}
		io.Copy(w, tr)
		w.Close()
		count++
	}
	return fmt.Sprintf("Extracted %d files", count), nil
}

func extractGz(gzFile, outputDir string) (string, error) {
	f, err := os.Open(gzFile)
	if err != nil {
		return "", err
	}
	defer f.Close()
	
	gr, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gr.Close()
	
	outFile := filepath.Join(outputDir, filepath.Base(gzFile)+".out")
	w, err := os.Create(outFile)
	if err != nil {
		return "", err
	}
	defer w.Close()
	
	io.Copy(w, gr)
	return fmt.Sprintf("Extracted to %s", outFile), nil
}
