package create_tasko_dto

import (
	"github.com/shubhamoys/todo-api/utils"
)

type CreateTaskDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"omitempty,oneof=Todo 'In Progress' Complete"`
}

// Validate the struct and return user-friendly error messages
func (dto *CreateTaskDTO) Validate() error {
	err := utils.Validator.Struct(dto)
	if err != nil {
		// Use the shared FormatValidationError function
		return utils.FormatValidationError(err)
	}
	return nil
}
