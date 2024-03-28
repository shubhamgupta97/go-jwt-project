package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shubhamgupta97/go-jwt-project/pkg/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": "Access granted for api-1"})
	})

	router.POST("/api-2", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(fmt.Sprintf(":%s", port))

}
