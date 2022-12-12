package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoCategories struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	CategoryName        string             `json:"category_name" bson:"category_name"`
	CategoryDescription string             `json:"category_desc" bson:"category_desc"`
	IsCategoryActive    bool               `json:"is_category_active" bson:"is_active"`
	CreatedAt           primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt           primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
