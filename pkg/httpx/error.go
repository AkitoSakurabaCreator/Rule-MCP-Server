package httpx

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"requestId,omitempty"`
	Timestamp string      `json:"timestamp"`
}

const (
	CodeValidation    = "validation_error"
	CodeUnauthorized  = "unauthorized"
	CodeForbidden     = "forbidden"
	CodeNotFound      = "not_found"
	CodeConflict      = "conflict"
	CodeUnprocessable = "unprocessable_entity"
	CodeInternal      = "internal_error"
)

func JSONError(c *gin.Context, status int, code string, message string, details interface{}) {
	requestID, _ := c.Get(ContextKeyRequestID)
	resp := ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	if rid, ok := requestID.(string); ok && rid != "" {
		resp.RequestID = rid
	}
	c.AbortWithStatusJSON(status, resp)
}

func MapError(err error) (int, string, string) {
	if err == nil {
		return http.StatusOK, "", ""
	}
	return http.StatusInternalServerError, CodeInternal, "サーバ内部でエラーが発生しました"
}

func JSONFromError(c *gin.Context, err error) {
	status, code, message, details := MapAppError(err)
	JSONError(c, status, code, message, details)
}
