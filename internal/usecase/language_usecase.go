package usecase

import (
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
)

type LanguageUseCase struct {
	languageRepo domain.LanguageRepository
}

func NewLanguageUseCase(languageRepo domain.LanguageRepository) *LanguageUseCase {
	return &LanguageUseCase{
		languageRepo: languageRepo,
	}
}

func (uc *LanguageUseCase) GetLanguages() ([]*domain.Language, error) {
	return uc.languageRepo.GetAll()
}

func (uc *LanguageUseCase) GetLanguage(code string) (*domain.Language, error) {
	return uc.languageRepo.GetByCode(code)
}

func (uc *LanguageUseCase) CreateLanguage(code, name, description, icon, color string, isActive bool) error {
	language := &domain.Language{
		Code:        code,
		Name:        name,
		Description: description,
		Icon:        icon,
		Color:       color,
		IsActive:    isActive,
	}
	return uc.languageRepo.Create(language)
}

func (uc *LanguageUseCase) UpdateLanguage(code, name, description, icon, color string, isActive bool) error {
	language := &domain.Language{
		Code:        code,
		Name:        name,
		Description: description,
		Icon:        icon,
		Color:       color,
		IsActive:    isActive,
	}
	return uc.languageRepo.Update(language)
}

func (uc *LanguageUseCase) DeleteLanguage(code string) error {
	return uc.languageRepo.Delete(code)
}
