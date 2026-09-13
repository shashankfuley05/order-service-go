package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shashank/order-service/internal/auth"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		fmt.Println("Authorization", authHeader)
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is missing !!!",
			})
			ctx.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token ",
			})

			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := auth.ParseToken(token)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})

			ctx.Abort()
			return
		}

		ctx.Set("userId", user)
		ctx.Next()

	}
}
