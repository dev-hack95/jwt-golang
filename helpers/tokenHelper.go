package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dev-hack95/jwt-golang/config"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Struct uses hasing mechanism to create a token from detils you have provided
type SignedDetails struct {
	Email     string
	Firstname string
	Lastname  string
	Uid       string
	UserType  string
	jwt.Claims
}

var userCollection *mongo.Collection = config.OpenCollection(config.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstname string, lastname string, userType string, uid string) (signedToken string, signedRefreshToken string) {
	claims := &SignedDetails{
		Email:     email,
		Firstname: firstname,
		Lastname:  lastname,
		Uid:       uid,
		UserType:  userType,
		Claims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(86400, 0)),
		},
	}

	refreshClaims := &SignedDetails{
		Claims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(432000, 0)),
		},
	}

	token, err1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("SECRET_KEY"))

	if err1 != nil {
		log.Panic(err1)
		return
	}

	refreshToken, err2 := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte("SECRET_KEY"))

	if err2 != nil {
		log.Panic(err2)
		return
	}

	return token, refreshToken
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if ok {
		msg = err.Error()
		return
	}

	return claims, msg

}
