package services

import (
	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/mashingan/smapping"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type VideoVerificationService interface {
	CreateVerification(verification dto.CreateVideoVerificationDTO) error
	GetAllVideosVerification() ([]dto.GetVideoVerificationDTO, error)

	CreatePublish(publish dto.CreatePublishDTO) error
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

		smapping.FillStruct(&temp, smapping.MapFields(res[i]))

		verificationData = append(verificationData, temp)
	}

	return verificationData, nil
}
func (ser *videoVerificationService) CreatePublish(publish dto.CreatePublishDTO) error {
	createPublish := models.VideoPublish{}

	createPublish.IsPublish = *publish.IsPublish
	createPublish.VideoId = publish.VideoId
	createPublish.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	createPublish.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return ser.verificationRepo.CreatePublish(createPublish)
}
func (ser *videoVerificationService) GetAllPublishData() ([]dto.GetPublishDTO, error) {

	res, err := ser.verificationRepo.GetAllPublishData()

	if err != nil {
		return nil, err
	}

	var publishData []dto.GetPublishDTO

	for i := range res {
		temp := dto.GetPublishDTO{}

		smapping.FillStruct(&temp, smapping.MapFields(res[i]))
		publishData = append(publishData, temp)
	}

	return publishData, nil
}

func (ser *videoVerificationService) CreateVerificationNotification(notification dto.CreateVerificationNotificationDTO) error {
	notificationToCreate := models.VideoVerificationNotification{}
	if smpErr := smapping.FillStruct(&notificationToCreate, smapping.MapFields(notification)); smpErr != nil {
		return smpErr
	}

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

		smapping.FillStruct(&temp, smapping.MapFields(res[i]))

		userNotifications = append(userNotifications, temp)
	}

	return userNotifications, nil
}
