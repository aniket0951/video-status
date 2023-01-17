package repositories

import (
	"context"
	dbconfig "github.com/aniket0951/Chatrapati-Maharaj/db-config"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	userVideoCollection = dbconfig.GetCollection(dbconfig.DB, "user_video")
)

type UserVideoRepository interface {
	Init() (context.Context, context.CancelFunc)
	AddUserVideo(userVideo models.UserVideos) error
}

type uservideorepository struct {
	userVideoConnection *mongo.Collection
}

func NewUserVideoRepository() UserVideoRepository {
	return &uservideorepository{
		userVideoConnection: userVideoCollection,
	}
}

func (db *uservideorepository) Init() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, cancel
}

func (db *uservideorepository) AddUserVideo(userVideo models.UserVideos) error {
	ctx, cancle := db.Init()

	defer cancle()

	_, err := db.userVideoConnection.InsertOne(ctx, userVideo)

	return err

}
