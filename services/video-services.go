package services

import (
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

	_, err := ser.repo.DuplicateCategory(category.CategoryName)

	if err != nil {
		return dto.GetVideoCategoriesDTO{}, err
	}

	result, resErr := ser.repo.UpdateCategory(categoryToUpdate)

	if resErr != nil {
		return dto.GetVideoCategoriesDTO{}, resErr
	}

	newCategory := dto.GetVideoCategoriesDTO{}
	smapping.FillStruct(&newCategory, smapping.MapFields(result))
	return newCategory, err
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
