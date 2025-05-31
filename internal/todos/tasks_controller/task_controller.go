package tasks_controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shubhamoys/todo-api/constants"
	create_tasko_dto "github.com/shubhamoys/todo-api/internal/todos/dtos/create_task_dto"
	"github.com/shubhamoys/todo-api/internal/todos/dtos/get_tasks_dto"
	"github.com/shubhamoys/todo-api/internal/todos/dtos/update_task_dto"
	"github.com/shubhamoys/todo-api/internal/todos/task_inputs"
	"github.com/shubhamoys/todo-api/internal/todos/tasks_service"
	"github.com/shubhamoys/todo-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskController handles task-related operations
type TaskController struct {
	TaskService *tasks_service.TasksService
}

// NewTaskController creates a new TaskController with the given TaskService
func NewTaskController(taskService *tasks_service.TasksService) *TaskController {
	return &TaskController{
		TaskService: taskService,
	}
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	// Get user Id from context set by RoleBasedAccess middleware. This will restrcict the user to create task only for themselves
	userId, _ := c.Get("userId")

	// Step 1: Bind and validate request body
	var createTaskDTO create_tasko_dto.CreateTaskDTO
	utils.Logger.Info("Task creation started :", createTaskDTO.Name)

	if err := c.ShouldBindJSON(&createTaskDTO); err != nil {
		utils.Logger.Warn("Invalid request payload :", err)

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.InvalidInput.Message.User,
			map[string]string{"message": err.Error()},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", errors.New(formattedMessage), nil)
		return
	}

	// Set 2: Validate struct fields
	if err := createTaskDTO.Validate(); err != nil {
		utils.Logger.Warn("Validation failed :", err)

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.ValidationError.Message.User,
			map[string]string{"message": err.Error()},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", errors.New(formattedMessage), nil)
		return
	}

	// Restrict user from creating task for other users
	userIdStr := userId.(string)
	userObjId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.InvalidField.Message.User,
			map[string]string{"field": "User ID"},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", errors.New(formattedMessage), nil)
		return
	}

	// Set default status if not provided
	status := createTaskDTO.Status
	if status == "" {
		status = constants.TaskStatuses.Todo
	} else if !constants.IsValidTaskStatus(status) {
		utils.Logger.Warn("Invalid task status:", status)

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.InvalidField.Message.User,
			map[string]string{"field": "status", "valid": "Todo, In Progress, Complete"},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task status", errors.New(formattedMessage), nil)
		return
	}

	createTaskInput := task_inputs.CreateTaskInput{
		Name:        createTaskDTO.Name,
		UserId:      userObjId,
		Description: createTaskDTO.Description,
		Status:      status,
	}

	// Step 3: Pass the validated DTO to the service layer
	task, statusCode, err := tc.TaskService.CreateTask(createTaskInput)
	if err != nil {
		utils.Logger.Error("Error occured when registering :", err)

		utils.ErrorResponse(c, statusCode, "Task creation failed", err, nil)
		return
	}

	// Step 4: Return response to the client
	utils.Logger.Info("Task creation successful :", task.Name)

	utils.SuccessResponse(c, statusCode, "Task created successfully", task)
}

func (tc *TaskController) GetTasks(c *gin.Context) {
	// Get user role and Id from context set by RoleBasedAccess middleware
	userRole, _ := c.Get("userRole")
	userId, _ := c.Get("userId")

	// Step 1: Bind and validate request query parameters
	var getTasksDTO get_tasks_dto.GetTasksDTO
	utils.Logger.Info("Fetching Tasks started")

	if err := c.ShouldBindQuery(&getTasksDTO); err != nil {
		utils.Logger.Warn("Invalid query parameters:", err)

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.InvalidInput.Message.User,
			map[string]string{"message": err.Error()},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", errors.New(formattedMessage), nil)
		return
	}

	// Step 2: Validate the DTO
	if err := getTasksDTO.Validate(); err != nil {
		utils.Logger.Warn("Validation failed:", err)

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.ValidationError.Message.User,
			map[string]string{"message": err.Error()},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", errors.New(formattedMessage), nil)
		return
	}

	// Step 3: Map DTO to service query input
	getTasksQuery := task_inputs.GetTasksQuery{
		TaskId:      getTasksDTO.TaskId,
		TaskIds:     getTasksDTO.TaskIds,
		UserId:      getTasksDTO.UserId,
		Name:        getTasksDTO.Name,
		Description: getTasksDTO.Description,
		Status:      getTasksDTO.Status,
		Search:      getTasksDTO.Search,
		Sort:        getTasksDTO.Sort,
		Limit:       getTasksDTO.Limit,
		Page:        getTasksDTO.Page,
		Fields:      getTasksDTO.Fields,
		Populate:    getTasksDTO.Populate,
	}

	// Restrict normal users from fetching other task of other users
	if userRole == constants.UserRoles.User {
		getTasksQuery.UserId = userId.(string)
	}

	// Step 4: Call the service layer
	foundTasks, statusCode, err := tc.TaskService.GetTasks(getTasksQuery)
	if err != nil {
		utils.Logger.Error("Error occurred while fetching tasks:", err)

		utils.ErrorResponse(c, statusCode, "Failed to fetch tasks", err, nil)
		return
	}

	// Step 5: Return the response
	utils.Logger.Info("Tasks fetched successfully")

	utils.SuccessResponse(c, http.StatusOK, "Tasks fetched successfully", foundTasks)
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	// Get user role and Id from context set by RoleBasedAccess middleware
	userRole, _ := c.Get("userRole")
	userId, _ := c.Get("userId")

	// Get taskId from URL parameter
	taskId := c.Param("id")
	if taskId == "" {
		utils.Logger.Warn("Task ID is required")

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.MissingField.Message.User,
			map[string]string{"field": "Task ID"},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Task ID is required", errors.New(formattedMessage), nil)
		return
	}

	// Validate if taskId is a valid MongoDB ObjectID
	taskObjId, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		utils.Logger.Warn("Invalid task ID format:", taskId)

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.InvalidField.Message.User,
			map[string]string{"field": "Task ID"},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID format", errors.New(formattedMessage), nil)
		return
	}

	// Additional validation: Check if task exists and belongs to user
	if userRole == constants.UserRoles.User {
		taskExists, err := tc.TaskService.VerifyTaskOwnership(taskObjId, userId.(string))
		if err != nil {
			utils.Logger.Error("Error verifying task ownership:", err)

			formattedMessage := constants.FormatErrorMessage(
				constants.ErrorConstants.DatabaseError.Message.User,
				map[string]string{"message": err.Error()},
			)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error verifying task ownership", errors.New(formattedMessage), nil)
			return
		}
		if !taskExists {
			formattedMessage := constants.FormatErrorMessage(
				constants.ErrorConstants.Unauthorized.Message.User,
				map[string]string{"resource": "task"},
			)
			utils.ErrorResponse(c, http.StatusNotFound, "Task not found or access denied", errors.New(formattedMessage), nil)
			return
		}

		utils.Logger.Debug("Task ownership verified for user:", userId)
	}

	// Step 1: Bind and validate request body
	var updateTaskDTO update_task_dto.UpdateTaskDTO
	utils.Logger.Info("Updating Task started for task:", taskId)

	if err := c.ShouldBindJSON(&updateTaskDTO); err != nil {
		utils.Logger.Warn("Invalid request payload:", err)

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.InvalidInput.Message.User,
			map[string]string{"message": err.Error()},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", errors.New(formattedMessage), nil)
		return
	}

	// Set 2: Validate struct fields
	if err := updateTaskDTO.Validate(); err != nil {
		utils.Logger.Warn("Validation failed :", err)

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.ValidationError.Message.User,
			map[string]string{"message": err.Error()},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", errors.New(formattedMessage), nil)
		return
	}

	// Set default status if not provided
	status := updateTaskDTO.Status
	if status == "" {
		status = constants.TaskStatuses.Todo
	} else if !constants.IsValidTaskStatus(status) {
		utils.Logger.Warn("Invalid task status:", status)

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.InvalidField.Message.User,
			map[string]string{"field": "status", "valid": "Todo, In Progress, Complete"},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task status", errors.New(formattedMessage), nil)
		return
	}

	// Conver userId to MongoDB ObjectID
	// userIdStr := userId.(string)
	// userObjId, err := primitive.ObjectIDFromHex(userIdStr)
	// if err != nil {
	// 	utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err, nil)
	// 	return
	// }

	updateTaskInput := task_inputs.UpdateTaskInput{
		Id: taskObjId,
		// UserId:      userObjId,
		Name:        updateTaskDTO.Name,
		Description: updateTaskDTO.Description,
		Status:      status,
	}

	// Step 3: Pass the validated DTO to the service layer
	task, statusCode, err := tc.TaskService.UpdateTask(updateTaskInput)
	if err != nil {
		utils.Logger.Error("Error occurred while updating task:", taskId, "error:", err)

		utils.ErrorResponse(c, statusCode, "Task update failed", err, nil)
		return
	}

	// Step 4: Return response to the client
	utils.Logger.Info("Task update successful:", task.Name)

	utils.SuccessResponse(c, statusCode, "Task updated successfully", task)
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	userRole, _ := c.Get("userRole")
	userId, _ := c.Get("userId")

	taskId := c.Param("id")
	if taskId == "" {
		utils.Logger.Warn("Task ID is required")
		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.MissingField.Message.User,
			map[string]string{"field": "Task ID"},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Task ID is required", errors.New(formattedMessage), nil)
		return
	}

	taskObjId, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		utils.Logger.Warn("Invalid task ID format:", taskId)
		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.InvalidField.Message.User,
			map[string]string{"field": "Task ID"},
		)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid task ID format", errors.New(formattedMessage), nil)
		return
	}

	// Check ownership for normal users
	if userRole == constants.UserRoles.User {
		taskExists, err := tc.TaskService.VerifyTaskOwnership(taskObjId, userId.(string))
		if err != nil {
			utils.Logger.Error("Error verifying task ownership:", err)
			formattedMessage := constants.FormatErrorMessage(
				constants.ErrorConstants.DatabaseError.Message.User,
				map[string]string{"message": err.Error()},
			)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error verifying task ownership", errors.New(formattedMessage), nil)
			return
		}
		if !taskExists {
			formattedMessage := constants.FormatErrorMessage(
				constants.ErrorConstants.Unauthorized.Message.User,
				map[string]string{"resource": "task"},
			)
			utils.ErrorResponse(c, http.StatusNotFound, "Task not found or access denied", errors.New(formattedMessage), nil)
			return
		}
	}

	// Call service to delete task
	statusCode, err := tc.TaskService.DeleteTask(taskObjId)
	if err != nil {
		utils.Logger.Error("Error occurred while deleting task:", taskId, "error:", err)
		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.DatabaseError.Message.User,
			map[string]string{"message": err.Error()},
		)
		utils.ErrorResponse(c, statusCode, "Task deletion failed", errors.New(formattedMessage), nil)
		return
	}

	utils.Logger.Info("Task deleted successfully:", taskId)
	utils.SuccessResponse(c, http.StatusOK, "Task deleted successfully", nil)
}
