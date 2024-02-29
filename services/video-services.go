package services

import (
	"errors"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/mashingan/smapping"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoService interface {
	CreateCategory(category dto.CreateVideoCategoriesDTO) (dto.GetVideoCategoriesDTO, error)
	UpdateCategory(category dto.CreateVideoCategoriesDTO) (dto.GetVideoCategoriesDTO, error)
	GetAllCategory() ([]dto.GetVideoCategoriesDTO, error)
	DeleteCategory(categoryId primitive.ObjectID) error
	DuplicateCategory(categoryName string) (bool, error)

	AddVideo(video dto.CreateVideosDTO, file multipart.File) error
	GetAllVideos() ([]dto.GetVideosDTO, error)
	UpdateVideo(video dto.UpdateVideoDTO) error
	DeleteVideo(videoId primitive.ObjectID) error

	FetchInActiveVideos() ([]dto.GetVideosDTO, error)
	ActiveVideo(video_id primitive.ObjectID) error
	IncreaseDownloadCount(video_id primitive.ObjectID) error
	GetVideoByID(videoId primitive.ObjectID) (dto.GetVideosDTO, error)
}

type videocategoriesservice struct {
	repo repositories.VideoRepository
}

func NewVideoCategoriesService(repo repositories.VideoRepository) VideoService {
	return &videocategoriesservice{
		repo: repo,
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

	smapping.FillStruct(&newCategory, smapping.MapFields(res))

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
	smapping.FillStruct(&newCategory, smapping.MapFields(result))
	return newCategory, nil
}

func (ser *videocategoriesservice) GetAllCategory() ([]dto.GetVideoCategoriesDTO, error) {
	res, err := ser.repo.GetAllCategory()

	if err != nil {
		return []dto.GetVideoCategoriesDTO{}, err
	}

	allCategory := []dto.GetVideoCategoriesDTO{}

	for i := range res {
		temp := dto.GetVideoCategoriesDTO{}
		smapping.FillStruct(&temp, smapping.MapFields(res[i]))

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

func (ser *videocategoriesservice) AddVideo(video dto.CreateVideosDTO, file multipart.File) error {
	videoToCreate := models.Videos{}

	if smpErr := smapping.FillStruct(&videoToCreate, smapping.MapFields(video)); smpErr != nil {
		return smpErr
	}

	_, isCatErr := ser.repo.GetCategoryById(videoToCreate.VideoCategoriesID)

	if isCatErr != nil {
		return isCatErr
	}

	err := ser.repo.AddVideo(videoToCreate, file)
	if err != nil {
		return err
	}

	return nil
}

func (ser *videocategoriesservice) GetAllVideos() ([]dto.GetVideosDTO, error) {
	res, err := ser.repo.GetAllVideos()

	if err != nil {
		return []dto.GetVideosDTO{}, nil
	}

	allVideos := []dto.GetVideosDTO{}

	if len(res) == 0 {
		return []dto.GetVideosDTO{}, errors.New("videos not availabel")
	}

	for i := range res {
		temp := dto.GetVideosDTO{}
		smapping.FillStruct(&temp, smapping.MapFields(res[i]))
		videoPath := "http://localhost:5000/static/" + temp.VideoPath
		temp.VideoPath = videoPath
		temp.DownloadCount = res[i].DownloadCount
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
		var videoPath string
		if !strings.Contains(temp.VideoPath, "static") {
			videoPath = "http://localhost:5000/static/" + temp.VideoPath
		} else {
			videoPath = "http://localhost:5000/" + temp.VideoPath
		}

		temp.VideoPath = videoPath
		temp.DownloadCount = result[i].DownloadCount
		inActiveVideos = append(inActiveVideos, temp)
	}

	return inActiveVideos, nil
}

func (ser *videocategoriesservice) ActiveVideo(video_id primitive.ObjectID) error {
	return ser.repo.ActiveVideo(video_id)
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
