package tools

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// HashTool represents a hashing tool
type HashTool struct{}

func (t *HashTool) Name() string { return "hash" }

func (t *HashTool) Description() string {
	return "Compute hash of data"
}

func (t *HashTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"data": map[string]interface{}{
				"type":        "string",
				"description": "Data to hash",
			},
			"algorithm": map[string]interface{}{
				"type":        "string",
				"description": "Algorithm: md5, sha256, sha512",
			},
		},
		"required": []string{"data"},
	}
}

func (t *HashTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Data      string `json:"data"`
		Algorithm string `json:"algorithm"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	switch params.Algorithm {
	case "md5":
		h := md5.Sum([]byte(params.Data))
		return hex.EncodeToString(h[:]), nil
	case "sha256":
		h := sha256.Sum256([]byte(params.Data))
		return hex.EncodeToString(h[:]), nil
	default:
		h := sha256.Sum256([]byte(params.Data))
		return hex.EncodeToString(h[:]), nil
	}
}
