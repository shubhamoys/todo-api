package tasks_service

import (
	"context"
	"errors"
	"net/http"

	"github.com/shubhamoys/todo-api/constants"
	"github.com/shubhamoys/todo-api/db"
	"github.com/shubhamoys/todo-api/internal/todos/models"
	"github.com/shubhamoys/todo-api/internal/todos/task_inputs"
	"github.com/shubhamoys/todo-api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskService handles task-related operations
type TasksService struct {
	collection *mongo.Collection
}

// NewTaskService creates a new TaskService with the required dependencies
func NewTasksService() *TasksService {
	return &TasksService{
		collection: db.GetCollection("todos"),
	}
}

func (s *TasksService) CreateTask(createTaskInput task_inputs.CreateTaskInput) (*models.Task, int, error) {
	// Create a new task
	newTask, err := models.NewTask(createTaskInput)
	if err != nil {
		utils.Logger.Error("Invalid input data :", err)

		return nil, http.StatusBadRequest, errors.New("invalid input data")
	}

	// Insert the task into the database
	result, err := s.collection.InsertOne(context.TODO(), newTask)
	if err != nil {
		utils.Logger.Error("Error when inserting task :", err)

		formattedMessage := constants.FormatErrorMessage(
			constants.ErrorConstants.DatabaseError.Message.User,
			map[string]string{"message": err.Error()},
		)

		return nil, http.StatusInternalServerError, errors.New(formattedMessage)
	}

	// Update the newTask object with the generated Id
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		newTask.Id = oid
	} else {
		utils.Logger.Error("Failed to cast inserted Id to ObjectID")
	}

	utils.Logger.Info("Task registered successfully :", newTask.Name)

	return newTask, http.StatusCreated, nil
}
