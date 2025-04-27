package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/shubhamoys/todo-api/config"
	"github.com/shubhamoys/todo-api/constants"
	"github.com/shubhamoys/todo-api/db"
	"github.com/shubhamoys/todo-api/internal/users/user_inputs"
	"github.com/shubhamoys/todo-api/internal/users/users_service"
	"github.com/shubhamoys/todo-api/utils"
)

func main() {
	// Initialize the logger
	utils.InitLogger()

	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		utils.Logger.Error("Error loading .env file: ", err)
	}

	// Load configuration
	config.LoadConfig()

	// Get superadmin details from config
	name := config.AppConfig.SuperadminName
	email := config.AppConfig.SuperadminEmail
	password := config.AppConfig.SuperadminPassword
	role := constants.UserRoles.Superadmin

	if name == "" || email == "" || password == "" {
		log.Println(name, email, password, role)

		log.Println("Super admin details are missing in the environment variables")
		return
	}

	// Connect to the database
	mongoURI := config.GetMongoURI()
	dbName := config.AppConfig.DBName
	db.Connect(mongoURI, dbName)
	defer db.Disconnect()

	// Create user service
	userService := users_service.NewUsersService()

	registerUserInput := user_inputs.CreateUserInput{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}

	user, statusCode, err := userService.RegisterUser(registerUserInput)
	if err != nil {
		utils.Logger.Error("Error occured when registering :", err)

		// If the error is because the user already exists, that's okay
		if statusCode == http.StatusConflict {
			log.Println("Super admin already exists.")
		} else {
			log.Printf("Failed to create super admin: %v", err)
		}
		return
	}

	log.Println("Super admin seeding process completed successfully.")
	utils.Logger.Info("Super admin created with email:", user.Email.Value)
}
