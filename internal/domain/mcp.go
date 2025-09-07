package domain

import (
	"encoding/json"
	"time"
)

// MCPRequest 標準的なMCPリクエストを表す
type MCPRequest struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// MCPResponse 標準的なMCPレスポンスを表す
type MCPResponse struct {
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *MCPError       `json:"error,omitempty"`
}

// MCPError MCPエラーを表す
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCPRuleRequest ルールリクエストを表す
type MCPRuleRequest struct {
	ProjectID string `json:"project_id"`
	Language  string `json:"language,omitempty"`
}

// MCPRuleResponse ルールを含むレスポンスを表す
type MCPRuleResponse struct {
	ProjectID    string       `json:"project_id"`
	Language     string       `json:"language"`
	Rules        []Rule       `json:"rules"`
	GlobalRules  []GlobalRule `json:"global_rules,omitempty"`
	AppliedRules []Rule       `json:"applied_rules"`
}

// MCPValidationRequest コード検証リクエストを表す
type MCPValidationRequest struct {
	ProjectID string `json:"project_id"`
	Code      string `json:"code"`
	Language  string `json:"language,omitempty"`
}

// MCPValidationResponse コード検証レスポンスを表す
type MCPValidationResponse struct {
	IsValid bool              `json:"is_valid"`
	Issues  []ValidationIssue `json:"issues,omitempty"`
	Rules   []Rule            `json:"applied_rules"`
}

// ValidationIssue コードで発見された検証問題を表す
type ValidationIssue struct {
	RuleID      string `json:"rule_id"`
	RuleName    string `json:"rule_name"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	LineNumber  int    `json:"line_number,omitempty"`
	ColumnStart int    `json:"column_start,omitempty"`
	ColumnEnd   int    `json:"column_end,omitempty"`
}

// MCPNotification サーバーからの通知を表す
type MCPNotification struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// MCPRuleUpdateNotification ルール更新通知を表す
type MCPRuleUpdateNotification struct {
	ProjectID string    `json:"project_id"`
	UpdatedAt time.Time `json:"updated_at"`
	RuleCount int       `json:"rule_count"`
}
