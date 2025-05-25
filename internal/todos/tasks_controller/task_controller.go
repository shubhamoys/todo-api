package tasks_controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shubhamoys/todo-api/constants"
	create_tasko_dto "github.com/shubhamoys/todo-api/internal/todos/dtos/create_task_dto"
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

		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err, nil)
		return
	}

	// Set 2: Validate struct fields
	if err := createTaskDTO.Validate(); err != nil {
		utils.Logger.Warn("Validation failed :", err)

		utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed. Incorrect task", err, nil)
		return
	}

	// Restrict user from creating task for other users
	userIdStr := userId.(string)
	userObjId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err, nil)
		return
	}

	// Set default status if not provided
	status := createTaskDTO.Status
	if status == "" {
		status = constants.TaskStatuses.Todo
	} else if !constants.IsValidTaskStatus(status) {
		utils.Logger.Warn("Invalid task status:", status)
		utils.ErrorResponse(c, http.StatusBadRequest,
			"Invalid task status. Must be one of: Todo, In Progress, Complete",
			errors.New("invalid task status"),
			nil)
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
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
}
