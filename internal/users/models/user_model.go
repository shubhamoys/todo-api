package models

import (
	"time"

	"github.com/shubhamoys/todo-api/constants"
	"github.com/shubhamoys/todo-api/internal/users/user_inputs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Email     Email              `bson:"email" json:"email"`
	Password  Password           `bson:"password" json:"password"`
	Role      string             `bson:"role" json:"role" validate:"required"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewUser constructor to initialize a user with default values
func NewUser(input user_inputs.CreateUserInput) (*User, error) {
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Set default role to "User" if not provided or empty
	role := input.Role
	if role == "" {
		role = constants.UserRoles.User
	}

	return &User{
		Name:      input.Name,
		Email:     NewEmail(input.Email),
		Password:  NewPassword(hashedPassword),
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// HashPassword hashes a plain password
func HashPassword(plainPassword string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// ComparePassword checks if the provided password matches the hashed password
func ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
