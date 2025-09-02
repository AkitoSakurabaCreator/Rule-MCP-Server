package usecase

import (
	"errors"

	"github.com/AkitoSakurabaCreator/RuleMCPServer/internal/domain"
)

type GlobalRuleUseCase struct {
	globalRuleRepo domain.GlobalRuleRepository
}

func NewGlobalRuleUseCase(globalRuleRepo domain.GlobalRuleRepository) *GlobalRuleUseCase {
	return &GlobalRuleUseCase{
		globalRuleRepo: globalRuleRepo,
	}
}

func (uc *GlobalRuleUseCase) CreateGlobalRule(language, ruleID, name, description, ruleType, severity, pattern, message string) error {
	if language == "" || ruleID == "" || name == "" {
		return errors.New("language, rule_id, and name are required")
	}

	rule := &domain.GlobalRule{
		Language:    language,
		RuleID:      ruleID,
		Name:        name,
		Description: description,
		Type:        ruleType,
		Severity:    severity,
		Pattern:     pattern,
		Message:     message,
		IsActive:    true,
	}

	return uc.globalRuleRepo.Create(rule)
}

func (uc *GlobalRuleUseCase) GetGlobalRules(language string) ([]*domain.GlobalRule, error) {
	return uc.globalRuleRepo.GetByLanguage(language)
}

func (uc *GlobalRuleUseCase) GetAllLanguages() ([]string, error) {
	return uc.globalRuleRepo.GetAllLanguages()
}

func (uc *GlobalRuleUseCase) DeleteGlobalRule(language, ruleID string) error {
	return uc.globalRuleRepo.Delete(language, ruleID)
}
