package fixtures

import (
	"time"

	"github.com/opm008077/RuleMCPServer/internal/domain"
)

// Users テスト用ユーザーデータ
var Users = []domain.User{
	{
		ID:        1,
		Username:  "admin",
		Email:     "admin@rulemcp.com",
		FullName:  "System Administrator",
		Role:      "admin",
		IsActive:  true,
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 14, 30, 15, 0, time.UTC),
	},
	{
		ID:        2,
		Username:  "user1",
		Email:     "user1@example.com",
		FullName:  "User One",
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 13, 45, 22, 0, time.UTC),
	},
	{
		ID:        3,
		Username:  "user2",
		Email:     "user2@example.com",
		FullName:  "User Two",
		Role:      "user",
		IsActive:  false,
		CreatedAt: time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 14, 16, 20, 10, 0, time.UTC),
	},
}

// AdminUser 管理者ユーザーを取得
func AdminUser() domain.User {
	return Users[0]
}

// ActiveUsers アクティブなユーザーのみを取得
func ActiveUsers() []domain.User {
	var activeUsers []domain.User
	for _, user := range Users {
		if user.IsActive {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers
}

// UserByID 指定されたIDのユーザーを取得
func UserByID(id int) *domain.User {
	for _, user := range Users {
		if user.ID == id {
			return &user
		}
	}
	return nil
}
