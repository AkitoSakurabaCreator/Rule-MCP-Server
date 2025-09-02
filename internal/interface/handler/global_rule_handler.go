package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/AkitoSakurabaCreator/RuleMCPServer/internal/usecase"
)

type GlobalRuleHandler struct {
	globalRuleUseCase *usecase.GlobalRuleUseCase
}

func NewGlobalRuleHandler(globalRuleUseCase *usecase.GlobalRuleUseCase) *GlobalRuleHandler {
	return &GlobalRuleHandler{
		globalRuleUseCase: globalRuleUseCase,
	}
}

func (h *GlobalRuleHandler) GetGlobalRules(c *gin.Context) {
	language := c.Param("language")
	if language == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "language parameter is required"})
		return
	}

	rules, err := h.globalRuleUseCase.GetGlobalRules(language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

func (h *GlobalRuleHandler) CreateGlobalRule(c *gin.Context) {
	var req struct {
		Language    string `json:"language" binding:"required"`
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

	err := h.globalRuleUseCase.CreateGlobalRule(req.Language, req.RuleID, req.Name, req.Description, req.Type, req.Severity, req.Pattern, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Global rule created successfully"})
}

func (h *GlobalRuleHandler) DeleteGlobalRule(c *gin.Context) {
	language := c.Param("language")
	ruleID := c.Param("rule_id")

	if language == "" || ruleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "language and rule_id are required"})
		return
	}

	err := h.globalRuleUseCase.DeleteGlobalRule(language, ruleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Global rule deleted successfully"})
}

func (h *GlobalRuleHandler) GetLanguages(c *gin.Context) {
	languages, err := h.globalRuleUseCase.GetAllLanguages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"languages": languages})
}
