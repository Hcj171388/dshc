package tools

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// EncodeTool represents an encoding tool
type EncodeTool struct{}

func (t *EncodeTool) Name() string { return "encode" }

func (t *EncodeTool) Description() string {
	return "Encode/decode data"
}

func (t *EncodeTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action: encode, decode",
			},
			"data": map[string]interface{}{
				"type":        "string",
				"description": "Data to encode/decode",
			},
			"encoding": map[string]interface{}{
				"type":        "string",
				"description": "Encoding: base64, hex",
			},
		},
		"required": []string{"action", "data"},
	}
}

func (t *EncodeTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Action   string `json:"action"`
		Data     string `json:"data"`
		Encoding string `json:"encoding"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	switch params.Encoding {
	case "base64":
		if params.Action == "encode" {
			return base64.StdEncoding.EncodeToString([]byte(params.Data)), nil
		}
		out, _ := base64.StdEncoding.DecodeString(params.Data); return string(out), nil
	case "hex":
		if params.Action == "encode" {
			return hex.EncodeToString([]byte(params.Data)), nil
		}
		out, _ := hex.DecodeString(params.Data); return string(out), nil
	default:
		return fmt.Sprintf("Unknown encoding: %s", params.Encoding), nil
	}
}
