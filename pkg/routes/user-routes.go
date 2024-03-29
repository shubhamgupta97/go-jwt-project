package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shubhamgupta97/go-jwt-project/pkg/handlers"
	"github.com/shubhamgupta97/go-jwt-project/pkg/middleware"
)

func UserRoutes(router *gin.Engine) {
	router.Use(middleware.Authenticate)
	router.GET("/users", handlers.GetUsers)
	router.GET("/users/:userId", handlers.GetUserById)
}
