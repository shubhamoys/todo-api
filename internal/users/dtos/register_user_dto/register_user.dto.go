package register_user_dto

import (
	"github.com/shubhamoys/todo-api/utils"
)

type RegisterUserDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Validate the struct and return user-friendly error messages
func (dto *RegisterUserDTO) Validate() error {
	err := utils.Validator.Struct(dto)
	if err != nil {
		// Use the shared FormatValidationError function
		return utils.FormatValidationError(err)
	}
	return nil
}
