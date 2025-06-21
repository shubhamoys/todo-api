package update_user_dto

import "github.com/shubhamoys/todo-api/utils"

type UpdateUserDTO struct {
	// UserId      primitive.ObjectID `json:"userId" validate:"required"`
	Name string `json:"name" validate:"required"`
	// Description string `json:"description" validate:"required"`
	// Status      string `json:"status" validate:"required"`
}

// Validate the struct and return user-friendly error messages
func (dto *UpdateUserDTO) Validate() error {
	err := utils.Validator.Struct(dto)
	if err != nil {
		// Use the shared FormatValidationError function
		return utils.FormatValidationError(err)
	}
	return nil
}
