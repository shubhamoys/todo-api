package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shubhamoys/todo-api/config"
)

// CorsMiddleware handles CORS for all incoming requests
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := config.AppConfig.AllowedOrigins
		allowedMethods := config.AppConfig.AllowedMethods
		allowedHeaders := config.AppConfig.AllowedHeaders
		allowCredentials := config.AppConfig.AllowCredentials

		origin := c.Request.Header.Get("Origin")

		// If allowedOrigins is "*", allow all origins (for development, no credentials)
		if allowedOrigins == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			c.Writer.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			// Do not set Allow-Credentials if using wildcard origin
		} else {
			origins := strings.Split(allowedOrigins, ",")
			for _, allowedOrigin := range origins {
				if origin == allowedOrigin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					if allowCredentials {
						c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
					}
					break
				}
			}
			c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			c.Writer.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
