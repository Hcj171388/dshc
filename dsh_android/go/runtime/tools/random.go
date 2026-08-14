package tools

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// RandomTool represents a random number generator tool
type RandomTool struct{}

func (t *RandomTool) Name() string { return "random" }

func (t *RandomTool) Description() string {
	return "Generate random numbers"
}

func (t *RandomTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"min": map[string]interface{}{
				"type":        "integer",
				"description": "Minimum value",
			},
			"max": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum value",
			},
			"count": map[string]interface{}{
				"type":        "integer",
				"description": "Number of random values",
			},
		},
		"required": []string{},
	}
}

func (t *RandomTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Min   int `json:"min"`
		Max   int `json:"max"`
		Count int `json:"count"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	if params.Count <= 0 {
		params.Count = 1
	}
	if params.Max <= params.Min {
		params.Max = params.Min + 100
	}
	
	rand.Seed(time.Now().UnixNano())
	var result []int
	for i := 0; i < params.Count; i++ {
		result = append(result, rand.Intn(params.Max-params.Min)+params.Min)
	}
	return fmt.Sprintf("%v", result), nil
}
