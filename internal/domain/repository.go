package domain

type ProjectRepository interface {
	Create(project *Project) error
	GetByID(projectID string) (*Project, error)
	GetAll() ([]*Project, error)
	Update(project *Project) error
	Delete(projectID string) error
}

type RuleRepository interface {
	Create(rule *Rule) error
	GetByProjectID(projectID string) ([]*Rule, error)
	Delete(projectID, ruleID string) error
}

type GlobalRuleRepository interface {
	Create(rule *GlobalRule) error
	GetByLanguage(language string) ([]*GlobalRule, error)
	GetAllLanguages() ([]string, error)
	Delete(language, ruleID string) error
}

type ValidationRepository interface {
	ValidateCode(projectID, code string) (*ValidationResult, error)
}
