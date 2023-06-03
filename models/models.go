package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Firstname    *string            `json:"firstname" validate:"required,min=2,max=100"`
	Lastname     *string            `json:"lastname" validate:"required,min=2,max=100"`
	Password     *string            `json:"Password" validate:"required,min=6"`
	Email        *string            `json:"email" validate:"email,required"`
	Phone        *string            `json:"phone" validate:"required"`
	Token        *string            `json:"token"`
	UserType     *string            `json:"usertype" validate:"required,eq=ADMIN|eq=USER"` // same concept as enum in javascript
	RefreshToken *string            `json:"refresh_token"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
	UserId       string             `json:"user_id"`
}
