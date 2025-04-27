package get_users_dto

import "github.com/shubhamoys/todo-api/utils"

type GetUsersDTO struct {
	UserId        string `json:"userId,omitempty" form:"userId"` // Add from tag to bind the query params with struct
	UserIds       string `json:"userIds,omitempty" form:"userIds"`
	Name          string `json:"name,omitempty" form:"name"`
	EmailValue    string `json:"emailValue,omitempty" form:"emailValue"` // Add form tag
	EmailVerified *bool  `json:"emailVerified,omitempty" form:"emailVerified"`
	Role          string `json:"role,omitempty" form:"role"`
	Search        string `json:"search,omitempty" form:"search"`
	Sort          string `json:"sort,omitempty" form:"sort"`
	Limit         int64  `json:"limit,omitempty" form:"limit"`
	Page          int64  `json:"page,omitempty" form:"page"`
	Fields        string `json:"fields,omitempty" form:"fields"`
}

// Validate the struct and return user-friendly error messages
func (dto *GetUsersDTO) Validate() error {
	err := utils.Validator.Struct(dto)
	if err != nil {
		// Use the shared FormatValidationError function
		return utils.FormatValidationError(err)
	}
	return nil
}
