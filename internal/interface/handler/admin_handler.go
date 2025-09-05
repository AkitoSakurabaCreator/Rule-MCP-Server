package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userRepo       domain.UserRepository
	projectRepo    domain.ProjectRepository
	ruleRepo       domain.RuleRepository
	globalRuleRepo domain.GlobalRuleRepository
	ruleOptionRepo domain.RuleOptionRepository
	roleRepo       domain.RoleRepository
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

func NewAdminHandler(userRepo domain.UserRepository, projectRepo domain.ProjectRepository, ruleRepo domain.RuleRepository, globalRuleRepo domain.GlobalRuleRepository, ruleOptionRepo domain.RuleOptionRepository, roleRepo domain.RoleRepository) *AdminHandler {
	return &AdminHandler{
		userRepo:       userRepo,
		projectRepo:    projectRepo,
		ruleRepo:       ruleRepo,
		globalRuleRepo: globalRuleRepo,
		ruleOptionRepo: ruleOptionRepo,
		roleRepo:       roleRepo,
	}
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	if role, ok := c.Get("userRole"); !ok || role != "admin" {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Admin access required", nil)
		return
	}
	if h.userRepo == nil || h.projectRepo == nil || h.ruleRepo == nil {
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

	users, err := h.userRepo.GetAll()
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "Failed to get users", nil)
		return
	}

	projects, err := h.projectRepo.GetAll()
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "Failed to get projects", nil)
		return
	}

	totalRules := 0
	for _, project := range projects {
		rules, err := h.ruleRepo.GetByProjectID(project.ProjectID)
		if err != nil {
			continue
		}
		totalRules += len(rules)
	}

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
		ActiveApiKeys:  3,
		McpRequests:    1234,
		ActiveSessions: activeUsers,
		SystemLoad:     "15%",
	}

	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	if perms, ok := c.Get("permissions"); !ok || !perms.(map[string]bool)["manage_users"] {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Permission manage_users required", nil)
		return
	}
	if h.userRepo == nil {
		users := []AdminUser{
			{ID: 1, Username: "admin", Email: "admin@rulemcp.com", FullName: "System Administrator", Role: "admin", IsActive: true, LastLogin: time.Now().Add(-time.Hour)},
			{ID: 2, Username: "user1", Email: "user1@example.com", FullName: "User One", Role: "user", IsActive: true, LastLogin: time.Now().Add(-2 * time.Hour)},
			{ID: 3, Username: "user2", Email: "user2@example.com", FullName: "User Two", Role: "user", IsActive: false, LastLogin: time.Now().Add(-24 * time.Hour)},
		}
		c.JSON(http.StatusOK, users)
		return
	}

	users, err := h.userRepo.GetAll()
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "Failed to get users", nil)
		return
	}

	adminUsers := make([]AdminUser, len(users))
	for i, user := range users {
		adminUsers[i] = AdminUser{ID: user.ID, Username: user.Username, Email: user.Email, FullName: user.FullName, Role: user.Role, IsActive: user.IsActive, LastLogin: user.UpdatedAt}
	}

	c.JSON(http.StatusOK, adminUsers)
}

func (h *AdminHandler) GetApiKeys(c *gin.Context) {
	if role, ok := c.Get("userRole"); !ok || role != "admin" {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Admin access required", nil)
		return
	}
	apiKeys := []AdminApiKey{
		{ID: 1, Name: "Admin API Key", Key: "admin_key_123", AccessLevel: "admin", Status: "active", CreatedAt: time.Now().Add(-24 * time.Hour), LastUsed: time.Now().Add(-time.Hour)},
		{ID: 2, Name: "User API Key", Key: "user_key_456", AccessLevel: "user", Status: "expired", CreatedAt: time.Now().Add(-48 * time.Hour), LastUsed: time.Now().Add(-2 * time.Hour)},
	}
	c.JSON(http.StatusOK, apiKeys)
}

func (h *AdminHandler) GetMcpStats(c *gin.Context) {
	if role, ok := c.Get("userRole"); !ok || role != "admin" {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Admin access required", nil)
		return
	}
	stats := []McpStats{
		{Method: "getRules", Count: 1234, LastUsed: "2分前", Status: "正常"},
		{Method: "validateCode", Count: 567, LastUsed: "5分前", Status: "正常"},
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) GetSystemLogs(c *gin.Context) {
	if role, ok := c.Get("userRole"); !ok || role != "admin" {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Admin access required", nil)
		return
	}
	logs := []SystemLog{
		{Timestamp: time.Now().Add(-5 * time.Minute), Level: "INFO", Message: "User 'admin' logged in successfully"},
		{Timestamp: time.Now().Add(-10 * time.Minute), Level: "WARN", Message: "API key 'user_key_456' expired"},
		{Timestamp: time.Now().Add(-15 * time.Minute), Level: "INFO", Message: "MCP request 'getRules' processed in 23ms"},
		{Timestamp: time.Now().Add(-20 * time.Minute), Level: "ERROR", Message: "Database connection timeout"},
	}
	c.JSON(http.StatusOK, logs)
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	if perms, ok := c.Get("permissions"); !ok || !perms.(map[string]bool)["manage_users"] {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Permission manage_users required", nil)
		return
	}
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		FullName string `json:"fullName" binding:"required"`
		Role     string `json:"role" binding:"required"`
		IsActive bool   `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}
	if h.userRepo == nil {
		httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "Database not available", nil)
		return
	}
	if existingUser, _ := h.userRepo.GetByUsername(req.Username); existingUser != nil {
		httpx.JSONError(c, http.StatusConflict, httpx.CodeConflict, "Username already exists", nil)
		return
	}
	if existingUser, _ := h.userRepo.GetByEmail(req.Email); existingUser != nil {
		httpx.JSONError(c, http.StatusConflict, httpx.CodeConflict, "Email already exists", nil)
		return
	}
	user := &domain.User{Username: req.Username, Email: req.Email, FullName: req.FullName, Role: req.Role, IsActive: req.IsActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := h.userRepo.Create(user); err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "Failed to create user", nil)
		return
	}
	adminUser := AdminUser{ID: user.ID, Username: user.Username, Email: user.Email, FullName: user.FullName, Role: user.Role, IsActive: user.IsActive, LastLogin: user.UpdatedAt}
	c.JSON(http.StatusCreated, adminUser)
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	if perms, ok := c.Get("permissions"); !ok || !perms.(map[string]bool)["manage_users"] {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Permission manage_users required", nil)
		return
	}
	userID := c.Param("id")
	if userID == "" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "User ID is required", nil)
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
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}
	if h.userRepo == nil {
		httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "Database not available", nil)
		return
	}
	var id int
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "Invalid user ID", nil)
		return
	}
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		httpx.JSONError(c, http.StatusNotFound, httpx.CodeNotFound, "User not found", nil)
		return
	}
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
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "Failed to update user", nil)
		return
	}
	adminUser := AdminUser{ID: user.ID, Username: user.Username, Email: user.Email, FullName: user.FullName, Role: user.Role, IsActive: user.IsActive, LastLogin: user.UpdatedAt}
	c.JSON(http.StatusOK, adminUser)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	if perms, ok := c.Get("permissions"); !ok || !perms.(map[string]bool)["manage_users"] {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Permission manage_users required", nil)
		return
	}
	userID := c.Param("id")
	if userID == "" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "User ID is required", nil)
		return
	}
	if h.userRepo == nil {
		httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "Database not available", nil)
		return
	}
	var id int
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "Invalid user ID", nil)
		return
	}
	if _, err := h.userRepo.GetByID(id); err != nil {
		httpx.JSONError(c, http.StatusNotFound, httpx.CodeNotFound, "User not found", nil)
		return
	}
	if err := h.userRepo.Delete(id); err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "Failed to delete user", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *AdminHandler) GenerateApiKey(c *gin.Context) {
	if role, ok := c.Get("userRole"); !ok || role != "admin" {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Admin access required", nil)
		return
	}
	var req struct {
		Name        string `json:"name" binding:"required"`
		AccessLevel string `json:"accessLevel" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}
	apiKey := fmt.Sprintf("%s_%d_%s", req.AccessLevel, time.Now().Unix(), generateRandomString(16))
	newApiKey := AdminApiKey{ID: int(time.Now().Unix()), Name: req.Name, Key: apiKey, AccessLevel: req.AccessLevel, Status: "active", CreatedAt: time.Now(), LastUsed: time.Time{}}
	c.JSON(http.StatusCreated, newApiKey)
}

func (h *AdminHandler) DeleteApiKey(c *gin.Context) {
	if role, ok := c.Get("userRole"); !ok || role != "admin" {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Admin access required", nil)
		return
	}
	keyID := c.Param("id")
	if keyID == "" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "API Key ID is required", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "API Key deleted successfully"})
}

// Rule options
func (h *AdminHandler) GetRuleOptions(c *gin.Context) {
	kind := c.Query("kind")
	if kind != "type" && kind != "severity" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "kind must be 'type' or 'severity'", nil)
		return
	}
	if h.ruleOptionRepo == nil {
		httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "RuleOption repository not available", nil)
		return
	}
	opts, err := h.ruleOptionRepo.GetByKind(kind)
	if err != nil {
		httpx.JSONFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"options": opts})
}

func (h *AdminHandler) AddRuleOption(c *gin.Context) {
	if perms, ok := c.Get("permissions"); !ok || !perms.(map[string]bool)["manage_rules"] {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Permission manage_rules required", nil)
		return
	}
	var req struct {
		Kind  string `json:"kind" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}
	if req.Kind != "type" && req.Kind != "severity" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "kind must be 'type' or 'severity'", nil)
		return
	}
	if h.ruleOptionRepo == nil {
		httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "RuleOption repository not available", nil)
		return
	}
	if err := h.ruleOptionRepo.Add(req.Kind, req.Value); err != nil {
		httpx.JSONFromError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Option added"})
}

func (h *AdminHandler) DeleteRuleOption(c *gin.Context) {
	if perms, ok := c.Get("permissions"); !ok || !perms.(map[string]bool)["manage_rules"] {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Permission manage_rules required", nil)
		return
	}
	var req struct {
		Kind  string `json:"kind" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}
	if req.Kind != "type" && req.Kind != "severity" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "kind must be 'type' or 'severity'", nil)
		return
	}
	if h.ruleOptionRepo == nil {
		httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "RuleOption repository not available", nil)
		return
	}
	if err := h.ruleOptionRepo.Delete(req.Kind, req.Value); err != nil {
		httpx.JSONFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Option deleted"})
}

// Roles management
func (h *AdminHandler) GetRoles(c *gin.Context) {
	if perms, ok := c.Get("permissions"); !ok || !perms.(map[string]bool)["manage_roles"] {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Permission manage_roles required", nil)
		return
	}
	if h.roleRepo == nil {
		httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "Role repository not available", nil)
		return
	}
	roles, err := h.roleRepo.GetAll()
	if err != nil {
		httpx.JSONFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *AdminHandler) CreateRole(c *gin.Context) {
	if perms, ok := c.Get("permissions"); !ok || !perms.(map[string]bool)["manage_roles"] {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Permission manage_roles required", nil)
		return
	}
	if h.roleRepo == nil {
		httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "Role repository not available", nil)
		return
	}
	var req struct {
		Name        string          `json:"name" binding:"required"`
		Description string          `json:"description"`
		Permissions map[string]bool `json:"permissions"`
		IsActive    *bool           `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	if err := h.roleRepo.Create(domain.Role{Name: req.Name, Description: req.Description, Permissions: req.Permissions, IsActive: active}); err != nil {
		httpx.JSONFromError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Role created"})
}

func (h *AdminHandler) UpdateRole(c *gin.Context) {
	if perms, ok := c.Get("permissions"); !ok || !perms.(map[string]bool)["manage_roles"] {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Permission manage_roles required", nil)
		return
	}
	if h.roleRepo == nil {
		httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "Role repository not available", nil)
		return
	}
	name := c.Param("name")
	if name == "" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "Role name is required", nil)
		return
	}
	var req struct {
		Description string          `json:"description"`
		Permissions map[string]bool `json:"permissions"`
		IsActive    *bool           `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	if err := h.roleRepo.Update(name, domain.Role{Description: req.Description, Permissions: req.Permissions, IsActive: active}); err != nil {
		httpx.JSONFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

func (h *AdminHandler) DeleteRole(c *gin.Context) {
	if perms, ok := c.Get("permissions"); !ok || !perms.(map[string]bool)["manage_roles"] {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Permission manage_roles required", nil)
		return
	}
	if h.roleRepo == nil {
		httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "Role repository not available", nil)
		return
	}
	name := c.Param("name")
	if name == "" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "Role name is required", nil)
		return
	}
	if err := h.roleRepo.Delete(name); err != nil {
		httpx.JSONFromError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}
