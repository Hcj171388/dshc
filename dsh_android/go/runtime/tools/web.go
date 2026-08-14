package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WebTools provides web-related tools
type WebTools struct {
	Timeout time.Duration
}

func NewWebTools(timeout time.Duration) *WebTools {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &WebTools{Timeout: timeout}
}

func (w *WebTools) RegisterAll(reg *Registry) {
	reg.Register(&WebSearchTool{web: w})
	reg.Register(&WebFetchTool{web: w})
}

type WebSearchTool struct {
	web *WebTools
}

func (t *WebSearchTool) Name() string { return "web_search" }

func (t *WebSearchTool) Description() string {
	return "Search the web for information"
}

func (t *WebSearchTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Search query",
			},
			"max_results": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum number of results",
			},
		},
		"required": []string{"query"},
	}
}

func (t *WebSearchTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Query      string `json:"query"`
		MaxResults int    `json:"max_results"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	if params.Query == "" {
		return "", fmt.Errorf("query is required")
	}
	if params.MaxResults == 0 {
		params.MaxResults = 5
	}
	
	// For now, return a placeholder - in production this would call a search API
	return fmt.Sprintf("Search results for: %s\n\n(Note: Web search requires API key configuration)", params.Query), nil
}

type WebFetchTool struct {
	web *WebTools
}

func (t *WebFetchTool) Name() string { return "web_fetch" }

func (t *WebFetchTool) Description() string {
	return "Fetch content from a URL"
}

func (t *WebFetchTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"url": map[string]interface{}{
				"type":        "string",
				"description": "URL to fetch",
			},
		},
		"required": []string{"url"},
	}
}

func (t *WebFetchTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	if params.URL == "" {
		return "", fmt.Errorf("url is required")
	}
	
	client := &http.Client{Timeout: t.web.Timeout}
	req, err := http.NewRequest("GET", params.URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "DeepSeekHarness/1.0")
	
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	// Extract text content (simplified HTML parsing)
	content := strings.TrimSpace(string(body))
	if len(content) > 10000 {
		content = content[:10000] + "...(truncated)"
	}
	
	return content, nil
}
