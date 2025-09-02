package domain

import "time"

type Project struct {
	ProjectID        string    `json:"project_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Language         string    `json:"language"`
	ApplyGlobalRules bool      `json:"apply_global_rules"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Rule struct {
	ID          int    `json:"id"`
	ProjectID   string `json:"project_id"`
	RuleID      string `json:"rule_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Pattern     string `json:"pattern"`
	Message     string `json:"message"`
	IsActive    bool   `json:"is_active"`
}

type GlobalRule struct {
	ID          int    `json:"id"`
	Language    string `json:"language"`
	RuleID      string `json:"rule_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Pattern     string `json:"pattern"`
	Message     string `json:"message"`
	IsActive    bool   `json:"is_active"`
}

type ValidationResult struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

type ProjectRules struct {
	ProjectID string `json:"project_id"`
	Rules     []Rule `json:"rules"`
}
