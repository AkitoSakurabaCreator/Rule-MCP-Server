package usecase

import (
	"errors"
	"time"

	"github.com/opm008077/RuleMCPServer/internal/domain"
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
		return errors.New("project_id and name are required")
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
