package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shubhamoys/todo-api/constants"
)

// Validator instance (shared across the application)
var Validator = validator.New()

// FormatValidationError converts validator errors to user-friendly messages
func FormatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string

		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())

			// Remove the DTO suffix from field names (optional)
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
