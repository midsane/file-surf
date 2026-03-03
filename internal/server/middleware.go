package server

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		log.Printf(
			"%s | %d | %v | %s %s",
			time.Now().Format(time.RFC3339),
			c.Writer.Status(),
			duration,
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}