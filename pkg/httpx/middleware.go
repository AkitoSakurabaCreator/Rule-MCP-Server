package httpx

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

const ContextKeyRequestID = "requestId"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf := make([]byte, 16)
		if _, err := rand.Read(buf); err != nil {
			buf = []byte("fallback-request-id-0000")
		}
		rid := hex.EncodeToString(buf)
		c.Set(ContextKeyRequestID, rid)
		c.Header("X-Request-ID", rid)
		c.Next()
	}
}

func RecoveryJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				JSONError(c, http.StatusInternalServerError, CodeInternal, "サーバ内部でエラーが発生しました", nil)
			}
		}()
		c.Next()
	}
}
