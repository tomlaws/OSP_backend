package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminBearerAuth(rootToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fail closed if not configured.
		if strings.TrimSpace(rootToken) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		auth := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(auth, prefix))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Constant-time compare to reduce timing leaks.
		if len(token) != len(rootToken) || subtle.ConstantTimeCompare([]byte(token), []byte(rootToken)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}
