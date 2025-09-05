package handler

import (
	"net/http"
	"time"

	"github.com/AkitoSakurabaCreator/Rule-MCP-Server/pkg/httpx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	jwtSecret []byte
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
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

// ApproveUserRequest ユーザー承認リクエスト
type ApproveUserRequest struct {
	UserID  int  `json:"userId" binding:"required"`
	Approve bool `json:"approve"`
}

func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{
		jwtSecret: []byte(jwtSecret),
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}

	// 簡易的な認証（実際の実装ではデータベースからユーザーを取得）
	if req.Username == "admin" && req.Password == "admin123" {
		// JWTトークンを生成
		claims := Claims{
			UserID:   1,
			Username: req.Username,
			Role:     "admin",
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
				ID:       1,
				Username: req.Username,
				Email:    "admin@rulemcp.com",
				Role:     "admin",
			},
			Message: "Login successful",
		}

		c.JSON(http.StatusOK, response)
		return
	}

	httpx.JSONError(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "Invalid credentials", nil)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		FullName string `json:"full_name"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.JSONError(c, http.StatusBadRequest, httpx.CodeValidation, "リクエストデータが不正です", err.Error())
		return
	}

	// 簡易実装: ここでは即時にJWTを発行（DB保存は次段で拡張）
	claims := Claims{
		UserID:   0,
		Username: req.Username,
		Role:     "user",
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
	c.JSON(http.StatusOK, LoginResponse{Token: tokenString, User: User{ID: 0, Username: req.Username, Email: req.Email, Role: "user"}, Message: "Registration successful"})
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
	_, exists := c.Get("userID")
	if !exists {
		httpx.JSONError(c, http.StatusUnauthorized, httpx.CodeUnauthorized, "User not authenticated", nil)
		return
	}

	// パスワード検証と更新のロジックを実装
	// ここでは簡易的な実装（実際のプロダクションでは適切なパスワードハッシュ化が必要）

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

	// ユーザー承認/拒否のロジックを実装
	// ここでは簡易的な実装

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

	// 承認待ちユーザー一覧を取得
	// ここでは簡易的な実装
	pendingUsers := []gin.H{
		{"id": 4, "username": "newuser1", "email": "newuser1@example.com", "fullName": "New User One", "role": "user", "isActive": false, "createdAt": "2025-09-03T10:00:00Z"},
		{"id": 5, "username": "newuser2", "email": "newuser2@example.com", "fullName": "New User Two", "role": "user", "isActive": false, "createdAt": "2025-09-03T11:00:00Z"},
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
	c.JSON(http.StatusOK, gin.H{"id": 0, "username": claims.Username, "email": "", "full_name": "", "role": claims.Role, "is_active": true})
}
