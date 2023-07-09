package models

type RuleInfo struct {
	Name        string         `json:"rule_name"`
	Constraints map[string]any `json:"constraints"`
}
