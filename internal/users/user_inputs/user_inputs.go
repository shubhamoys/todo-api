package user_inputs

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetUsersQuery struct {
	UserId        string
	UserIds       string
	Name          string
	EmailValue    string
	EmailVerified *bool
	Role          string
	Search        string
	Sort          string
	Limit         int64
	Page          int64
	Fields        string
	Populate      string
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	Role     string
}

type UpdateUserInput struct {
	Id       primitive.ObjectID
	Name     string
	Email    string
	Password string
}
