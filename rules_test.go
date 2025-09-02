package main

import (
	"os"
	"testing"
)

func TestLoadRules(t *testing.T) {
	// テスト用の一時的なrules.jsonを作成
	testRules := `{
		"test-project": {
			"project_id": "test-project",
			"rules": [
				{
					"id": "test-rule",
					"name": "Test Rule",
					"description": "A test rule",
					"type": "test",
					"severity": "warning",
					"pattern": "test_pattern",
					"message": "Test pattern detected"
				}
			]
		}
	}`

	err := os.WriteFile("rules.json", []byte(testRules), 0644)
	if err != nil {
		t.Fatalf("Failed to create test rules file: %v", err)
	}
	defer os.Remove("rules.json")

	// ルール読み込みテスト
	rules, err := loadRules("test-project")
	if err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	if rules.ProjectID != "test-project" {
		t.Errorf("Expected project ID 'test-project', got '%s'", rules.ProjectID)
	}

	if len(rules.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(rules.Rules))
	}

	if rules.Rules[0].ID != "test-rule" {
		t.Errorf("Expected rule ID 'test-rule', got '%s'", rules.Rules[0].ID)
	}
}

func TestLoadRulesProjectNotFound(t *testing.T) {
	// 存在しないプロジェクトのテスト
	testRules := `{
		"existing-project": {
			"project_id": "existing-project",
			"rules": []
		}
	}`

	err := os.WriteFile("rules.json", []byte(testRules), 0644)
	if err != nil {
		t.Fatalf("Failed to create test rules file: %v", err)
	}
	defer os.Remove("rules.json")

	_, err = loadRules("non-existent-project")
	if err == nil {
		t.Error("Expected error for non-existent project, got nil")
	}
}

func TestValidateCode(t *testing.T) {
	// テスト用のルールを作成
	testRules := `{
		"validation-test": {
			"project_id": "validation-test",
			"rules": [
				{
					"id": "no-console-log",
					"name": "No Console Log",
					"description": "Console.log statements should not be in production code",
					"type": "style",
					"severity": "warning",
					"pattern": "console.log",
					"message": "Console.log detected. Use proper logging framework in production."
				},
				{
					"id": "no-hardcoded-secrets",
					"name": "No Hardcoded Secrets",
					"description": "API keys should not be hardcoded",
					"type": "security",
					"severity": "error",
					"pattern": "api_key",
					"message": "Hardcoded API key detected. Use environment variables instead."
				}
			]
		}
	}`

	err := os.WriteFile("rules.json", []byte(testRules), 0644)
	if err != nil {
		t.Fatalf("Failed to create test rules file: %v", err)
	}
	defer os.Remove("rules.json")

	// ルール違反なしのコード
	result := validateCode("validation-test", "function testFunction() { return true; }")
	if !result.Valid {
		t.Errorf("Expected valid code, got invalid with errors: %v", result.Errors)
	}

	// 警告レベルのルール違反
	result = validateCode("validation-test", "console.log('test')")
	if !result.Valid {
		t.Error("Expected valid code with warnings, got invalid")
	}
	if len(result.Warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(result.Warnings))
	}

	// エラーレベルのルール違反
	result = validateCode("validation-test", "const api_key = 'secret123'")
	if result.Valid {
		t.Error("Expected invalid code due to hardcoded API key, got valid")
	}
	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}
}
