package users_service

import (
	"context"
	"errors"
	"net/http"

	"github.com/shubhamoys/todo-api/constants"
	"github.com/shubhamoys/todo-api/db"
	"github.com/shubhamoys/todo-api/internal/users/dtos"
	"github.com/shubhamoys/todo-api/internal/users/models"
	"github.com/shubhamoys/todo-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService handles user-related operations
type UserService struct {
	collection *mongo.Collection
}

// NewUserService creates a new UserService with the required dependencies
func NewUserService() *UserService {
	return &UserService{
		collection: db.GetCollection("users"),
	}
}

func (s *UserService) RegisterUser(registerUserDTO dtos.RegisterUserDTO) (*models.User, int, error) {
	// Check if the user already exists
	var existingUser models.User
	err := s.collection.FindOne(context.TODO(), bson.M{"email.value": registerUserDTO.Email}).Decode(&existingUser)
	if err == nil {
		utils.Logger.Warn("Email already exists :", registerUserDTO.Email)

		// Format the error message for duplicate entity
		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.DuplicateEntity.Message.User,
			map[string]string{"entity": "email"},
		)

		return nil, http.StatusConflict, errors.New(formattedMessage)
	} else if err != mongo.ErrNoDocuments {
		utils.Logger.Error("Database error :", err)

		return nil, http.StatusInternalServerError, errors.New("database error") // Other database error
	}

	// Create a new user
	newUser, err := models.NewUser(registerUserDTO.Name, registerUserDTO.Email, registerUserDTO.Password)
	if err != nil {
		utils.Logger.Error("Invalid input data :", err)

		return nil, http.StatusBadRequest, errors.New("invalid input data")
	}

	// Insert the user into the database
	result, err := s.collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		utils.Logger.Error("Error when inserting user :", err)

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.DatabaseError.Message.User,
			map[string]string{"message": err.Error()},
		)

		return nil, http.StatusInternalServerError, errors.New(formattedMessage)
	}

	// Update the newUser object with the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		newUser.ID = oid
	} else {
		utils.Logger.Error("Failed to cast inserted ID to ObjectID")
	}

	utils.Logger.Info("User registered successfully :", newUser.Email.Value)

	return newUser, http.StatusCreated, nil
}
