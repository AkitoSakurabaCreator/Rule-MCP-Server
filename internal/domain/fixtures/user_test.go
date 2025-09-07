package fixtures

import (
	"testing"
)

func TestUsers_AdminUserExists(t *testing.T) {
	adminUser := AdminUser()
	if adminUser.Role != "admin" {
		t.Errorf("Expected admin user role to be 'admin', got %s", adminUser.Role)
	}
}

func TestUsers_AdminUserIsActive(t *testing.T) {
	adminUser := AdminUser()
	if !adminUser.IsActive {
		t.Errorf("Expected admin user to be active")
	}
}

func TestUsers_AdminUserEmailFormat(t *testing.T) {
	adminUser := AdminUser()
	if adminUser.Email == "" {
		t.Errorf("Expected admin user to have an email")
	}

	// シンプルなメール形式チェック
	if len(adminUser.Email) < 5 || !contains(adminUser.Email, "@") {
		t.Errorf("Expected valid email format, got %s", adminUser.Email)
	}
}

func TestUsers_AdminUserCreatedAtBeforeUpdatedAt(t *testing.T) {
	adminUser := AdminUser()
	if adminUser.CreatedAt.After(adminUser.UpdatedAt) {
		t.Errorf("CreatedAt should be before or equal to UpdatedAt")
	}
}

func TestUsers_AdminUserUsernameNotEmpty(t *testing.T) {
	adminUser := AdminUser()
	if adminUser.Username == "" {
		t.Errorf("Expected username to not be empty")
	}
}

func TestUsers_ActiveUsersCount(t *testing.T) {
	activeUsers := ActiveUsers()
	expectedCount := 2 // admin + user1
	if len(activeUsers) != expectedCount {
		t.Errorf("Expected %d active users, got %d", expectedCount, len(activeUsers))
	}
}

func TestUsers_UserByID(t *testing.T) {
	user := UserByID(1)
	if user == nil {
		t.Errorf("Expected to find user with ID 1")
	}
	if user.Username != "admin" {
		t.Errorf("Expected user with ID 1 to be admin, got %s", user.Username)
	}
}

func TestUsers_UserByIDNotFound(t *testing.T) {
	user := UserByID(999)
	if user != nil {
		t.Errorf("Expected nil for non-existent user ID, got %v", user)
	}
}

// ヘルパー関数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			contains(s[1:], substr))))
}
