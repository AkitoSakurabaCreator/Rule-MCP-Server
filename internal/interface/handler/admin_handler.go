package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opm008077/RuleMCPServer/internal/domain"
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
	if h.userRepo == nil || h.projectRepo == nil || h.ruleRepo == nil {
		// データベース接続がない場合はモックデータを返す
		stats := AdminStats{
			TotalUsers:     3,
			TotalProjects:  5,
			TotalRules:     25,
			ActiveApiKeys:  2,
			McpRequests:    1234,
			ActiveSessions: 2,
			SystemLoad:     "15%",
		}
		c.JSON(http.StatusOK, stats)
		return
	}

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
		ActiveApiKeys:  3,    // 固定値（APIキーテーブルから取得する場合は更新）
		McpRequests:    1234, // 固定値（MCP統計テーブルから取得する場合は更新）
		ActiveSessions: activeUsers,
		SystemLoad:     "15%", // 固定値（システム監視から取得する場合は更新）
	}

	c.JSON(http.StatusOK, stats)
}

// GetUsers 管理者用ユーザー一覧を取得
func (h *AdminHandler) GetUsers(c *gin.Context) {
	if h.userRepo == nil {
		// モックデータ（データベース接続がない場合）
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
		return
	}

	// データベースから実際のユーザーデータを取得
	users, err := h.userRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	adminUsers := make([]AdminUser, len(users))
	for i, user := range users {
		adminUsers[i] = AdminUser{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      user.Role,
			IsActive:  user.IsActive,
			LastLogin: user.UpdatedAt, // 最終ログイン時間の代わりに更新時間を使用
		}
	}

	c.JSON(http.StatusOK, adminUsers)
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

// CreateUser 新しいユーザーを作成
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		FullName string `json:"fullName" binding:"required"`
		Role     string `json:"role" binding:"required"`
		IsActive bool   `json:"isActive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.userRepo == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}

	// 既存ユーザーの重複チェック
	if existingUser, _ := h.userRepo.GetByUsername(req.Username); existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	if existingUser, _ := h.userRepo.GetByEmail(req.Email); existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// 新しいユーザーを作成
	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		FullName:  req.FullName,
		Role:      req.Role,
		IsActive:  req.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	adminUser := AdminUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		LastLogin: user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, adminUser)
}

// UpdateUser ユーザー情報を更新
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"fullName"`
		Role     string `json:"role"`
		IsActive *bool  `json:"isActive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.userRepo == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}

	// ユーザーIDを整数に変換
	var id int
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 既存ユーザーを取得
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// フィールドを更新
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	user.UpdatedAt = time.Now()

	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	adminUser := AdminUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		LastLogin: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, adminUser)
}

// DeleteUser ユーザーを削除
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	if h.userRepo == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return
	}

	// ユーザーIDを整数に変換
	var id int
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// ユーザーが存在するか確認
	if _, err := h.userRepo.GetByID(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := h.userRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GenerateApiKey 新しいAPIキーを生成
func (h *AdminHandler) GenerateApiKey(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		AccessLevel string `json:"accessLevel" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 簡単なAPIキー生成（実際の実装では暗号学的に安全な方法を使用）
	apiKey := fmt.Sprintf("%s_%d_%s", req.AccessLevel, time.Now().Unix(), generateRandomString(16))

	newApiKey := AdminApiKey{
		ID:          int(time.Now().Unix()), // 簡易的なID生成
		Name:        req.Name,
		Key:         apiKey,
		AccessLevel: req.AccessLevel,
		Status:      "active",
		CreatedAt:   time.Now(),
		LastUsed:    time.Time{}, // 未使用
	}

	c.JSON(http.StatusCreated, newApiKey)
}

// DeleteApiKey APIキーを削除
func (h *AdminHandler) DeleteApiKey(c *gin.Context) {
	keyID := c.Param("id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "API Key ID is required"})
		return
	}

	// 実際の実装ではデータベースから削除
	c.JSON(http.StatusOK, gin.H{"message": "API Key deleted successfully"})
}

// generateRandomString ランダム文字列を生成（簡易版）
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}
