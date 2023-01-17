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
	IsVerified        bool               `json:"is_verified" bson:"is_verified"`
	IsPublished       bool               `json:"is_published" bson:"is_published"`
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

type VideoVerification struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VideoId            primitive.ObjectID `json:"video_id" bson:"video_id"`
	UserId             primitive.ObjectID `json:"user_id" bson:"user_id"`
	VerificationStatus string             `json:"verification_status" bson:"verification_status"`
	Reason             string             `json:"reason" bson:"reason"`
	VideoInfo          []Videos           `json:"video_info,omitempty" bson:"video_info"`
	CreatedAt          primitive.DateTime `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt          primitive.DateTime `json:"updated_at" bson:"updated_at,omitempty"`
}

type VideoPublish struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VideoId   primitive.ObjectID `json:"video_id" bson:"video_id"`
	IsPublish bool               `json:"is_publish" bson:"is_publish"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

type VideoVerificationNotification struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	IsApproved  string             `json:"isApproved" bson:"isApproved"`
	VideoId     primitive.ObjectID `json:"video_id" bson:"video_id"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt   primitive.DateTime `json:"updated_at" bson:"updated_at"`
}
