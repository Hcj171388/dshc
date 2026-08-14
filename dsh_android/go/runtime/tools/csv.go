package tools

import (
	"encoding/csv"
	"encoding/json"
	"strings"
)

// CsvTool represents a CSV parsing tool
type CsvTool struct{}

func (t *CsvTool) Name() string { return "csv_parse" }

func (t *CsvTool) Description() string {
	return "Parse CSV data"
}

func (t *CsvTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"data": map[string]interface{}{
				"type":        "string",
				"description": "CSV string",
			},
		},
		"required": []string{"data"},
	}
}

func (t *CsvTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	reader := csv.NewReader(strings.NewReader(params.Data))
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}
	
	var result []map[string]string
	if len(records) > 0 {
		headers := records[0]
		for _, row := range records[1:] {
			record := make(map[string]string)
			for i, h := range headers {
				if i < len(row) {
					record[h] = row[i]
				}
			}
			result = append(result, record)
		}
	}
	
	out, _ := json.Marshal(result)
	return string(out), nil
}
