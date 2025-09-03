package domain

import (
	"encoding/json"
	"time"
)

// MCPRequest represents a standard MCP request
type MCPRequest struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// MCPResponse represents a standard MCP response
type MCPResponse struct {
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *MCPError       `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCPRuleRequest represents a request for rules
type MCPRuleRequest struct {
	ProjectID string `json:"project_id"`
	Language  string `json:"language,omitempty"`
}

// MCPRuleResponse represents a response containing rules
type MCPRuleResponse struct {
	ProjectID    string       `json:"project_id"`
	Language     string       `json:"language"`
	Rules        []Rule       `json:"rules"`
	GlobalRules  []GlobalRule `json:"global_rules,omitempty"`
	AppliedRules []Rule       `json:"applied_rules"`
}

// MCPValidationRequest represents a code validation request
type MCPValidationRequest struct {
	ProjectID string `json:"project_id"`
	Code      string `json:"code"`
	Language  string `json:"language,omitempty"`
}

// MCPValidationResponse represents a code validation response
type MCPValidationResponse struct {
	IsValid bool              `json:"is_valid"`
	Issues  []ValidationIssue `json:"issues,omitempty"`
	Rules   []Rule            `json:"applied_rules"`
}

// ValidationIssue represents a validation issue found in code
type ValidationIssue struct {
	RuleID      string `json:"rule_id"`
	RuleName    string `json:"rule_name"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	LineNumber  int    `json:"line_number,omitempty"`
	ColumnStart int    `json:"column_start,omitempty"`
	ColumnEnd   int    `json:"column_end,omitempty"`
}

// MCPNotification represents a notification from the server
type MCPNotification struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// MCPRuleUpdateNotification represents a rule update notification
type MCPRuleUpdateNotification struct {
	ProjectID string    `json:"project_id"`
	UpdatedAt time.Time `json:"updated_at"`
	RuleCount int       `json:"rule_count"`
}
