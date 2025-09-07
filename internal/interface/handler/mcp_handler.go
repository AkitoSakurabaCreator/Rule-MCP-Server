package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/usecase"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/mcpx"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type MCPHandler struct {
	ruleUseCase       *usecase.RuleUseCase
	globalRuleUseCase *usecase.GlobalRuleUseCase
	projectDetector   *usecase.ProjectDetector
	metricsRepo       domain.MetricsRepository
	metricsHandler    *MetricsHandler
}

func NewMCPHandler(ruleUseCase *usecase.RuleUseCase, globalRuleUseCase *usecase.GlobalRuleUseCase, projectDetector *usecase.ProjectDetector) *MCPHandler {
	return &MCPHandler{
		ruleUseCase:       ruleUseCase,
		globalRuleUseCase: globalRuleUseCase,
		projectDetector:   projectDetector,
	}
}

// SetMetricsRepo メトリクスリポジトリを注入
func (h *MCPHandler) SetMetricsRepo(repo domain.MetricsRepository) {
	h.metricsRepo = repo
}

// SetMetricsHandler メトリクスハンドラーを注入
func (h *MCPHandler) SetMetricsHandler(handler *MetricsHandler) {
	h.metricsHandler = handler
}

func (h *MCPHandler) withMetrics(method string, handler func() error) {
	start := time.Now()
	status := "ok"
	if err := handler(); err != nil {
		status = "error"
	}
	duration := time.Since(start)

	// データベースメトリクス
	if h.metricsRepo != nil {
		_ = h.metricsRepo.RecordMCP(method, status, int(duration/time.Millisecond))
	}

	// Prometheusメトリクス
	if h.metricsHandler != nil {
		h.metricsHandler.RecordMCPRequest(method, status, duration)
	}
}

// HandleMCPRequest MCPプロトコルリクエストを処理
func (h *MCPHandler) HandleMCPRequest(c *gin.Context) {
	var req domain.MCPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendMCPError(c, "", mcpx.CodeValidation, "Invalid request format")
		return
	}

	switch req.Method {
	case "tools/list":
		h.withMetrics("tools/list", func() error { h.handleToolsList(c, req); return nil })
	case "getRules":
		h.withMetrics("getRules", func() error { h.handleGetRules(c, req); return nil })
	case "validateCode":
		h.withMetrics("validateCode", func() error { h.handleValidateCode(c, req); return nil })
	case "getProjectInfo":
		h.withMetrics("getProjectInfo", func() error { h.handleGetProjectInfo(c, req); return nil })
	case "autoDetectProject":
		h.withMetrics("autoDetectProject", func() error { h.handleAutoDetectProject(c, req); return nil })
	case "scanLocalProjects":
		h.withMetrics("scanLocalProjects", func() error { h.handleScanLocalProjects(c, req); return nil })
	default:
		h.sendMCPError(c, req.ID, mcpx.CodeNotFound, "Method not found: "+req.Method)
	}
}

// handleToolsList tools/list MCPメソッドを処理
func (h *MCPHandler) handleToolsList(c *gin.Context, req domain.MCPRequest) {
	tools := []map[string]interface{}{
		{
			"name":        "getRules",
			"description": "Get coding rules for a specific project",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "The project ID to get rules for",
					},
					"language": map[string]interface{}{
						"type":        "string",
						"description": "Programming language (optional)",
					},
				},
				"required": []string{"project_id"},
			},
		},
		{
			"name":        "validateCode",
			"description": "Validate code against project rules",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "The project ID to validate against",
					},
					"code": map[string]interface{}{
						"type":        "string",
						"description": "The code to validate",
					},
					"language": map[string]interface{}{
						"type":        "string",
						"description": "Programming language (optional)",
					},
				},
				"required": []string{"project_id", "code"},
			},
		},
		{
			"name":        "getProjectInfo",
			"description": "Get information about a specific project",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "The project ID to get info for",
					},
				},
				"required": []string{"project_id"},
			},
		},
		{
			"name":        "autoDetectProject",
			"description": "Automatically detect project from path and get appropriate rules",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "The path to detect project from",
					},
				},
				"required": []string{"path"},
			},
		},
		{
			"name":        "scanLocalProjects",
			"description": "Scan local directory to detect multiple projects",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"base_path": map[string]interface{}{
						"type":        "string",
						"description": "The base path to scan for projects (optional, defaults to /)",
					},
				},
			},
		},
	}

	h.sendMCPResponse(c, req.ID, map[string]interface{}{"tools": tools})
}

// handleGetRules getRules MCPメソッドを処理
func (h *MCPHandler) handleGetRules(c *gin.Context, req domain.MCPRequest) {
	var params domain.MCPRuleRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendMCPError(c, req.ID, mcpx.CodeValidation, "Invalid parameters")
		return
	}

	if params.ProjectID == "" {
		h.sendMCPError(c, req.ID, mcpx.CodeValidation, "Project ID is required")
		return
	}

	// プロジェクトルールを取得
	projectRules, err := h.ruleUseCase.GetProjectRules(params.ProjectID)
	if err != nil {
		code, msg := mcpx.MapAppErrorToMCP(err)
		h.sendMCPError(c, req.ID, code, "Failed to get project rules: "+msg)
		return
	}

	// 言語が指定されている場合はグローバルルールを取得
	var globalRules []domain.GlobalRule
	if params.Language != "" {
		globalRulesPtr, err := h.globalRuleUseCase.GetGlobalRules(params.Language)
		if err == nil {
			// ポインタスライスを値スライスに変換
			globalRules = make([]domain.GlobalRule, len(globalRulesPtr))
			for i, gr := range globalRulesPtr {
				globalRules[i] = *gr
			}
		}
	}

	// プロジェクトルールとグローバルルールを結合
	appliedRules := make([]domain.Rule, 0, len(projectRules.Rules)+len(globalRules))
	appliedRules = append(appliedRules, projectRules.Rules...)

	// グローバルルールをプロジェクトルール形式に変換
	for _, gr := range globalRules {
		rule := domain.Rule{
			RuleID:      gr.RuleID,
			Name:        gr.Name,
			Description: gr.Description,
			Type:        gr.Type,
			Severity:    gr.Severity,
			Pattern:     gr.Pattern,
			Message:     gr.Message,
			IsActive:    gr.IsActive,
		}
		appliedRules = append(appliedRules, rule)
	}

	response := domain.MCPRuleResponse{
		ProjectID:    params.ProjectID,
		Language:     params.Language,
		Rules:        projectRules.Rules,
		GlobalRules:  globalRules,
		AppliedRules: appliedRules,
	}

	h.sendMCPResponse(c, req.ID, response)
}

// handleValidateCode validateCode MCPメソッドを処理
func (h *MCPHandler) handleValidateCode(c *gin.Context, req domain.MCPRequest) {
	var params domain.MCPValidationRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendMCPError(c, req.ID, mcpx.CodeValidation, "Invalid parameters")
		return
	}

	if params.ProjectID == "" || params.Code == "" {
		h.sendMCPError(c, req.ID, mcpx.CodeValidation, "Project ID and code are required")
		return
	}

	// プロジェクトルールに対してコードを検証
	validationResult, err := h.ruleUseCase.ValidateCode(params.ProjectID, params.Code)
	if err != nil {
		code, msg := mcpx.MapAppErrorToMCP(err)
		h.sendMCPError(c, req.ID, code, "Failed to validate code: "+msg)
		return
	}

	// 検証結果をMCP形式に変換
	issues := make([]domain.ValidationIssue, 0)

	// エラーを検証問題に変換
	for _, errorMsg := range validationResult.Errors {
		issue := domain.ValidationIssue{
			RuleID:   "validation-error",
			RuleName: "Code Validation Error",
			Severity: "error",
			Message:  errorMsg,
		}
		issues = append(issues, issue)
	}

	// 警告を検証問題に変換
	for _, warningMsg := range validationResult.Warnings {
		issue := domain.ValidationIssue{
			RuleID:   "validation-warning",
			RuleName: "Code Validation Warning",
			Severity: "warning",
			Message:  warningMsg,
		}
		issues = append(issues, issue)
	}

	// コンテキスト用に適用されたルールを取得
	projectRules, err := h.ruleUseCase.GetProjectRules(params.ProjectID)
	if err != nil {
		projectRules = &domain.ProjectRules{Rules: []domain.Rule{}}
	}

	response := domain.MCPValidationResponse{
		IsValid: validationResult.Valid,
		Issues:  issues,
		Rules:   projectRules.Rules,
	}

	h.sendMCPResponse(c, req.ID, response)
}

// handleGetProjectInfo getProjectInfo MCPメソッドを処理
func (h *MCPHandler) handleGetProjectInfo(c *gin.Context, req domain.MCPRequest) {
	var params struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendMCPError(c, req.ID, mcpx.CodeValidation, "Invalid parameters")
		return
	}

	if params.ProjectID == "" {
		h.sendMCPError(c, req.ID, mcpx.CodeValidation, "Project ID is required")
		return
	}

	// プロジェクト情報を取得 - このメソッドはusecaseで実装する必要がある
	// 今のところ、エラーを返す
	h.sendMCPError(c, req.ID, mcpx.CodeInternal, "getProjectInfo not yet implemented")
}

// handleAutoDetectProject autoDetectProject MCPメソッドを処理
func (h *MCPHandler) handleAutoDetectProject(c *gin.Context, req domain.MCPRequest) {
	var params struct {
		Path string `json:"path"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendMCPError(c, req.ID, mcpx.CodeValidation, "Invalid parameters")
		return
	}

	if params.Path == "" {
		h.sendMCPError(c, req.ID, mcpx.CodeValidation, "Path is required")
		return
	}

	// プロジェクトを自動検出
	result, err := h.projectDetector.AutoDetectProject(params.Path)
	if err != nil {
		code, msg := mcpx.MapAppErrorToMCP(err)
		h.sendMCPError(c, req.ID, code, "Project not found: "+msg)
		return
	}

	// 結果をJSONに変換
	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.sendMCPError(c, req.ID, mcpx.CodeInternal, "Failed to serialize result")
		return
	}

	response := domain.MCPResponse{
		ID:     req.ID,
		Result: resultJSON,
	}

	c.JSON(http.StatusOK, response)
}

// handleScanLocalProjects scanLocalProjects MCPメソッドを処理
func (h *MCPHandler) handleScanLocalProjects(c *gin.Context, req domain.MCPRequest) {
	var params struct {
		BasePath string `json:"base_path"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendMCPError(c, req.ID, mcpx.CodeValidation, "Invalid parameters")
		return
	}

	if params.BasePath == "" {
		params.BasePath = "/" // デフォルトはルートディレクトリ
	}

	// ローカルプロジェクトをスキャン
	results, err := h.projectDetector.ScanLocalProjects(params.BasePath)
	if err != nil {
		code, msg := mcpx.MapAppErrorToMCP(err)
		h.sendMCPError(c, req.ID, code, "Failed to scan local projects: "+msg)
		return
	}

	// 結果をJSONに変換
	responseData := gin.H{
		"projects": results,
		"count":    len(results),
	}

	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		h.sendMCPError(c, req.ID, mcpx.CodeInternal, "Failed to serialize results")
		return
	}

	response := domain.MCPResponse{
		ID:     req.ID,
		Result: responseJSON,
	}

	c.JSON(http.StatusOK, response)
}

// sendMCPResponse 成功したMCPレスポンスを送信
func (h *MCPHandler) sendMCPResponse(c *gin.Context, id string, result interface{}) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.sendMCPError(c, id, mcpx.CodeInternal, "Failed to marshal response")
		return
	}

	response := domain.MCPResponse{
		ID:     id,
		Result: resultJSON,
	}

	c.JSON(http.StatusOK, response)
}

// sendMCPError MCPエラーレスポンスを送信
func (h *MCPHandler) sendMCPError(c *gin.Context, id string, code int, message string) {
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
func (h *MCPHandler) HandleWebSocket(c *gin.Context) {
	// WebSocket接続にアップグレード
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 開発用にすべてのオリジンを許可
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}
	defer conn.Close()

	// WebSocketメッセージを処理
	for {
		var req domain.MCPRequest
		err := conn.ReadJSON(&req)
		if err != nil {
			break
		}

		// MCPリクエストを処理
		h.processWebSocketRequest(conn, req)
	}
}

// processWebSocketRequest WebSocket経由でMCPリクエストを処理
func (h *MCPHandler) processWebSocketRequest(conn *websocket.Conn, req domain.MCPRequest) {
	switch req.Method {
	case "getRules":
		h.handleWebSocketGetRules(conn, req)
	case "validateCode":
		h.handleWebSocketValidateCode(conn, req)
	default:
		h.sendWebSocketError(conn, req.ID, 404, "Method not found: "+req.Method)
	}
}

// handleWebSocketGetRules WebSocket経由でgetRulesを処理
func (h *MCPHandler) handleWebSocketGetRules(conn *websocket.Conn, req domain.MCPRequest) {
	var params domain.MCPRuleRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendWebSocketError(conn, req.ID, 400, "Invalid parameters")
		return
	}

	if params.ProjectID == "" {
		h.sendWebSocketError(conn, req.ID, 400, "Project ID is required")
		return
	}

	// プロジェクトルールを取得
	projectRules, err := h.ruleUseCase.GetProjectRules(params.ProjectID)
	if err != nil {
		h.sendWebSocketError(conn, req.ID, 500, "Failed to get project rules: "+err.Error())
		return
	}

	// 言語が指定されている場合はグローバルルールを取得
	var globalRules []domain.GlobalRule
	if params.Language != "" {
		globalRulesPtr, err := h.globalRuleUseCase.GetGlobalRules(params.Language)
		if err == nil {
			// ポインタスライスを値スライスに変換
			globalRules = make([]domain.GlobalRule, len(globalRulesPtr))
			for i, gr := range globalRulesPtr {
				globalRules[i] = *gr
			}
		}
	}

	// プロジェクトルールとグローバルルールを結合
	appliedRules := make([]domain.Rule, 0, len(projectRules.Rules)+len(globalRules))
	appliedRules = append(appliedRules, projectRules.Rules...)

	// グローバルルールをプロジェクトルール形式に変換
	for _, gr := range globalRules {
		rule := domain.Rule{
			RuleID:      gr.RuleID,
			Name:        gr.Name,
			Description: gr.Description,
			Type:        gr.Type,
			Severity:    gr.Severity,
			Pattern:     gr.Pattern,
			Message:     gr.Message,
			IsActive:    gr.IsActive,
		}
		appliedRules = append(appliedRules, rule)
	}

	response := domain.MCPRuleResponse{
		ProjectID:    params.ProjectID,
		Language:     params.Language,
		Rules:        projectRules.Rules,
		GlobalRules:  globalRules,
		AppliedRules: appliedRules,
	}

	h.sendWebSocketResponse(conn, req.ID, response)
}

// handleWebSocketValidateCode WebSocket経由でvalidateCodeを処理
func (h *MCPHandler) handleWebSocketValidateCode(conn *websocket.Conn, req domain.MCPRequest) {
	var params domain.MCPValidationRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendWebSocketError(conn, req.ID, 400, "Invalid parameters")
		return
	}

	if params.ProjectID == "" || params.Code == "" {
		h.sendWebSocketError(conn, req.ID, 400, "Project ID and code are required")
		return
	}

	// プロジェクトルールに対してコードを検証
	validationResult, err := h.ruleUseCase.ValidateCode(params.ProjectID, params.Code)
	if err != nil {
		h.sendWebSocketError(conn, req.ID, 500, "Failed to validate code: "+err.Error())
		return
	}

	// 検証結果をMCP形式に変換
	issues := make([]domain.ValidationIssue, 0)

	// エラーを検証問題に変換
	for _, errorMsg := range validationResult.Errors {
		issue := domain.ValidationIssue{
			RuleID:   "validation-error",
			RuleName: "Code Validation Error",
			Severity: "error",
			Message:  errorMsg,
		}
		issues = append(issues, issue)
	}

	// 警告を検証問題に変換
	for _, warningMsg := range validationResult.Warnings {
		issue := domain.ValidationIssue{
			RuleID:   "validation-warning",
			RuleName: "Code Validation Warning",
			Severity: "warning",
			Message:  warningMsg,
		}
		issues = append(issues, issue)
	}

	// コンテキスト用に適用されたルールを取得
	projectRules, err := h.ruleUseCase.GetProjectRules(params.ProjectID)
	if err != nil {
		projectRules = &domain.ProjectRules{Rules: []domain.Rule{}}
	}

	response := domain.MCPValidationResponse{
		IsValid: validationResult.Valid,
		Issues:  issues,
		Rules:   projectRules.Rules,
	}

	h.sendWebSocketResponse(conn, req.ID, response)
}

// sendWebSocketResponse WebSocket経由でレスポンスを送信
func (h *MCPHandler) sendWebSocketResponse(conn *websocket.Conn, id string, result interface{}) {
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

// sendWebSocketError WebSocket経由でエラーを送信
func (h *MCPHandler) sendWebSocketError(conn *websocket.Conn, id string, code int, message string) {
	response := domain.MCPResponse{
		ID: id,
		Error: &domain.MCPError{
			Code:    code,
			Message: message,
		},
	}

	conn.WriteJSON(response)
}
