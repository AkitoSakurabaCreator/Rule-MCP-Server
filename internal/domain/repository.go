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

type Language struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	IsActive    bool   `json:"is_active"`
}

type LanguageRepository interface {
	Create(language *Language) error
	GetByCode(code string) (*Language, error)
	GetAll() ([]*Language, error)
	Update(language *Language) error
	Delete(code string) error
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

type Role struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Permissions map[string]bool `json:"permissions"`
	IsActive    bool            `json:"is_active"`
}

type RoleRepository interface {
	GetAll() ([]Role, error)
	GetByName(name string) (Role, error)
	Create(role Role) error
	Update(name string, role Role) error
	Delete(name string) error
}

type MCPMethodStat struct {
	Method   string `json:"method"`
	Count    int    `json:"count"`
	LastUsed string `json:"lastUsed"`
	Status   string `json:"status"`
}

type MetricsRepository interface {
	RecordMCP(method string, status string, durationMs int) error
	GetMCPStatsLast24h() ([]MCPMethodStat, error)
	GetMCPRequestsCountLast24h() (int, error)
}
