package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username        *string            `json:"username" validate:"required,min=3,max=50"`
	Email           *string            `json:"email" validate:"required,email"`
	Password        *string            `json:"password" validate:"required,min=6,max=50"`
	ConfirmPassword *string            `json:"confirmPassword" validate:"required,min=6,max=50"`
	Gender          *string            `json:"gender"`
	Token           *string            `json:"token"`
	User_type       *string            `json:"user_type" validate:"required,oneof=USER ADMIN"`
	Refresh_token   *string            `json:"refresh_token"`
	Created_at      time.Time          `json:"created_at"`
	Updated_at      time.Time          `json:"updated_at"`
}
