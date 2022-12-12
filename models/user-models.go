package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	UserName     string             `json:"username" bson:"username" `
	MobileNumber string             `json:"mobile" bson:"contact"`
	UserType     string             `json:"user_type" bson:"user_type"`
	CreatedAt    primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt    primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

type AdminUser struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	UserName     string             `json:"username" bson:"username" `
	MobileNumber string             `json:"mobile" bson:"contact"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	UserType     string             `json:"user_type" bson:"user_type"`
	CreatedAt    primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt    primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
