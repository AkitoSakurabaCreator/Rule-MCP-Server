package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Rule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Pattern     string `json:"pattern"`
	Message     string `json:"message"`
}

type ProjectRules struct {
	ProjectID string `json:"project_id"`
	Rules     []Rule `json:"rules"`
}

type ValidationResult struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

func loadRules(projectID string) (*ProjectRules, error) {
	log.Printf("Loading rules for project: %s", projectID)

	// データベースが利用可能な場合はデータベースから読み込み
	if db != nil {
		log.Printf("Loading rules from database")
		return db.LoadRulesFromDB(projectID)
	}

	// データベースが利用できない場合はJSONファイルから読み込み
	log.Printf("Loading rules from JSON file")

	// 現在の作業ディレクトリを確認
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get working directory: %v", err)
	} else {
		log.Printf("Current working directory: %s", wd)
	}

	data, err := os.ReadFile("rules.json")
	if err != nil {
		log.Printf("Failed to read rules.json: %v", err)
		return nil, fmt.Errorf("failed to read rules file: %w", err)
	}

	var allRules map[string]ProjectRules
	if err := json.Unmarshal(data, &allRules); err != nil {
		log.Printf("Failed to parse rules.json: %v", err)
		return nil, fmt.Errorf("failed to parse rules file: %w", err)
	}

	log.Printf("Available projects: %v", getProjectIDs(allRules))

	projectRules, exists := allRules[projectID]
	if !exists {
		return nil, fmt.Errorf("project %s not found", projectID)
	}

	log.Printf("Successfully loaded %d rules for project %s", len(projectRules.Rules), projectID)
	return &projectRules, nil
}

func getProjectIDs(rules map[string]ProjectRules) []string {
	ids := make([]string, 0, len(rules))
	for id := range rules {
		ids = append(ids, id)
	}
	return ids
}

func validateCode(projectID, code string) ValidationResult {
	rules, err := loadRules(projectID)
	if err != nil {
		return ValidationResult{
			Valid:  false,
			Errors: []string{fmt.Sprintf("Failed to load rules: %v", err)},
		}
	}

	var errors, warnings []string

	for _, rule := range rules.Rules {
		if strings.Contains(code, rule.Pattern) {
			switch rule.Severity {
			case "error":
				errors = append(errors, rule.Message)
			case "warning":
				warnings = append(warnings, rule.Message)
			}
		}
	}

	return ValidationResult{
		Valid:    len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}
}
