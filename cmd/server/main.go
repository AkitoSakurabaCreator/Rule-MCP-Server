package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/infrastructure/database"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/interface/handler"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/usecase"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/config"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ActiveTracker struct {
	mu   sync.Mutex
	last map[string]time.Time
}

func NewActiveTracker() *ActiveTracker     { return &ActiveTracker{last: map[string]time.Time{}} }
func (t *ActiveTracker) Touch(user string) { t.mu.Lock(); t.last[user] = time.Now(); t.mu.Unlock() }
func (t *ActiveTracker) CountSince(d time.Duration) int {
	now := time.Now()
	c := 0
	t.mu.Lock()
	for _, ts := range t.last {
		if now.Sub(ts) <= d {
			c++
		}
	}
	t.mu.Unlock()
	return c
}

// generateRandomString セキュアな乱数で英数字文字列を生成
func generateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			// フォールバック（低品質）
			b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
			continue
		}
		b[i] = letters[nBig.Int64()]
	}
	return string(b)
}

func main() {
	cfg := config.LoadConfig()

	var projectRepo domain.ProjectRepository
	var ruleRepo domain.RuleRepository
	var globalRuleRepo domain.GlobalRuleRepository
	var userRepo domain.UserRepository
	var ruleOptionRepo domain.RuleOptionRepository
	var roleRepo domain.RoleRepository
	var metricsRepo domain.MetricsRepository
	activeTracker := NewActiveTracker()

	db, err := database.NewPostgresDatabase(
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Printf("Falling back to JSON file mode")
	} else {
		defer db.Close()
		log.Printf("Successfully connected to database")

		projectRepo = db
		ruleRepo = database.NewPostgresRuleRepository(db.DB)
		globalRuleRepo = database.NewPostgresGlobalRuleRepository(db.DB)
		ruleOptionRepo = database.NewPostgresRuleOptionRepository(db.DB)
		userRepo = database.NewPostgresUserRepository(db.DB)
		roleRepo = database.NewPostgresRoleRepository(db.DB)
		metricsRepo = database.NewPostgresMetricsRepository(db.DB)
	}

	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(httpx.RecoveryJSON())
	r.Use(httpx.RequestID())

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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key-change-in-production"
	}
	r.Use(func(c *gin.Context) {
		c.Set("userRole", "public")
		c.Set("permissions", map[string]bool{"manage_users": false, "manage_rules": false, "manage_roles": false})
		auth := c.GetHeader("Authorization")
		if len(auth) > 7 && (auth[:7] == "Bearer " || auth[:7] == "bearer ") {
			tokenStr := auth[7:]
			claims := &handler.Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) { return []byte(jwtSecret), nil })
			if err == nil && token.Valid {
				if claims.Role != "" {
					c.Set("userRole", claims.Role)
				}
				// permissions lookup with fallback
				perm := map[string]bool{"manage_users": false, "manage_rules": false, "manage_roles": false}
				loaded := false
				if roleRepo != nil {
					if role, err := roleRepo.GetByName(claims.Role); err == nil && role.Permissions != nil {
						perm = role.Permissions
						loaded = true
					}
				}
				if !loaded {
					switch claims.Role {
					case "admin":
						perm["manage_users"] = true
						perm["manage_rules"] = true
						perm["manage_roles"] = true
					case "user":
						perm["manage_rules"] = true
					}
				}
				c.Set("permissions", perm)
				// track active session
				activeTracker.Touch(claims.Username)
			}
		}
		c.Next()
	})

	healthHandler := handler.NewHealthHandler()
	r.GET("/api/v1/health", healthHandler.HealthCheck)

	authHandler := handler.NewAuthHandler(jwtSecret, userRepo, roleRepo)
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.GET("/validate", authHandler.ValidateToken)
		auth.GET("/me", authHandler.Me)
		auth.POST("/change-password", authHandler.ChangePassword)
		auth.GET("/pending-users", authHandler.GetPendingUsers)
		auth.POST("/approve-user", authHandler.ApproveUser)
	}

	var adminHandler *handler.AdminHandler
	if projectRepo != nil {
		adminHandler = handler.NewAdminHandler(userRepo, projectRepo, ruleRepo, globalRuleRepo, ruleOptionRepo, roleRepo)
	} else {
		adminHandler = handler.NewAdminHandler(nil, nil, nil, nil, nil, nil)
	}
	admin := r.Group("/api/v1/admin")
	{
		admin.GET("/stats", func(c *gin.Context) {
			// 実データ化: totalUsers/totalProjects/totalRules + mcpRequests + activeSessions + systemLoad
			if projectRepo == nil || ruleRepo == nil {
				adminHandler.GetStats(c)
				return
			}
			// reuse existing but override some fields
			// call original to compute counts
			// For simplicity, call handler then mutate response is complex; instead, implement small inline
			users, _ := userRepo.GetAll()
			projects, _ := projectRepo.GetAll()
			totalRules := 0
			for _, p := range projects {
				rs, _ := ruleRepo.GetByProjectID(p.ProjectID)
				totalRules += len(rs)
			}
			mcpCount := 0
			if metricsRepo != nil {
				mcpCount, _ = metricsRepo.GetMCPRequestsCountLast24h()
			}
			active := activeTracker.CountSince(10 * time.Minute)
			// System load (approximate percent): loadavg1 / NumCPU * 100
			sysLoad := ""
			if b, err := os.ReadFile("/proc/loadavg"); err == nil {
				parts := strings.Fields(string(b))
				if len(parts) > 0 {
					if load1, err := strconv.ParseFloat(parts[0], 64); err == nil {
						cores := float64(runtime.NumCPU())
						pct := int((load1/cores)*100 + 0.5)
						sysLoad = strings.TrimSpace((func(v int) string { return fmt.Sprintf("%d%%", v) })(pct))
					}
				}
			}
			c.JSON(http.StatusOK, handler.AdminStats{TotalUsers: len(users), TotalProjects: len(projects), TotalRules: totalRules, ActiveApiKeys: 0, McpRequests: mcpCount, ActiveSessions: active, SystemLoad: sysLoad})
		})
		admin.GET("/users", adminHandler.GetUsers)
		admin.POST("/users", adminHandler.CreateUser)
		admin.PUT("/users/:id", adminHandler.UpdateUser)
		admin.DELETE("/users/:id", adminHandler.DeleteUser)
		// Real API keys from DB
		admin.GET("/api-keys", func(c *gin.Context) {
			if db == nil || db.DB == nil {
				c.JSON(http.StatusOK, []gin.H{})
				return
			}
			rows, err := db.DB.Query(`SELECT id, name, key_hash, access_level, is_active, created_at, updated_at FROM api_keys ORDER BY created_at DESC LIMIT 100`)
			if err != nil {
				httpx.JSONFromError(c, err)
				return
			}
			defer rows.Close()
			var keys []gin.H
			for rows.Next() {
				var id int
				var name, keyHash, accessLevel string
				var isActive bool
				var createdAt, updatedAt time.Time
				if err := rows.Scan(&id, &name, &keyHash, &accessLevel, &isActive, &createdAt, &updatedAt); err != nil {
					continue
				}
				status := "active"
				if !isActive {
					status = "inactive"
				}
				keys = append(keys, gin.H{"id": id, "name": name, "key": keyHash, "accessLevel": accessLevel, "status": status, "createdAt": createdAt.Format(time.RFC3339), "lastUsed": updatedAt.Format(time.RFC3339)})
			}
			c.JSON(http.StatusOK, keys)
		})
		admin.POST("/api-keys", func(c *gin.Context) {
			var req struct {
				Name        string `json:"name"`
				AccessLevel string `json:"accessLevel"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
				return
			}
			if db == nil || db.DB == nil {
				httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "Database not available", nil)
				return
			}
			apiKey := fmt.Sprintf("%s_%d_%s", req.AccessLevel, time.Now().Unix(), generateRandomString(16))
			_, err := db.DB.Exec(`INSERT INTO api_keys (key_hash, name, access_level, is_active, created_by, created_at, updated_at) VALUES ($1, $2, $3, true, $4, NOW(), NOW())`, apiKey, req.Name, req.AccessLevel, "admin")
			if err != nil {
				httpx.JSONFromError(c, err)
				return
			}
			c.JSON(http.StatusCreated, gin.H{"id": time.Now().Unix(), "name": req.Name, "key": apiKey, "accessLevel": req.AccessLevel, "status": "active", "createdAt": time.Now().Format(time.RFC3339), "lastUsed": ""})
		})
		admin.DELETE("/api-keys/:id", func(c *gin.Context) {
			if db == nil || db.DB == nil {
				httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "Database not available", nil)
				return
			}
			id := c.Param("id")
			_, err := db.DB.Exec(`DELETE FROM api_keys WHERE id = $1`, id)
			if err != nil {
				httpx.JSONFromError(c, err)
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "API Key deleted successfully"})
		})
		// Update API key (name/description/is_active)
		admin.PUT("/api-keys/:id", func(c *gin.Context) {
			if db == nil || db.DB == nil {
				httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "Database not available", nil)
				return
			}
			id := c.Param("id")
			var req struct {
				Name        *string `json:"name"`
				Description *string `json:"description"`
				IsActive    *bool   `json:"isActive"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
				return
			}
			q := "UPDATE api_keys SET updated_at = NOW()"
			args := []interface{}{}
			idx := 1
			if req.Name != nil {
				q += ", name = $" + fmt.Sprint(idx)
				args = append(args, *req.Name)
				idx++
			}
			if req.Description != nil {
				q += ", description = $" + fmt.Sprint(idx)
				args = append(args, *req.Description)
				idx++
			}
			if req.IsActive != nil {
				q += ", is_active = $" + fmt.Sprint(idx)
				args = append(args, *req.IsActive)
				idx++
			}
			q += " WHERE id = $" + fmt.Sprint(idx)
			args = append(args, id)
			if _, err := db.DB.Exec(q, args...); err != nil {
				httpx.JSONFromError(c, err)
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "API Key updated"})
		})
		// Settings: simple key-value store
		admin.GET("/settings", func(c *gin.Context) {
			if db == nil || db.DB == nil {
				c.JSON(http.StatusOK, gin.H{"defaultAccessLevel": "public", "requestsPerMinute": 100})
				return
			}
			rows, err := db.DB.Query(`SELECT key, value FROM settings`)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"defaultAccessLevel": "public", "requestsPerMinute": 100})
				return
			}
			defer rows.Close()
			conf := map[string]interface{}{}
			for rows.Next() {
				var k string
				var v string
				if err := rows.Scan(&k, &v); err == nil {
					conf[k] = v
				}
			}
			c.JSON(http.StatusOK, conf)
		})
		admin.PUT("/settings", func(c *gin.Context) {
			if db == nil || db.DB == nil {
				httpx.JSONError(c, http.StatusServiceUnavailable, httpx.CodeInternal, "Database not available", nil)
				return
			}
			var payload map[string]interface{}
			if err := c.ShouldBindJSON(&payload); err != nil {
				httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
				return
			}
			_, _ = db.DB.Exec(`CREATE TABLE IF NOT EXISTS settings (key TEXT PRIMARY KEY, value TEXT, updated_at TIMESTAMP DEFAULT NOW())`)
			for k, v := range payload {
				_, _ = db.DB.Exec(`INSERT INTO settings(key, value, updated_at) VALUES ($1, $2, NOW()) ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = NOW()`, k, fmt.Sprint(v))
			}
			c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
		})
		admin.GET("/mcp-stats", func(c *gin.Context) {
			if metricsRepo == nil {
				c.JSON(http.StatusOK, []domain.MCPMethodStat{})
				return
			}
			s, err := metricsRepo.GetMCPStatsLast24h()
			if err != nil {
				c.JSON(http.StatusOK, []domain.MCPMethodStat{})
				return
			}
			c.JSON(http.StatusOK, s)
		})
		// MCP performance aggregates (avg, success rate, error rate, p95) over last 24h
		admin.GET("/mcp-performance", func(c *gin.Context) {
			if db == nil || db.DB == nil {
				c.JSON(http.StatusOK, gin.H{"avgMs": 0, "successRate": 0, "errorRate": 0, "p95Ms": 0})
				return
			}
			row := db.DB.QueryRow(`
				SELECT 
					COALESCE(AVG(duration_ms)::INT,0) AS avg_ms,
					COALESCE(SUM(CASE WHEN status='ok' THEN 1 ELSE 0 END),0) AS ok_cnt,
					COALESCE(SUM(CASE WHEN status<>'ok' THEN 1 ELSE 0 END),0) AS err_cnt,
					COALESCE(PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY duration_ms),0) AS p95
				FROM mcp_requests 
				WHERE created_at > NOW() - INTERVAL '24 hours'
			`)
			var avgMs int
			var okCnt, errCnt int
			var p95 float64
			if err := row.Scan(&avgMs, &okCnt, &errCnt, &p95); err != nil {
				c.JSON(http.StatusOK, gin.H{"avgMs": 0, "successRate": 0, "errorRate": 0, "p95Ms": 0})
				return
			}
			total := okCnt + errCnt
			var succRate, errRate float64
			if total > 0 {
				succRate = float64(okCnt) * 100.0 / float64(total)
				errRate = float64(errCnt) * 100.0 / float64(total)
			}
			c.JSON(http.StatusOK, gin.H{"avgMs": avgMs, "successRate": succRate, "errorRate": errRate, "p95Ms": int(p95 + 0.5)})
		})
		admin.GET("/system-logs", func(c *gin.Context) {
			// Return recent MCP request logs as system logs
			if db == nil || db.DB == nil {
				c.JSON(http.StatusOK, []gin.H{})
				return
			}
			rows, err := db.DB.Query(`SELECT created_at, method, status, duration_ms FROM mcp_requests ORDER BY created_at DESC LIMIT 100`)
			if err != nil {
				c.JSON(http.StatusOK, []gin.H{})
				return
			}
			defer rows.Close()
			var logs []gin.H
			for rows.Next() {
				var ts time.Time
				var method, status string
				var dur int
				if err := rows.Scan(&ts, &method, &status, &dur); err != nil {
					continue
				}
				level := "INFO"
				if status != "ok" {
					level = "ERROR"
				} else if dur > 500 {
					level = "WARN"
				}
				msg := fmt.Sprintf("MCP %s %s in %dms", method, status, dur)
				logs = append(logs, gin.H{"timestamp": ts.Format(time.RFC3339), "level": level, "message": msg})
			}
			c.JSON(http.StatusOK, logs)
		})
		admin.GET("/rule-options", adminHandler.GetRuleOptions)
		admin.POST("/rule-options", adminHandler.AddRuleOption)
		admin.DELETE("/rule-options", adminHandler.DeleteRuleOption)
		admin.GET("/roles", adminHandler.GetRoles)
		admin.POST("/roles", adminHandler.CreateRole)
		admin.PUT("/roles/:name", adminHandler.UpdateRole)
		admin.DELETE("/roles/:name", adminHandler.DeleteRole)
	}

	if projectRepo != nil {
		projectUseCase := usecase.NewProjectUseCase(projectRepo)
		ruleUseCase := usecase.NewRuleUseCase(ruleRepo, globalRuleRepo, projectRepo)
		globalRuleUseCase := usecase.NewGlobalRuleUseCase(globalRuleRepo)
		projectHandler := handler.NewProjectHandler(projectUseCase)
		ruleHandler := handler.NewRuleHandler(ruleUseCase)
		languageRepo := database.NewPostgresLanguageRepository(db.DB)
		globalRuleHandler := handler.NewGlobalRuleHandler(globalRuleUseCase, languageRepo)
		api := r.Group("/api/v1")
		{
			api.GET("/projects", projectHandler.GetProjects)
			api.POST("/projects", projectHandler.CreateProject)
			api.PUT("/projects/:project_id", projectHandler.UpdateProject)
			api.DELETE("/projects/:project_id", projectHandler.DeleteProject)
			api.GET("/rules", ruleHandler.GetRules)
			api.GET("/rules/:project_id/:rule_id", ruleHandler.GetRule)
			api.POST("/rules", ruleHandler.CreateRule)
			api.PUT("/rules/:project_id/:rule_id", ruleHandler.UpdateRule)
			api.DELETE("/rules/:project_id/:rule_id", ruleHandler.DeleteRule)
			api.POST("/rules/validate", ruleHandler.ValidateCode)
			api.POST("/rules/export", ruleHandler.ExportRules)
			api.POST("/rules/import", ruleHandler.ImportRules)
			api.GET("/global-rules/:language", globalRuleHandler.GetGlobalRules)
			api.POST("/global-rules", globalRuleHandler.CreateGlobalRule)
			api.DELETE("/global-rules/:language/:rule_id", globalRuleHandler.DeleteGlobalRule)
			api.GET("/languages", globalRuleHandler.GetLanguages)
			api.POST("/languages", globalRuleHandler.CreateLanguage)
			api.PUT("/languages/:code", globalRuleHandler.UpdateLanguage)
			api.DELETE("/languages/:code", globalRuleHandler.DeleteLanguage)
			api.POST("/global-rules/export", globalRuleHandler.ExportGlobalRules)
			api.POST("/global-rules/import", globalRuleHandler.ImportGlobalRules)
		}

		// MCP endpoints
		projectDetector := usecase.NewProjectDetector(projectRepo, ruleRepo)
		mcpHandler := handler.NewMCPHandler(ruleUseCase, globalRuleUseCase, projectDetector)
		// inject metrics repo via setter
		mcpHandler.SetMetricsRepo(metricsRepo)
		mcp := r.Group("/mcp")
		{
			mcp.POST("/request", mcpHandler.HandleMCPRequest)
			mcp.GET("/ws", mcpHandler.HandleWebSocket)
		}
	} else {
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
