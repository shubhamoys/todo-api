package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shubhamoys/todo-api/pkg/auth"
	"github.com/shubhamoys/todo-api/utils"
)

// AuthMiddleware checks if the token is present and valid
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized,
				"Authorization header is missing",
				errors.New("no authorization header provided"),
				nil)
			c.Abort()
			return
		}

		// Check if the token is in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid Authorization header format", errors.New("invalid Authorization header format"), nil)
			c.Abort()
			return
		}

		token := parts[1]

		// Validate the token
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", err, nil)
			c.Abort()
			return
		}

		// Store the claims in the context for use in handlers
		c.Set("user", claims)

		// Proceed to the next handler
		c.Next()
	}
}
