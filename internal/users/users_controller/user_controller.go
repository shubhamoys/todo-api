package users_controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shubhamoys/todo-api/constants"
	"github.com/shubhamoys/todo-api/internal/users/dtos/get_users_dto"
	"github.com/shubhamoys/todo-api/internal/users/dtos/login_user_dto"
	"github.com/shubhamoys/todo-api/internal/users/dtos/register_user_dto"
	"github.com/shubhamoys/todo-api/internal/users/models"
	"github.com/shubhamoys/todo-api/internal/users/user_inputs"
	"github.com/shubhamoys/todo-api/internal/users/users_service"
	"github.com/shubhamoys/todo-api/pkg/auth"
	"github.com/shubhamoys/todo-api/utils"
)

// UserController handles user-related operations
type UserController struct {
	UserService *users_service.UsersService
}

// NewUserController creates a new UserController with the given UserService
func NewUserController(userService *users_service.UsersService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) Register(c *gin.Context) {
	// Step 1: Bind and validate request body
	var registerUserDTO register_user_dto.RegisterUserDTO
	utils.Logger.Info("User registration started for :", registerUserDTO.Email)

	if err := c.ShouldBindJSON(&registerUserDTO); err != nil {
		utils.Logger.Warn("Invalid request payload :", err)

		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err, nil)
		return
	}

	// Set 2: Validate struct fields
	if err := registerUserDTO.Validate(); err != nil {
		utils.Logger.Warn("Validation failed :", err)

		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed. Incorrect email or password", err, nil)
		return
	}

	registerUserInput := user_inputs.CreateUserInput{
		Name:     registerUserDTO.Name,
		Email:    registerUserDTO.Email,
		Password: registerUserDTO.Password,
	}

	// Step 3: Pass the validated DTO to the service layer
	user, statusCode, err := uc.UserService.RegisterUser(registerUserInput)
	if err != nil {
		utils.Logger.Error("Error occured when registering :", err)

		utils.ErrorResponse(c, statusCode, "User registration failed", err, nil)
		return
	}

	// Generate JWT token for the user
	token, err := auth.GenerateJWT(user.Id.Hex())
	if err != nil {
		utils.Logger.Error("Unexpected error occured when generating token :", err)

		utils.ErrorResponse(c, http.StatusInternalServerError, "Unexpected error occured", err, nil)
	}

	// Step 4: Return response to the client
	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	utils.Logger.Info("User registration and token generation successful for :", user.Email.Value)

	utils.SuccessResponse(c, statusCode, "User registered successfully", response)
}

func (uc *UserController) Login(c *gin.Context) {
	// Step 1: Bind and validate request body
	var loginUserDTO login_user_dto.LoginUserDTO
	utils.Logger.Info("User login started for :", loginUserDTO.Email)

	if err := c.ShouldBindJSON(&loginUserDTO); err != nil {
		utils.Logger.Warn("Invalid request payload :", err)

		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err, nil)
		return
	}

	// Set 2: Validate struct fields
	if err := loginUserDTO.Validate(); err != nil {
		utils.Logger.Warn("Validation failed :", err)

		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed. Incorrect email or password", err, nil)
		return
	}

	getUserInput := user_inputs.GetUsersQuery{
		EmailValue: loginUserDTO.Email,
		// Fields:     "password,email.value",
	}

	// Step 3: Pass the validated DTO to the service layer
	foundUser, statusCode, err := uc.UserService.GetUsers(getUserInput)
	if err != nil {
		utils.Logger.Error("Error occured during user login :", err)

		utils.ErrorResponse(c, statusCode, "User login failed", err, nil)
		return
	}

	// Check if any users were found
	users, ok := foundUser["users"].([]models.User)
	if !ok || len(users) == 0 {
		utils.Logger.Warn("No user found with email:", loginUserDTO.Email)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", errors.New("invalid email or password"), nil)
		return
	}
	passwordHash := users[0].Password.Hash

	if !models.ComparePassword(passwordHash, loginUserDTO.Password) {
		utils.Logger.Warn("Invalid password for user with email:", loginUserDTO.Email)
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", errors.New("invalid email or password"), nil)
		return
	}

	// Generate JWT token for the user
	token, err := auth.GenerateJWT(users[0].Id.Hex())
	if err != nil {
		utils.Logger.Error("Unexpected error occured when generating token :", err)

		utils.ErrorResponse(c, http.StatusInternalServerError, "Unexpected error occured", err, nil)
	}

	// Step 4: Return response to the client
	response := map[string]interface{}{
		"user":  users[0],
		"token": token,
	}

	utils.Logger.Info("User login and token generation successful for :", users[0].Email.Value)

	utils.SuccessResponse(c, statusCode, "User login successful", response)
}

func (uc *UserController) GetUsers(c *gin.Context) {
	// Get user role and Id from context set by RoleBasedAccess middleware
	userRole, _ := c.Get("userRole")
	userId, _ := c.Get("userId")

	// Step 1: Bind and validate request query parameters
	var getUsersDTO get_users_dto.GetUsersDTO
	utils.Logger.Info("Fetching Users started")

	if err := c.ShouldBindQuery(&getUsersDTO); err != nil {
		utils.Logger.Warn("Invalid query parameters:", err)

		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err, nil)
		return
	}

	// Step 2: Validate the DTO
	if err := getUsersDTO.Validate(); err != nil {
		utils.Logger.Warn("Validation failed:", err)

		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err, nil)
		return
	}

	// Step 3: Map DTO to service query input
	getUsersQuery := user_inputs.GetUsersQuery{
		UserId:        getUsersDTO.UserId,
		UserIds:       getUsersDTO.UserIds,
		Name:          getUsersDTO.Name,
		EmailValue:    getUsersDTO.EmailValue,
		EmailVerified: getUsersDTO.EmailVerified,
		Search:        getUsersDTO.Search,
		Sort:          getUsersDTO.Sort,
		Limit:         getUsersDTO.Limit,
		Page:          getUsersDTO.Page,
		Fields:        getUsersDTO.Fields,
	}

	// Restrict normal users from fetching other users
	if userRole == constants.UserRoles.User {
		getUsersQuery.UserId = userId.(string)
	}

	// Step 4: Call the service layer
	foundUsers, statusCode, err := uc.UserService.GetUsers(getUsersQuery)
	if err != nil {
		utils.Logger.Error("Error occurred while fetching users:", err)

		utils.ErrorResponse(c, statusCode, "Failed to fetch users", err, nil)
		return
	}

	// Step 5: Return the response
	utils.Logger.Info("Users fetched successfully")

	utils.SuccessResponse(c, http.StatusOK, "Users fetched successfully", foundUsers)
}
