package config

import (
	"fmt"
	"os"

	"github.com/shubhamoys/todo-api/utils"
)

type Config struct {
	AppURL             string
	AppPort            string
	AppEnv             string
	AllowedOrigins     string
	AllowedMethods     string
	AllowedHeaders     string
	AllowCredentials   bool
	DBHost             string
	DBPort             string
	DBName             string
	DBUser             string
	DBPassword         string
	JWTSecretKey       []byte
	SuperadminName     string
	SuperadminEmail    string
	SuperadminPassword string
}

var AppConfig Config

func LoadConfig() {
	// Get the MongoDB connection details from environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		utils.Logger.Error("JWT_SECRET environment variable not set")
	}

	AppConfig = Config{
		AppURL:             os.Getenv("APP_URL"),
		AppEnv:             os.Getenv("APP_ENV"),
		AppPort:            os.Getenv("APP_PORT"),
		AllowedOrigins:     getEnvOrDefault("CORS_ALLOWED_ORIGINS", "*"),
		AllowedMethods:     getEnvOrDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		AllowedHeaders:     getEnvOrDefault("CORS_ALLOWED_HEADERS", "Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,accept,origin,Cache-Control,X-Requested-With"),
		AllowCredentials:   os.Getenv("CORS_ALLOW_CREDENTIALS") != "false", // defaults to true
		DBHost:             getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:             getEnvOrDefault("DB_PORT", "27017"),
		DBName:             getEnvOrDefault("DB_NAME", "todo_app"),
		DBUser:             getEnvOrDefault("DB_USERNAME", ""),
		DBPassword:         getEnvOrDefault("DB_PASSWORD", ""),
		JWTSecretKey:       []byte(jwtSecret),
		SuperadminName:     os.Getenv("SUPERADMIN_NAME"),
		SuperadminEmail:    os.Getenv("SUPERADMIN_EMAIL"),
		SuperadminPassword: os.Getenv("SUPERADMIN_PASSWORD"),
	}
}

// Helper function to get environment variable with fallback
func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GetMongoURI() string {
	// Create the MongoDB URI
	var mongoURI string
	if AppConfig.DBUser != "" && AppConfig.DBPassword != "" {
		// Use authentication if credentials exist
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", AppConfig.DBUser, AppConfig.DBPassword, AppConfig.DBHost, AppConfig.DBPort, AppConfig.DBName)
	} else {
		// No authentication, but still specify the database
		mongoURI = fmt.Sprintf("mongodb://%s:%s/%s", AppConfig.DBHost, AppConfig.DBPort, AppConfig.DBName)
	}

	return mongoURI
}
