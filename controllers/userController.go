package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dev-hack95/jwt-golang/config"
	helper "github.com/dev-hack95/jwt-golang/helpers"
	"github.com/dev-hack95/jwt-golang/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = config.OpenCollection(config.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Panic(err)
	}
	return string(hashedPassword)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = err.Error()
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		}

		// validate models.user
		validationError := validate.Struct(user)

		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": validationError.Error()})
		}

		count1, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		defer cancel()

		if err != nil {

			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for user email"})
		}

		if count1 > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This email is already exist"})
		}

		count2, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for phone number"})
		}

		if count2 > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This phone number is already exsist"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()
		token, refreshToken := helper.GenerateAllTokens(*user.Email, *user.Firstname, *user.Lastname, *user.UserType, user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshToken

		// insering token to database

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			msg := err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User // person want to login should be present in your database

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or Password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		//VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}

		token, refreshToken := helper.GenerateAllTokens(*foundUser.Email, *foundUser.Firstname, *foundUser.Lastname, *foundUser.UserType, foundUser.UserId)

		helper.UpdateAllTokens(token, refreshToken, foundUser.UserId)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.UserId}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		_, err := userCollection.DeleteOne(ctx, bson.M{"user_id": userId})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		//var user models.User

		err := helper.CheckUserType(c, "ADMIN")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		// For first page i want 10 records
		page, err_1 := strconv.Atoi(c.Query("page"))

		if err_1 != nil || page < 1 {
			page = 1
		}

		startIndex, err := strconv.Atoi(c.Query("startIndex"))

		if err != nil {
			startIndex = 0
		}

		matchStage := bson.D{
			{Key: "$match", Value: bson.D{{}}},
		}

		// Counting Total number of records
		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "_id", Value: "null"},
				}},
				{Key: "total_count", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
				{Key: "data", Value: bson.D{
					{Key: "$push", Value: "$$ROUTE"},
				}}},
			},
		}

		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_item", Value: bson.D{
					{Key: "$slice", Value: []interface{}{
						"$data",
						startIndex,
						recordPerPage,
					}},
				}},
			}},
		}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			groupStage,
			projectStage,
		})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing items"})
		}

		var allusers []bson.M

		err = result.All(ctx, &allusers)

		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allusers[0])
	}
}
