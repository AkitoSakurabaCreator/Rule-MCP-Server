package usecase

import (
	"regexp"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/apperr"
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
		missing := []string{}
		if projectID == "" {
			missing = append(missing, "project_id")
		}
		if ruleID == "" {
			missing = append(missing, "rule_id")
		}
		if name == "" {
			missing = append(missing, "name")
		}
		return apperr.WrapWithDetails(apperr.ErrValidation, "入力値が不正です", map[string]interface{}{"missing": missing})
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

func (uc *RuleUseCase) GetRule(projectID, ruleID string) (*domain.Rule, error) {
	if projectID == "" || ruleID == "" {
		return nil, apperr.WrapWithDetails(apperr.ErrValidation, "入力値が不正です", map[string]interface{}{"missing": []string{"project_id", "rule_id"}})
	}
	return uc.ruleRepo.GetByID(projectID, ruleID)
}

func (uc *RuleUseCase) UpdateRule(projectID, ruleID, name, description, ruleType, severity, pattern, message string, isActive *bool) error {
	if projectID == "" || ruleID == "" {
		return apperr.WrapWithDetails(apperr.ErrValidation, "入力値が不正です", map[string]interface{}{"missing": []string{"project_id", "rule_id"}})
	}
	existing, err := uc.ruleRepo.GetByID(projectID, ruleID)
	if err != nil {
		return err
	}
	if name != "" {
		existing.Name = name
	}
	existing.Description = description
	if ruleType != "" {
		existing.Type = ruleType
	}
	if severity != "" {
		existing.Severity = severity
	}
	existing.Pattern = pattern
	existing.Message = message
	if isActive != nil {
		existing.IsActive = *isActive
	}
	return uc.ruleRepo.Update(existing)
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
		if rule.Pattern == "" {
			continue
		}

		matched, err := regexp.MatchString(rule.Pattern, code)
		if err != nil {
			continue
		}

		if matched {
			msg := rule.Message
			if msg == "" {
				if rule.Name != "" {
					msg = rule.Name
				} else {
					msg = rule.Description
				}
			}
			if rule.Severity == "error" {
				result.Errors = append(result.Errors, msg)
				result.Valid = false
			} else if rule.Severity == "warning" {
				result.Warnings = append(result.Warnings, msg)
			}
		}
	}

	return result, nil
}
