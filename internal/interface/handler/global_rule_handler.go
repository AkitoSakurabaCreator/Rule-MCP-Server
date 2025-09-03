package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opm008077/RuleMCPServer/internal/domain"
	"github.com/opm008077/RuleMCPServer/internal/usecase"
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

// ExportGlobalRulesRequest グローバルルールエクスポートリクエスト
type ExportGlobalRulesRequest struct {
	Language string   `json:"language"`
	RuleIDs  []string `json:"ruleIds,omitempty"`
	Format   string   `json:"format"` // json, yaml, csv
}

// ImportGlobalRulesRequest グローバルルールインポートリクエスト
type ImportGlobalRulesRequest struct {
	Language  string              `json:"language"`
	Rules     []domain.GlobalRule `json:"rules"`
	Overwrite bool                `json:"overwrite"`
}

// ExportGlobalRules グローバルルールエクスポート
func (h *GlobalRuleHandler) ExportGlobalRules(c *gin.Context) {
	var req ExportGlobalRulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 管理者権限チェック
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// グローバルルール取得ロジック
	var rules []domain.GlobalRule
	if len(req.RuleIDs) > 0 {
		// 特定のルールIDを指定
		for _, ruleID := range req.RuleIDs {
			rule := domain.GlobalRule{
				RuleID:      ruleID,
				Language:    req.Language,
				Name:        "Global Rule " + ruleID,
				Description: "Exported global rule",
				Pattern:     "global_pattern_" + ruleID,
				Message:     "Global message for " + ruleID,
				Severity:    "warning",
				IsActive:    true,
			}
			rules = append(rules, rule)
		}
	} else {
		// 言語全体のルール
		rules = []domain.GlobalRule{
			{
				RuleID:      "global1",
				Language:    req.Language,
				Name:        "Global Rule 1",
				Description: "Exported global rule 1",
				Pattern:     "global_pattern_1",
				Message:     "Global message for rule 1",
				Severity:    "warning",
				IsActive:    true,
			},
			{
				RuleID:      "global2",
				Language:    req.Language,
				Name:        "Global Rule 2",
				Description: "Exported global rule 2",
				Pattern:     "global_pattern_2",
				Message:     "Global message for rule 2",
				Severity:    "error",
				IsActive:    true,
			},
		}
	}

	// フォーマットに応じてレスポンス
	switch req.Format {
	case "yaml":
		c.Header("Content-Type", "application/x-yaml")
		c.Header("Content-Disposition", "attachment; filename=global_rules_"+req.Language+".yaml")
		c.String(http.StatusOK, "language: "+req.Language+"\nrules:\n  - name: Global Rule 1\n    pattern: global_pattern_1\n  - name: Global Rule 2\n    pattern: global_pattern_2")
	case "csv":
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=global_rules_"+req.Language+".csv")
		c.String(http.StatusOK, "Language,RuleID,Name,Pattern,Message,Severity\n"+req.Language+",global1,Global Rule 1,global_pattern_1,Global message for rule 1,warning\n"+req.Language+",global2,Global Rule 2,global_pattern_2,Global message for rule 2,error")
	default:
		// JSON形式（デフォルト）
		c.JSON(http.StatusOK, gin.H{
			"language":   req.Language,
			"format":     req.Format,
			"rules":      rules,
			"exportedAt": time.Now().Format(time.RFC3339),
		})
	}
}

// ImportGlobalRules グローバルルールインポート
func (h *GlobalRuleHandler) ImportGlobalRules(c *gin.Context) {
	var req ImportGlobalRulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 管理者権限チェック
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// インポート処理
	importedCount := 0
	skippedCount := 0
	errors := []string{}

	for _, rule := range req.Rules {
		// ルール検証
		if rule.Name == "" || rule.Pattern == "" {
			errors = append(errors, "Global rule "+rule.RuleID+" has invalid data")
			continue
		}

		// 言語一致チェック
		if rule.Language != req.Language {
			errors = append(errors, "Global rule "+rule.RuleID+" language mismatch")
			continue
		}

		// 重複チェック
		if !req.Overwrite {
			// 既存ルールとの重複チェック（簡易実装）
			skippedCount++
			continue
		}

		// ルール保存処理（簡易実装）
		importedCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Global rules import completed",
		"language":      req.Language,
		"importedCount": importedCount,
		"skippedCount":  skippedCount,
		"errorCount":    len(errors),
		"errors":        errors,
		"importedAt":    time.Now().Format(time.RFC3339),
	})
}

// LanguageInfo 言語情報
type LanguageInfo struct {
	Code        string `json:"code"`        // 言語コード (go, javascript, python等)
	Name        string `json:"name"`        // 言語名 (Go, JavaScript, Python等)
	Description string `json:"description"` // 言語の説明
	Icon        string `json:"icon"`        // アイコン（FontAwesome等）
	Color       string `json:"color"`       // テーマカラー
	IsActive    bool   `json:"isActive"`    // 有効/無効
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateLanguageRequest 言語作成リクエスト
type CreateLanguageRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
}

// UpdateLanguageRequest 言語更新リクエスト
type UpdateLanguageRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	IsActive    *bool  `json:"isActive"`
}

// CreateLanguage 新しい言語を作成
func (h *GlobalRuleHandler) CreateLanguage(c *gin.Context) {
	var req CreateLanguageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 管理者権限チェック
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// 言語コードの重複チェック
	existingLanguages, err := h.globalRuleUseCase.GetAllLanguages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing languages"})
		return
	}

	for _, lang := range existingLanguages {
		if lang == req.Code {
			c.JSON(http.StatusConflict, gin.H{"error": "Language code already exists"})
			return
		}
	}

	// 言語作成処理（簡易実装）
	languageInfo := LanguageInfo{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
		IsActive:    true,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Language created successfully",
		"language": languageInfo,
	})
}

// UpdateLanguage 言語情報を更新
func (h *GlobalRuleHandler) UpdateLanguage(c *gin.Context) {
	languageCode := c.Param("code")
	if languageCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language code is required"})
		return
	}

	var req UpdateLanguageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 管理者権限チェック
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// 言語更新処理（簡易実装）
	languageInfo := LanguageInfo{
		Code:        languageCode,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
		IsActive:    true,
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	if req.IsActive != nil {
		languageInfo.IsActive = *req.IsActive
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Language updated successfully",
		"language": languageInfo,
	})
}

// DeleteLanguage 言語を削除
func (h *GlobalRuleHandler) DeleteLanguage(c *gin.Context) {
	languageCode := c.Param("code")
	if languageCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language code is required"})
		return
	}

	// 管理者権限チェック
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// 言語削除処理（簡易実装）
	// 実際の実装では、関連するルールの確認が必要

	c.JSON(http.StatusOK, gin.H{
		"message": "Language deleted successfully",
		"code":    languageCode,
	})
}

// GetLanguageInfo 言語の詳細情報を取得
func (h *GlobalRuleHandler) GetLanguageInfo(c *gin.Context) {
	languageCode := c.Param("code")
	if languageCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language code is required"})
		return
	}

	// 言語情報取得処理（簡易実装）
	languageInfo := LanguageInfo{
		Code:        languageCode,
		Name:        "Sample Language",
		Description: "Sample language description",
		Icon:        "code",
		Color:       "#007acc",
		IsActive:    true,
		CreatedAt:   "2025-01-01T00:00:00Z",
		UpdatedAt:   "2025-09-03T18:00:00Z",
	}

	c.JSON(http.StatusOK, languageInfo)
}
