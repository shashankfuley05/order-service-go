package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		start := time.Now()

		fmt.Printf("Request Recieved at %v\n", start)
		fmt.Printf("Received Request method %v\n", ctx.Request.Method)
		fmt.Printf("Received Request for path %v\n", ctx.Request.URL.Path)

		ctx.Next()

		responseTime := time.Since(start)

		fmt.Printf("Response time %v\n", responseTime)

		fmt.Printf("Response status %v\n", ctx.Writer.Status())

	}
}
