package util

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func MatchUserTypeToUid(ctx *gin.Context, userId string) error {
	userType := ctx.GetString("userType")
	receivedId := ctx.GetString("userId")

	if userType == "USER" && receivedId != userId {
		return errors.New("unauthorized access")
	}

	if err := CheckUserType(ctx, userType); err != nil {
		return err
	}

	return nil
}

func CheckUserType(ctx *gin.Context, role string) error {
	userType := ctx.GetString("userType")

	if userType != role {
		return errors.New("unauthorized access")
	}

	return nil
}
