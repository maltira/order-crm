package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCode := c.MustGet("role_code")

		currentRole := roleCode.(string)

		for _, allowed := range allowedRoles {
			if currentRole == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Доступ к контенту запрещён. Недостаточно прав.",
		})
	}
}

func AdminOnly() gin.HandlerFunc {
	return RoleMiddleware("admin")
}

func ManagerOrHigher() gin.HandlerFunc {
	return RoleMiddleware("admin", "manager")
}
