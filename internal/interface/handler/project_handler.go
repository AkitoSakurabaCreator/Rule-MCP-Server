package handler

import (
	"net/http"
	"strings"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/usecase"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectUseCase *usecase.ProjectUseCase
}

func NewProjectHandler(projectUseCase *usecase.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{
		projectUseCase: projectUseCase,
	}
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
	projects, err := h.projectUseCase.GetProjects()
	if err != nil {
		httpx.JSONFromError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectID := c.Param("project_id")
	if projectID == "" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "project_id is required", nil)
		return
	}

	project, err := h.projectUseCase.GetByID(projectID)
	if err != nil {
		httpx.JSONFromError(c, err)
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req struct {
		ProjectID        string `json:"project_id" binding:"required"`
		Name             string `json:"name" binding:"required"`
		Description      string `json:"description"`
		Language         string `json:"language"`
		ApplyGlobalRules bool   `json:"apply_global_rules"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}

	err := h.projectUseCase.CreateProject(req.ProjectID, req.Name, req.Description, req.Language, req.ApplyGlobalRules)
	if err != nil {
		if strings.Contains(err.Error(), "一意制約") {
			httpx.JSONError(c, http.StatusConflict, httpx.CodeConflict, "このプロジェクトIDは既に使用されています。別のプロジェクトIDを指定してください。", nil)
			return
		}
		httpx.JSONFromError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	projectID := c.Param("project_id")
	if projectID == "" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "project_id is required", nil)
		return
	}

	var req struct {
		Name             string `json:"name" binding:"required"`
		Description      string `json:"description"`
		Language         string `json:"language"`
		ApplyGlobalRules bool   `json:"apply_global_rules"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}

	err := h.projectUseCase.UpdateProject(projectID, req.Name, req.Description, req.Language, req.ApplyGlobalRules)
	if err != nil {
		httpx.JSONFromError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID := c.Param("project_id")
	if projectID == "" {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "project_id is required", nil)
		return
	}

	err := h.projectUseCase.DeleteProject(projectID)
	if err != nil {
		httpx.JSONFromError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
