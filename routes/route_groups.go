package routes

import (
	"github.com/gin-gonic/gin"
)

func RouteGroups(router *gin.Engine) {
	UserRoutes(router)

	TaskRoutes(router)
}
