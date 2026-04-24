package middleware

import (
	"net/http"
	"order-crm/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // не было префикса Bearer
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format. Use Bearer <token>"})
			return
		}

		claims, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token: " + err.Error()})
			return
		}

		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("role_id", int(claims["role_id"].(float64)))
		c.Set("role_code", claims["role_code"].(string))
		c.Next()
	}
}
