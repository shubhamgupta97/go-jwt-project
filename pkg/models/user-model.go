package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    string             `bson:"firstName" json:"firstName" validate:"required,min=2,max=100"`
	LastName     string             `bson:"lastName" json:"lastName" validate:"required,min=2,max=100"`
	Password     string             `bson:"password" json:"password" validate:"required,min=6"`
	Email        string             `bson:"email" json:"email" validate:"required,email"`
	Phone        string             `bson:"phone" json:"phone" validate:"required,min=10,max=10"`
	Token        string             `bson:"token" json:"token"`
	UserType     string             `bson:"userType" json:"userType" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken string             `bson:"refreshToken" json:"refreshToken"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	UserId       string             `bson:"userId" json:"userId"`
}
