package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// TodoTool manages a todo list
type TodoTool struct {
	Filepath string
}

func (t *TodoTool) Name() string { return "todo_write" }

func (t *TodoTool) Description() string {
	return "Write/update the todo list"
}

func (t *TodoTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"todos": map[string]interface{}{
				"type":        "array",
				"description": "List of todos",
			},
		},
		"required": []string{"todos"},
	}
}

func (t *TodoTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Todos []string `json:"todos"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	if t.Filepath == "" {
		t.Filepath = "TODO.md"
	}
	
	var content string
	for i, todo := range params.Todos {
		content += fmt.Sprintf("%d. %s\n", i+1, todo)
	}
	
	os.MkdirAll(filepath.Dir(t.Filepath), 0755)
	return "", os.WriteFile(t.Filepath, []byte(content), 0644)
}
