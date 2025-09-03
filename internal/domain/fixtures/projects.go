package fixtures

import (
	"time"

	"github.com/opm008077/RuleMCPServer/internal/domain"
)

// Projects テスト用プロジェクトデータ
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

// ProjectByID 指定されたIDのプロジェクトを取得
func ProjectByID(projectID string) *domain.Project {
	for _, project := range Projects {
		if project.ProjectID == projectID {
			return &project
		}
	}
	return nil
}

// ProjectsByLanguage 指定された言語のプロジェクトを取得
func ProjectsByLanguage(language string) []domain.Project {
	var filteredProjects []domain.Project
	for _, project := range Projects {
		if project.Language == language {
			filteredProjects = append(filteredProjects, project)
		}
	}
	return filteredProjects
}

// PublicProjects パブリックアクセスのプロジェクトを取得
func PublicProjects() []domain.Project {
	var publicProjects []domain.Project
	for _, project := range Projects {
		if project.AccessLevel == "public" {
			publicProjects = append(publicProjects, project)
		}
	}
	return publicProjects
}
