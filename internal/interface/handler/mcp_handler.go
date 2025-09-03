package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/opm008077/RuleMCPServer/internal/domain"
	"github.com/opm008077/RuleMCPServer/internal/usecase"
)

type MCPHandler struct {
	ruleUseCase       *usecase.RuleUseCase
	globalRuleUseCase *usecase.GlobalRuleUseCase
}

func NewMCPHandler(ruleUseCase *usecase.RuleUseCase, globalRuleUseCase *usecase.GlobalRuleUseCase) *MCPHandler {
	return &MCPHandler{
		ruleUseCase:       ruleUseCase,
		globalRuleUseCase: globalRuleUseCase,
	}
}

// HandleMCPRequest handles MCP protocol requests
func (h *MCPHandler) HandleMCPRequest(c *gin.Context) {
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

// handleGetRules handles the getRules MCP method
func (h *MCPHandler) handleGetRules(c *gin.Context, req domain.MCPRequest) {
	var params domain.MCPRuleRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendMCPError(c, req.ID, 400, "Invalid parameters")
		return
	}

	if params.ProjectID == "" {
		h.sendMCPError(c, req.ID, 400, "Project ID is required")
		return
	}

	// Get project rules
	projectRules, err := h.ruleUseCase.GetProjectRules(params.ProjectID)
	if err != nil {
		h.sendMCPError(c, req.ID, 500, "Failed to get project rules: "+err.Error())
		return
	}

	// Get global rules if language is specified
	var globalRules []domain.GlobalRule
	if params.Language != "" {
		globalRulesPtr, err := h.globalRuleUseCase.GetGlobalRules(params.Language)
		if err != nil {
			// Log error but continue without global rules
			c.Error(err)
		} else {
			// Convert pointer slice to value slice
			globalRules = make([]domain.GlobalRule, len(globalRulesPtr))
			for i, gr := range globalRulesPtr {
				globalRules[i] = *gr
			}
		}
	}

	// Combine project rules and global rules
	appliedRules := make([]domain.Rule, 0, len(projectRules.Rules)+len(globalRules))
	appliedRules = append(appliedRules, projectRules.Rules...)

	// Convert global rules to project rules format
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

// handleValidateCode handles the validateCode MCP method
func (h *MCPHandler) handleValidateCode(c *gin.Context, req domain.MCPRequest) {
	var params domain.MCPValidationRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendMCPError(c, req.ID, 400, "Invalid parameters")
		return
	}

	if params.ProjectID == "" || params.Code == "" {
		h.sendMCPError(c, req.ID, 400, "Project ID and code are required")
		return
	}

	// Validate code against project rules
	validationResult, err := h.ruleUseCase.ValidateCode(params.ProjectID, params.Code)
	if err != nil {
		h.sendMCPError(c, req.ID, 500, "Failed to validate code: "+err.Error())
		return
	}

	// Convert validation result to MCP format
	issues := make([]domain.ValidationIssue, 0)

	// Convert errors to validation issues
	for _, errorMsg := range validationResult.Errors {
		issue := domain.ValidationIssue{
			RuleID:   "validation-error",
			RuleName: "Code Validation Error",
			Severity: "error",
			Message:  errorMsg,
		}
		issues = append(issues, issue)
	}

	// Convert warnings to validation issues
	for _, warningMsg := range validationResult.Warnings {
		issue := domain.ValidationIssue{
			RuleID:   "validation-warning",
			RuleName: "Code Validation Warning",
			Severity: "warning",
			Message:  warningMsg,
		}
		issues = append(issues, issue)
	}

	// Get applied rules for context
	projectRules, err := h.ruleUseCase.GetProjectRules(params.ProjectID)
	if err != nil {
		c.Error(err)
		projectRules = &domain.ProjectRules{Rules: []domain.Rule{}}
	}

	response := domain.MCPValidationResponse{
		IsValid: validationResult.Valid,
		Issues:  issues,
		Rules:   projectRules.Rules,
	}

	h.sendMCPResponse(c, req.ID, response)
}

// handleGetProjectInfo handles the getProjectInfo MCP method
func (h *MCPHandler) handleGetProjectInfo(c *gin.Context, req domain.MCPRequest) {
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

	// Get project information - this method needs to be implemented in usecase
	// For now, return an error
	h.sendMCPError(c, req.ID, 501, "getProjectInfo not yet implemented")
}

// sendMCPResponse sends a successful MCP response
func (h *MCPHandler) sendMCPResponse(c *gin.Context, id string, result interface{}) {
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

// sendMCPError sends an MCP error response
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

// HandleWebSocket handles WebSocket connections for real-time MCP communication
func (h *MCPHandler) HandleWebSocket(c *gin.Context) {
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

// handleWebSocketGetRules handles getRules over WebSocket
func (h *MCPHandler) handleWebSocketGetRules(conn *websocket.Conn, req domain.MCPRequest) {
	var params domain.MCPRuleRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendWebSocketError(conn, req.ID, 400, "Invalid parameters")
		return
	}

	// Implementation similar to HTTP handler
	// ... (same logic as handleGetRules)
}

// handleWebSocketValidateCode handles validateCode over WebSocket
func (h *MCPHandler) handleWebSocketValidateCode(conn *websocket.Conn, req domain.MCPRequest) {
	var params domain.MCPValidationRequest
	if err := json.Unmarshal(req.Params, &params); err != nil {
		h.sendWebSocketError(conn, req.ID, 400, "Invalid parameters")
		return
	}

	// Implementation similar to HTTP handler
	// ... (same logic as handleValidateCode)
}

// sendWebSocketResponse sends a response over WebSocket
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

// sendWebSocketError sends an error over WebSocket
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
