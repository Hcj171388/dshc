package tools

import (
	"encoding/json"
	"strings"
)

// XmlTool represents an XML parsing tool
type XmlTool struct{}

func (t *XmlTool) Name() string { return "xml_parse" }

func (t *XmlTool) Description() string {
	return "Parse XML data"
}

func (t *XmlTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"data": map[string]interface{}{
				"type":        "string",
				"description": "XML string",
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "XPath expression",
			},
		},
		"required": []string{"data"},
	}
}

func (t *XmlTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Data string `json:"data"`
		Path string `json:"path"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	// Simple XML parsing
	if strings.Contains(params.Data, "<") {
		lines := strings.Split(params.Data, "\n")
		var result strings.Builder
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) > 0 && !strings.HasPrefix(line, "<?") && !strings.HasPrefix(line, "<!") {
				result.WriteString(line + "\n")
			}
		}
		return result.String(), nil
	}
	return params.Data, nil
}
