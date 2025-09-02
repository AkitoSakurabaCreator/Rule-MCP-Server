package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opm008077/RuleMCPServer/internal/usecase"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.projectUseCase.CreateProject(req.ProjectID, req.Name, req.Description, req.Language, req.ApplyGlobalRules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	projectID := c.Param("project_id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	var req struct {
		Name             string `json:"name" binding:"required"`
		Description      string `json:"description"`
		Language         string `json:"language"`
		ApplyGlobalRules bool   `json:"apply_global_rules"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.projectUseCase.UpdateProject(projectID, req.Name, req.Description, req.Language, req.ApplyGlobalRules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID := c.Param("project_id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	err := h.projectUseCase.DeleteProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
