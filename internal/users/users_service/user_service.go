package users_service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/shubhamoys/todo-api/constants"
	"github.com/shubhamoys/todo-api/db"
	"github.com/shubhamoys/todo-api/internal/users/models"
	"github.com/shubhamoys/todo-api/internal/users/user_inputs"
	"github.com/shubhamoys/todo-api/pkg/mongodb"
	"github.com/shubhamoys/todo-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserService handles user-related operations
type UsersService struct {
	collection *mongo.Collection
}

// NewUserService creates a new UserService with the required dependencies
func NewUsersService() *UsersService {
	return &UsersService{
		collection: db.GetCollection("users"),
	}
}

func (s *UsersService) GetUsers(query user_inputs.GetUsersQuery) (map[string]interface{}, int, error) {
	readQuery := bson.M{}

	// Filters
	if query.UserId != "" {
		id, _ := primitive.ObjectIDFromHex(query.UserId)
		readQuery["_id"] = id
	}

	if query.UserIds != "" {
		ids := strings.Split(query.UserIds, ",")
		var objIds []primitive.ObjectID
		for _, id := range ids {
			if oid, err := primitive.ObjectIDFromHex(id); err == nil {
				objIds = append(objIds, oid)
			}
		}
		readQuery["_id"] = bson.M{"$in": objIds}
	}

	if query.Name != "" {
		readQuery["name"] = bson.M{"$regex": primitive.Regex{Pattern: query.Name, Options: "i"}}
	}

	if query.EmailValue != "" {
		readQuery["email.value"] = query.EmailValue
	}

	if query.EmailVerified != nil {
		readQuery["email.verified"] = *query.EmailVerified
	}

	if query.Role != "" {
		readQuery["role"] = query.Role
	}

	if query.Search != "" {
		// Create case-insensitive regex pattern with partial matching
		searchPattern := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", regexp.QuoteMeta(query.Search)),
			Options: "i",
		}

		// Search in both name and email
		readQuery["$or"] = []bson.M{
			{"name": bson.M{"$regex": searchPattern}},
			{"email.value": bson.M{"$regex": searchPattern}},
		}
	}

	// Create pipeline builder
	pipelineBuilder := mongodb.NewPipelineBuilder()

	// Add match stage with filters
	pipelineBuilder.AddMatch(readQuery)

	// Sorting
	sort := bson.D{{Key: "timestamp.createdAt", Value: -1}}
	switch query.Sort {
	case "mrc":
		sort = bson.D{{Key: "timestamp.createdAt", Value: -1}}
	case "mru":
		sort = bson.D{{Key: "timestamp.updatedAt", Value: -1}}
	case "namea":
		sort = bson.D{{Key: "name", Value: 1}}
	case "named":
		sort = bson.D{{Key: "name", Value: -1}}
	}
	pipelineBuilder.AddSort(sort)

	// Pagination
	limit := query.Limit
	if limit == 0 {
		limit = 20
	}
	page := query.Page
	if page == 0 {
		page = 1
	}
	skip := (page - 1) * limit
	pipelineBuilder.AddPagination(skip, limit)

	// Projection
	projection := bson.M{}
	if query.Fields != "" {
		fields := strings.Split(query.Fields, ",")
		for _, field := range fields {
			projection[field] = 1
		}
	}
	pipelineBuilder.AddProjection(projection)

	// Execute pipeline
	cursor, err := s.collection.Aggregate(context.TODO(), pipelineBuilder.Build())
	if err != nil {
		utils.Logger.Error("Failed to execute aggregation pipeline:", err)
		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.DatabaseError.Message.User,
			map[string]string{"message": err.Error()},
		)
		return nil, http.StatusInternalServerError, errors.New(formattedMessage)
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.DatabaseError.Message.User,
			map[string]string{"message": err.Error()},
		)
		return nil, http.StatusInternalServerError, errors.New(formattedMessage)
	}

	// Get total count (without pagination)
	count, err := s.collection.CountDocuments(context.TODO(), readQuery)
	if err != nil {
		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.DatabaseError.Message.User,
			map[string]string{"message": "Failed to get total count"},
		)
		return nil, http.StatusInternalServerError, errors.New(formattedMessage)
	}

	return map[string]interface{}{
		"totalCount":   count,
		"currentCount": len(users),
		"users":        users,
	}, http.StatusOK, nil
}

func (s *UsersService) RegisterUser(createUserInput user_inputs.CreateUserInput) (*models.User, int, error) {
	// Check if the user already exists
	var existingUser models.User
	err := s.collection.FindOne(context.TODO(), bson.M{"email.value": createUserInput.Email}).Decode(&existingUser)
	if err == nil {
		utils.Logger.Warn("Email already exists :", createUserInput.Email)

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
	newUser, err := models.NewUser(createUserInput)
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

	// Update the newUser object with the generated Id
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		newUser.Id = oid
	} else {
		utils.Logger.Error("Failed to cast inserted Id to ObjectID")
	}

	utils.Logger.Info("User registered successfully :", newUser.Email.Value)

	return newUser, http.StatusCreated, nil
}

func (s *UsersService) UpdateUser(updateUserInput user_inputs.UpdateUserInput) (*models.User, int, error) {
	utils.Logger.Info("Updating user:", updateUserInput.Id.Hex())

	// Build update document
	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	// Only add fields that are provided
	if updateUserInput.Name != "" {
		update["$set"].(bson.M)["name"] = updateUserInput.Name
	}
	// if updateUserInput.Description != "" {
	// 	update["$set"].(bson.M)["description"] = updateUserInput.Description
	// }
	// if updateUserInput.Status != "" {
	// 	update["$set"].(bson.M)["status"] = updateUserInput.Status
	// }

	// Find and update the user
	var updatedUser models.User
	err := s.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": updateUserInput.Id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedUser)

	// Handle errors
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.Logger.Warn("User not found:", updateUserInput.Id.Hex())
			formattedMessage := constants.FormatErrorMessage(
				constants.ErrorConstants.EntityNotFound.Message.User,
				map[string]string{"entity": "User"},
			)
			return nil, http.StatusNotFound, errors.New(formattedMessage)
		}

		utils.Logger.Error("Database error while updating user:", err)
		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.DatabaseError.Message.User,
			map[string]string{"message": err.Error()},
		)
		return nil, http.StatusInternalServerError, errors.New(formattedMessage)
	}

	utils.Logger.Info("User updated successfully:", updatedUser.Id.Hex())
	return &updatedUser, http.StatusOK, nil
}
