package user_inputs

type GetUsersQuery struct {
	UserID        string
	UserIDs       string
	Name          string
	EmailValue    string
	EmailVerified *bool
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
}

type UpdateUserInput struct {
	Id       string
	Name     string
	Email    string
	Password string
}
