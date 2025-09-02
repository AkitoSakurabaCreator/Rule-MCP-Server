package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var db *Database

func main() {
	config := LoadConfig()

	// データベース接続を試行
	var err error
	db, err = NewDatabase()
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Printf("Falling back to JSON file mode")
		db = nil
	} else {
		defer db.Close()
		log.Printf("Successfully connected to database")
	}

	// 本番環境ではGinのデバッグモードを無効化
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/health", healthCheck)
		api.GET("/rules", getRules)
		api.POST("/rules/validate", validateRules)

		// データベース管理用のエンドポイント
		if db != nil {
			api.GET("/projects", getProjects)
			api.POST("/projects", createProject)
			api.POST("/rules", createRule)
			api.DELETE("/rules/:project_id/:rule_id", deleteRule)

			// グローバルルール管理
			api.GET("/global-rules/:language", getGlobalRules)
			api.POST("/global-rules", createGlobalRule)
			api.DELETE("/global-rules/:language/:rule_id", deleteGlobalRule)
			api.GET("/languages", getLanguages)
		}
	}

	log.Printf("Rule MCP Server starting on %s", config.GetAddress())
	log.Printf("Environment: %s, Log Level: %s", config.Environment, config.LogLevel)
	if db != nil {
		log.Printf("Database: Connected")
	} else {
		log.Printf("Database: JSON file mode")
	}

	if err := r.Run(config.GetAddress()); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "rule-mcp-server",
	})
}

func getRules(c *gin.Context) {
	projectID := c.Query("project")
	if projectID == "" {
		c.JSON(400, gin.H{"error": "project parameter is required"})
		return
	}

	rules, err := loadRules(projectID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to load rules"})
		return
	}

	c.JSON(200, rules)
}

func validateRules(c *gin.Context) {
	var req struct {
		ProjectID string `json:"project_id" binding:"required"`
		Code      string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	validation := validateCode(req.ProjectID, req.Code)
	c.JSON(200, validation)
}
