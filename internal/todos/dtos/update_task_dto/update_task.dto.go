package update_todo_dto

import (
	"github.com/shubhamoys/todo-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateTaskDTO struct {
	TaskId      primitive.ObjectID `json:"taskId" validate:"required"`
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description" validate:"required"`
	Status      string             `json:"status" validate:"required"`
}

// Validate the struct and return user-friendly error messages
func (dto *UpdateTaskDTO) Validate() error {
	err := utils.Validator.Struct(dto)
	if err != nil {
		// Use the shared FormatValidationError function
		return utils.FormatValidationError(err)
	}
	return nil
}
