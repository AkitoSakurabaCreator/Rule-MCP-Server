package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/internal/domain"
	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	jwtSecret []byte
	userRepo  domain.UserRepository
	roleRepo  domain.RoleRepository
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	User    User   `json:"user"`
	Message string `json:"message"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// ChangePasswordRequest パスワード変更リクエスト
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=12"`
}

// ApproveUserRequest ユーザー承認リクエスト
type ApproveUserRequest struct {
	UserID  int  `json:"userId" binding:"required"`
	Approve bool `json:"approve"`
}

func NewAuthHandler(jwtSecret string, userRepo domain.UserRepository, roleRepo domain.RoleRepository) *AuthHandler {
	return &AuthHandler{
		jwtSecret: []byte(jwtSecret),
		userRepo:  userRepo,
		roleRepo:  roleRepo,
	}
}

// validatePasswordStrength パスワードの複雑性要件を検証
func validatePasswordStrength(password string) error {
	if len(password) < 12 {
		return fmt.Errorf("パスワードは12文字以上である必要があります")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("パスワードには大文字を含める必要があります")
	}
	if !hasLower {
		return fmt.Errorf("パスワードには小文字を含める必要があります")
	}
	if !hasDigit {
		return fmt.Errorf("パスワードには数字を含める必要があります")
	}
	if !hasSpecial {
		return fmt.Errorf("パスワードには特殊文字を含める必要があります")
	}

	// Check for common patterns
	commonPatterns := []string{
		"password", "123456", "qwerty", "admin", "user",
		"letmein", "welcome", "monkey", "dragon", "master",
	}
	lowerPassword := strings.ToLower(password)
	for _, pattern := range commonPatterns {
		if strings.Contains(lowerPassword, pattern) {
			return fmt.Errorf("パスワードには一般的な単語やパターンを含めることはできません")
		}
	}

	return nil
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}

	// データベースからユーザーを取得
	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		httpx.JSONError(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "Invalid credentials", nil)
		return
	}

	// パスワードを検証
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		httpx.JSONError(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "Invalid credentials", nil)
		return
	}

	// ユーザーがアクティブかチェック
	if !user.IsActive {
		httpx.JSONError(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "Account is inactive", nil)
		return
	}

	// JWTトークンを生成
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "トークン生成に失敗しました", nil)
		return
	}

	response := LoginResponse{
		Token: tokenString,
		User: User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
		Message: "Login successful",
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		FullName string `json:"full_name"`
		Password string `json:"password" binding:"required,min=12"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}

	// ユーザー名の重複チェック
	existingUser, err := h.userRepo.GetByUsername(req.Username)
	if err == nil && existingUser != nil {
		httpx.JSONError(c, http.StatusConflict, httpx.CodeValidation, "Username already exists", nil)
		return
	}

	// メールアドレスの重複チェック
	existingEmail, err := h.userRepo.GetByEmail(req.Email)
	if err == nil && existingEmail != nil {
		httpx.JSONError(c, http.StatusConflict, httpx.CodeValidation, "Email already exists", nil)
		return
	}

	// パスワードの強度を検証
	if err := validatePasswordStrength(req.Password); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, err.Error(), nil)
		return
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "Password hashing failed", nil)
		return
	}

	// ユーザーを作成
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,
		Role:         "user",
		IsActive:     true,
		PasswordHash: string(hashedPassword),
	}

	err = h.userRepo.Create(user)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "User creation failed", nil)
		return
	}

	// JWTトークンを生成
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "トークン生成に失敗しました", nil)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: tokenString,
		User: User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
		Message: "Registration successful",
	})
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	// トークン検証の実装（必要に応じて）
	httpx.JSONError(c, http.StatusNotImplemented, httpx.CodeUnprocessable, "Token validation not implemented yet", nil)
}

// ChangePassword パスワード変更
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}

	// 現在のユーザーIDを取得（JWTから）
	userID, exists := c.Get("userID")
	if !exists {
		httpx.JSONError(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "User not authenticated", nil)
		return
	}

	// ユーザーを取得
	user, err := h.userRepo.GetByID(userID.(int))
	if err != nil {
		httpx.JSONError(c, http.StatusNotFound, httpx.CodeNotFound, "User not found", nil)
		return
	}

	// 現在のパスワードを検証
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "Current password is incorrect", nil)
		return
	}

	// 新しいパスワードの強度を検証
	if err := validatePasswordStrength(req.NewPassword); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, err.Error(), nil)
		return
	}

	// 新しいパスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "Password hashing failed", nil)
		return
	}

	// パスワードを更新
	user.PasswordHash = string(hashedPassword)
	err = h.userRepo.Update(user)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "Password update failed", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ApproveUser ユーザー承認/拒否
func (h *AuthHandler) ApproveUser(c *gin.Context) {
	var req ApproveUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}

	// 管理者権限チェック
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Admin access required", nil)
		return
	}

	// ユーザーを取得
	user, err := h.userRepo.GetByID(req.UserID)
	if err != nil {
		httpx.JSONError(c, http.StatusNotFound, httpx.CodeNotFound, "User not found", nil)
		return
	}

	// ユーザーのアクティブ状態を更新
	user.IsActive = req.Approve
	err = h.userRepo.Update(user)
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "User approval failed", nil)
		return
	}

	action := "approved"
	if !req.Approve {
		action = "rejected"
	}

	c.JSON(http.StatusOK, gin.H{"message": "User " + action + " successfully"})
}

// GetPendingUsers 承認待ちユーザー一覧取得
func (h *AuthHandler) GetPendingUsers(c *gin.Context) {
	// 管理者権限チェック
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		httpx.JSONError(c, http.StatusForbidden, httpx.CodeForbidden, "Admin access required", nil)
		return
	}

	// 非アクティブユーザーを取得
	users, err := h.userRepo.GetAll()
	if err != nil {
		httpx.JSONError(c, http.StatusInternalServerError, httpx.CodeInternal, "Failed to get users", nil)
		return
	}

	var pendingUsers []gin.H
	for _, user := range users {
		if !user.IsActive {
			pendingUsers = append(pendingUsers, gin.H{
				"id":        user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"fullName":  user.FullName,
				"role":      user.Role,
				"isActive":  user.IsActive,
				"createdAt": user.CreatedAt,
			})
		}
	}

	c.JSON(http.StatusOK, pendingUsers)
}

// Me 現在のユーザー情報
func (h *AuthHandler) Me(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if len(auth) <= 7 {
		httpx.JSONError(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "Unauthorized", nil)
		return
	}
	tokenStr := auth[7:]
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) { return h.jwtSecret, nil })
	if err != nil || !token.Valid {
		httpx.JSONError(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "Invalid token", nil)
		return
	}

	// データベースから最新のユーザー情報を取得
	user, err := h.userRepo.GetByID(claims.UserID)
	if err != nil {
		httpx.JSONError(c, http.StatusNotFound, httpx.CodeNotFound, "User not found", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"full_name": user.FullName,
		"role":      user.Role,
		"is_active": user.IsActive,
	})
}
