package fixtures

import (
	"time"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
)

// Projects test project data
var Projects = []domain.Project{
	{
		ProjectID:        "default",
		Name:             "Default Project",
		Description:      "Default project with common rules",
		Language:         "general",
		ApplyGlobalRules: true,
		AccessLevel:      "public",
		CreatedBy:        "system",
		CreatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		ProjectID:        "web-app",
		Name:             "Web Application",
		Description:      "Web application specific rules",
		Language:         "javascript",
		ApplyGlobalRules: true,
		AccessLevel:      "public",
		CreatedBy:        "admin",
		CreatedAt:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	},
	{
		ProjectID:        "api-service",
		Name:             "API Service",
		Description:      "API service specific rules",
		Language:         "go",
		ApplyGlobalRules: true,
		AccessLevel:      "public",
		CreatedBy:        "admin",
		CreatedAt:        time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
	},
	{
		ProjectID:        "team-project",
		Name:             "Team Project",
		Description:      "Team collaboration project",
		Language:         "typescript",
		ApplyGlobalRules: true,
		AccessLevel:      "user",
		CreatedBy:        "admin",
		CreatedAt:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	},
}

// ProjectByID retrieves a project with the specified ID
func ProjectByID(projectID string) *domain.Project {
	for _, project := range Projects {
		if project.ProjectID == projectID {
			return &project
		}
	}
	return nil
}

// ProjectsByLanguage retrieves projects for the specified language
func ProjectsByLanguage(language string) []domain.Project {
	var filteredProjects []domain.Project
	for _, project := range Projects {
		if project.Language == language {
			filteredProjects = append(filteredProjects, project)
		}
	}
	return filteredProjects
}

// PublicProjects retrieves projects with public access
func PublicProjects() []domain.Project {
	var publicProjects []domain.Project
	for _, project := range Projects {
		if project.AccessLevel == "public" {
			publicProjects = append(publicProjects, project)
		}
	}
	return publicProjects
}
