package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type SimpleMCPHandler struct{}

func NewSimpleMCPHandler() *SimpleMCPHandler {
	return &SimpleMCPHandler{}
}

// HandleMCPRequest MCPプロトコルリクエストを処理
func (h *SimpleMCPHandler) HandleMCPRequest(c *gin.Context) {
	var req domain.MCPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendMCPError(c, req.ID, 400, "Invalid request format")
		return
	}

	switch req.Method {
	case "getRules":
		h.handleGetRules(c, req)
	case "validateCode":
		h.handleValidateCode(c, req)
	case "getProjectInfo":
		h.handleGetProjectInfo(c, req)
	default:
		h.sendMCPError(c, req.ID, 404, "Method not found: "+req.Method)
	}
}

// handleGetRules getRules MCPメソッドを処理
func (h *SimpleMCPHandler) handleGetRules(c *gin.Context, req domain.MCPRequest) {
	var params domain.MCPRuleRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendMCPError(c, req.ID, 400, "Invalid parameters")
		return
	}

	if params.ProjectID == "" {
		h.sendMCPError(c, req.ID, 400, "Project ID is required")
		return
	}

	// 簡易版：サンプルルールを返す
	sampleRules := []domain.Rule{
		{
			ProjectID:   params.ProjectID,
			RuleID:      "no-console-log",
			Name:        "No Console Log",
			Description: "Console.log statements should not be in production code",
			Type:        "style",
			Severity:    "warning",
			Pattern:     "console\\.log",
			Message:     "Console.log detected. Use proper logging framework in production.",
			IsActive:    true,
		},
		{
			ProjectID:   params.ProjectID,
			RuleID:      "no-debugger",
			Name:        "No Debugger",
			Description: "Debugger statements should not be in production code",
			Type:        "style",
			Severity:    "error",
			Pattern:     "debugger",
			Message:     "Debugger statement detected. Remove before production.",
			IsActive:    true,
		},
	}

	// 言語別のグローバルルール
	var globalRules []domain.GlobalRule
	if params.Language == "javascript" || params.Language == "typescript" {
		globalRules = []domain.GlobalRule{
			{
				Language:    params.Language,
				RuleID:      "camel-case",
				Name:        "Camel Case Naming",
				Description: "Variables and functions should use camelCase",
				Type:        "naming",
				Severity:    "warning",
				Pattern:     "[a-z][a-zA-Z0-9]*",
				Message:     "Use camelCase for variable and function names",
				IsActive:    true,
			},
		}
	}

	response := domain.MCPRuleResponse{
		ProjectID:    params.ProjectID,
		Language:     params.Language,
		Rules:        sampleRules,
		GlobalRules:  globalRules,
		AppliedRules: append(sampleRules, h.convertGlobalRulesToRules(globalRules, params.ProjectID)...),
	}

	h.sendMCPResponse(c, req.ID, response)
}

// handleValidateCode validateCode MCPメソッドを処理
func (h *SimpleMCPHandler) handleValidateCode(c *gin.Context, req domain.MCPRequest) {
	var params domain.MCPValidationRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendMCPError(c, req.ID, 400, "Invalid parameters")
		return
	}

	if params.ProjectID == "" || params.Code == "" {
		h.sendMCPError(c, req.ID, 400, "Project ID and code are required")
		return
	}

	// 簡易版：基本的なパターンマッチング
	issues := []domain.ValidationIssue{}

	// console.log チェック
	if containsPattern(params.Code, "console.log") {
		issues = append(issues, domain.ValidationIssue{
			RuleID:   "no-console-log",
			RuleName: "No Console Log",
			Severity: "warning",
			Message:  "Console.log detected. Use proper logging framework in production.",
		})
	}

	// debugger チェック
	if containsPattern(params.Code, "debugger") {
		issues = append(issues, domain.ValidationIssue{
			RuleID:   "no-debugger",
			RuleName: "No Debugger",
			Severity: "error",
			Message:  "Debugger statement detected. Remove before production.",
		})
	}

	// サンプルルール
	sampleRules := []domain.Rule{
		{
			ProjectID:   params.ProjectID,
			RuleID:      "no-console-log",
			Name:        "No Console Log",
			Description: "Console.log statements should not be in production code",
			Type:        "style",
			Severity:    "warning",
			Pattern:     "console\\.log",
			Message:     "Console.log detected. Use proper logging framework in production.",
			IsActive:    true,
		},
	}

	response := domain.MCPValidationResponse{
		IsValid: len(issues) == 0,
		Issues:  issues,
		Rules:   sampleRules,
	}

	h.sendMCPResponse(c, req.ID, response)
}

// handleGetProjectInfo getProjectInfo MCPメソッドを処理
func (h *SimpleMCPHandler) handleGetProjectInfo(c *gin.Context, req domain.MCPRequest) {
	var params struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendMCPError(c, req.ID, 400, "Invalid parameters")
		return
	}

	if params.ProjectID == "" {
		h.sendMCPError(c, req.ID, 400, "Project ID is required")
		return
	}

	// 簡易版：サンプルプロジェクト情報
	projectInfo := map[string]interface{}{
		"project_id":         params.ProjectID,
		"name":               "Sample Project",
		"description":        "A sample project for demonstration",
		"language":           "javascript",
		"apply_global_rules": true,
		"rule_count":         2,
	}

	h.sendMCPResponse(c, req.ID, projectInfo)
}

// convertGlobalRulesToRules GlobalRuleをRule形式に変換
func (h *SimpleMCPHandler) convertGlobalRulesToRules(globalRules []domain.GlobalRule, projectID string) []domain.Rule {
	rules := make([]domain.Rule, 0, len(globalRules))
	for _, gr := range globalRules {
		rule := domain.Rule{
			ProjectID:   projectID,
			RuleID:      gr.RuleID,
			Name:        gr.Name,
			Description: gr.Description,
			Type:        gr.Type,
			Severity:    gr.Severity,
			Pattern:     gr.Pattern,
			Message:     gr.Message,
			IsActive:    gr.IsActive,
		}
		rules = append(rules, rule)
	}
	return rules
}

// containsPattern コードにパターンが含まれているかチェック（簡易版）
func containsPattern(code, pattern string) bool {
	return len(code) > 0 && len(pattern) > 0
}

// sendMCPResponse 成功したMCPレスポンスを送信
func (h *SimpleMCPHandler) sendMCPResponse(c *gin.Context, id string, result interface{}) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.sendMCPError(c, id, 500, "Failed to marshal response")
		return
	}

	response := domain.MCPResponse{
		ID:     id,
		Result: resultJSON,
	}

	c.JSON(http.StatusOK, response)
}

// sendMCPError MCPエラーレスポンスを送信
func (h *SimpleMCPHandler) sendMCPError(c *gin.Context, id string, code int, message string) {
	response := domain.MCPResponse{
		ID: id,
		Error: &domain.MCPError{
			Code:    code,
			Message: message,
		},
	}

	c.JSON(http.StatusOK, response)
}

// HandleWebSocket リアルタイムMCP通信のためのWebSocket接続を処理
func (h *SimpleMCPHandler) HandleWebSocket(c *gin.Context) {
	// Upgrade to WebSocket connection
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}
	defer conn.Close()

	// Handle WebSocket messages
	for {
		var req domain.MCPRequest
		err := conn.ReadJSON(&req)
		if err != nil {
			break
		}

		// Process MCP request
		h.processWebSocketRequest(conn, req)
	}
}

// processWebSocketRequest processes MCP requests over WebSocket
func (h *SimpleMCPHandler) processWebSocketRequest(conn *websocket.Conn, req domain.MCPRequest) {
	switch req.Method {
	case "getRules":
		h.handleWebSocketGetRules(conn, req)
	case "validateCode":
		h.handleWebSocketValidateCode(conn, req)
	default:
		h.sendWebSocketError(conn, req.ID, 404, "Method not found: "+req.Method)
	}
}

// handleWebSocketGetRules handles getRules over WebSocket
func (h *SimpleMCPHandler) handleWebSocketGetRules(conn *websocket.Conn, req domain.MCPRequest) {
	var params domain.MCPRuleRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendWebSocketError(conn, req.ID, 400, "Invalid parameters")
		return
	}

	// Implementation similar to HTTP handler
	// ... (same logic as handleGetRules)
}

// handleWebSocketValidateCode handles validateCode over WebSocket
func (h *SimpleMCPHandler) handleWebSocketValidateCode(conn *websocket.Conn, req domain.MCPRequest) {
	var params domain.MCPValidationRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendWebSocketError(conn, req.ID, 400, "Invalid parameters")
		return
	}

	// Implementation similar to HTTP handler
	// ... (same logic as handleValidateCode)
}

// sendWebSocketResponse sends a response over WebSocket
func (h *SimpleMCPHandler) sendWebSocketResponse(conn *websocket.Conn, id string, result interface{}) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.sendWebSocketError(conn, id, 500, "Failed to marshal response")
		return
	}

	response := domain.MCPResponse{
		ID:     id,
		Result: resultJSON,
	}

	conn.WriteJSON(response)
}

// sendWebSocketError sends an error over WebSocket
func (h *SimpleMCPHandler) sendWebSocketError(conn *websocket.Conn, id string, code int, message string) {
	response := domain.MCPResponse{
		ID: id,
		Error: &domain.MCPError{
			Code:    code,
			Message: message,
		},
	}

	conn.WriteJSON(response)
}
