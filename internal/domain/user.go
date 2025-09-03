package domain

import "time"

// User ユーザーエンティティ
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository ユーザーリポジトリインターフェース
type UserRepository interface {
	GetByID(id int) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAll() ([]User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
	GetActiveUsers() ([]User, error)
	GetUsersByRole(role string) ([]User, error)
}
