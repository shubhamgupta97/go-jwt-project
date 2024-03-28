package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shubhamgupta97/go-jwt-project/pkg/handlers"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/users/signup", handlers.SignUp)
	router.POST("/users/login", handlers.Login)
}
