package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shubhamoys/todo-api/utils"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the incoming request
		startTime := time.Now()
		utils.Logger.Infof("➡️ Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)

		// Process the request
		c.Next()

		// Log the response details
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		utils.Logger.Infof("⬅️ Completed request: %s %s | Status: %d | Duration: %v",
			c.Request.Method, c.Request.URL.Path, statusCode, duration)
	}
}
