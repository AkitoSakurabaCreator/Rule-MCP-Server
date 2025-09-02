package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getProjects(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}

	projects, err := db.GetProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func createProject(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}

	var req struct {
		ProjectID        string `json:"project_id" binding:"required"`
		Name             string `json:"name" binding:"required"`
		Description      string `json:"description"`
		Language         string `json:"language"`
		ApplyGlobalRules bool   `json:"apply_global_rules"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// デフォルト値を設定
	if req.Language == "" {
		req.Language = "general"
	}

	err := db.CreateProject(req.ProjectID, req.Name, req.Description, req.Language, req.ApplyGlobalRules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}

func createRule(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}

	var req struct {
		ProjectID   string `json:"project_id" binding:"required"`
		RuleID      string `json:"rule_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Type        string `json:"type" binding:"required"`
		Severity    string `json:"severity" binding:"required"`
		Pattern     string `json:"pattern" binding:"required"`
		Message     string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 有効性チェック
	if req.Severity != "error" && req.Severity != "warning" && req.Severity != "info" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Severity must be 'error', 'warning', or 'info'"})
		return
	}

	err := db.CreateRule(req.ProjectID, req.RuleID, req.Name, req.Description, req.Type, req.Severity, req.Pattern, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Rule created successfully"})
}

func deleteRule(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}

	projectID := c.Param("project_id")
	ruleID := c.Param("rule_id")

	if projectID == "" || ruleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID and Rule ID are required"})
		return
	}

	err := db.DeleteRule(projectID, ruleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rule deleted successfully"})
}

func getGlobalRules(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}

	language := c.Param("language")
	if language == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language parameter is required"})
		return
	}

	rules, err := db.GetGlobalRules(language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get global rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"language": language, "rules": rules})
}

func createGlobalRule(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}

	var req struct {
		Language    string `json:"language" binding:"required"`
		RuleID      string `json:"rule_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Type        string `json:"type" binding:"required"`
		Severity    string `json:"severity" binding:"required"`
		Pattern     string `json:"pattern" binding:"required"`
		Message     string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 有効性チェック
	if req.Severity != "error" && req.Severity != "warning" && req.Severity != "info" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Severity must be 'error', 'warning', or 'info'"})
		return
	}

	err := db.CreateGlobalRule(req.Language, req.RuleID, req.Name, req.Description, req.Type, req.Severity, req.Pattern, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create global rule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Global rule created successfully"})
}

func deleteGlobalRule(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}

	language := c.Param("language")
	ruleID := c.Param("rule_id")

	if language == "" || ruleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language and Rule ID are required"})
		return
	}

	err := db.DeleteGlobalRule(language, ruleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete global rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Global rule deleted successfully"})
}

func getLanguages(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}

	languages, err := db.GetLanguages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get languages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"languages": languages})
}
