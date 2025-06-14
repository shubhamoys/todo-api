package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shubhamoys/todo-api/config"
	"github.com/shubhamoys/todo-api/db"
	"github.com/shubhamoys/todo-api/pkg/middleware"
	"github.com/shubhamoys/todo-api/routes"
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

	mongoURI := config.GetMongoURI()
	// jwtKey := config.LoadJWTKey()
	dbName := config.AppConfig.DBName

	// log.Printf("Mongo URI :%s", mongoURI)
	// log.Printf("JWT Key :%s", config.AppConfig.JWTSecretKey)

	// Connect to the database
	db.Connect(mongoURI, dbName)

	// Ensure the database connection is closed when the main function exits
	defer db.Disconnect()

	// Set Gin mode based on environment
	switch config.AppConfig.AppEnv {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "staging":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode) // Default to debug mode
	}

	port := config.AppConfig.AppPort
	host := config.AppConfig.AppURL
	if host == "" {
		host = "127.0.0.1" // Default to localhost
	}
	if port == "" {
		port = "8080" // Default port if not specified
	}

	router := gin.Default()

	// Register the middleware logger
	router.Use(middleware.RequestLogger())

	routes.RouteGroups(router)

	// Run the server on the specified HOST and PORT
	address := fmt.Sprintf("%s:%s", host, port)
	utils.Logger.Info("Starting server at :", address)
	utils.Logger.Warn(router.Run(address))
}
