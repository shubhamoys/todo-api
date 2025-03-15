package routes

import (
	"github.com/gin-gonic/gin"
	usersController "github.com/shubhamoys/todo-api/internal/users/users_controller"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")

	userGroup.POST("/register", usersController.Register)
	userGroup.POST("/login", usersController.Login)
}

func TodoRoutes(router *gin.Engine) {
	// todoGroup := router.Group("/todos")

	// todoGroup.POST("/", todoController.CreateTodo)    // Create a new todo
	// todoGroup.GET("/", todoController.GetTodos)       // Get all todos
	// todoGroup.GET("/:id", todoController.GetTodo)     // Get a specific todo by ID
	// todoGroup.PUT("/:id", todoController.UpdateTodo)  // Update a specific todo by ID
	// todoGroup.DELETE("/:id", todoController.DeleteTodo) // Delete a specific todo by ID
}
