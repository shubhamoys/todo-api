package config

import (
	"fmt"
	"os"

	"github.com/shubhamoys/todo-api/utils"
)

type Config struct {
	AppPort      string
	AppEnv       string
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPassword   string
	JWTSecretKey []byte
}

var AppConfig Config

func LoadConfig() {
	// Get the MongoDB connection details from environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		utils.Logger.Error("JWT_SECRET environment variable not set")
	}

	AppConfig = Config{
		AppEnv:       os.Getenv("APP_ENV"),
		AppPort:      os.Getenv("APP_PORT"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBName:       os.Getenv("DB_NAME"),
		DBUser:       os.Getenv("DB_USERNAME"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		JWTSecretKey: []byte(jwtSecret),
	}
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
