package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shubhamoys/todo-api/internal/todos/tasks_controller"
	"github.com/shubhamoys/todo-api/internal/todos/tasks_service"
	"github.com/shubhamoys/todo-api/internal/users/users_controller"
	"github.com/shubhamoys/todo-api/internal/users/users_service"
	"github.com/shubhamoys/todo-api/pkg/middleware"
)

func UserRoutes(router *gin.Engine) {
	usersService := users_service.NewUsersService()
	usersController := users_controller.NewUserController(usersService)
	userGroup := router.Group("/users")

	userGroup.POST("/register", usersController.Register)
	userGroup.POST("/login", usersController.Login)

	// Protected routes
	userGroup.Use(middleware.AuthMiddleware())
	userGroup.Use(middleware.RoleBasedAccess())
	userGroup.GET("/", usersController.GetUsers)
	userGroup.PUT("/:id", usersController.UpdateUser)
}

func TaskRoutes(router *gin.Engine) {
	tasksService := tasks_service.NewTasksService()
	tasksController := tasks_controller.NewTaskController(tasksService)
	taskGroup := router.Group("/tasks")
	taskGroup.Use(middleware.AuthMiddleware())
	taskGroup.Use(middleware.RoleBasedAccess())

	taskGroup.POST("/", tasksController.CreateTask)      // Create a new todo
	taskGroup.GET("/", tasksController.GetTasks)         // Get all todos
	taskGroup.PUT("/:id", tasksController.UpdateTask)    // Update a specific todo by Id
	taskGroup.DELETE("/:id", tasksController.DeleteTask) // Delete a specific todo by Id
}
