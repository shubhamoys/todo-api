package users_service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/shubhamoys/todo-api/constants"
	"github.com/shubhamoys/todo-api/db"
	"github.com/shubhamoys/todo-api/internal/users/models"
	"github.com/shubhamoys/todo-api/internal/users/user_inputs"
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
	if query.UserID != "" {
		id, _ := primitive.ObjectIDFromHex(query.UserID)
		readQuery["_id"] = id
	}

	if query.UserIDs != "" {
		ids := strings.Split(query.UserIDs, ",")
		var objIDs []primitive.ObjectID
		for _, id := range ids {
			if oid, err := primitive.ObjectIDFromHex(id); err == nil {
				objIDs = append(objIDs, oid)
			}
		}
		readQuery["_id"] = bson.M{"$in": objIDs}
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

	if query.Search != "" {
		readQuery["$text"] = bson.M{"$search": query.Search}
	}

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

	// Projection
	projection := bson.M{}
	if query.Fields != "" {
		fields := strings.Split(query.Fields, ",")
		for _, field := range fields {
			projection[field] = 1
		}
	}

	findOptions := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(sort).
		SetProjection(projection)

	cursor, err := s.collection.Find(context.TODO(), readQuery, findOptions)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	count, err := s.collection.CountDocuments(context.TODO(), readQuery)
	if err != nil {
		return nil, http.StatusInternalServerError, err
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

	// Update the newUser object with the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		newUser.ID = oid
	} else {
		utils.Logger.Error("Failed to cast inserted ID to ObjectID")
	}

	utils.Logger.Info("User registered successfully :", newUser.Email.Value)

	return newUser, http.StatusCreated, nil
}
