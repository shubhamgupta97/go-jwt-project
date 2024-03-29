package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shubhamgupta97/go-jwt-project/pkg/models"

	"github.com/shubhamgupta97/go-jwt-project/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

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

	hashedPassword := HashPassword(user.Password)
	user.Password = hashedPassword

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
		msg := "user item was not created"
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, resultInsertionNumber)

}

func Login(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	var foundUser models.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := userCollection.FindOne(c, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email or password is incorrect"})
		return
	}

	isValidPassword, msg := VerifyPassword(user.Password, foundUser.Password)
	if !isValidPassword {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	if foundUser.Email == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
	}

	token, refreshToken, _ := util.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.UserType, foundUser.UserId)
	util.UpdateAllTokens(token, refreshToken, foundUser.UserId)

	err = userCollection.FindOne(c, bson.M{"userId": foundUser.UserId}).Decode(&foundUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, foundUser)

}

func GetUsers(ctx *gin.Context) {
	if err := util.CheckUserType(ctx, "ADMIN"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))

	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, _ = strconv.Atoi(ctx.Query("startIndex"))

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}},
	}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "total_count", Value: 1},
		{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []any{"$data", startIndex, recordPerPage}}}}}},
	}

	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
	}

	var allUsers []bson.M
	if err = result.All(c, &allUsers); err != nil {
		log.Fatal(err)
	}

	ctx.JSON(http.StatusOK, allUsers[0])
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

func HashPassword(password string) string {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(hashedPasswordBytes)
}

func VerifyPassword(password string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(providedPassword))

	chk := true
	msg := ""

	if err != nil {
		msg = "email or password is incorrect"
		chk = false
	}

	return chk, msg
}
