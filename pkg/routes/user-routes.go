package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shubhamgupta97/go-jwt-project/pkg/handlers"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/users", handlers.GetUsers)
	router.GET("/users/:userId", handlers.GetUserById)
}
