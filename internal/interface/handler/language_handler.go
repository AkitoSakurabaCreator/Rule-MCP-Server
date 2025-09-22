package handler

import (
	"net/http"
	"strings"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/usecase"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type LanguageHandler struct {
	languageUseCase *usecase.LanguageUseCase
}

func NewLanguageHandler(languageUseCase *usecase.LanguageUseCase) *LanguageHandler {
	return &LanguageHandler{
		languageUseCase: languageUseCase,
	}
}

func (h *LanguageHandler) GetLanguages(c *gin.Context) {
	languages, err := h.languageUseCase.GetLanguages()
	if err != nil {
		httpx.JSONFromError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"languages": languages})
}

func (h *LanguageHandler) GetLanguage(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "language code is required", nil)
		return
	}

	language, err := h.languageUseCase.GetLanguage(code)
	if err != nil {
		httpx.JSONFromError(c, err)
		return
	}

	c.JSON(http.StatusOK, language)
}

func (h *LanguageHandler) CreateLanguage(c *gin.Context) {
	var req struct {
		Code        string `json:"code" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Color       string `json:"color"`
		IsActive    bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}

	// デフォルト値の設定
	if req.Icon == "" {
		req.Icon = strings.ToLower(req.Code)
	}
	if req.Color == "" {
		req.Color = "#666666"
	}

	err := h.languageUseCase.CreateLanguage(req.Code, req.Name, req.Description, req.Icon, req.Color, req.IsActive)
	if err != nil {
		if strings.Contains(err.Error(), "一意制約") {
			httpx.JSONError(c, http.StatusConflict, httpx.CodeConflict, "この言語コードは既に使用されています。別の言語コードを指定してください。", nil)
			return
		}
		httpx.JSONFromError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Language created successfully"})
}

func (h *LanguageHandler) UpdateLanguage(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "language code is required", nil)
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Color       string `json:"color"`
		IsActive    bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}

	err := h.languageUseCase.UpdateLanguage(code, req.Name, req.Description, req.Icon, req.Color, req.IsActive)
	if err != nil {
		httpx.JSONFromError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language updated successfully"})
}

func (h *LanguageHandler) DeleteLanguage(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "language code is required", nil)
		return
	}

	err := h.languageUseCase.DeleteLanguage(code)
	if err != nil {
		httpx.JSONFromError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Language deleted successfully"})
}
