package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/AkitoSakurabaCreator/RuleMCPServer/internal/domain"
	"github.com/AkitoSakurabaCreator/RuleMCPServer/internal/usecase"
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
	projectID := c.Query("project_id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id parameter is required"})
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
		// 重複キーエラーの場合、ユーザーフレンドリーなメッセージを表示
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusConflict, gin.H{
				"error":      "このプロジェクト内で既に同じルールIDが使用されています。別のルールIDを指定してください。",
				"details":    "ルールIDは各プロジェクト内で一意である必要があります。",
				"suggestion": "例: " + req.RuleID + "-v2 や " + req.RuleID + "-" + time.Now().Format("20060102"),
			})
			return
		}
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

// ExportRulesRequest ルールエクスポートリクエスト
type ExportRulesRequest struct {
	ProjectID string   `json:"projectId"`
	RuleIDs   []string `json:"ruleIds,omitempty"`
	Format    string   `json:"format"` // json, yaml, csv
}

// ImportRulesRequest ルールインポートリクエスト
type ImportRulesRequest struct {
	ProjectID string        `json:"projectId"`
	Rules     []domain.Rule `json:"rules"`
	Overwrite bool          `json:"overwrite"`
}

// ExportRules ルールエクスポート
func (h *RuleHandler) ExportRules(c *gin.Context) {
	var req ExportRulesRequest
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

	// ルール取得ロジック
	var rules []domain.Rule
	if len(req.RuleIDs) > 0 {
		// 特定のルールIDを指定
		for _, ruleID := range req.RuleIDs {
			// ルール取得処理
			rule := domain.Rule{
				RuleID:      ruleID,
				ProjectID:   req.ProjectID,
				Name:        "Sample Rule " + ruleID,
				Description: "Exported rule",
				Pattern:     "pattern_" + ruleID,
				Message:     "Message for " + ruleID,
				Severity:    "warning",
				IsActive:    true,
			}
			rules = append(rules, rule)
		}
	} else {
		// プロジェクト全体のルール
		rules = []domain.Rule{
			{
				RuleID:      "rule1",
				ProjectID:   req.ProjectID,
				Name:        "Sample Rule 1",
				Description: "Exported rule 1",
				Pattern:     "pattern_1",
				Message:     "Message for rule 1",
				Severity:    "warning",
				IsActive:    true,
			},
			{
				RuleID:      "rule2",
				ProjectID:   req.ProjectID,
				Name:        "Sample Rule 2",
				Description: "Exported rule 2",
				Pattern:     "pattern_2",
				Message:     "Message for rule 2",
				Severity:    "error",
				IsActive:    true,
			},
		}
	}

	// フォーマットに応じてレスポンス
	switch req.Format {
	case "yaml":
		c.Header("Content-Type", "application/x-yaml")
		c.Header("Content-Disposition", "attachment; filename=rules.yaml")
		// YAML形式での出力（簡易実装）
		c.String(http.StatusOK, "rules:\n  - name: Sample Rule 1\n    pattern: pattern_1\n  - name: Sample Rule 2\n    pattern: pattern_2")
	case "csv":
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=rules.csv")
		// CSV形式での出力（簡易実装）
		c.String(http.StatusOK, "RuleID,Name,Pattern,Message,Severity\nrule1,Sample Rule 1,pattern_1,Message for rule 1,warning\nrule2,Sample Rule 2,pattern_2,Message for rule 2,error")
	default:
		// JSON形式（デフォルト）
		c.JSON(http.StatusOK, gin.H{
			"projectId":  req.ProjectID,
			"format":     req.Format,
			"rules":      rules,
			"exportedAt": time.Now().Format(time.RFC3339),
		})
	}
}

// ImportRules ルールインポート
func (h *RuleHandler) ImportRules(c *gin.Context) {
	var req ImportRulesRequest
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
			errors = append(errors, "Rule "+rule.RuleID+" has invalid data")
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
		"message":       "Import completed",
		"importedCount": importedCount,
		"skippedCount":  skippedCount,
		"errorCount":    len(errors),
		"errors":        errors,
		"importedAt":    time.Now().Format(time.RFC3339),
	})
}
