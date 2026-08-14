package tools

import (
	"encoding/json"
	"fmt"
)

// SkillTool represents a skill management tool
type SkillTool struct{}

func (t *SkillTool) Name() string { return "skill_manage" }

func (t *SkillTool) Description() string {
	return "Manage registered skills"
}

func (t *SkillTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action to perform: list, register, invoke",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Skill name",
			},
			"data": map[string]interface{}{
				"type":        "object",
				"description": "Skill data",
			},
		},
		"required": []string{"action"},
	}
}

func (t *SkillTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Action string          `json:"action"`
		Name   string          `json:"name"`
		Data   json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	switch params.Action {
	case "list":
		skills := ListSkills()
		return fmt.Sprintf("%v", skills), nil
	case "register":
		RegisterSkill(params.Name, params.Data)
		return fmt.Sprintf("Registered skill: %s", params.Name), nil
	case "invoke":
		_, ok := GetSkill(params.Name)
		if !ok {
			return fmt.Sprintf("Skill not found: %s", params.Name), nil
		}
		return fmt.Sprintf("Invoked skill: %s", params.Name), nil
	default:
		return fmt.Sprintf("Unknown action: %s", params.Action), nil
	}
}
