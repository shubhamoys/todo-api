package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shubhamoys/todo-api/config"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := strings.Split(config.AppConfig.AllowedOrigins, ",")
		allowedMethods := config.AppConfig.AllowedMethods
		allowedHeaders := config.AppConfig.AllowedHeaders
		allowCredentials := "false"
		if config.AppConfig.AllowCredentials {
			allowCredentials = "true"
		}

		// Get the origin from the request
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is allowed
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", allowCredentials)
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Methods", allowedMethods)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
