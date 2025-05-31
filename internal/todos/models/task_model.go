package models

import (
	"time"

	"github.com/shubhamoys/todo-api/internal/todos/task_inputs"
	"github.com/shubhamoys/todo-api/internal/users/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId      primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	User        *models.User       `bson:"user" json:"user,omitempty"` // Reference to User, excluded from BSON
	Name        string             `bson:"name" json:"name" validate:"required"`
	Description string             `bson:"description" json:"description" validate:"required"`
	Status      string             `bson:"status" json:"status" validate:"required"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func NewTask(input task_inputs.CreateTaskInput) (*Task, error) {
	return &Task{
		UserId:      input.UserId,
		Name:        input.Name,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
