package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"reflect"

	log "github.com/sirupsen/logrus"

	"firebase.google.com/go/messaging"
	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	notificationmanager "github.com/aniket0951/Chatrapati-Maharaj/notification_manager"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/aniket0951/Chatrapati-Maharaj/s3"
	"github.com/mashingan/smapping"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"strings"
)

var (
	userVideoRepo = repositories.NewUserVideoRepository()
)

type VideoService interface {
	CreateCategory(category dto.CreateVideoCategoriesDTO) (dto.GetVideoCategoriesDTO, error)
	UpdateCategory(category dto.CreateVideoCategoriesDTO) (dto.GetVideoCategoriesDTO, error)
	GetAllCategory() ([]dto.GetVideoCategoriesDTO, error)
	DeleteCategory(categoryId primitive.ObjectID) error
	DuplicateCategory(categoryName string) (bool, error)

	AddVideo(video dto.CreateVideosDTO, file, thumbnail multipart.File) error

	AddVideo(video dto.CreateVideosDTO, file multipart.File) (primitive.ObjectID, error)

	GetAllVideos() ([]dto.GetVideosDTO, error)
	UpdateVideo(video dto.UpdateVideoDTO) error
	UpdateVideoVerification(video models.Videos) error
	DeleteVideo(videoId primitive.ObjectID) error


	FetchInActiveVideos() ([]dto.GetVideosDTO, error)
	ActiveVideo(video_id primitive.ObjectID, isActive bool) error
	IncreaseDownloadCount(video_id primitive.ObjectID) error
	GetVideoByID(videoId primitive.ObjectID) (dto.GetVideosDTO, error)

	// after share the link to user
	GetShareVideo(fileKey string) error
}

type videocategoriesservice struct {
	repo                repositories.VideoRepository
	notificationService notificationmanager.NotificationManager
	//GetVideoById()

	VideoFullDetails(videoId primitive.ObjectID) (interface{}, error)
}

type videocategoriesservice struct {
	repo          repositories.VideoRepository
	userVideoRepo repositories.UserVideoRepository

}

func NewVideoCategoriesService(repo repositories.VideoRepository, notificationManager notificationmanager.NotificationManager) VideoService {
	return &videocategoriesservice{

		repo:                repo,
		notificationService: notificationManager,

		repo:          repo,
		userVideoRepo: userVideoRepo,
	}
}

func (ser *videocategoriesservice) CreateCategory(category dto.CreateVideoCategoriesDTO) (dto.GetVideoCategoriesDTO, error) {

	_, err := ser.repo.DuplicateCategory(category.CategoryName)

	if err != nil {
		return dto.GetVideoCategoriesDTO{}, err
	}

	categoryToCreate := models.VideoCategories{}

	if smpErr := smapping.FillStruct(&categoryToCreate, smapping.MapFields(category)); smpErr != nil {
		return dto.GetVideoCategoriesDTO{}, smpErr
	}

	res, err := ser.repo.CreateCategory(categoryToCreate)

	if err != nil {
		return dto.GetVideoCategoriesDTO{}, err
	}

	newCategory := dto.GetVideoCategoriesDTO{}

	_ = smapping.FillStruct(&newCategory, smapping.MapFields(res))

	return newCategory, nil
}

func (ser *videocategoriesservice) UpdateCategory(category dto.CreateVideoCategoriesDTO) (dto.GetVideoCategoriesDTO, error) {
	categoryToUpdate := models.VideoCategories{}

	if smpErr := smapping.FillStruct(&categoryToUpdate, smapping.MapFields(category)); smpErr != nil {
		return dto.GetVideoCategoriesDTO{}, smpErr
	}

	result, resErr := ser.repo.UpdateCategory(categoryToUpdate)

	if resErr != nil {
		return dto.GetVideoCategoriesDTO{}, resErr
	}

	newCategory := dto.GetVideoCategoriesDTO{}
	_ = smapping.FillStruct(&newCategory, smapping.MapFields(result))
	return newCategory, nil
}

func (ser *videocategoriesservice) GetAllCategory() ([]dto.GetVideoCategoriesDTO, error) {
	res, err := ser.repo.GetAllCategory()

	if err != nil {
		return []dto.GetVideoCategoriesDTO{}, err
	}

	var allCategory []dto.GetVideoCategoriesDTO

	for i := range res {
		temp := dto.GetVideoCategoriesDTO{}
		_ = smapping.FillStruct(&temp, smapping.MapFields(res[i]))

		allCategory = append(allCategory, temp)
	}

	return allCategory, nil
}

func (ser *videocategoriesservice) DeleteCategory(categoryId primitive.ObjectID) error {
	if err := ser.repo.DeleteCategory(categoryId); err != nil {
		return err
	}
	return nil
}

func (ser *videocategoriesservice) DuplicateCategory(categoryName string) (bool, error) {
	return ser.repo.DuplicateCategory(categoryName)
}


func (ser *videocategoriesservice) AddVideo(video dto.CreateVideosDTO, file, thumbnailFile multipart.File) error {

func (ser *videocategoriesservice) AddVideo(video dto.CreateVideosDTO, file multipart.File) (primitive.ObjectID, error) {
	_, isCatErr := ser.repo.GetCategoryById(video.VideoCategoriesID)

	if isCatErr != nil {
		return primitive.NewObjectID(), isCatErr
	}


	videoToCreate := models.Videos{}

	if smpErr := smapping.FillStruct(&videoToCreate, smapping.MapFields(video)); smpErr != nil {
		return primitive.NewObjectID(), smpErr
	}

	if _, isCatErr := ser.repo.GetCategoryById(videoToCreate.VideoCategoriesID); isCatErr != nil {
		return isCatErr
	}

	// save the video thumbnail locally
	tFileKey, tFilePath, s_err := helper.LocalFileWrite(thumbnailFile, "static/thumbnail", "thumbnail-*.png")

	if s_err != nil {
		return s_err
	}

	// upload a thumbnail to s3
	if err := s3.UploadFileToS3(tFilePath, tFileKey, "image/png"); err != nil {
		if err := os.Remove(tFilePath); err != nil {
			return err
		}
		return err
	}

	if err := os.Remove(tFilePath); err != nil {
		log.Println("Error while removing local file :", err)
	}

	videoToCreate.VideoThumbnail = tFileKey

	// prepare video file
	var fileContent = "video/mp4"

	fileKey, filePath, err := helper.LocalFileWrite(file, "static", "upload-*.mp4")

	res, err := ser.repo.AddVideo(videoToCreate, file)

	if err != nil {
		return primitive.NewObjectID(), err
	}


	if err := s3.UploadFileToS3(filePath, fileKey, fileContent); err != nil {
		log.Println("File Upload Error : ", err)
		return err
	}

	// remove local file
	if err := os.Remove(filePath); err != nil {
		log.Println("Error occured while remove local file : ", err)
	}

	videoToCreate.VideoPath = fileKey

	err = ser.repo.AddVideo2(videoToCreate)
	return err

	return res, nil

}

func (ser *videocategoriesservice) GetAllVideos() ([]dto.GetVideosDTO, error) {
	res, err := ser.repo.GetAllVideos()

	if err != nil {
		return []dto.GetVideosDTO{}, nil
	}

	var allVideos []dto.GetVideosDTO

	if len(res) == 0 {
		return []dto.GetVideosDTO{}, errors.New("videos not availabel")
	}

	for i := range res {
		temp := dto.GetVideosDTO{}

		smapping.FillStruct(&temp, smapping.MapFields(res[i]))

		url, err := s3.GetVideoObjectUrl(temp.VideoPath)
		if err != nil {
			log.Error("Video Object URL : ", err)
		}

		thumbnailUrl, err := s3.GetVideoThumbnailObjectUrl(temp.VideoThumbnail)

		if err != nil {
			log.Error("Thumbnail Error : ", err)
		}

		temp.VideoPath = url
		temp.VideoThumbnail = thumbnailUrl
		temp.DownloadCount = res[i].DownloadCount

		_ = smapping.FillStruct(&temp, smapping.MapFields(res[i]))
		videoPath := ""
		if strings.Contains(temp.VideoPath, "static") {
			videoPath = "http://localhost:5000/" + temp.VideoPath
		} else {
			videoPath = "http://localhost:5000/static/" + temp.VideoPath
		}

		temp.VideoPath = videoPath

		allVideos = append(allVideos, temp)
	}

	return allVideos, nil
}

func (ser *videocategoriesservice) UpdateVideo(video dto.UpdateVideoDTO) error {
	videoToUpdate := models.Videos{}

	if smpErr := smapping.FillStruct(&videoToUpdate, smapping.MapFields(video)); smpErr != nil {
		return smpErr
	}

	return ser.repo.UpdateVideo(videoToUpdate)
}

func (ser *videocategoriesservice) UpdateVideoVerification(video models.Videos) error {
	return ser.repo.UpdateVideoVerification(video)
}

func (ser *videocategoriesservice) DeleteVideo(videoId primitive.ObjectID) error {
	err := ser.repo.DeleteVideo(videoId)

	return err
}


func (ser *videocategoriesservice) FetchInActiveVideos() ([]dto.GetVideosDTO, error) {
	result, err := ser.repo.FetchInActiveVideos()

	if err != nil {
		return nil, err
	}

	var inActiveVideos []dto.GetVideosDTO

	if len(result) == 0 {
		return []dto.GetVideosDTO{}, errors.New("videos not availabel")
	}

	for i := range result {
		temp := dto.GetVideosDTO{}
		smapping.FillStruct(&temp, smapping.MapFields(result[i]))
		url, err := s3.GetVideoObjectUrl(temp.VideoPath)
		if err != nil {
			log.Error("Error while generate signed URL : ", err)
		}

		temp.VideoPath = url
		temp.DownloadCount = result[i].DownloadCount
		inActiveVideos = append(inActiveVideos, temp)
	}

	return inActiveVideos, nil
}

func (ser *videocategoriesservice) ActiveVideo(video_id primitive.ObjectID, isActive bool) error {
	err := ser.repo.ActiveVideo(video_id, isActive)

	if err != nil {
		return err
	}

	// if video is get active
	if isActive {
		// fetch video by id
		video, err := ser.repo.GetVideoByID(video_id)

		if err != nil {
			fmt.Println("Fetch Video By ID Error : ", err)
			return nil
		}
		thumbnail := "http://192.168.0.109:5000/" + video.VideoThumbnail
		notificationMessage := messaging.Message{
			Notification: &messaging.Notification{
				Title:    "Jay Bhavani !",
				Body:     video.VideoDescription,
				ImageURL: thumbnail,
			},
		}

		// notify all user
		ser.notificationService.NotifyAllUser(&notificationMessage)
	}

	return nil
}

func (ser *videocategoriesservice) IncreaseDownloadCount(video_id primitive.ObjectID) error {
	return ser.repo.IncreaseDownloadCount(video_id)
}

func (ser *videocategoriesservice) GetVideoByID(videoId primitive.ObjectID) (dto.GetVideosDTO, error) {
	result, err := ser.repo.GetVideoByID(videoId)

	if err != nil {
		return dto.GetVideosDTO{}, err
	}

	if (reflect.DeepEqual(dto.GetVideosDTO{}, result)) {
		return dto.GetVideosDTO{}, errors.New("video not found")
	}

	var video dto.GetVideosDTO
	video.ID = result.ID
	video.VideoTitle = result.VideoTitle
	video.VideoDescription = result.VideoDescription
	video.IsVideoActive = result.IsVideoActive
	videoPath := "http://localhost:5000/static/" + result.VideoPath
	video.VideoPath = videoPath
	video.DownloadCount = result.DownloadCount
	video.CreatedAt = result.CreatedAt
	video.UpdatedAt = result.UpdatedAt

	return video, nil
}

func (ser *videocategoriesservice) GetShareVideo(fileKey string) error {
	// check the file key exists or not
	isExists, err := ser.repo.IsFileKeyExists(fileKey)
	log.Infof("Is Exists : %+v and err %+v", isExists, err)
	if err != nil || !isExists {
		return err
	}
	return nil

func (ser *videocategoriesservice) VideoFullDetails(videoId primitive.ObjectID) (interface{}, error) {
	return ser.repo.VideoFullDetails(videoId)

}
