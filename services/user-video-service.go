package services

import (
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserVideoService interface {
	AddUserVideo(videoId primitive.ObjectID) error
}

type uservideoservice struct {
	repo repositories.UserVideoRepository
}

func NewUserVideoService(repo repositories.UserVideoRepository) UserVideoService {
	return &uservideoservice{
		repo: repo,
	}
}

func (ser *uservideoservice) AddUserVideo(videoId primitive.ObjectID) error {
	userVideo := models.UserVideos{}

	userId := helper.USER_ID

	objId, _ := primitive.ObjectIDFromHex(userId)

	userVideo.VideoId = videoId
	userVideo.UserId = objId
	userVideo.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	userVideo.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return ser.repo.AddUserVideo(userVideo)
}
