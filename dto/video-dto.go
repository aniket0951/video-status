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
	IsVerified        bool               `json:"is_verified" `
	IsPublished       bool               `json:"is_published" `
	CreatedAt         primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt         primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

type CreateVideoVerificationDTO struct {
	VideoId            primitive.ObjectID `json:"video_id" validate:"required"`
	UserId             primitive.ObjectID `json:"user_id" validate:"required"`
	VerificationStatus string             `json:"verification_status" validate:"required"`
}

type UpdateVideoVerificationDTO struct {
	ID                 primitive.ObjectID `json:"id" validate:"required"`
	VideoId            primitive.ObjectID `json:"video_id"`
	UserId             primitive.ObjectID `json:"user_id"`
	VerificationStatus string             `json:"verification_status"`
	Reason             string             `json:"reason"`
}

type GetVideoVerificationDTO struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VideoId            primitive.ObjectID `json:"video_id" bson:"video_id"`
	UserId             primitive.ObjectID `json:"user_id" bson:"user_id"`
	VerificationStatus string             `json:"verification_status" bson:"verification_status"`
	Reason             string             `json:"reason,omitempty" bson:"reason"`
	CreatedAt          primitive.DateTime `json:"created_at"`
	UpdatedAt          primitive.DateTime `json:"updated_at"`
}

type CreatePublishDTO struct {
	VideoId   primitive.ObjectID `json:"video_id" validate:"required"`
	IsPublish *bool              `json:"is_publish" validate:"required"`
}

type UpdatePublishDTO struct {
	ID        primitive.ObjectID `json:"id" validate:"required"`
	VideoId   primitive.ObjectID `json:"video_id" validate:"required"`
	IsPublish *bool              `json:"is_publish" validate:"required"`
}

type GetPublishDTO struct {
	ID        primitive.ObjectID `json:"id" `
	VideoId   primitive.ObjectID `json:"video_id" `
	IsPublish *bool              `json:"is_publish"`
	CreatedAt primitive.DateTime `json:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at"`
}

type CreateVerificationNotificationDTO struct {
	Title        string             `json:"title" validate:"required"`
	Description  string             `json:"description" validate:"required"`
	IsApproved   *bool              `json:"isApproved" validate:"required"`
	VideoId      primitive.ObjectID `json:"video_id" validate:"required"`
	UserId       primitive.ObjectID `json:"user_id" validate:"required"`
	UploadedDate primitive.DateTime `json:"uploadedDate" validate:"required"`
}

type GetVerificationNotificationDTO struct {
	ID           primitive.ObjectID `json:"id" `
	Title        string             `json:"title" `
	Description  string             `json:"description" `
	IsApproved   bool               `json:"isApproved"`
	VideoId      primitive.ObjectID `json:"video_id"`
	UserId       primitive.ObjectID `json:"user_id"`
	UploadedDate primitive.DateTime `json:"uploadedDate"`
	CreatedAt    primitive.DateTime `json:"created_at"`
	UpdatedAt    primitive.DateTime `json:"updated_at"`
}
