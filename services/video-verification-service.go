package services

import (
	"errors"
	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/mashingan/smapping"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type VideoVerificationService interface {
	CreateVerification(verification dto.CreateVideoVerificationDTO) error
	GetAllVideosVerification() ([]dto.GetVideoVerificationDTO, error)
	ApproveOrDeniedVideo(videoId primitive.ObjectID, verificationStatus string) error
	VideosForVerification(tag string) ([]dto.GetVideoVerificationDTO, error)

	PublishedVideo(publish dto.CreatePublishDTO) error
	GetAllPublishData() ([]dto.GetPublishDTO, error)

	CreateVerificationNotification(notification dto.CreateVerificationNotificationDTO) error
	GetUserVerificationNotification(userId primitive.ObjectID) ([]dto.GetVerificationNotificationDTO, error)
}

type videoVerificationService struct {
	verificationRepo repositories.VideoVerificationRepository
}

func NewVideoVerificationService(repo repositories.VideoVerificationRepository) VideoVerificationService {
	return &videoVerificationService{
		verificationRepo: repo,
	}
}

func (ser *videoVerificationService) CreateVerification(verification dto.CreateVideoVerificationDTO) error {
	verificationToCreate := models.VideoVerification{}

	if smpErr := smapping.FillStruct(&verificationToCreate, smapping.MapFields(&verification)); smpErr != nil {
		return smpErr
	}

	verificationToCreate.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	verificationToCreate.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return ser.verificationRepo.CreateVerification(verificationToCreate)

}
func (ser *videoVerificationService) GetAllVideosVerification() ([]dto.GetVideoVerificationDTO, error) {
	res, err := ser.verificationRepo.GetAllVideosVerification()

	if err != nil {
		return nil, err
	}

	var verificationData []dto.GetVideoVerificationDTO

	for i := range res {
		temp := dto.GetVideoVerificationDTO{}

		_ = smapping.FillStruct(&temp, smapping.MapFields(res[i]))

		verificationData = append(verificationData, temp)
	}

	return verificationData, nil
}
func (ser *videoVerificationService) ApproveOrDeniedVideo(videoId primitive.ObjectID, verificationStatus string) error {

	err := ser.verificationRepo.ApproveOrDeniedVideo(videoId, verificationStatus)
	return err

}
func (ser *videoVerificationService) VideosForVerification(tag string) ([]dto.GetVideoVerificationDTO, error) {
	res, err := ser.verificationRepo.VideosForVerification(tag)

	if err != nil {
		return nil, err
	}

	if len(res) <= 0 {
		return nil, errors.New("no video available at this time to publish")
	}

	var approvedVideos []dto.GetVideoVerificationDTO

	for i := range res {
		temp := dto.GetVideoVerificationDTO{}
		smapping.FillStruct(&temp, smapping.MapFields(res[i]))

		if len(temp.VideoInfo) > 0 {
			videoPath := ""
			if strings.Contains(temp.VideoInfo[0].VideoPath, "static") {
				videoPath = "http://localhost:5000/" + temp.VideoInfo[0].VideoPath
			} else {
				videoPath = "http://localhost:5000/static/" + temp.VideoInfo[0].VideoPath
			}

			temp.VideoInfo[0].VideoPath = videoPath
		}

		approvedVideos = append(approvedVideos, temp)
	}

	return approvedVideos, nil
}

func (ser *videoVerificationService) PublishedVideo(publish dto.CreatePublishDTO) error {
	createPublish := models.VideoPublish{}

	createPublish.IsPublish = *publish.IsPublish
	createPublish.VideoId = publish.VideoId
	createPublish.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	createPublish.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return ser.verificationRepo.PublishedVideo(createPublish)
}
func (ser *videoVerificationService) GetAllPublishData() ([]dto.GetPublishDTO, error) {

	res, err := ser.verificationRepo.GetAllPublishData()

	if err != nil {
		return nil, err
	}

	var publishData []dto.GetPublishDTO

	for i := range res {
		temp := dto.GetPublishDTO{}

		_ = smapping.FillStruct(&temp, smapping.MapFields(res[i]))
		publishData = append(publishData, temp)
	}

	return publishData, nil
}

func (ser *videoVerificationService) CreateVerificationNotification(notification dto.CreateVerificationNotificationDTO) error {
	notificationToCreate := models.VideoVerificationNotification{}
	if smpErr := smapping.FillStruct(&notificationToCreate, smapping.MapFields(notification)); smpErr != nil {
		return smpErr
	}

	notificationToCreate.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	notificationToCreate.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return ser.verificationRepo.CreateVerificationNotification(notificationToCreate)
}
func (ser *videoVerificationService) GetUserVerificationNotification(userId primitive.ObjectID) ([]dto.GetVerificationNotificationDTO, error) {
	res, err := ser.verificationRepo.GetUserVerificationNotification(userId)

	if err != nil {
		return nil, err
	}

	var userNotifications []dto.GetVerificationNotificationDTO

	for i := range res {
		temp := dto.GetVerificationNotificationDTO{}

		_ = smapping.FillStruct(&temp, smapping.MapFields(res[i]))

		userNotifications = append(userNotifications, temp)
	}

	return userNotifications, nil
}
