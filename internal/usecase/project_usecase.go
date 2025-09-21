package usecase

import (
	"time"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/apperr"
)

type ProjectUseCase struct {
	projectRepo domain.ProjectRepository
}

func NewProjectUseCase(projectRepo domain.ProjectRepository) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepo: projectRepo,
	}
}

func (uc *ProjectUseCase) CreateProject(projectID, name, description, language string, applyGlobalRules bool) error {
	if projectID == "" || name == "" {
		return apperr.WrapWithDetails(apperr.ErrValidation, "入力値が不正です", map[string]interface{}{"missing": []string{"project_id", "name"}})
	}

	project := &domain.Project{
		ProjectID:        projectID,
		Name:             name,
		Description:      description,
		Language:         language,
		ApplyGlobalRules: applyGlobalRules,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return uc.projectRepo.Create(project)
}

func (uc *ProjectUseCase) GetProjects() ([]*domain.Project, error) {
	return uc.projectRepo.GetAll()
}

func (uc *ProjectUseCase) GetByID(projectID string) (*domain.Project, error) {
	return uc.projectRepo.GetByID(projectID)
}

func (uc *ProjectUseCase) UpdateProject(projectID, name, description, language string, applyGlobalRules bool) error {
	project, err := uc.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}

	project.Name = name
	project.Description = description
	project.Language = language
	project.ApplyGlobalRules = applyGlobalRules
	project.UpdatedAt = time.Now()

	return uc.projectRepo.Update(project)
}

func (uc *ProjectUseCase) DeleteProject(projectID string) error {
	return uc.projectRepo.Delete(projectID)
}
