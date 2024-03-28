package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shubhamgupta97/go-jwt-project/pkg/models"
	_ "github.com/shubhamgupta97/go-jwt-project/pkg/models"

	"github.com/shubhamgupta97/go-jwt-project/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	_ "golang.org/x/crypto/bcrypt"

	"github.com/shubhamgupta97/go-jwt-project/pkg/config"
)

var userCollection *mongo.Collection = config.OpenCollection(config.DBInstance(), "user")
var validate = validator.New()

func SignUp(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	count, err := userCollection.CountDocuments(c, bson.M{"email": user.Email})
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking email"})
	}

	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "this email or phone number already exists"})
	}

	count, err = userCollection.CountDocuments(c, bson.M{"phone": user.Phone})
	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for phone number"})
	}

	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "this email or phone number already exists"})
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt = user.CreatedAt
	user.ID = primitive.NewObjectID()
	user.UserId = user.ID.Hex()
	token, refreshToken, _ := util.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.UserType, user.UserId)
	user.Token = token
	user.RefreshToken = refreshToken

	resultInsertionNumber, insertErr := userCollection.InsertOne(c, user)
	if insertErr != nil {
		msg := fmt.Sprintf("user item was not created")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, resultInsertionNumber)

}

func Login(ctx *gin.Context) {

}

func GetUsers(ctx *gin.Context) {

}

func GetUserById(ctx *gin.Context) {
	userId := ctx.Param("userId")

	if err := util.MatchUserTypeToUid(ctx, userId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(c, bson.M{"userId": userId}).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)

}

func HashPassword() {}

func VerifyPassword() {}
