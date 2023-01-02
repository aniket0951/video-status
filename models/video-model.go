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

type Videos struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	VideoTitle        string             `json:"video_title" bson:"video_title"`
	VideoDescription  string             `json:"video_desc" bson:"video_desc"`
	IsVideoActive     bool               `json:"is_active" bson:"is_active"`
	VideoCategoriesID primitive.ObjectID `json:"video_cat_id" bson:"v_cat_id"`
	VideoPath         string             `json:"video_path" bson:"video_path"`
	CreatedAt         primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt         primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

type UserVideos struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VideoId   primitive.ObjectID `json:"videoId" bson:"video_id"`
	UserId    primitive.ObjectID `json:"userId" bson:"user_id"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
