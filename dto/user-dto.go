package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type RegisterEndUserDTO struct {
	ID           primitive.ObjectID `json:"id,omitempty"`
	MobileNumber string             `json:"mobile"`
}

type CreateAdminUserDTO struct {
	UserName     string `json:"username" bson:"username" validate:"required"`
	MobileNumber string `json:"mobile" validate:"required,e164"`
	Email        string `json:"email"  validate:"required,email"`
	Password     string `json:"password"  validate:"required"`
	UserType     string `json:"user_type" `
}

type AdminLoginDTO struct {
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password"  validate:"required"`
}

type GetAdminUserDTO struct {
	ID           primitive.ObjectID  `json:"id,omitempty"`
	UserName     string              `json:"username"`
	MobileNumber string              `json:"mobile"`
	Email        string              `json:"email" `
	UserType     string              `json:"user_type"`
	Token        string              `json:"token,omitempty"`
	UserAddress  GetAdminUserAddress `json:"user_adddress,omitempty"`
	CreatedAt    primitive.DateTime  `json:"created_at"`
	UpdatedAt    primitive.DateTime  `json:"updated_at"`
}

type CreateAdminUserAddress struct {
	State       string             `json:"state" validate:"required"`
	City        string             `json:"city" validate:"required"`
	Addressline string             `json:"address" validate:"required"`
	UserID      primitive.ObjectID `json:"userId" validate:"required"`
}

type UpdateAdminAddressDTO struct {
	ID          primitive.ObjectID `json:"id" validate:"required"`
	State       string             `json:"state" validate:"required"`
	City        string             `json:"city" validate:"required"`
	Addressline string             `json:"address" validate:"required"`
}

type GetAdminUserAddress struct {
	ID          primitive.ObjectID `json:"id"`
	State       string             `json:"state"`
	City        string             `json:"city"`
	Addressline string             `json:"address" `
	UserID      primitive.ObjectID `json:"userId" `
	CreatedAt   primitive.DateTime `json:"created_at" `
	UpdatedAt   primitive.DateTime `json:"updated_at" `
}
