package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequiredRole(requiredRoles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawRoles, ok := ctx.Get("roles")

		if !ok {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "User not authorized",
			})

			ctx.Abort()
			return
		}

		roles, ok := rawRoles.([]string)

		if !ok {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid user roles",
			})

			ctx.Abort()
			return
		}

		for _, role := range roles {
			for _, requiredRole := range requiredRoles {
				if role == requiredRole {
					ctx.Next()
					return
				}
			}
		}

		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Required user role not found",
		})

		ctx.Abort()

	}
}
