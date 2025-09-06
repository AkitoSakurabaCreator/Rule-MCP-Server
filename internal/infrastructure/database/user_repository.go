package database

import (
	"database/sql"
	"time"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) domain.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetByID(id int) (*domain.User, error) {
	query := `
		SELECT id, username, email, full_name, role, is_active, password_hash, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.FullName,
		&user.Role, &user.IsActive, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, full_name, role, is_active, password_hash, created_at, updated_at
		FROM users WHERE username = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.FullName,
		&user.Role, &user.IsActive, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, full_name, role, is_active, password_hash, created_at, updated_at
		FROM users WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.FullName,
		&user.Role, &user.IsActive, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetAll() ([]domain.User, error) {
	query := `
		SELECT id, username, email, full_name, role, is_active, password_hash, created_at, updated_at
		FROM users ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.FullName,
			&user.Role, &user.IsActive, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (username, email, full_name, role, is_active, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query, user.Username, user.Email, user.FullName, user.Role, user.IsActive, user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *PostgresUserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET username = $1, email = $2, full_name = $3, role = $4, is_active = $5, password_hash = $6, updated_at = $7
		WHERE id = $8
	`

	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		query, user.Username, user.Email, user.FullName, user.Role, user.IsActive, user.PasswordHash, user.UpdatedAt, user.ID,
	)
	return err
}

func (r *PostgresUserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PostgresUserRepository) GetActiveUsers() ([]domain.User, error) {
	query := `
		SELECT id, username, email, full_name, role, is_active, password_hash, created_at, updated_at
		FROM users WHERE is_active = true ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.FullName,
			&user.Role, &user.IsActive, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *PostgresUserRepository) GetUsersByRole(role string) ([]domain.User, error) {
	query := `
		SELECT id, username, email, full_name, role, is_active, password_hash, created_at, updated_at
		FROM users WHERE role = $1 ORDER BY id
	`

	rows, err := r.db.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.FullName,
			&user.Role, &user.IsActive, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
