package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SkillTool manages skill definitions
type SkillTool struct {
	SkillDir string
}

func (t *SkillTool) Name() string { return "skill_manage" }

func (t *SkillTool) Description() string {
	return "Manage skills and skill definitions"
}

func (t *SkillTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"list", "load", "execute"},
				"description": "Action to perform",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Skill name",
			},
		},
		"required": []string{"action"},
	}
}

func (t *SkillTool) Execute(input json.RawMessage) (string, error) {
	var params struct {
		Action string `json:"action"`
		Name   string `json:"name"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return "", err
	}
	
	if t.SkillDir == "" {
		t.SkillDir = "./skills"
	}
	
	switch params.Action {
	case "list":
		return t.listSkills()
	case "load":
		return t.loadSkill(params.Name)
	case "execute":
		return t.executeSkill(params.Name)
	default:
		return "", fmt.Errorf("unknown action: %s", params.Action)
	}
}

func (t *SkillTool) listSkills() (string, error) {
	entries, err := os.ReadDir(t.SkillDir)
	if err != nil {
		return "No skills found", nil
	}
	
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return strings.Join(names, ", "), nil
}

func (t *SkillTool) loadSkill(name string) (string, error) {
	path := filepath.Join(t.SkillDir, name, "skill.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("skill not found: %s", name)
	}
	return string(data), nil
}

func (t *SkillTool) executeSkill(name string) (string, error) {
	return fmt.Sprintf("Executed skill: %s", name), nil
}
