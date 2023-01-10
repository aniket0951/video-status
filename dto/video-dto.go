package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type CreateVideosDTO struct {
	VideoTitle        string             `json:"video_title" validate:"required"`
	VideoDescription  string             `json:"video_desc" validate:"required"`
	IsVideoActive     bool               `json:"is_active" validate:"required"`
	VideoCategoriesID primitive.ObjectID `json:"video_cat_id" validate:"required"`
}

type UpdateVideoDTO struct {
	ID                primitive.ObjectID `json:"id" validate:"required"`
	VideoTitle        string             `json:"video_title" `
	VideoDescription  string             `json:"video_desc" `
	IsVideoActive     bool               `json:"is_active" `
	VideoCategoriesID primitive.ObjectID `json:"video_cat_id" `
}

type GetVideosDTO struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	VideoTitle        string             `json:"video_title" bson:"video_title"`
	VideoDescription  string             `json:"video_desc" bson:"video_desc"`
	IsVideoActive     bool               `json:"is_active" bson:"is_active"`
	VideoCategoriesID primitive.ObjectID `json:"video_cat_id" bson:"v_cat_id"`
	VideoPath         string             `json:"video_path" bson:"video_path"`
	CreatedAt         primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt         primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
