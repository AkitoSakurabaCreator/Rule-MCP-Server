package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/opm008077/RuleMCPServer/internal/domain"
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

// Ensure PostgresDatabase implements ProjectRepository
var _ domain.ProjectRepository = (*PostgresDatabase)(nil)
// Ensure PostgresRuleRepository implements RuleRepository
var _ domain.RuleRepository = (*PostgresRuleRepository)(nil)
// Ensure PostgresGlobalRuleRepository implements GlobalRuleRepository
var _ domain.GlobalRuleRepository = (*PostgresGlobalRuleRepository)(nil)

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

func (d *PostgresDatabase) Close() error {
	return d.DB.Close()
}

func (d *PostgresDatabase) Create(project *domain.Project) error {
	query := `INSERT INTO projects (project_id, name, description, language, apply_global_rules, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := d.DB.Exec(query, project.ProjectID, project.Name, project.Description, project.Language, project.ApplyGlobalRules, project.CreatedAt, project.UpdatedAt)
	return err
}

func (d *PostgresDatabase) GetByID(projectID string) (*domain.Project, error) {
	query := `SELECT project_id, name, description, language, apply_global_rules, created_at, updated_at 
			  FROM projects WHERE project_id = $1`
	
	var project domain.Project
	err := d.DB.QueryRow(query, projectID).Scan(
		&project.ProjectID, &project.Name, &project.Description, &project.Language,
		&project.ApplyGlobalRules, &project.CreatedAt, &project.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (d *PostgresDatabase) GetAll() ([]*domain.Project, error) {
	query := `SELECT project_id, name, description, language, apply_global_rules, created_at, updated_at 
			  FROM projects ORDER BY created_at DESC`
	
	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var project domain.Project
		err := rows.Scan(
			&project.ProjectID, &project.Name, &project.Description, &project.Language,
			&project.ApplyGlobalRules, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}

	return projects, nil
}

func (d *PostgresDatabase) Update(project *domain.Project) error {
	query := `UPDATE projects SET name = $2, description = $3, language = $4, apply_global_rules = $5, updated_at = $6 
			  WHERE project_id = $1`
	_, err := d.DB.Exec(query, project.ProjectID, project.Name, project.Description, project.Language, project.ApplyGlobalRules, project.UpdatedAt)
	return err
}

func (d *PostgresDatabase) Delete(projectID string) error {
	query := `DELETE FROM projects WHERE project_id = $1`
	_, err := d.DB.Exec(query, projectID)
	return err
}

// RuleRepository implementation
func (d *PostgresRuleRepository) Create(rule *domain.Rule) error {
	query := `INSERT INTO rules (project_id, rule_id, name, description, type, severity, pattern, message, is_active) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := d.DB.Exec(query, rule.ProjectID, rule.RuleID, rule.Name, rule.Description, rule.Type, rule.Severity, rule.Pattern, rule.Message, rule.IsActive)
	return err
}

func (d *PostgresRuleRepository) GetByProjectID(projectID string) ([]*domain.Rule, error) {
	query := `SELECT id, project_id, rule_id, name, description, type, severity, pattern, message, is_active 
			  FROM rules WHERE project_id = $1 AND is_active = true ORDER BY severity DESC, name ASC`
	
	rows, err := d.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*domain.Rule
	for rows.Next() {
		var rule domain.Rule
		err := rows.Scan(
			&rule.ID, &rule.ProjectID, &rule.RuleID, &rule.Name, &rule.Description,
			&rule.Type, &rule.Severity, &rule.Pattern, &rule.Message, &rule.IsActive)
		if err != nil {
			return nil, err
		}
		rules = append(rules, &rule)
	}

	return rules, nil
}

func (d *PostgresRuleRepository) Delete(projectID, ruleID string) error {
	query := `DELETE FROM rules WHERE project_id = $1 AND rule_id = $2`
	_, err := d.DB.Exec(query, projectID, ruleID)
	return err
}

// GlobalRuleRepository implementation
func (d *PostgresGlobalRuleRepository) Create(rule *domain.GlobalRule) error {
	query := `INSERT INTO global_rules (language, rule_id, name, description, type, severity, pattern, message, is_active) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := d.DB.Exec(query, rule.Language, rule.RuleID, rule.Name, rule.Description, rule.Type, rule.Severity, rule.Pattern, rule.Message, rule.IsActive)
	return err
}

func (d *PostgresGlobalRuleRepository) GetByLanguage(language string) ([]*domain.GlobalRule, error) {
	query := `SELECT id, language, rule_id, name, description, type, severity, pattern, message, is_active 
			  FROM global_rules WHERE language = $1 AND is_active = true ORDER BY severity DESC, name ASC`
	
	rows, err := d.DB.Query(query, language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*domain.GlobalRule
	for rows.Next() {
		var rule domain.GlobalRule
		err := rows.Scan(
			&rule.ID, &rule.Language, &rule.RuleID, &rule.Name, &rule.Description,
			&rule.Type, &rule.Severity, &rule.Pattern, &rule.Message, &rule.IsActive)
		if err != nil {
			return nil, err
		}
		rules = append(rules, &rule)
	}

	return rules, nil
}

func (d *PostgresGlobalRuleRepository) GetAllLanguages() ([]string, error) {
	query := `SELECT DISTINCT language FROM global_rules WHERE is_active = true ORDER BY language`
	
	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languages []string
	for rows.Next() {
		var language string
		err := rows.Scan(&language)
		if err != nil {
			return nil, err
		}
		languages = append(languages, language)
	}

	return languages, nil
}

func (d *PostgresGlobalRuleRepository) Delete(language, ruleID string) error {
	query := `DELETE FROM global_rules WHERE language = $1 AND rule_id = $2`
	_, err := d.DB.Exec(query, language, ruleID)
	return err
}
