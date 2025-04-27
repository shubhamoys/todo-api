package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shubhamoys/todo-api/constants"
	"github.com/shubhamoys/todo-api/internal/users/models"
	"github.com/shubhamoys/todo-api/internal/users/user_inputs"
	"github.com/shubhamoys/todo-api/internal/users/users_service"
	"github.com/shubhamoys/todo-api/pkg/auth"
	"github.com/shubhamoys/todo-api/utils"
)

func RoleBasedAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims, exists := c.Get("user")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized access", nil, nil)
			c.Abort()
			return
		}

		claims := userClaims.(*auth.Claims)
		userId := claims.UserID

		// Get user details including role from database
		userService := users_service.NewUsersService()
		foundUser, statusCode, err := userService.GetUsers(user_inputs.GetUsersQuery{
			UserID: userId,
		})

		if err != nil {
			utils.ErrorResponse(c, statusCode, "Failed to verify user access", err, nil)
			c.Abort()
			return
		}

		users, ok := foundUser["users"].([]models.User)
		if !ok || len(users) == 0 {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User not found", nil, nil)
			c.Abort()
			return
		}

		user := users[0]
		userRole := user.Role

		// Store user role and ID in context for use in controllers
		c.Set("userRole", userRole)
		c.Set("userId", userId)

		// Handle role-based access
		switch userRole {
		case constants.UserRoles.Superadmin:
			// Superadmin can access everything
			c.Next()
			return
		case constants.UserRoles.Admin:
			// Admin can access everything except superadmin data
			c.Next()
			return
		case constants.UserRoles.User:
			// Regular users can only access their own data
			requestedUserId := c.Query("userId")
			if requestedUserId == "" {
				// If no userId is provided in query, assume they're trying to access their own data
				c.Set("userId", userId)
				c.Next()
				return
			}

			// If userId is provided, it must match their own ID
			if requestedUserId != userId {
				utils.ErrorResponse(c, http.StatusForbidden, "Access denied. You can only access your own data", errors.New("unauthorized access"), nil)
				c.Abort()
				return
			}
			c.Next()
			return
		default:
			utils.ErrorResponse(c, http.StatusForbidden, "Invalid user role", errors.New("invalid role"), nil)
			c.Abort()
			return
		}
	}
}
