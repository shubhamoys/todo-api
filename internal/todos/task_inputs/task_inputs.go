package task_inputs

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetTasksQuery struct {
	TaskId      string
	TaskIds     string
	UserId      string
	Name        string
	Description string
	Status      string
	Search      string
	Sort        string
	Limit       int64
	Page        int64
	Fields      string
	Populate    string
}

type CreateTaskInput struct {
	UserId      primitive.ObjectID
	Name        string
	Description string
	Status      string
}

type UpdateTaskInput struct {
	Id          primitive.ObjectID
	Name        string
	Description string
	Status      string
}

type DeleteTaskInput struct {
	Id     primitive.ObjectID
	UserId primitive.ObjectID
}
