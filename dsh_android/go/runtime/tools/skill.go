package tools

import (
)

// SkillRegistry stores registered skills
type SkillRegistry map[string]interface{}

var skillRegistry = make(SkillRegistry)

// RegisterSkill registers a skill
func RegisterSkill(name string, data interface{}) {
	skillRegistry[name] = data
}

// GetSkill returns a skill
func GetSkill(name string) (interface{}, bool) {
	v, ok := skillRegistry[name]
	return v, ok
}

// ListSkills returns all skills
func ListSkills() []string {
	skills := make([]string, 0, len(skillRegistry))
	for k := range skillRegistry {
		skills = append(skills, k)
	}
	return skills
}
