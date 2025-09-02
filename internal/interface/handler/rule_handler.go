package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opm008077/RuleMCPServer/internal/usecase"
)

type RuleHandler struct {
	ruleUseCase *usecase.RuleUseCase
}

func NewRuleHandler(ruleUseCase *usecase.RuleUseCase) *RuleHandler {
	return &RuleHandler{
		ruleUseCase: ruleUseCase,
	}
}

func (h *RuleHandler) GetRules(c *gin.Context) {
	projectID := c.Query("project")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project parameter is required"})
		return
	}

	rules, err := h.ruleUseCase.GetProjectRules(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rules)
}

func (h *RuleHandler) CreateRule(c *gin.Context) {
	var req struct {
		ProjectID   string `json:"project_id" binding:"required"`
		RuleID      string `json:"rule_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Type        string `json:"type"`
		Severity    string `json:"severity"`
		Pattern     string `json:"pattern"`
		Message     string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ruleUseCase.CreateRule(req.ProjectID, req.RuleID, req.Name, req.Description, req.Type, req.Severity, req.Pattern, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Rule created successfully"})
}

func (h *RuleHandler) DeleteRule(c *gin.Context) {
	projectID := c.Param("project_id")
	ruleID := c.Param("rule_id")

	if projectID == "" || ruleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id and rule_id are required"})
		return
	}

	err := h.ruleUseCase.DeleteRule(projectID, ruleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rule deleted successfully"})
}

func (h *RuleHandler) ValidateCode(c *gin.Context) {
	var req struct {
		ProjectID string `json:"project_id" binding:"required"`
		Code      string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.ruleUseCase.ValidateCode(req.ProjectID, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
