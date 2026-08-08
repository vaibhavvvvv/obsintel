// API key middleware
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	validKeys map[string]bool
}

func NewAuthMiddleware(keys []string) *AuthMiddleware {
	keyMap := make(map[string]bool, len(keys))
	for _, key := range keys {
		trimmedKey := strings.TrimSpace(key)
		if trimmedKey != "" {
			keyMap[trimmedKey] = true
		}
	}
	return &AuthMiddleware{validKeys: keyMap}
}

func (a *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := strings.TrimSpace(c.GetHeader("X-API-Key"))
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required."})
			return
		}

		isValid := a.validKeys[key]
		if !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key is invalid."})
			return
		}
		//api key isnt invalid so...
		c.Set("api_key", key)
		c.Next()
	}
}
