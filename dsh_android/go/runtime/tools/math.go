package tools

import (
	"encoding/json"
	"fmt"
	"math"
)

// MathTool represents a math tool
type MathTool struct{}

func (t *MathTool) Name() string { return "math" }

func (t *MathTool) Description() string {
	return "Perform mathematical operations"
}

func (t *MathTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"operation": map[string]interface{}{
				"type":        "string",
				"description": "Operation: add, sub, mul, div, pow, sqrt, abs",
			},
			"a": map[string]interface{}{
				"type":        "number",
				"description": "First operand",
			},
			"b": map[string]interface{}{
				"type":        "number",
				"description": "Second operand",
			},
		},
		"required": []string{"operation", "a"},
	}
}

func (t *MathTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Operation string  `json:"operation"`
		A         float64 `json:"a"`
		B         float64 `json:"b"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	var result float64
	switch params.Operation {
	case "add":
		result = params.A + params.B
	case "sub":
		result = params.A - params.B
	case "mul":
		result = params.A * params.B
	case "div":
		if params.B == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = params.A / params.B
	case "pow":
		result = math.Pow(params.A, params.B)
	case "sqrt":
		if params.A < 0 {
			return "", fmt.Errorf("square root of negative number")
		}
		result = math.Sqrt(params.A)
	case "abs":
		result = math.Abs(params.A)
	default:
		return "", fmt.Errorf("unknown operation: %s", params.Operation)
	}
	
	return fmt.Sprintf("%v", result), nil
}
