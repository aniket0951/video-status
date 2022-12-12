package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateVideoCategoriesDTO struct {
	ID                  primitive.ObjectID `json:"id"`
	CategoryName        string             `json:"category_name" validate:"required"`
	CategoryDescription string             `json:"category_desc" validate:"required"`
	IsCategoryActive    bool               `json:"is_category_active"`
}

type GetVideoCategoriesDTO struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	CategoryName        string             `json:"category_name"`
	CategoryDescription string             `json:"category_desc"`
	IsCategoryActive    bool               `json:"is_category_active"`
	CreatedAt           primitive.DateTime `json:"created_at"`
	UpdatedAt           primitive.DateTime `json:"updated_at"`
}
