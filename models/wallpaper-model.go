package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type WallPaper struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Title         string             `bson:"title"`
	FilePath      string             `bson:"file_path"`
	IsActive      bool               `bson:"is_active"`
	DownloadCount int32              `bson:"download_count"`
	ShareCount    int32              `bson:"share_count"`
	Category      string             `bson:"category"`
	CreatedAt     primitive.DateTime `bson:"created_at"`
	UpdatedAt     primitive.DateTime `bson:"updated_at"`
}
