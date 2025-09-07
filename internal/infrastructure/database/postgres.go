package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/apperr"
	pq "github.com/lib/pq"
)

type PostgresDatabase struct {
	DB *sql.DB
}

type PostgresRuleRepository struct {
	DB *sql.DB
}

type PostgresGlobalRuleRepository struct {
	DB *sql.DB
}

type PostgresRuleOptionRepository struct {
	DB *sql.DB
}

type PostgresRoleRepository struct {
	DB *sql.DB
}

type PostgresMetricsRepository struct {
	DB *sql.DB
}

// Ensure implementations
var _ domain.ProjectRepository = (*PostgresDatabase)(nil)
var _ domain.RuleRepository = (*PostgresRuleRepository)(nil)
var _ domain.GlobalRuleRepository = (*PostgresGlobalRuleRepository)(nil)
var _ domain.RuleOptionRepository = (*PostgresRuleOptionRepository)(nil)
var _ domain.RoleRepository = (*PostgresRoleRepository)(nil)
var _ domain.MetricsRepository = (*PostgresMetricsRepository)(nil)

func NewPostgresDatabase(host, port, user, password, dbname string) (*PostgresDatabase, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return &PostgresDatabase{DB: db}, nil
}

func NewPostgresRuleRepository(db *sql.DB) *PostgresRuleRepository {
	return &PostgresRuleRepository{DB: db}
}

func NewPostgresGlobalRuleRepository(db *sql.DB) *PostgresGlobalRuleRepository {
	return &PostgresGlobalRuleRepository{DB: db}
}

func NewPostgresRuleOptionRepository(db *sql.DB) *PostgresRuleOptionRepository {
	return &PostgresRuleOptionRepository{DB: db}
}

func NewPostgresRoleRepository(db *sql.DB) *PostgresRoleRepository {
	return &PostgresRoleRepository{DB: db}
}

func NewPostgresMetricsRepository(db *sql.DB) *PostgresMetricsRepository {
	return &PostgresMetricsRepository{DB: db}
}

func (d *PostgresDatabase) Close() error {
	return d.DB.Close()
}

func (d *PostgresDatabase) Create(project *domain.Project) error {
	query := `INSERT INTO projects (project_id, name, description, language, apply_global_rules, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := d.DB.Exec(query, project.ProjectID, project.Name, project.Description, project.Language, project.ApplyGlobalRules, project.CreatedAt, project.UpdatedAt)
	return mapDBError(err)
}

func (d *PostgresDatabase) GetByID(projectID string) (*domain.Project, error) {
	query := `SELECT project_id, name, description, language, apply_global_rules, created_at, updated_at 
			  FROM projects WHERE project_id = $1`

	var project domain.Project
	err := d.DB.QueryRow(query, projectID).Scan(
		&project.ProjectID, &project.Name, &project.Description, &project.Language,
		&project.ApplyGlobalRules, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return nil, mapDBError(err)
	}
	return &project, nil
}

func (d *PostgresDatabase) GetAll() ([]*domain.Project, error) {
	query := `SELECT project_id, name, description, language, apply_global_rules, created_at, updated_at 
			  FROM projects ORDER BY created_at DESC`

	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, mapDBError(err)
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var project domain.Project
		err := rows.Scan(
			&project.ProjectID, &project.Name, &project.Description, &project.Language,
			&project.ApplyGlobalRules, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, mapDBError(err)
		}
		projects = append(projects, &project)
	}

	return projects, nil
}

func (d *PostgresDatabase) Update(project *domain.Project) error {
	query := `UPDATE projects SET name = $2, description = $3, language = $4, apply_global_rules = $5, updated_at = $6 
			  WHERE project_id = $1`
	_, err := d.DB.Exec(query, project.ProjectID, project.Name, project.Description, project.Language, project.ApplyGlobalRules, project.UpdatedAt)
	return mapDBError(err)
}

func (d *PostgresDatabase) Delete(projectID string) error {
	query := `DELETE FROM projects WHERE project_id = $1`
	_, err := d.DB.Exec(query, projectID)
	return mapDBError(err)
}

// RuleRepository implementation
func (d *PostgresRuleRepository) Create(rule *domain.Rule) error {
	query := `INSERT INTO rules (project_id, rule_id, name, description, type, severity, pattern, message, is_active) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := d.DB.Exec(query, rule.ProjectID, rule.RuleID, rule.Name, rule.Description, rule.Type, rule.Severity, rule.Pattern, rule.Message, rule.IsActive)
	return mapDBError(err)
}

func (d *PostgresRuleRepository) GetByProjectID(projectID string) ([]*domain.Rule, error) {
	query := `SELECT id, project_id, rule_id, name, description, type, severity, pattern, message, is_active 
			  FROM rules WHERE project_id = $1 AND is_active = true ORDER BY severity DESC, name ASC`

	rows, err := d.DB.Query(query, projectID)
	if err != nil {
		return nil, mapDBError(err)
	}
	defer rows.Close()

	var rules []*domain.Rule
	for rows.Next() {
		var rule domain.Rule
		err := rows.Scan(
			&rule.ID, &rule.ProjectID, &rule.RuleID, &rule.Name, &rule.Description,
			&rule.Type, &rule.Severity, &rule.Pattern, &rule.Message, &rule.IsActive)
		if err != nil {
			return nil, mapDBError(err)
		}
		rules = append(rules, &rule)
	}

	return rules, nil
}

func (d *PostgresRuleRepository) GetByID(projectID, ruleID string) (*domain.Rule, error) {
	query := `SELECT id, project_id, rule_id, name, description, type, severity, pattern, message, is_active
              FROM rules WHERE project_id = $1 AND rule_id = $2`
	var rule domain.Rule
	err := d.DB.QueryRow(query, projectID, ruleID).Scan(
		&rule.ID, &rule.ProjectID, &rule.RuleID, &rule.Name, &rule.Description,
		&rule.Type, &rule.Severity, &rule.Pattern, &rule.Message, &rule.IsActive,
	)
	if err != nil {
		return nil, mapDBError(err)
	}
	return &rule, nil
}

func (d *PostgresRuleRepository) Update(rule *domain.Rule) error {
	query := `UPDATE rules SET name=$3, description=$4, type=$5, severity=$6, pattern=$7, message=$8, is_active=$9, project_id=$2
              WHERE project_id=$1 AND rule_id=$10`
	_, err := d.DB.Exec(query, rule.ProjectID, rule.ProjectID, rule.Name, rule.Description, rule.Type, rule.Severity, rule.Pattern, rule.Message, rule.IsActive, rule.RuleID)
	return mapDBError(err)
}

func (d *PostgresRuleRepository) Delete(projectID, ruleID string) error {
	query := `DELETE FROM rules WHERE project_id = $1 AND rule_id = $2`
	_, err := d.DB.Exec(query, projectID, ruleID)
	return mapDBError(err)
}

// GlobalRuleRepository implementation
func (d *PostgresGlobalRuleRepository) Create(rule *domain.GlobalRule) error {
	query := `INSERT INTO global_rules (language, rule_id, name, description, type, severity, pattern, message, is_active) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := d.DB.Exec(query, rule.Language, rule.RuleID, rule.Name, rule.Description, rule.Type, rule.Severity, rule.Pattern, rule.Message, rule.IsActive)
	return mapDBError(err)
}

func (d *PostgresGlobalRuleRepository) GetByLanguage(language string) ([]*domain.GlobalRule, error) {
	query := `SELECT id, language, rule_id, name, description, type, severity, pattern, message, is_active 
			  FROM global_rules WHERE language = $1 AND is_active = true ORDER BY severity DESC, name ASC`

	rows, err := d.DB.Query(query, language)
	if err != nil {
		return nil, mapDBError(err)
	}
	defer rows.Close()

	var rules []*domain.GlobalRule
	for rows.Next() {
		var rule domain.GlobalRule
		err := rows.Scan(
			&rule.ID, &rule.Language, &rule.RuleID, &rule.Name, &rule.Description,
			&rule.Type, &rule.Severity, &rule.Pattern, &rule.Message, &rule.IsActive)
		if err != nil {
			return nil, mapDBError(err)
		}
		rules = append(rules, &rule)
	}

	return rules, nil
}

func (d *PostgresGlobalRuleRepository) GetAllLanguages() ([]string, error) {
	query := `SELECT DISTINCT language FROM global_rules WHERE is_active = true ORDER BY language`

	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, mapDBError(err)
	}
	defer rows.Close()

	var languages []string
	for rows.Next() {
		var language string
		err := rows.Scan(&language)
		if err != nil {
			return nil, mapDBError(err)
		}
		languages = append(languages, language)
	}

	return languages, nil
}

func (d *PostgresGlobalRuleRepository) Delete(language, ruleID string) error {
	query := `DELETE FROM global_rules WHERE language = $1 AND rule_id = $2`
	_, err := d.DB.Exec(query, language, ruleID)
	return mapDBError(err)
}

// RuleOptionRepository implementation
func (r *PostgresRuleOptionRepository) GetByKind(kind string) ([]domain.RuleOption, error) {
	query := `SELECT id, kind, value, is_active FROM rule_options WHERE kind = $1 AND is_active = true ORDER BY value`
	rows, err := r.DB.Query(query, kind)
	if err != nil {
		return nil, mapDBError(err)
	}
	defer rows.Close()
	var opts []domain.RuleOption
	for rows.Next() {
		var o domain.RuleOption
		if err := rows.Scan(&o.ID, &o.Kind, &o.Value, &o.IsActive); err != nil {
			return nil, mapDBError(err)
		}
		opts = append(opts, o)
	}
	return opts, nil
}

func (r *PostgresRuleOptionRepository) Add(kind, value string) error {
	query := `INSERT INTO rule_options (kind, value, is_active) VALUES ($1, $2, true) ON CONFLICT (kind, value) DO UPDATE SET is_active = EXCLUDED.is_active`
	_, err := r.DB.Exec(query, kind, value)
	return mapDBError(err)
}

func (r *PostgresRuleOptionRepository) Delete(kind, value string) error {
	query := `DELETE FROM rule_options WHERE kind = $1 AND value = $2`
	_, err := r.DB.Exec(query, kind, value)
	return mapDBError(err)
}

// RoleRepository implementation
func (r *PostgresRoleRepository) GetAll() ([]domain.Role, error) {
	query := `SELECT id, name, description, permissions, is_active FROM roles ORDER BY name`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, mapDBError(err)
	}
	defer rows.Close()
	var roles []domain.Role
	for rows.Next() {
		var id int
		var name, description string
		var isActive bool
		var permRaw []byte
		if err := rows.Scan(&id, &name, &description, &permRaw, &isActive); err != nil {
			return nil, mapDBError(err)
		}
		var perms map[string]bool
		if len(permRaw) > 0 {
			_ = json.Unmarshal(permRaw, &perms)
		}
		roles = append(roles, domain.Role{ID: id, Name: name, Description: description, Permissions: perms, IsActive: isActive})
	}
	return roles, nil
}

func (r *PostgresRoleRepository) GetByName(name string) (domain.Role, error) {
	query := `SELECT id, name, description, permissions, is_active FROM roles WHERE name=$1`
	var id int
	var description string
	var isActive bool
	var permRaw []byte
	if err := r.DB.QueryRow(query, name).Scan(&id, &name, &description, &permRaw, &isActive); err != nil {
		return domain.Role{}, mapDBError(err)
	}
	var perms map[string]bool
	if len(permRaw) > 0 {
		_ = json.Unmarshal(permRaw, &perms)
	}
	return domain.Role{ID: id, Name: name, Description: description, Permissions: perms, IsActive: isActive}, nil
}

func (r *PostgresRoleRepository) Create(role domain.Role) error {
	permJSON, _ := json.Marshal(role.Permissions)
	query := `INSERT INTO roles (name, description, permissions, is_active) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, role.Name, role.Description, string(permJSON), role.IsActive)
	return mapDBError(err)
}

func (r *PostgresRoleRepository) Update(name string, role domain.Role) error {
	permJSON, _ := json.Marshal(role.Permissions)
	query := `UPDATE roles SET description=$2, permissions=$3, is_active=$4, updated_at=NOW() WHERE name=$1`
	_, err := r.DB.Exec(query, name, role.Description, string(permJSON), role.IsActive)
	return mapDBError(err)
}

func (r *PostgresRoleRepository) Delete(name string) error {
	query := `DELETE FROM roles WHERE name = $1`
	_, err := r.DB.Exec(query, name)
	return mapDBError(err)
}

// GetByLanguage 言語別にプロジェクトを取得
func (d *PostgresDatabase) GetByLanguage(language string) ([]*domain.Project, error) {
	query := `SELECT project_id, name, description, language, apply_global_rules, access_level, created_by, created_at, updated_at 
			  FROM projects WHERE language = $1 ORDER BY created_at DESC`

	rows, err := d.DB.Query(query, language)
	if err != nil {
		return nil, mapDBError(err)
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var project domain.Project
		err := rows.Scan(
			&project.ProjectID,
			&project.Name,
			&project.Description,
			&project.Language,
			&project.ApplyGlobalRules,
			&project.AccessLevel,
			&project.CreatedBy,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, mapDBError(err)
		}
		projects = append(projects, &project)
	}

	if err = rows.Err(); err != nil {
		return nil, mapDBError(err)
	}

	return projects, nil
}

// LanguageRepository implementation
type PostgresLanguageRepository struct {
	DB *sql.DB
}

func NewPostgresLanguageRepository(db *sql.DB) domain.LanguageRepository {
	return &PostgresLanguageRepository{DB: db}
}

func (r *PostgresLanguageRepository) Create(language *domain.Language) error {
	query := `INSERT INTO languages (code, name, description, icon, color, is_active) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.DB.Exec(query, language.Code, language.Name, language.Description, language.Icon, language.Color, language.IsActive)
	return mapDBError(err)
}

func (r *PostgresLanguageRepository) GetByCode(code string) (*domain.Language, error) {
	query := `SELECT code, name, description, icon, color, is_active FROM languages WHERE code = $1`
	var language domain.Language
	err := r.DB.QueryRow(query, code).Scan(
		&language.Code, &language.Name, &language.Description,
		&language.Icon, &language.Color, &language.IsActive,
	)
	if err != nil {
		return nil, mapDBError(err)
	}
	return &language, nil
}

func (r *PostgresLanguageRepository) GetAll() ([]*domain.Language, error) {
	query := `SELECT code, name, description, icon, color, is_active FROM languages ORDER BY name`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, mapDBError(err)
	}
	defer rows.Close()

	var languages []*domain.Language
	for rows.Next() {
		var language domain.Language
		err := rows.Scan(
			&language.Code, &language.Name, &language.Description,
			&language.Icon, &language.Color, &language.IsActive,
		)
		if err != nil {
			return nil, mapDBError(err)
		}
		languages = append(languages, &language)
	}
	return languages, nil
}

func (r *PostgresLanguageRepository) Update(language *domain.Language) error {
	query := `UPDATE languages SET name = $1, description = $2, icon = $3, color = $4, is_active = $5 WHERE code = $6`
	_, err := r.DB.Exec(query, language.Name, language.Description, language.Icon, language.Color, language.IsActive, language.Code)
	return mapDBError(err)
}

func (r *PostgresLanguageRepository) Delete(code string) error {
	query := `DELETE FROM languages WHERE code = $1`
	_, err := r.DB.Exec(query, code)
	return mapDBError(err)
}

// MetricsRepository implementation
func (m *PostgresMetricsRepository) RecordMCP(method string, status string, durationMs int) error {
	q := `INSERT INTO mcp_requests (method, status, duration_ms, created_at) VALUES ($1, $2, $3, $4)`
	_, err := m.DB.Exec(q, method, status, durationMs, time.Now())
	return mapDBError(err)
}

func (m *PostgresMetricsRepository) GetMCPStatsLast24h() ([]domain.MCPMethodStat, error) {
	q := `SELECT method, COUNT(*), MAX(created_at) FROM mcp_requests WHERE created_at > NOW() - INTERVAL '24 hours' GROUP BY method ORDER BY COUNT(*) DESC`
	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, mapDBError(err)
	}
	defer rows.Close()
	stats := []domain.MCPMethodStat{}
	for rows.Next() {
		var method string
		var cnt int
		var last time.Time
		if err := rows.Scan(&method, &cnt, &last); err != nil {
			return nil, mapDBError(err)
		}
		stats = append(stats, domain.MCPMethodStat{Method: method, Count: cnt, LastUsed: last.Format("2006-01-02 15:04:05"), Status: "ok"})
	}
	return stats, nil
}

func (m *PostgresMetricsRepository) GetMCPRequestsCountLast24h() (int, error) {
	q := `SELECT COUNT(*) FROM mcp_requests WHERE created_at > NOW() - INTERVAL '24 hours'`
	var cnt int
	if err := m.DB.QueryRow(q).Scan(&cnt); err != nil {
		return 0, mapDBError(err)
	}
	return cnt, nil
}

func mapDBError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return apperr.Wrap(apperr.ErrNotFound, "対象が見つかりません")
	}
	if pqErr, ok := err.(*pq.Error); ok {
		switch string(pqErr.Code) {
		case "23505": // unique violation
			return apperr.WrapWithDetails(apperr.ErrConflict, "一意制約に違反しています", map[string]string{"constraint": pqErr.Constraint})
		case "23503": // foreign key violation
			return apperr.Wrap(apperr.ErrUnprocessable, "関連データが存在しないため処理できません")
		}
	}
	return err
}
