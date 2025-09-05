package domain

type ProjectRepository interface {
	Create(project *Project) error
	GetByID(projectID string) (*Project, error)
	GetAll() ([]*Project, error)
	GetByLanguage(language string) ([]*Project, error)
	Update(project *Project) error
	Delete(projectID string) error
}

type RuleRepository interface {
	Create(rule *Rule) error
	GetByProjectID(projectID string) ([]*Rule, error)
	GetByID(projectID, ruleID string) (*Rule, error)
	Update(rule *Rule) error
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

type RuleOption struct {
	ID       int    `json:"id"`
	Kind     string `json:"kind"` // type | severity
	Value    string `json:"value"`
	IsActive bool   `json:"is_active"`
}

type RuleOptionRepository interface {
	GetByKind(kind string) ([]RuleOption, error)
	Add(kind, value string) error
	Delete(kind, value string) error
}
