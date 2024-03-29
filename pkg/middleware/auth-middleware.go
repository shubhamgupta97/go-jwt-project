package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shubhamgupta97/go-jwt-project/pkg/util"
)

func Authenticate(ctx *gin.Context) {
	clientToken := ctx.Request.Header.Get("token")
	if clientToken == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no authorization header provied"})
		ctx.Abort()
		return
	}

	claims, errMsg := util.ValidateToken(clientToken)
	if errMsg != "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		ctx.Abort()
		return
	}

	ctx.Set("email", claims.Email)
	ctx.Set("firstName", claims.FirstName)
	ctx.Set("lastName", claims.LastName)
	ctx.Set("userId", claims.UserId)
	ctx.Set("userType", claims.UserType)
	ctx.Next()
}
