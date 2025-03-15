package users_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shubhamoys/todo-api/internal/users/dtos"
	"github.com/shubhamoys/todo-api/internal/users/users_service"
	"github.com/shubhamoys/todo-api/pkg/auth"
	"github.com/shubhamoys/todo-api/utils"
)

// Create a service instance
var userService *users_service.UserService

func Register(c *gin.Context) {
	if userService == nil {
		userService = users_service.NewUserService()
	}

	// Step 1: Bind and validate request body
	var registerUserDTO dtos.RegisterUserDTO
	utils.Logger.Info("User registration started for :", registerUserDTO.Email)

	if err := c.ShouldBindJSON(&registerUserDTO); err != nil {
		utils.Logger.Warn("Invalid request payload :", err)

		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err, nil)
		return
	}

	// Set 2: Validate struct fields
	if err := registerUserDTO.Validate(); err != nil {
		utils.Logger.Warn("Validation failed :", err)

		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err, nil)
		return
	}

	// Step 3: Pass the validated DTO to the service layer
	user, statusCode, err := userService.RegisterUser(registerUserDTO)
	if err != nil {
		utils.Logger.Error("Error occured when registering :", err)

		utils.ErrorResponse(c, statusCode, "User registration failed", err, nil)
		return
	}

	// Generate JWT token for the user
	token, err := auth.GenerateJWT(user.Email.Value)
	if err != nil {
		utils.Logger.Error("Unexpected error occured when generating token :", err)

		utils.ErrorResponse(c, http.StatusInternalServerError, "Unexpected error occured", err, nil)
	}

	// Step 4: Return response to the client
	utils.Logger.Info("User registration and token generation successful for :", user.Email.Value)

	utils.SuccessResponse(c, statusCode, "User registered successfully", token)
}

func Login(c *gin.Context) {
}
