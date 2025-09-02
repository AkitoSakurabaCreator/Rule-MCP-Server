package usecase

import (
	"errors"
	"regexp"

	"github.com/AkitoSakurabaCreator/RuleMCPServer/internal/domain"
)

type RuleUseCase struct {
	ruleRepo       domain.RuleRepository
	globalRuleRepo domain.GlobalRuleRepository
	projectRepo    domain.ProjectRepository
}

func NewRuleUseCase(ruleRepo domain.RuleRepository, globalRuleRepo domain.GlobalRuleRepository, projectRepo domain.ProjectRepository) *RuleUseCase {
	return &RuleUseCase{
		ruleRepo:       ruleRepo,
		globalRuleRepo: globalRuleRepo,
		projectRepo:    projectRepo,
	}
}

func (uc *RuleUseCase) CreateRule(projectID, ruleID, name, description, ruleType, severity, pattern, message string) error {
	if projectID == "" || ruleID == "" || name == "" {
		return errors.New("project_id, rule_id, and name are required")
	}

	rule := &domain.Rule{
		ProjectID:   projectID,
		RuleID:      ruleID,
		Name:        name,
		Description: description,
		Type:        ruleType,
		Severity:    severity,
		Pattern:     pattern,
		Message:     message,
		IsActive:    true,
	}

	return uc.ruleRepo.Create(rule)
}

func (uc *RuleUseCase) GetProjectRules(projectID string) (*domain.ProjectRules, error) {
	project, err := uc.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}

	rules, err := uc.ruleRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	projectRules := &domain.ProjectRules{
		ProjectID: projectID,
		Rules:     make([]domain.Rule, 0, len(rules)),
	}

	for _, r := range rules {
		projectRules.Rules = append(projectRules.Rules, *r)
	}

	if project.ApplyGlobalRules {
		globalRules, err := uc.globalRuleRepo.GetByLanguage(project.Language)
		if err != nil {
			return nil, err
		}

		for _, globalRule := range globalRules {
			rule := domain.Rule{
				ProjectID:   projectID,
				RuleID:      globalRule.RuleID,
				Name:        globalRule.Name,
				Description: globalRule.Description,
				Type:        globalRule.Type,
				Severity:    globalRule.Severity,
				Pattern:     globalRule.Pattern,
				Message:     globalRule.Message,
				IsActive:    globalRule.IsActive,
			}
			projectRules.Rules = append(projectRules.Rules, rule)
		}
	}

	return projectRules, nil
}

func (uc *RuleUseCase) DeleteRule(projectID, ruleID string) error {
	return uc.ruleRepo.Delete(projectID, ruleID)
}

func (uc *RuleUseCase) ValidateCode(projectID, code string) (*domain.ValidationResult, error) {
	projectRules, err := uc.GetProjectRules(projectID)
	if err != nil {
		return nil, err
	}

	result := &domain.ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	for _, rule := range projectRules.Rules {
		if !rule.IsActive {
			continue
		}

		matched, err := regexp.MatchString(rule.Pattern, code)
		if err != nil {
			continue
		}

		if matched {
			if rule.Severity == "error" {
				result.Errors = append(result.Errors, rule.Message)
				result.Valid = false
			} else if rule.Severity == "warning" {
				result.Warnings = append(result.Warnings, rule.Message)
			}
		}
	}

	return result, nil
}
