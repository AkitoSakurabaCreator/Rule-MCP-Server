package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	// データベース接続情報を環境変数から取得
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "rule_mcp_db")
	dbUser := getEnv("DB_USER", "rule_mcp_user")
	dbPassword := getEnv("DB_PASSWORD", "rule_mcp_password")

	// 接続文字列を作成
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// データベースに接続
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 接続をテスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL database: %s:%s/%s", dbHost, dbPort, dbName)

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) LoadRulesFromDB(projectID string) (*ProjectRules, error) {
	// プロジェクト情報を取得
	projectQuery := `
		SELECT language, apply_global_rules
		FROM projects
		WHERE project_id = $1
	`

	var language string
	var applyGlobalRules bool
	err := d.db.QueryRow(projectQuery, projectID).Scan(&language, &applyGlobalRules)
	if err != nil {
		return nil, fmt.Errorf("failed to get project info: %w", err)
	}

	// プロジェクト固有のルールを取得
	projectRulesQuery := `
		SELECT r.rule_id, r.name, r.description, r.type, r.severity, r.pattern, r.message
		FROM rules r
		WHERE r.project_id = $1
		ORDER BY r.severity DESC, r.name ASC
	`

	projectRules, err := d.queryRules(projectRulesQuery, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project rules: %w", err)
	}

	var allRules []Rule
	allRules = append(allRules, projectRules...)

	// グローバルルールが適用される場合、言語別のグローバルルールを追加
	if applyGlobalRules && language != "" {
		globalRulesQuery := `
			SELECT rule_id, name, description, type, severity, pattern, message
			FROM global_rules
			WHERE language = $1 AND is_active = true
			ORDER BY severity DESC, name ASC
		`

		globalRules, err := d.queryRules(globalRulesQuery, language)
		if err != nil {
			return nil, fmt.Errorf("failed to get global rules: %w", err)
		}

		allRules = append(allRules, globalRules...)
	}

	if len(allRules) == 0 {
		return nil, fmt.Errorf("project %s has no rules", projectID)
	}

	return &ProjectRules{
		ProjectID: projectID,
		Rules:     allRules,
	}, nil
}

func (d *Database) queryRules(query string, args ...interface{}) ([]Rule, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []Rule
	for rows.Next() {
		var rule Rule
		err := rows.Scan(
			&rule.ID,
			&rule.Name,
			&rule.Description,
			&rule.Type,
			&rule.Severity,
			&rule.Pattern,
			&rule.Message,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rule: %w", err)
		}
		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rules: %w", err)
	}

	return rules, nil
}

func (d *Database) CreateProject(projectID, name, description, language string, applyGlobalRules bool) error {
	query := `
		INSERT INTO projects (project_id, name, description, language, apply_global_rules)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (project_id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			language = EXCLUDED.language,
			apply_global_rules = EXCLUDED.apply_global_rules,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := d.db.Exec(query, projectID, name, description, language, applyGlobalRules)
	if err != nil {
		return fmt.Errorf("failed to create/update project: %w", err)
	}

	return nil
}

func (d *Database) CreateRule(projectID, ruleID, name, description, ruleType, severity, pattern, message string) error {
	query := `
		INSERT INTO rules (project_id, rule_id, name, description, type, severity, pattern, message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (project_id, rule_id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			type = EXCLUDED.type,
			severity = EXCLUDED.severity,
			pattern = EXCLUDED.pattern,
			message = EXCLUDED.message,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := d.db.Exec(query, projectID, ruleID, name, description, ruleType, severity, pattern, message)
	if err != nil {
		return fmt.Errorf("failed to create/update rule: %w", err)
	}

	return nil
}

func (d *Database) DeleteRule(projectID, ruleID string) error {
	query := `DELETE FROM rules WHERE project_id = $1 AND rule_id = $2`

	result, err := d.db.Exec(query, projectID, ruleID)
	if err != nil {
		return fmt.Errorf("failed to delete rule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rule %s not found in project %s", ruleID, projectID)
	}

	return nil
}

func (d *Database) GetProjects() ([]map[string]interface{}, error) {
	query := `
		SELECT p.project_id, p.name, p.description, p.language, p.apply_global_rules, p.created_at, p.updated_at,
		       COUNT(r.id) as rule_count
		FROM projects p
		LEFT JOIN rules r ON p.project_id = r.project_id
		GROUP BY p.id, p.project_id, p.name, p.description, p.language, p.apply_global_rules, p.created_at, p.updated_at
		ORDER BY p.name ASC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []map[string]interface{}
	for rows.Next() {
		var projectID, name, description, language string
		var applyGlobalRules bool
		var createdAt, updatedAt sql.NullTime
		var ruleCount int

		err := rows.Scan(&projectID, &name, &description, &language, &applyGlobalRules, &createdAt, &updatedAt, &ruleCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		project := map[string]interface{}{
			"project_id":         projectID,
			"name":               name,
			"description":        description,
			"language":           language,
			"apply_global_rules": applyGlobalRules,
			"rule_count":         ruleCount,
		}

		if createdAt.Valid {
			project["created_at"] = createdAt.Time
		}
		if updatedAt.Valid {
			project["updated_at"] = updatedAt.Time
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over projects: %w", err)
	}

	return projects, nil
}

func (d *Database) GetGlobalRules(language string) ([]Rule, error) {
	query := `
		SELECT rule_id, name, description, type, severity, pattern, message
		FROM global_rules
		WHERE language = $1 AND is_active = true
		ORDER BY severity DESC, name ASC
	`

	return d.queryRules(query, language)
}

func (d *Database) CreateGlobalRule(language, ruleID, name, description, ruleType, severity, pattern, message string) error {
	query := `
		INSERT INTO global_rules (language, rule_id, name, description, type, severity, pattern, message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (language, rule_id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			type = EXCLUDED.type,
			severity = EXCLUDED.severity,
			pattern = EXCLUDED.pattern,
			message = EXCLUDED.message,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := d.db.Exec(query, language, ruleID, name, description, ruleType, severity, pattern, message)
	if err != nil {
		return fmt.Errorf("failed to create/update global rule: %w", err)
	}

	return nil
}

func (d *Database) DeleteGlobalRule(language, ruleID string) error {
	query := `DELETE FROM global_rules WHERE language = $1 AND rule_id = $2`

	result, err := d.db.Exec(query, language, ruleID)
	if err != nil {
		return fmt.Errorf("failed to delete global rule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("global rule %s not found for language %s", ruleID, language)
	}

	return nil
}

func (d *Database) GetLanguages() ([]string, error) {
	query := `SELECT DISTINCT language FROM global_rules WHERE is_active = true ORDER BY language`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query languages: %w", err)
	}
	defer rows.Close()

	var languages []string
	for rows.Next() {
		var language string
		if err := rows.Scan(&language); err != nil {
			return nil, fmt.Errorf("failed to scan language: %w", err)
		}
		languages = append(languages, language)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over languages: %w", err)
	}

	return languages, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
