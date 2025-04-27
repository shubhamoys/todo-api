package get_tasks_dto

import "github.com/shubhamoys/todo-api/utils"

type GetTasksDTO struct {
	TaskId      string `json:"taskId,omitempty" form:"taskId"` // Add from tag to bind the query params with struct
	TaskIds     string `json:"taskIds,omitempty" form:"taskIds"`
	UserId      string `json:"userId,omitempty" form:"userId"`
	Name        string `json:"name,omitempty" form:"name"`
	Description string `json:"description,omitempty" form:"description"` // Add form tag
	Status      string `json:"status,omitempty" form:"status"`
	Search      string `json:"search,omitempty" form:"search"`
	Sort        string `json:"sort,omitempty" form:"sort"`
	Limit       int64  `json:"limit,omitempty" form:"limit"`
	Page        int64  `json:"page,omitempty" form:"page"`
	Fields      string `json:"fields,omitempty" form:"fields"`
}

// Validate the struct and return user-friendly error messages
func (dto *GetTasksDTO) Validate() error {
	err := utils.Validator.Struct(dto)
	if err != nil {
		// Use the shared FormatValidationError function
		return utils.FormatValidationError(err)
	}
	return nil
}
