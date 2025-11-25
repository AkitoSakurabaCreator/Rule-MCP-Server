package httpx

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

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
				log.Printf("PANIC RECOVERED: %v", r)
				log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
				log.Printf("RemoteAddr: %s", c.Request.RemoteAddr)
				// スタックトレースを出力
				log.Printf("Stack trace:\n%s", debug.Stack())
				// エラー詳細を出力
				if err, ok := r.(error); ok {
					log.Printf("Error details: %+v", err)
				}
				JSONError(c, http.StatusInternalServerError, CodeInternal, fmt.Sprintf("サーバ内部でエラーが発生しました: %v", r), nil)
			}
		}()
		c.Next()
	}
}
