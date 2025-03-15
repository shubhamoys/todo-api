package dtos

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shubhamoys/todo-api/constants"
)

type RegisterUserDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Validator instance
var validate = validator.New()

// Validate the struct and return user-friendly error messages
func (dto *RegisterUserDTO) Validate() error {
	err := validate.Struct(dto)
	if err != nil {
		// Convert the validation errors to user-friendly messages
		return formatValidationError(err)
	}
	return nil
}

// formatValidationError converts validator errors to user-friendly messages
func formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			
			// Remove the DTO suffix from field names
			field = strings.TrimSuffix(field, "dto")
			
			switch e.Tag() {
			case "required":
				message := constants.FormatErrorMessage(
					constants.ErrorConstants.MissingField.Message.User,
					map[string]string{"field": field},
				)
				errorMessages = append(errorMessages, message)
			case "email":
				message := constants.FormatErrorMessage(
					constants.ErrorConstants.InvalidField.Message.User,
					map[string]string{"field": field},
				)
				errorMessages = append(errorMessages, message)
			default:
				message := constants.FormatErrorMessage(
					constants.ErrorConstants.InvalidField.Message.User,
					map[string]string{"field": field},
				)
				errorMessages = append(errorMessages, message)
			}
		}
		
		return errors.New(strings.Join(errorMessages, "; "))
	}
	
	return err
}