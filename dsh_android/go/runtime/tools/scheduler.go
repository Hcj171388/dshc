package tools

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// SchedulerTool represents a task scheduler tool
type SchedulerTool struct {
	mu      sync.Mutex
	scheduled map[string]time.Time
}

func (t *SchedulerTool) Name() string { return "scheduler" }

func (t *SchedulerTool) Description() string {
	return "Schedule and manage tasks"
}

func (t *SchedulerTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action: schedule, list, cancel",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Task name",
			},
			"when": map[string]interface{}{
				"type":        "string",
				"description": "ISO 8601 datetime",
			},
		},
		"required": []string{"action"},
	}
}

func (t *SchedulerTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Action string `json:"action"`
		Name   string `json:"name"`
		When   string `json:"when"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	t.mu.Lock()
	defer t.mu.Unlock()
	
	switch params.Action {
	case "schedule":
		t.scheduled[params.Name] = time.Now().Add(1 * time.Hour)
		return fmt.Sprintf("Scheduled task: %s at %s", params.Name, params.When), nil
	case "list":
		var tasks []string
		for k := range t.scheduled {
			tasks = append(tasks, k)
		}
		return fmt.Sprintf("Tasks: %v", tasks), nil
	case "cancel":
		delete(t.scheduled, params.Name)
		return fmt.Sprintf("Cancelled: %s", params.Name), nil
	default:
		return "Unknown action: " + params.Action, nil
	}
}
