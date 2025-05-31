package tasks_service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/shubhamoys/todo-api/constants"
	"github.com/shubhamoys/todo-api/db"
	"github.com/shubhamoys/todo-api/internal/todos/models"
	"github.com/shubhamoys/todo-api/internal/todos/task_inputs"
	"github.com/shubhamoys/todo-api/pkg/mongodb"
	"github.com/shubhamoys/todo-api/utils"
	"go.mongodb.org/mongo-driver/bson"
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

func (s *TasksService) GetTasks(query task_inputs.GetTasksQuery) (map[string]interface{}, int, error) {
	readQuery := bson.M{}

	// Filters
	if query.TaskId != "" {
		id, _ := primitive.ObjectIDFromHex(query.TaskId)
		readQuery["_id"] = id
	}

	if query.TaskIds != "" {
		ids := strings.Split(query.TaskIds, ",")
		var objIds []primitive.ObjectID
		for _, id := range ids {
			if oid, err := primitive.ObjectIDFromHex(id); err == nil {
				objIds = append(objIds, oid)
			}
		}
		readQuery["_id"] = bson.M{"$in": objIds}
	}

	if query.UserId != "" {
		id, _ := primitive.ObjectIDFromHex(query.UserId)
		readQuery["userId"] = id
	}

	if query.Name != "" {
		readQuery["name"] = bson.M{"$regex": primitive.Regex{Pattern: query.Name, Options: "i"}}
	}

	if query.Description != "" {
		readQuery["description"] = query.Description
	}

	if query.Status != "" {
		readQuery["status"] = query.Status
	}

	if query.Search != "" {
		readQuery["$text"] = bson.M{"$search": query.Search}
	}

	// Create pipeline builder
	pipelineBuilder := mongodb.NewPipelineBuilder()

	// Add match stage with filters
	pipelineBuilder.AddMatch(readQuery)

	// Add sort stage
	sort := bson.D{{Key: "createdAt", Value: -1}}
	switch query.Sort {
	case "mrc":
		sort = bson.D{{Key: "createdAt", Value: -1}}
	case "mru":
		sort = bson.D{{Key: "updatedAt", Value: -1}}
	case "namea":
		sort = bson.D{{Key: "name", Value: 1}}
	case "named":
		sort = bson.D{{Key: "name", Value: -1}}
	}
	pipelineBuilder.AddSort(sort)

	// Add pagination
	limit := query.Limit
	if limit == 0 {
		limit = 20
	}
	page := query.Page
	if page == 0 {
		page = 1
	}
	skip := (page - 1) * limit
	pipelineBuilder.AddPagination(skip, limit)

	// Add projection
	projection := bson.M{}
	if query.Fields != "" {
		fields := strings.Split(query.Fields, ",")
		for _, field := range fields {
			projection[field] = 1
		}
	}
	pipelineBuilder.AddProjection(projection)

	// Add lookup if populate is requested
	if query.Populate != "" {
		fields := strings.Split(query.Populate, ",")
		for _, field := range fields {
			switch field {
			case "user":
				pipelineBuilder.AddLookup(mongodb.LookupConfig{
					From:         "users",
					LocalField:   "userId",
					ForeignField: "_id",
					As:           "user",
				})
				// Example of how to add more lookups
				// case "example":
				// 	pipelineBuilder.AddLookup(mongodb.LookupConfig{
				// 		From:         "examples",
				// 		LocalField:   "exampleId",
				// 		ForeignField: "_id",
				// 		As:           "example",
				// 	})
			}
		}
	}

	// Execute pipeline
	cursor, err := s.collection.Aggregate(context.TODO(), pipelineBuilder.Build())
	if err != nil {
		utils.Logger.Error("Failed to execute aggregation pipeline:", err)
		return nil, http.StatusInternalServerError, errors.New("failed to fetch tasks")
	}
	defer cursor.Close(context.TODO())

	var tasks []models.Task
	if err = cursor.All(context.TODO(), &tasks); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Get total count (without pagination)
	count, err := s.collection.CountDocuments(context.TODO(), readQuery)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return map[string]interface{}{
		"totalCount":   count,
		"currentCount": len(tasks),
		"tasks":        tasks,
	}, http.StatusOK, nil
}
