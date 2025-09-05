package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opm008077/RuleMCPServer/internal/domain"
	"github.com/opm008077/RuleMCPServer/internal/infrastructure/database"
	"github.com/opm008077/RuleMCPServer/internal/interface/handler"
	"github.com/opm008077/RuleMCPServer/internal/usecase"
	"github.com/opm008077/RuleMCPServer/pkg/config"
)

func main() {
	// 設定の読み込み
	cfg := config.LoadConfig()

	// データベース接続を試行
	var projectRepo domain.ProjectRepository
	var ruleRepo domain.RuleRepository
	var globalRuleRepo domain.GlobalRuleRepository
	var userRepo domain.UserRepository

	db, err := database.NewPostgresDatabase(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Printf("Falling back to JSON file mode")
	} else {
		defer db.Close()
		log.Printf("Successfully connected to database")

		// リポジトリの初期化
		projectRepo = db
		ruleRepo = database.NewPostgresRuleRepository(db.DB)
		globalRuleRepo = database.NewPostgresGlobalRuleRepository(db.DB)
		userRepo = database.NewPostgresUserRepository(db.DB)
	}

	// 本番環境ではGinのデバッグモードを無効化
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORSミドルウェアを追加
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// ヘルスチェック
	healthHandler := handler.NewHealthHandler()
	r.GET("/api/v1/health", healthHandler.HealthCheck)

	// 認証エンドポイント（データベース接続なしでも利用可能）
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-in-production"
	}
	authHandler := handler.NewAuthHandler(jwtSecret)
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.GET("/validate", authHandler.ValidateToken)
		auth.POST("/change-password", authHandler.ChangePassword)
		auth.GET("/pending-users", authHandler.GetPendingUsers)
		auth.POST("/approve-user", authHandler.ApproveUser)
	}

	// 管理者用APIエンドポイント
	var adminHandler *handler.AdminHandler
	if projectRepo != nil {
		adminHandler = handler.NewAdminHandler(userRepo, projectRepo, ruleRepo, globalRuleRepo)
	} else {
		// データベースが利用できない場合はモックハンドラーを使用
		adminHandler = handler.NewAdminHandler(nil, nil, nil, nil)
	}
	admin := r.Group("/api/v1/admin")
	{
		admin.GET("/stats", adminHandler.GetStats)
		admin.GET("/users", adminHandler.GetUsers)
		admin.POST("/users", adminHandler.CreateUser)
		admin.PUT("/users/:id", adminHandler.UpdateUser)
		admin.DELETE("/users/:id", adminHandler.DeleteUser)
		admin.GET("/api-keys", adminHandler.GetApiKeys)
		admin.POST("/api-keys", adminHandler.GenerateApiKey)
		admin.DELETE("/api-keys/:id", adminHandler.DeleteApiKey)
		admin.GET("/mcp-stats", adminHandler.GetMcpStats)
		admin.GET("/system-logs", adminHandler.GetSystemLogs)
	}

	// データベースが利用可能な場合のみ、管理用エンドポイントを有効化
	if projectRepo != nil {
		// ユースケースの初期化
		projectUseCase := usecase.NewProjectUseCase(projectRepo)
		ruleUseCase := usecase.NewRuleUseCase(ruleRepo, globalRuleRepo, projectRepo)
		globalRuleUseCase := usecase.NewGlobalRuleUseCase(globalRuleRepo)

		// ハンドラーの初期化
		projectHandler := handler.NewProjectHandler(projectUseCase)
		ruleHandler := handler.NewRuleHandler(ruleUseCase)
		globalRuleHandler := handler.NewGlobalRuleHandler(globalRuleUseCase)

		// APIルートの設定
		api := r.Group("/api/v1")
		{
			// プロジェクト管理
			api.GET("/projects", projectHandler.GetProjects)
			api.POST("/projects", projectHandler.CreateProject)
			api.PUT("/projects/:project_id", projectHandler.UpdateProject)
			api.DELETE("/projects/:project_id", projectHandler.DeleteProject)

			// ルール管理
			api.GET("/rules", ruleHandler.GetRules)
			api.POST("/rules", ruleHandler.CreateRule)
			api.DELETE("/rules/:project_id/:rule_id", ruleHandler.DeleteRule)
			api.POST("/rules/validate", ruleHandler.ValidateCode)
			api.POST("/rules/export", ruleHandler.ExportRules)
			api.POST("/rules/import", ruleHandler.ImportRules)

			// グローバルルール管理
			api.GET("/global-rules/:language", globalRuleHandler.GetGlobalRules)
			api.POST("/global-rules", globalRuleHandler.CreateGlobalRule)
			api.DELETE("/global-rules/:language/:rule_id", globalRuleHandler.DeleteGlobalRule)
			api.GET("/languages", globalRuleHandler.GetLanguages)
			api.POST("/global-rules/export", globalRuleHandler.ExportGlobalRules)
			api.POST("/global-rules/import", globalRuleHandler.ImportGlobalRules)
		}
	}

	// MCP プロトコルエンドポイント（データベース接続なしでも利用可能）
	if projectRepo != nil {
		// データベースが利用可能な場合は完全版のMCPハンドラーを使用
		ruleUseCase := usecase.NewRuleUseCase(ruleRepo, globalRuleRepo, projectRepo)
		globalRuleUseCase := usecase.NewGlobalRuleUseCase(globalRuleRepo)
		projectDetector := usecase.NewProjectDetector(projectRepo, ruleRepo)

		mcpHandler := handler.NewMCPHandler(ruleUseCase, globalRuleUseCase, projectDetector)

		mcp := r.Group("/mcp")
		{
			mcp.POST("/request", mcpHandler.HandleMCPRequest)
			mcp.GET("/ws", mcpHandler.HandleWebSocket)
		}
	} else {
		// データベースが利用できない場合は簡易版のMCPハンドラーを使用
		simpleMCPHandler := handler.NewSimpleMCPHandler()

		mcp := r.Group("/mcp")
		{
			mcp.POST("/request", simpleMCPHandler.HandleMCPRequest)
			mcp.GET("/ws", simpleMCPHandler.HandleWebSocket)
		}
	}

	log.Printf("Rule MCP Server starting on %s", cfg.GetAddress())
	log.Printf("Environment: %s, Log Level: %s", cfg.Environment, cfg.LogLevel)
	if projectRepo != nil {
		log.Printf("Database: Connected")
	} else {
		log.Printf("Database: JSON file mode")
	}
	log.Printf("MCP Endpoints: /mcp/request, /mcp/ws")

	if err := r.Run(cfg.GetAddress()); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
