package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/AkitoSakurabaCreator/RuleMCPServer/internal/domain"
)

type AdminHandler struct {
	userRepo       domain.UserRepository
	projectRepo    domain.ProjectRepository
	ruleRepo       domain.RuleRepository
	globalRuleRepo domain.GlobalRuleRepository
}

type AdminStats struct {
	TotalUsers     int    `json:"totalUsers"`
	TotalProjects  int    `json:"totalProjects"`
	TotalRules     int    `json:"totalRules"`
	ActiveApiKeys  int    `json:"activeApiKeys"`
	McpRequests    int    `json:"mcpRequests"`
	ActiveSessions int    `json:"activeSessions"`
	SystemLoad     string `json:"systemLoad"`
}

type AdminUser struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullName"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"isActive"`
	LastLogin time.Time `json:"lastLogin"`
}

type AdminApiKey struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	AccessLevel string    `json:"accessLevel"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	LastUsed    time.Time `json:"lastUsed"`
}

type McpStats struct {
	Method   string `json:"method"`
	Count    int    `json:"count"`
	LastUsed string `json:"lastUsed"`
	Status   string `json:"status"`
}

type SystemLog struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

func NewAdminHandler(userRepo domain.UserRepository, projectRepo domain.ProjectRepository, ruleRepo domain.RuleRepository, globalRuleRepo domain.GlobalRuleRepository) *AdminHandler {
	return &AdminHandler{
		userRepo:       userRepo,
		projectRepo:    projectRepo,
		ruleRepo:       ruleRepo,
		globalRuleRepo: globalRuleRepo,
	}
}

// GetStats 管理者統計データを取得
func (h *AdminHandler) GetStats(c *gin.Context) {
	// データベースから実際のデータを取得
	users, err := h.userRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	projects, err := h.projectRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get projects"})
		return
	}

	// 全プロジェクトのルール数を計算
	totalRules := 0
	for _, project := range projects {
		rules, err := h.ruleRepo.GetByProjectID(project.ProjectID)
		if err != nil {
			continue // エラーが発生しても他のプロジェクトは処理を続行
		}
		totalRules += len(rules)
	}

	// アクティブユーザー数を計算
	activeUsers := 0
	for _, user := range users {
		if user.IsActive {
			activeUsers++
		}
	}

	stats := AdminStats{
		TotalUsers:     len(users),
		TotalProjects:  len(projects),
		TotalRules:     totalRules,
		ActiveApiKeys:  3, // 固定値（APIキーテーブルから取得する場合は更新）
		McpRequests:    0, // 固定値（MCP統計テーブルから取得する場合は更新）
		ActiveSessions: activeUsers,
		SystemLoad:     "15%", // 固定値（システム監視から取得する場合は更新）
	}

	c.JSON(http.StatusOK, stats)
}

// GetUsers 管理者用ユーザー一覧を取得
func (h *AdminHandler) GetUsers(c *gin.Context) {
	// モックデータ（実際の実装ではデータベースから取得）
	users := []AdminUser{
		{
			ID:        1,
			Username:  "admin",
			Email:     "admin@rulemcp.com",
			FullName:  "System Administrator",
			Role:      "admin",
			IsActive:  true,
			LastLogin: time.Now().Add(-time.Hour),
		},
		{
			ID:        2,
			Username:  "user1",
			Email:     "user1@example.com",
			FullName:  "User One",
			Role:      "user",
			IsActive:  true,
			LastLogin: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        3,
			Username:  "user2",
			Email:     "user2@example.com",
			FullName:  "User Two",
			Role:      "user",
			IsActive:  false,
			LastLogin: time.Now().Add(-24 * time.Hour),
		},
	}

	c.JSON(http.StatusOK, users)
}

// GetApiKeys 管理者用APIキー一覧を取得
func (h *AdminHandler) GetApiKeys(c *gin.Context) {
	// モックデータ（実際の実装ではデータベースから取得）
	apiKeys := []AdminApiKey{
		{
			ID:          1,
			Name:        "Admin API Key",
			Key:         "admin_key_123",
			AccessLevel: "admin",
			Status:      "active",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
			LastUsed:    time.Now().Add(-time.Hour),
		},
		{
			ID:          2,
			Name:        "User API Key",
			Key:         "user_key_456",
			AccessLevel: "user",
			Status:      "expired",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
			LastUsed:    time.Now().Add(-2 * time.Hour),
		},
	}

	c.JSON(http.StatusOK, apiKeys)
}

// GetMcpStats MCP統計データを取得
func (h *AdminHandler) GetMcpStats(c *gin.Context) {
	// モックデータ（実際の実装ではデータベースから取得）
	stats := []McpStats{
		{
			Method:   "getRules",
			Count:    1234,
			LastUsed: "2分前",
			Status:   "正常",
		},
		{
			Method:   "validateCode",
			Count:    567,
			LastUsed: "5分前",
			Status:   "正常",
		},
	}

	c.JSON(http.StatusOK, stats)
}

// GetSystemLogs システムログを取得
func (h *AdminHandler) GetSystemLogs(c *gin.Context) {
	// モックデータ（実際の実装ではデータベースから取得）
	logs := []SystemLog{
		{
			Timestamp: time.Now().Add(-5 * time.Minute),
			Level:     "INFO",
			Message:   "User 'admin' logged in successfully",
		},
		{
			Timestamp: time.Now().Add(-10 * time.Minute),
			Level:     "WARN",
			Message:   "API key 'user_key_456' expired",
		},
		{
			Timestamp: time.Now().Add(-15 * time.Minute),
			Level:     "INFO",
			Message:   "MCP request 'getRules' processed in 23ms",
		},
		{
			Timestamp: time.Now().Add(-20 * time.Minute),
			Level:     "ERROR",
			Message:   "Database connection timeout",
		},
	}

	c.JSON(http.StatusOK, logs)
}
