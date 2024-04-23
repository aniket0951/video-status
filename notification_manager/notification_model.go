package notificationmanager

import "go.mongodb.org/mongo-driver/bson/primitive"

type TokenData struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token     string             `json:"token" bson:"token"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
}

type TokenRequestDTO struct {
	Token string `json:"token" validate:"required"`
}
