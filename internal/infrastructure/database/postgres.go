package database

import (
	"database/sql"
	"fmt"
	"log"

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

// Ensure PostgresDatabase implements ProjectRepository
var _ domain.ProjectRepository = (*PostgresDatabase)(nil)

// Ensure PostgresRuleRepository implements RuleRepository
var _ domain.RuleRepository = (*PostgresRuleRepository)(nil)

// Ensure PostgresGlobalRuleRepository implements GlobalRuleRepository
var _ domain.GlobalRuleRepository = (*PostgresGlobalRuleRepository)(nil)

// Ensure PostgresRuleOptionRepository implements RuleOptionRepository
var _ domain.RuleOptionRepository = (*PostgresRuleOptionRepository)(nil)

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
