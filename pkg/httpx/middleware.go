package httpx

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
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
				// パニックの詳細をログに出力
				log.Printf("PANIC RECOVERED: %v\nRequest: %s %s\n", r, c.Request.Method, c.Request.URL.Path)
				log.Printf("Stack trace would be helpful here")
				// スタックトレースも出力できるように
				if err, ok := r.(error); ok {
					log.Printf("Error details: %+v", err)
				}
				JSONError(c, http.StatusInternalServerError, CodeInternal, fmt.Sprintf("サーバ内部でエラーが発生しました: %v", r), nil)
			}
		}()
		c.Next()
	}
}
