package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

// Client connects to an MCP server and provides tools
type Client struct {
	serverName  string
	cmd         *exec.Cmd
	stdin       io.WriteCloser
	stdout      *bufio.Scanner
	tools       map[string]Tool
	mu          sync.Mutex
	isConnected bool
}

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// CallResult represents a tool call result
type CallResult struct {
	Content []map[string]interface{} `json:"content"`
	IsError bool                     `json:"isError"`
}

// NewClient creates a new MCP client
func NewClient(serverName, command string, args []string) (*Client, error) {
	cmd := exec.Command(command, args...)
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	
	client := &Client{
		serverName: serverName,
		cmd:        cmd,
		stdin:      stdin,
		stdout:     bufio.NewScanner(stdout),
		tools:      make(map[string]Tool),
	}
	
	if err := cmd.Start(); err != nil {
		stdin.Close()
		return nil, fmt.Errorf("start server: %w", err)
	}
	
	go client.readLoop()
	
	// Initialize connection
	if err := client.initialize(); err != nil {
		cmd.Wait()
		stdin.Close()
		return nil, fmt.Errorf("initialize: %w", err)
	}
	
	client.mu.Lock()
	client.isConnected = true
	client.mu.Unlock()
	
	return client, nil
}

func (c *Client) readLoop() {
	for c.stdout.Scan() {
		line := c.stdout.Text()
		var msg map[string]interface{}
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}
		
		method, _ := msg["method"].(string)
		switch method {
		case "tools/list":
			c.handleToolsList(msg)
		}
	}
}

func (c *Client) handleToolsList(msg map[string]interface{}) {
	// Wait for tools/list result
	// This is a simplified implementation
}

func (c *Client) initialize() error {
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "dsh-android",
				"version": "1.0.0",
			},
		},
	}
	
	reqJSON, _ := json.Marshal(req)
	reqJSON = append(reqJSON, '\n')
	
	_, err := c.stdin.Write(reqJSON)
	return err
}

func (c *Client) discoverTools() error {
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
	}
	
	reqJSON, _ := json.Marshal(req)
	reqJSON = append(reqJSON, '\n')
	
	_, err := c.stdin.Write(reqJSON)
	return err
}

func (c *Client) CallTool(name string, arguments json.RawMessage) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if !c.isConnected {
		return "", fmt.Errorf("not connected to MCP server")
	}
	
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      3,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      name,
			"arguments": json.RawMessage(arguments),
		},
	}
	
	reqJSON, _ := json.Marshal(req)
	reqJSON = append(reqJSON, '\n')
	
	_, err := c.stdin.Write(reqJSON)
	if err != nil {
		return "", err
	}
	
	// Read response
	var result string
	for c.stdout.Scan() {
		line := c.stdout.Text()
		var resp map[string]interface{}
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			continue
		}
		
		resultID, _ := resp["id"].(float64)
		if int(resultID) == 3 {
			content, _ := resp["result"].(map[string]interface{})
			items, _ := content["content"].([]interface{})
			for _, item := range items {
				textItem, ok := item.(map[string]interface{})
				if ok && textItem["type"] == "text" {
					text, _ := textItem["text"].(string)
					result += text
				}
			}
			break
		}
	}
	
	return result, nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.isConnected = false
	c.stdin.Close()
	return c.cmd.Wait()
}

func (c *Client) Name() string {
	return c.serverName
}
